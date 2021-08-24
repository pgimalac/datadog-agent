package snmp

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/snmp/checkconfig"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/snmp/valuestore"
	"regexp"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/snmp/common"
)

type metricSender struct {
	sender           aggregator.Sender
	submittedMetrics int
}

// ReportMetrics reports metrics using sender
func (ms *metricSender) ReportMetrics(metrics []checkconfig.MetricsConfig, values *valuestore.ResultValueStore, tags []string) {
	for _, metric := range metrics {
		if metric.IsScalar() {
			ms.reportScalarMetrics(metric, values, tags)
		} else if metric.IsColumn() {
			ms.reportColumnMetrics(metric, values, tags)
		}
	}
}

func (ms *metricSender) getCheckInstanceMetricTags(metricTags []checkconfig.MetricTagConfig, values *valuestore.ResultValueStore) []string {
	var globalTags []string

	for _, metricTag := range metricTags {
		value, err := values.GetScalarValue(metricTag.OID)
		if err != nil {
			log.Debugf("metric tags: error getting scalar value: %v", err)
			continue
		}
		strValue, err := value.ToString()
		if err != nil {
			log.Debugf("error converting value (%#v) to string : %v", value, err)
			continue
		}
		globalTags = append(globalTags, metricTag.GetTags(strValue)...)
	}
	return globalTags
}

func (ms *metricSender) reportScalarMetrics(metric checkconfig.MetricsConfig, values *valuestore.ResultValueStore, tags []string) {
	value, err := values.GetScalarValue(metric.Symbol.OID)
	if err != nil {
		log.Debugf("report scalar: error getting scalar value: %v", err)
		return
	}

	scalarTags := common.CopyStrings(tags)
	scalarTags = append(scalarTags, metric.GetSymbolTags()...)
	ms.sendMetric(metric.Symbol.Name, value, scalarTags, metric.ForcedType, metric.Options, metric.Symbol.ExtractValuePattern)
}

func (ms *metricSender) reportColumnMetrics(metricConfig checkconfig.MetricsConfig, values *valuestore.ResultValueStore, tags []string) {
	rowTagsCache := make(map[string][]string)
	for _, symbol := range metricConfig.Symbols {
		metricValues, err := values.GetColumnValues(symbol.OID)
		if err != nil {
			log.Debugf("report column: error getting column value: %v", err)
			continue
		}
		for fullIndex, value := range metricValues {
			// cache row tags by fullIndex to avoid rebuilding it for every column rows
			if _, ok := rowTagsCache[fullIndex]; !ok {
				rowTagsCache[fullIndex] = append(common.CopyStrings(tags), metricConfig.GetTags(fullIndex, values)...)
				log.Debugf("report column: caching tags `%v` for fullIndex `%s`", rowTagsCache[fullIndex], fullIndex)
			}
			rowTags := rowTagsCache[fullIndex]
			ms.sendMetric(symbol.Name, value, rowTags, metricConfig.ForcedType, metricConfig.Options, symbol.ExtractValuePattern)
			ms.trySendBandwidthUsageMetric(symbol, fullIndex, values, rowTags)
		}
	}
}

func (ms *metricSender) sendMetric(metricName string, value valuestore.ResultValue, tags []string, forcedType string, options checkconfig.MetricsConfigOption, extractValuePattern *regexp.Regexp) {
	if extractValuePattern != nil {
		extractedValue, err := value.ExtractStringValue(extractValuePattern)
		if err != nil {
			log.Debugf("error extracting value from `%v` with pattern `%v`: %v", value, extractValuePattern, err)
			return
		}
		value = extractedValue
	}

	metricFullName := "snmp." + metricName
	if forcedType == "" {
		if value.SubmissionType != "" {
			forcedType = value.SubmissionType
		} else {
			forcedType = "gauge"
		}
	} else if forcedType == "flag_stream" {
		strValue, err := value.ToString()
		if err != nil {
			log.Debugf("error converting value (%#v) to string : %v", value, err)
			return
		}
		floatValue, err := getFlagStreamValue(options.Placement, strValue)
		if err != nil {
			log.Debugf("metric `%s`: failed to get flag stream value: %s", metricFullName, err)
			return
		}
		metricFullName = metricFullName + "." + options.MetricSuffix
		value = valuestore.ResultValue{Value: floatValue}
		forcedType = "gauge"
	}

	floatValue, err := value.ToFloat64()
	if err != nil {
		log.Debugf("metric `%s`: failed to convert to float64: %s", metricFullName, err)
		return
	}

	switch forcedType {
	case "gauge":
		ms.gauge(metricFullName, floatValue, "", tags)
		ms.submittedMetrics++
	case "counter":
		ms.rate(metricFullName, floatValue, "", tags)
		ms.submittedMetrics++
	case "percent":
		ms.rate(metricFullName, floatValue*100, "", tags)
		ms.submittedMetrics++
	case "monotonic_count":
		ms.monotonicCount(metricFullName, floatValue, "", tags)
		ms.submittedMetrics++
	case "monotonic_count_and_rate":
		ms.monotonicCount(metricFullName, floatValue, "", tags)
		ms.rate(metricFullName+".rate", floatValue, "", tags)
		ms.submittedMetrics += 2
	default:
		log.Debugf("metric `%s`: unsupported forcedType: %s", metricFullName, forcedType)
		return
	}
}

func (ms *metricSender) gauge(metric string, value float64, hostname string, tags []string) {
	// we need copy tags before using sender due to https://github.com/DataDog/datadog-agent/issues/7159
	ms.sender.Gauge(metric, value, hostname, common.CopyStrings(tags))
}

func (ms *metricSender) rate(metric string, value float64, hostname string, tags []string) {
	// we need copy tags before using sender due to https://github.com/DataDog/datadog-agent/issues/7159
	ms.sender.Rate(metric, value, hostname, common.CopyStrings(tags))
}

func (ms *metricSender) monotonicCount(metric string, value float64, hostname string, tags []string) {
	// we need copy tags before using sender due to https://github.com/DataDog/datadog-agent/issues/7159
	ms.sender.MonotonicCount(metric, value, hostname, common.CopyStrings(tags))
}

func (ms *metricSender) serviceCheck(checkName string, status metrics.ServiceCheckStatus, hostname string, tags []string, message string) {
	// we need copy tags before using sender due to https://github.com/DataDog/datadog-agent/issues/7159
	ms.sender.ServiceCheck(checkName, status, hostname, common.CopyStrings(tags), message)
}

func getFlagStreamValue(placement uint, strValue string) (float64, error) {
	index := placement - 1
	if int(index) >= len(strValue) {
		return 0, fmt.Errorf("flag stream index `%d` not found in `%s`", index, strValue)
	}
	charAtIndex := strValue[index]
	floatValue := 0.0
	if charAtIndex == '1' {
		floatValue = 1.0
	}
	return floatValue, nil
}
