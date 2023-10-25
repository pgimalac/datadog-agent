// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.

//go:build test

package server

import (
	"context"
	"encoding/json"
	"github.com/DataDog/datadog-agent/comp/netflow/goflowlib"
	"github.com/DataDog/datadog-agent/comp/netflow/goflowlib/netflowstate"
	"github.com/DataDog/datadog-agent/comp/netflow/payload"
	"github.com/netsampler/goflow2/decoders/netflow/templates"
	"github.com/netsampler/goflow2/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/pkg/epforwarder"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"

	"github.com/DataDog/datadog-agent/comp/ndmtmp/forwarder"

	ndmtestutils "github.com/DataDog/datadog-agent/pkg/networkdevice/testutils"

	"github.com/DataDog/datadog-agent/comp/netflow/common"
	nfconfig "github.com/DataDog/datadog-agent/comp/netflow/config"
	"github.com/DataDog/datadog-agent/comp/netflow/flowaggregator"
	"github.com/DataDog/datadog-agent/comp/netflow/testutil"
)

func singleListenerConfig(flowType common.FlowType, port uint16) *nfconfig.NetflowConfig {
	return &nfconfig.NetflowConfig{
		Enabled:                 true,
		AggregatorFlushInterval: 1,
		Listeners: []nfconfig.ListenerConfig{{
			FlowType: flowType,
			BindHost: "127.0.0.1",
			Port:     port,
		}},
	}
}

var flushTime, _ = time.Parse(time.RFC3339, "2019-02-18T16:00:06Z")

var setTimeNow = fx.Invoke(func(c Component) {
	c.(*Server).FlowAgg.TimeNowFunction = func() time.Time {
		return flushTime
	}
})

func assertFlowEventsCount(t *testing.T, port uint16, srv *Server, packetData []byte, expectedEvents uint64) bool {
	return assert.EventuallyWithT(t, func(c *assert.CollectT) {
		err := testutil.SendUDPPacket(port, packetData)
		assert.NoError(c, err, "error sending udp packet")
		if err != nil {
			return
		}

		netflowEvents, err := flowaggregator.WaitForFlowsToBeFlushed(srv.FlowAgg, 1*time.Second, 2)
		assert.Equal(c, expectedEvents, netflowEvents)
		assert.NoError(c, err)
	}, 10*time.Second, 10*time.Millisecond)
}

func TestNetFlow_IntegrationTest_NetFlow5(t *testing.T) {
	port, err := ndmtestutils.GetFreePort()
	require.NoError(t, err)
	var epForwarder forwarder.MockComponent
	srv := fxutil.Test[Component](t, fx.Options(
		testOptions,
		fx.Populate(&epForwarder),
		fx.Replace(
			singleListenerConfig("netflow5", port),
		),
		setTimeNow,
	)).(*Server)

	// Set expectations
	testutil.ExpectNetflow5Payloads(t, epForwarder)
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), "network-devices-metadata").Return(nil).Times(1)

	// Send netflowV5Data twice to test aggregator
	// Flows will have 2x bytes/packets after aggregation
	packetData, err := testutil.GetNetFlow5Packet()
	require.NoError(t, err, "error getting packet")

	assertFlowEventsCount(t, port, srv, packetData, 2)
}

func TestNetFlow_IntegrationTest_NetFlow9(t *testing.T) {
	port, err := ndmtestutils.GetFreePort()
	require.NoError(t, err)
	var epForwarder forwarder.MockComponent
	srv := fxutil.Test[Component](t, fx.Options(
		testOptions,
		fx.Populate(&epForwarder),
		fx.Replace(
			singleListenerConfig("netflow9", port),
		),
		setTimeNow,
	)).(*Server)

	// Test later content of payloads if needed for more precise test.
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), epforwarder.EventTypeNetworkDevicesNetFlow).Return(nil).Times(29)
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), "network-devices-metadata").Return(nil).Times(1)

	packetData, err := testutil.GetNetFlow9Packet()
	require.NoError(t, err, "error getting packet")

	assertFlowEventsCount(t, port, srv, packetData, 29)
}

func TestNetFlow_IntegrationTest_SFlow5(t *testing.T) {
	port, err := ndmtestutils.GetFreePort()
	require.NoError(t, err)
	var epForwarder forwarder.MockComponent
	srv := fxutil.Test[Component](t, fx.Options(
		testOptions,
		fx.Populate(&epForwarder),
		fx.Replace(
			singleListenerConfig("sflow5", port),
		),
		setTimeNow,
	)).(*Server)

	// Test later content of payloads if needed for more precise test.
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), epforwarder.EventTypeNetworkDevicesNetFlow).Return(nil).Times(7)
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), "network-devices-metadata").Return(nil).Times(1)

	packetData, err := testutil.GetSFlow5Packet()
	require.NoError(t, err, "error getting sflow data")

	assertFlowEventsCount(t, port, srv, packetData, 7)
}

