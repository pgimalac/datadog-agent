// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package client implements helpers APIs to interact with a [fakeintake server] from go tests
// Helpers fetch fakeintake endpoints, unpackage payloads and store parsed data in aggregators
//
// # Using fakeintake in go tests
//
// In this example we assert that a fakeintake running at localhost on port 8080 received
// "system.uptime" metrics with tags "app:system" and values in range 4226000 < value < 4226050.
//
//	client := NewClient("http://localhost:8080")
//	metrics, err := client.FilterMetrics("system.uptime",
//			WithTags[*aggregator.MetricSeries]([]string{"app:system"}),
//			WithMetricValueInRange(4226000, 4226050))
//	assert.NoError(t, err)
//	assert.NotEmpty(t, metrics)
//
// In this example we assert that a fakeintake running at localhost on port 8080 received
// logs by service "system" with tags "app:system" and content containing "totoro"
//
//	client := NewClient("http://localhost:8080")
//	logs, err := client.FilterLogs("system",
//			WithTags[*aggregator.Log]([]string{"totoro"}),
//	assert.NoError(t, err)
//	assert.NotEmpty(t, logs)
//
// In this example we assert that a fakeintake running at localhost on port 8080 received
// check runs by name "totoro" with tags "status:ok"
//
//	client := NewClient("http://localhost:8080")
//	logs, err := client.GetCheckRun("totoro")
//	assert.NoError(t, err)
//	assert.NotEmpty(t, logs)
//
// [fakeintake server]: https://pkg.go.dev/github.com/DataDog/datadog-agent@main/test/fakeintake/server
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"

	agentmodel "github.com/DataDog/agent-payload/v5/process"
	"github.com/DataDog/datadog-agent/test/fakeintake/aggregator"
	"github.com/DataDog/datadog-agent/test/fakeintake/api"
	"github.com/DataDog/datadog-agent/test/fakeintake/client/flare"
)

const (
	metricsEndpoint              = "/api/v2/series"
	checkRunsEndpoint            = "/api/v1/check_run"
	logsEndpoint                 = "/api/v2/logs"
	connectionsEndpoint          = "/api/v1/connections"
	processesEndpoint            = "/api/v1/collector"
	containersEndpoint           = "/api/v1/container"
	processDiscoveryEndpoint     = "/api/v1/discovery"
	flareEndpoint                = "/support/flare"
	orchestratorEndpoint         = "/api/v2/orch"
	orchestratorManifestEndpoint = "/api/v2/orchmanif"
)

// ErrNoFlareAvailable is returned when no flare is available
var ErrNoFlareAvailable = errors.New("no flare available")

//nolint:revive // TODO(APL) Fix revive linter
type Client struct {
	fakeIntakeURL string

	metricAggregator               aggregator.MetricAggregator
	checkRunAggregator             aggregator.CheckRunAggregator
	logAggregator                  aggregator.LogAggregator
	connectionAggregator           aggregator.ConnectionsAggregator
	processAggregator              aggregator.ProcessAggregator
	containerAggregator            aggregator.ContainerAggregator
	processDiscoveryAggregator     aggregator.ProcessDiscoveryAggregator
	orchestratorAggregator         aggregator.OrchestratorAggregator
	orchestratorManifestAggregator aggregator.OrchestratorManifestAggregator
}

// NewClient creates a new fake intake client
// fakeIntakeURL: the host of the fake Datadog intake server
func NewClient(fakeIntakeURL string) *Client {
	return &Client{
		fakeIntakeURL:                  strings.TrimSuffix(fakeIntakeURL, "/"),
		metricAggregator:               aggregator.NewMetricAggregator(),
		checkRunAggregator:             aggregator.NewCheckRunAggregator(),
		logAggregator:                  aggregator.NewLogAggregator(),
		connectionAggregator:           aggregator.NewConnectionsAggregator(),
		processAggregator:              aggregator.NewProcessAggregator(),
		containerAggregator:            aggregator.NewContainerAggregator(),
		processDiscoveryAggregator:     aggregator.NewProcessDiscoveryAggregator(),
		orchestratorAggregator:         aggregator.NewOrchestratorAggregator(),
		orchestratorManifestAggregator: aggregator.NewOrchestratorManifestAggregator(),
	}
}

