package propagation

import (
	"errors"

	"github.com/DataDog/datadog-agent/pkg/trace/sampler"
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const defaultPriority sampler.SamplingPriority = sampler.PriorityNone

var (
	errorUnsupportedExtractionType = errors.New("Unsupported event type for trace context extraction")
	errorNoContextFound            = errors.New("No trace context found")
	errorNoSQSRecordFound          = errors.New("No sqs message records found for trace context extraction")
)

type Extractor struct {
	propagator tracer.Propagator
}

type TraceContext struct {
	TraceID  uint64
	ParentID uint64
	Priority sampler.SamplingPriority
}

type TraceContexter interface {
	GetTraceContext() *TraceContext
}

func (e Extractor) Extract(events ...interface{}) (*TraceContext, error) {
	for _, event := range events {
		if tc, err := e.extract(event); err == nil {
			return tc, nil
		}
	}
	return nil, errorNoContextFound
}

func (e Extractor) extract(event interface{}) (*TraceContext, error) {
	var carrier tracer.TextMapReader
	var err error

	switch ev := event.(type) {
	case []byte:
		carrier, err = rawPayloadCarrier(ev)
	case events.SQSEvent:
		// look for context in just the first message
		if len(ev.Records) > 0 {
			return e.extract(ev.Records[0])
		}
		return nil, errorNoSQSRecordFound
	case events.SQSMessage:
		if attr, ok := ev.Attributes[awsTraceHeader]; ok {
			if tc, err := extractTraceContextfromAWSTraceHeader(attr); err == nil {
				// Return early if AWSTraceHeader contains trace context
				return tc, nil
			}
		}
		carrier, err = sqsMessageCarrier(ev)
	case TraceContexter:
		// TODO: only look for datadog headers
		if tc := ev.GetTraceContext(); tc != nil {
			return tc, nil
		}
		return nil, errorNoContextFound
	default:
		err = errorUnsupportedExtractionType
	}

	if err != nil {
		return nil, err
	}
	if e.propagator == nil {
		e.propagator = tracer.NewPropagator(nil)
	}
	sc, err := e.propagator.Extract(carrier)
	if err != nil {
		return nil, err
	}
	return &TraceContext{
		TraceID:  sc.TraceID(),
		ParentID: sc.SpanID(),
		Priority: getPriority(sc),
	}, nil
}

func getPriority(sc ddtrace.SpanContext) (priority sampler.SamplingPriority) {
	priority = defaultPriority
	if pc, ok := sc.(interface{ SamplingPriority() (int, bool) }); ok {
		if p, ok := pc.SamplingPriority(); ok {
			priority = sampler.SamplingPriority(p)
		}
	}
	return
}