func TestNetFlow_IntegrationTest_CustomFields(t *testing.T) {
	port, err := ndmtestutils.GetFreePort()
	require.NoError(t, err)
	var epForwarder forwarder.MockComponent
	srv := fxutil.Test[Component](t, fx.Options(
		testOptions,
		fx.Populate(&epForwarder),
		fx.Replace(
			&nfconfig.NetflowConfig{
				Enabled:                 true,
				AggregatorFlushInterval: 1,
				Listeners: []nfconfig.ListenerConfig{{
					FlowType: common.TypeNetFlow9,
					BindHost: "127.0.0.1",
					Port:     port,
					Mapping: []nfconfig.Mapping{
						{
							Field:       11,
							Destination: "source.port", // Inverting source and destination port to test
							Type:        common.Varint,
						},
						{
							Field:       7,
							Destination: "destination.port",
							Type:        common.Varint,
						},
						{
							Field:       32,
							Destination: "icmp_type",
							Type:        common.Bytes,
						},
					},
				}},
			},
		),
		setTimeNow,
	)).(*Server)

	flowData, err := testutil.GetNetFlow9Packet()
	require.NoError(t, err, "error getting packet")

	// Set expectations
	testutil.ExpectPayloadWithCustomFields(t, epForwarder)
	epForwarder.EXPECT().SendEventPlatformEventBlocking(gomock.Any(), "network-devices-metadata").Return(nil).Times(1)

	assertFlowEventsCount(t, port, srv, flowData, 29)
}

func BenchmarkNetflowCustomFields(b *testing.B) {
	flowChan := make(chan *common.Flow, 10)
	listenerFlowCount := atomic.NewInt64(0)

	go func() {
		for {
			// Consume chan while benchmarking
			<-flowChan
		}
	}()

	formatDriver := goflowlib.NewAggregatorFormatDriver(flowChan, "bench", listenerFlowCount)
	logrusLogger := logrus.StandardLogger()
	ctx := context.Background()

	templateSystem, err := templates.FindTemplateSystem(ctx, "memory")
	if err != nil {
		require.NoError(b, err, "error with template")
	}
	defer templateSystem.Close(ctx)

	goflowState := utils.NewStateNetFlow()
	goflowState.Format = formatDriver
	goflowState.Logger = logrusLogger
	goflowState.TemplateSystem = templateSystem

	customStateWithoutFields := netflowstate.NewStateNetFlow(nil)
	customStateWithoutFields.Format = formatDriver
	customStateWithoutFields.Logger = logrusLogger
	customStateWithoutFields.TemplateSystem = templateSystem

	customState := netflowstate.NewStateNetFlow([]nfconfig.Mapping{
		{
			Field:       11,
			Destination: "source.port",
			Type:        common.Varint,
		},
		{
			Field:       7,
			Destination: "destination.port",
			Type:        common.Varint,
		},
		{
			Field:       32,
			Destination: "icmp_type",
			Type:        common.Bytes,
		},
	})

	customState.Format = formatDriver
	customState.Logger = logrusLogger
	customState.TemplateSystem = templateSystem

	flowData, err := testutil.GetNetFlow9Packet()
	require.NoError(b, err, "error getting cflow flow data")

	flowPacket := utils.BaseMessage{
		Src:      net.ParseIP("127.0.0.1"),
		Port:     3000,
		Payload:  flowData,
		SetTime:  false,
		RecvTime: time.Now(),
	}

	b.Run("goflow2 default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = goflowState.DecodeFlow(flowPacket)
			require.NoError(b, err, "error processing packet")
		}
	})

	b.Run("goflow2 netflow custom state without custom fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = customStateWithoutFields.DecodeFlow(flowPacket)
			require.NoError(b, err, "error processing packet")
		}
	})

	b.Run("goflow2 netflow custom state with custom fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = customState.DecodeFlow(flowPacket)
			require.NoError(b, err, "error processing packet")
		}
	})
}

func BenchmarkNetflowPayloadMarshalling(b *testing.B) {
	type FlowPayload_ payload.FlowPayload
	flowPayload := payload.FlowPayload{
		FlushTimestamp: 1000,
		FlowType:       "netflow9",
		SamplingRate:   1,
		Direction:      "ingress",
		Start:          13213234,
		End:            24342343,
		Bytes:          23423423,
		Packets:        13212,
		EtherType:      "a",
		IPProtocol:     "a",
		Device: payload.Device{
			Namespace: "default",
		},
		Exporter: payload.Exporter{
			IP: "10.0.0.3",
		},
		Source: payload.Endpoint{
			IP:   "10.0.1.13",
			Port: "4567",
			Mac:  "00:01:qq:02",
			Mask: "24",
		},
		Destination: payload.Endpoint{
			IP:   "10.0.1.14",
			Port: "22",
			Mac:  "00:01:qq:03",
			Mask: "24",
		},
		Ingress: payload.ObservationPoint{
			Interface: payload.Interface{
				Index: 5,
			},
		},
		Egress: payload.ObservationPoint{
			Interface: payload.Interface{
				Index: 12,
			},
		},
		Host:     "127.0.0.1",
		TCPFlags: nil,
		NextHop: payload.NextHop{
			IP: "10.4.5.6",
		},
		AdditionalFields: map[string]any{
			"test_fields": "bonjour",
		},
	}

	b.Run("Without custom marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(FlowPayload_(flowPayload))
			require.NoError(b, err, "error processing packet")
		}
	})
	b.Run("With Thibaud custom marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := flowPayload.MarshalWithAdditionalFields()
			require.NoError(b, err, "error processing packet")
		}
	})
	b.Run("With Alex custom marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := flowPayload.MarshalWithAdditionalFieldsLessMarshall()
			require.NoError(b, err, "error processing packet")
		}
	})
	b.Run("With Reflection custom marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := flowPayload.MarshalWithAdditionalFieldsReflect()
			require.NoError(b, err, "error processing packet")
		}
	})
	b.Run("With Manual custom marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := flowPayload.MarshalJSON()
			require.NoError(b, err, "error processing packet")
		}
	})
}