func (c *Client) getMetrics() error {
	payloads, err := c.getFakePayloads(metricsEndpoint)
	if err != nil {
		return err
	}
	return c.metricAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getCheckRuns() error {
	payloads, err := c.getFakePayloads(checkRunsEndpoint)
	if err != nil {
		return err
	}
	return c.checkRunAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getLogs() error {
	payloads, err := c.getFakePayloads(logsEndpoint)
	if err != nil {
		return err
	}
	return c.logAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getConnections() error {
	payloads, err := c.getFakePayloads(connectionsEndpoint)
	if err != nil {
		return err
	}
	return c.connectionAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getProcesses() error {
	payloads, err := c.getFakePayloads(processesEndpoint)
	if err != nil {
		return err
	}
	return c.processAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getContainers() error {
	payloads, err := c.getFakePayloads(containersEndpoint)
	if err != nil {
		return err
	}
	return c.containerAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getProcessDiscoveries() error {
	payloads, err := c.getFakePayloads(processDiscoveryEndpoint)
	if err != nil {
		return err
	}
	return c.processDiscoveryAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getOrchestrators() error {
	payloads, err := c.getFakePayloads(orchestratorEndpoint)
	if err != nil {
		return err
	}
	return c.orchestratorAggregator.UnmarshallPayloads(payloads)
}

func (c *Client) getOrchestratorManifests() error {
	payloads, err := c.getFakePayloads(orchestratorManifestEndpoint)
	if err != nil {
		return err
	}
	return c.orchestratorManifestAggregator.UnmarshallPayloads(payloads)
}

// GetLatestFlare queries the Fake Intake to fetch flares that were sent by a Datadog Agent and returns the latest flare as a Flare struct
// TODO: handle multiple flares / flush when returning latest flare
func (c *Client) GetLatestFlare() (flare.Flare, error) {
	payloads, err := c.getFakePayloads(flareEndpoint)
	if err != nil {
		return flare.Flare{}, err
	}

	if len(payloads) == 0 {
		return flare.Flare{}, ErrNoFlareAvailable
	}

	return flare.ParseRawFlare(payloads[len(payloads)-1])
}

func (c *Client) getFakePayloads(endpoint string) (rawPayloads []api.Payload, err error) {
	body, err := c.get(fmt.Sprintf("fakeintake/payloads?endpoint=%s", endpoint))
	if err != nil {
		return nil, err
	}
	var response api.APIFakeIntakePayloadsRawGETResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Payloads, nil
}

// GetServerHealth fetches fakeintake health status and returns an error if
// fakeintake is unhealthy
func (c *Client) GetServerHealth() error {
	resp, err := http.Get(fmt.Sprintf("%s/fakeintake/health", c.fakeIntakeURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error code %v", resp.StatusCode)
	}
	return nil
}

// ConfigureOverride sets a response override on the fakeintake server
func (c *Client) ConfigureOverride(override api.ResponseOverride) error {
	route := fmt.Sprintf("%s/fakeintake/configure/override", c.fakeIntakeURL)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(override)
	if err != nil {
		return err
	}

	resp, err := http.Post(route, "application/json", buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error code %v", resp.StatusCode)
	}
	return nil
}

func (c *Client) getMetric(name string) ([]*aggregator.MetricSeries, error) {
	err := c.getMetrics()
	if err != nil {
		return nil, err
	}
	return c.metricAggregator.GetPayloadsByName(name), nil
}

// A MatchOpt to filter fakeintake payloads
type MatchOpt[P aggregator.PayloadItem] func(payload P) (bool, error)

// GetMetricNames fetches fakeintake on `/api/v2/series` endpoint and returns
// all received metric names
func (c *Client) GetMetricNames() ([]string, error) {
	err := c.getMetrics()
	if err != nil {
		return nil, err
	}
	return c.metricAggregator.GetNames(), nil
}

// FilterMetrics fetches fakeintake on `/api/v2/series` endpoint and returns
// metrics matching `name` and any [MatchOpt](#MatchOpt) options
func (c *Client) FilterMetrics(name string, options ...MatchOpt[*aggregator.MetricSeries]) ([]*aggregator.MetricSeries, error) {
	metrics, err := c.getMetric(name)
	if err != nil {
		return nil, err
	}
	// apply filters one after the other
	filteredMetrics := []*aggregator.MetricSeries{}
	for _, metric := range metrics {
		matchCount := 0
		for _, matchOpt := range options {
			isMatch, err := matchOpt(metric)
			if err != nil {
				return nil, err
			}
			if !isMatch {
				break
			}
			matchCount++
		}
		if matchCount == len(options) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}
	return filteredMetrics, nil
}

// WithTags filters by `tags`
func WithTags[P aggregator.PayloadItem](tags []string) MatchOpt[P] {
	return func(payload P) (bool, error) {
		if aggregator.AreTagsSubsetOfOtherTags(tags, payload.GetTags()) {
			return true, nil
		}
		// TODO return similarity error score
		return false, nil
	}
}

// WithMetricValueInRange filters metrics with values in range `minValue < value < maxValue`
func WithMetricValueInRange(minValue float64, maxValue float64) MatchOpt[*aggregator.MetricSeries] {
	return func(metric *aggregator.MetricSeries) (bool, error) {
		isMatch, err := WithMetricValueHigherThan(minValue)(metric)
		if !isMatch || err != nil {
			return isMatch, err
		}
		return WithMetricValueLowerThan(maxValue)(metric)
	}
}

// WithMetricValueLowerThan filters metrics with values lower than `maxValue`
func WithMetricValueLowerThan(maxValue float64) MatchOpt[*aggregator.MetricSeries] {
	return func(metric *aggregator.MetricSeries) (bool, error) {
		for _, point := range metric.Points {
			if point.Value < maxValue {
				return true, nil
			}
		}
		// TODO return similarity error score
		return false, nil
	}
}

// WithMetricValueLowerThan filters metrics with values higher than `minValue`
//
//nolint:revive // TODO(APL) Fix revive linter
func WithMetricValueHigherThan(minValue float64) MatchOpt[*aggregator.MetricSeries] {
	return func(metric *aggregator.MetricSeries) (bool, error) {
		for _, point := range metric.Points {
			if point.Value > minValue {
				return true, nil
			}
		}
		// TODO return similarity error score
		return false, nil
	}
}

func (c *Client) getLog(service string) ([]*aggregator.Log, error) {
	err := c.getLogs()
	if err != nil {
		return nil, err
	}
	return c.logAggregator.GetPayloadsByName(service), nil
}

// GetLogNames fetches fakeintake on `/api/v2/logs` endpoint and returns
// all received log service names
//
//nolint:revive // TODO(APL) Fix revive linter
func (c *Client) GetLogServiceNames() ([]string, error) {
	err := c.getLogs()
	if err != nil {
		return nil, err
	}
	return c.logAggregator.GetNames(), nil
}

// FilterLogs fetches fakeintake on `/api/v2/logs` endpoint, unpackage payloads and returns
// logs matching `service` and any [MatchOpt](#MatchOpt) options
func (c *Client) FilterLogs(service string, options ...MatchOpt[*aggregator.Log]) ([]*aggregator.Log, error) {
	logs, err := c.getLog(service)
	if err != nil {
		return nil, err
	}
	// apply filters one after the other
	filteredLogs := []*aggregator.Log{}
	for _, log := range logs {
		matchCount := 0
		for _, matchOpt := range options {
			isMatch, err := matchOpt(log)
			if err != nil {
				return nil, err
			}
			if !isMatch {
				break
			}
			matchCount++
		}
		if matchCount == len(options) {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs, nil
}

// WithMessageContaining filters logs by message containing `content`
func WithMessageContaining(content string) MatchOpt[*aggregator.Log] {
	return func(log *aggregator.Log) (bool, error) {
		if strings.Contains(log.Message, content) {
			return true, nil
		}
		// TODO return similarity score in error
		return false, nil
	}
}

// WithMessageMatching filters logs by message matching [regexp](https://pkg.go.dev/regexp) `pattern`
func WithMessageMatching(pattern string) MatchOpt[*aggregator.Log] {
	return func(log *aggregator.Log) (bool, error) {
		matched, err := regexp.MatchString(pattern, log.Message)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
		// TODO return similarity score in error
		return false, nil
	}
}

// GetCheckRunNames fetches fakeintake on `/api/v1/check_run` endpoint and returns
// all received check run names
func (c *Client) GetCheckRunNames() ([]string, error) {
	err := c.getCheckRuns()
	if err != nil {
		return nil, err
	}
	return c.checkRunAggregator.GetNames(), nil
}

// FilterLogs fetches fakeintake on `/api/v1/check_run` endpoint, unpackage payloads and returns
// checks matching `name`
//
//nolint:revive // TODO(APL) Fix revive linter
func (c *Client) GetCheckRun(name string) ([]*aggregator.CheckRun, error) {
	err := c.getCheckRuns()
	if err != nil {
		return nil, err
	}
	return c.checkRunAggregator.GetPayloadsByName(name), nil
}

// FlushServerAndResetAggregators sends a request to delete any stored payload
// and resets client's  aggregators
// Call this in between tests to reset the fakeintake status on both client and server side
func (c *Client) FlushServerAndResetAggregators() error {
	err := c.flushPayloads()
	if err != nil {
		return err
	}
	c.checkRunAggregator.Reset()
	c.metricAggregator.Reset()
	c.logAggregator.Reset()
	return nil
}

func (c *Client) flushPayloads() error {
	resp, err := http.Get(fmt.Sprintf("%s/fakeintake/flushPayloads", c.fakeIntakeURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error code %v", resp.StatusCode)
	}
	return nil
}

// GetConnections fetches fakeintake on `/api/v1/connections` endpoint and returns
// all received connections
func (c *Client) GetConnections() (conns *aggregator.ConnectionsAggregator, err error) {
	err = c.getConnections()
	if err != nil {
		return nil, err
	}
	return &c.connectionAggregator, nil
}

// GetConnectionsNames fetches fakeintake on `/api/v1/connections` endpoint and returns
// all received connections from hostname+network_id
func (c *Client) GetConnectionsNames() ([]string, error) {
	err := c.getConnections()
	if err != nil {
		return []string{}, err
	}
	return c.connectionAggregator.GetNames(), nil
}

// URL returns the client's URL
func (c *Client) URL() string {
	return c.fakeIntakeURL
}

// GetProcesses fetches fakeintake on `/api/v1/collector` endpoint and returns
// all received process payloads
func (c *Client) GetProcesses() ([]*aggregator.ProcessPayload, error) {
	err := c.getProcesses()
	if err != nil {
		return nil, err
	}

	var procs []*aggregator.ProcessPayload
	for _, name := range c.processAggregator.GetNames() {
		procs = append(procs, c.processAggregator.GetPayloadsByName(name)...)
	}

	return procs, nil
}

// GetContainers fetches fakeintake on `/api/v1/container` endpoint and returns
// all received container payloads
func (c *Client) GetContainers() ([]*aggregator.ContainerPayload, error) {
	err := c.getContainers()
	if err != nil {
		return nil, err
	}

	var containers []*aggregator.ContainerPayload
	for _, name := range c.containerAggregator.GetNames() {
		containers = append(containers, c.containerAggregator.GetPayloadsByName(name)...)
	}

	return containers, nil
}

// GetProcessDiscoveries fetches fakeintake on `/api/v1/discovery` endpoint and returns
// all received process discovery payloads
func (c *Client) GetProcessDiscoveries() ([]*aggregator.ProcessDiscoveryPayload, error) {
	err := c.getProcessDiscoveries()
	if err != nil {
		return nil, err
	}

	var discs []*aggregator.ProcessDiscoveryPayload
	for _, name := range c.processDiscoveryAggregator.GetNames() {
		discs = append(discs, c.processDiscoveryAggregator.GetPayloadsByName(name)...)
	}

	return discs, nil
}

// GetOrchestrators fetches fakeintake on `/api/v2/orch` endpoint and returns
// all received process payloads
func (c *Client) GetOrchestrators(filter *PayloadFilter) ([]*aggregator.OrchestratorPayload, error) {
	err := c.getOrchestrators()
	if err != nil {
		return nil, err
	}

	var orchs []*aggregator.OrchestratorPayload
	for _, name := range c.orchestratorAggregator.GetNames() {
		if filter != nil && filter.Name != "" && filter.Name != name {
			continue
		}
		for _, payload := range c.orchestratorAggregator.GetPayloadsByName(name) {
			if filter != nil && filter.ResourceType != 0 && filter.ResourceType != payload.Type {
				continue
			}
			orchs = append(orchs, payload)
		}
	}
	return orchs, nil
}

// GetOrchestratorManifests fetches fakeintake on `/api/v2/orchmanif` endpoint and returns
// all received process payloads
func (c *Client) GetOrchestratorManifests() ([]*aggregator.OrchestratorManifestPayload, error) {
	err := c.getOrchestratorManifests()
	if err != nil {
		return nil, err
	}

	var manifs []*aggregator.OrchestratorManifestPayload
	for _, name := range c.orchestratorManifestAggregator.GetNames() {
		manifs = append(manifs, c.orchestratorManifestAggregator.GetPayloadsByName(name)...)
	}

	return manifs, nil
}

func (c *Client) get(route string) ([]byte, error) {
	var body []byte
	err := backoff.Retry(func() error {
		tmpResp, err := http.Get(fmt.Sprintf("%s/%s", c.fakeIntakeURL, route))
		if err != nil {
			return err
		}
		defer tmpResp.Body.Close()
		if tmpResp.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected %d got %d", http.StatusOK, tmpResp.StatusCode)
		}
		body, err = io.ReadAll(tmpResp.Body)
		return err
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 4))
	return body, err
}

// RouteStats queries the routestats fakeintake endpoint to get statistics about each route.
// It only returns statistics about endpoint which store some payloads.
func (c *Client) RouteStats() (map[string]int, error) {
	body, err := c.get("fakeintake/routestats")
	if err != nil {
		return nil, err
	}

	var routestats api.APIFakeIntakeRouteStatsGETResponse
	err = json.Unmarshal(body, &routestats)
	if err != nil {
		return nil, err
	}

	routes := map[string]int{}
	for endpoint, stats := range routestats.Routes {
		// the count of a given endpoint can be zero when old payloads are periodically removed
		if stats.Count != 0 {
			routes[endpoint] = stats.Count
		}
	}

	return routes, nil
}

type PayloadFilter struct {
	Name         string
	ResourceType agentmodel.MessageType
}
