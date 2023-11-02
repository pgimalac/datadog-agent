// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package testdns

import (
	"net"
	"os/exec"
	"sync"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

var globalServer *server
var globalServerError error
var serverOnce sync.Once

// GetServerIP returns the IP address of the test DNS server. The test DNS server returns canned responses for several
// known domains that are used in integration tests.
//
// see server#start to see which domains are handled.
func GetServerIP(t *testing.T) net.IP {
	serverOnce.Do(func() {
		globalServer, globalServerError = NewServer()
		globalServer.Start("tcp")
		globalServer.Start("udp")
	})

	require.NoError(t, globalServerError)
	return net.ParseIP("10.10.10.10")
}

type server struct{}

func NewServer() (*server, error) {
	// ignore errors as device might not exist from prior test run
	exec.Command("ip", "link", "del", "dev", "dnstestdummy").Run()

	err := exec.Command("ip", "link", "add", "dnstestdummy", "type", "dummy").Run()
	if err != nil {
		return nil, err
	}

	err = exec.Command("ip", "addr", "add", "dev", "dnstestdummy", "10.10.10.10", "broadcast", "+").Run()
	if err != nil {
		return nil, err
	}

	err = exec.Command("ip", "link", "set", "dnstestdummy", "up").Run()
	if err != nil {
		return nil, err
	}

	return &server{}, nil
}

func (s *server) Start(transport string) {
	started := make(chan struct{}, 1)
	srv := dns.Server{
		Addr: "10.10.10.10:53",
		Net:  transport,
		Handler: dns.HandlerFunc(func(writer dns.ResponseWriter, msg *dns.Msg) {
			switch msg.Question[0].Name {
			case "good.com.":
				respond(msg, writer, "good.com. 30 IN A  10.0.0.1")
			case "golang.org.":
				respond(msg, writer, "golang.org. 30 IN A  10.0.0.2")
			case "google.com.":
				respond(msg, writer, "google.com. 30 IN A  10.0.0.3")
			case "acm.org.":
				respond(msg, writer, "acm.org. 30 IN A  10.0.0.4")
			case "nonexistenent.net.com.":
				resp := &dns.Msg{}
				resp.SetReply(msg)
				resp.Rcode = dns.RcodeNameError
				writer.WriteMsg(resp)
			case "missingdomain.com.":
				resp := &dns.Msg{}
				resp.SetReply(msg)
				resp.Rcode = dns.RcodeNameError
				writer.WriteMsg(resp)
			default:
				resp := &dns.Msg{}
				resp.SetReply(msg)
				resp.Rcode = dns.RcodeServerFailure
				writer.WriteMsg(resp)
			}
		}),
		NotifyStartedFunc: func() {
			started <- struct{}{}
		},
	}
	go srv.ListenAndServe()
	<-started
}

func respond(req *dns.Msg, writer dns.ResponseWriter, record string) {
	resp := &dns.Msg{}
	resp.SetReply(req)

	rr, err := dns.NewRR(record)
	if err != nil {
		panic(err)
	}
	resp.Answer = []dns.RR{rr}
	writer.WriteMsg(resp)
}
