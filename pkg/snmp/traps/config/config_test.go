// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020-present Datadog, Inc.

package config

import (
	"strings"
	"testing"

	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/comp/core/log"
	"github.com/DataDog/datadog-agent/comp/core/log/logimpl"
	ddconf "github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const mockedHostname = "VeryLongHostnameThatDoesNotFitIntoTheByteArray"

var expectedEngineID = "\x80\xff\xff\xff\xff\x67\xb2\x0f\xe4\xdf\x73\x7a\xce\x28\x47\x03\x8f\x57\xe6\x5c\x98"

var expectedEngineIDs = map[string]string{
	"VeryLongHostnameThatDoesNotFitIntoTheByteArray": "\x80\xff\xff\xff\xff\x67\xb2\x0f\xe4\xdf\x73\x7a\xce\x28\x47\x03\x8f\x57\xe6\x5c\x98",
	"VeryLongHostnameThatIsDifferent":                "\x80\xff\xff\xff\xff\xe7\x21\xcc\xd7\x0b\xe1\x60\xc5\x18\xd7\xde\x17\x86\xb0\x7d\x36",
}

var usersV3 = []UserV3{
	{
		Username:     "user",
		AuthKey:      "password",
		AuthProtocol: "MD5",
		PrivKey:      "password",
		PrivProtocol: "AES",
	},
	{
		Username:     "user",
		AuthKey:      "password",
		AuthProtocol: "SHA",
		PrivKey:      "password",
		PrivProtocol: "DES",
	},
	{
		Username:     "user2",
		AuthKey:      "password",
		AuthProtocol: "MD5",
		PrivKey:      "password",
		PrivProtocol: "AES",
	},
}

var usmUsers = []*gosnmp.UsmSecurityParameters{
	{
		UserName:                 "user",
		AuthenticationProtocol:   gosnmp.MD5,
		AuthenticationPassphrase: "password",
		PrivacyProtocol:          gosnmp.AES,
		PrivacyPassphrase:        "password",
	},
	{
		UserName:                 "user",
		AuthenticationProtocol:   gosnmp.SHA,
		AuthenticationPassphrase: "password",
		PrivacyProtocol:          gosnmp.DES,
		PrivacyPassphrase:        "password",
	},
	{
		UserName:                 "user2",
		AuthenticationProtocol:   gosnmp.MD5,
		AuthenticationPassphrase: "password",
		PrivacyProtocol:          gosnmp.AES,
		PrivacyPassphrase:        "password",
	},
}

func makeConfig(t *testing.T, trapConfig TrapsConfig) config.Component {
	return makeConfigWithGlobalNamespace(t, trapConfig, "")
}

func makeConfigWithGlobalNamespace(t *testing.T, trapConfig TrapsConfig, globalNamespace string) config.Component {
	trapConfig.Enabled = true
	conf := ddconf.SetupConf()
	if globalNamespace != "" {
		conf.SetWithoutSource("network_devices.namespace", globalNamespace)
	}
	conf.SetConfigType("yaml")
	yamlData := map[string]map[string]interface{}{
		"network_devices": {
			"snmp_traps": trapConfig,
		},
	}
	out, err := yaml.Marshal(yamlData)
	require.NoError(t, err)
	err = conf.ReadConfig(strings.NewReader(string(out)))
	require.NoError(t, err)
	return conf
}

func TestFullConfig(t *testing.T) {
	logger := fxutil.Test[log.Component](t, logimpl.MockModule())
	rootConfig := makeConfig(t, TrapsConfig{
		Port:             1234,
		Users:            usersV3,
		BindHost:         "127.0.0.1",
		CommunityStrings: []string{"public"},
		StopTimeout:      12,
		Namespace:        "foo",
	})
	config, err := ReadConfig(mockedHostname, rootConfig)
	assert.NoError(t, err)
	assert.Equal(t, uint16(1234), config.Port)
	assert.Equal(t, 12, config.StopTimeout)
	assert.Equal(t, []string{"public"}, config.CommunityStrings)
	assert.Equal(t, "127.0.0.1", config.BindHost)
	assert.Equal(t, "foo", config.Namespace)
	assert.Equal(t, usersV3, config.Users)

	params, err := config.BuildSNMPParams(logger)
	assert.NoError(t, err)
	assert.Equal(t, uint16(1234), params.Port)
	assert.Equal(t, gosnmp.Version3, params.Version)
	assert.Equal(t, "udp", params.Transport)
	assert.NotNil(t, params.Logger)
	assert.Equal(t, gosnmp.UserSecurityModel, params.SecurityModel)
	assert.Equal(t, &gosnmp.UsmSecurityParameters{AuthoritativeEngineID: expectedEngineID}, params.SecurityParameters)

	table := gosnmp.NewSnmpV3SecurityParametersTable()
	for _, usmUser := range usmUsers {
		table.Add(usmUser.UserName, usmUser)
	}
	var usmConfigTests = []struct {
		name       string
		identifier string
	}{
		{
			"identifier: user has 2 entries",
			"user",
		},
		{
			"identifier: user2 has 1 entry",
			"user2",
		},
	}
	for _, usmConfigTest := range usmConfigTests {
		// Compare the security params after initializing the security keys (happens in the add to table)
		expected, _ := table.Get(usmConfigTest.identifier)
		actual, _ := params.TrapSecurityParametersTable.Get(usmConfigTest.identifier)
		assert.ElementsMatch(t, expected, actual)
	}
}

func TestMinimalConfig(t *testing.T) {
	logger := fxutil.Test[log.Component](t, logimpl.MockModule())
	config, err := ReadConfig("", makeConfig(t, TrapsConfig{}))
	assert.NoError(t, err)
	assert.Equal(t, uint16(9162), config.Port)
	assert.Equal(t, 5, config.StopTimeout)
	assert.Empty(t, config.CommunityStrings)
	assert.Equal(t, "0.0.0.0", config.BindHost)
	assert.Empty(t, config.Users)
	assert.Equal(t, "default", config.Namespace)

	params, err := config.BuildSNMPParams(logger)
	assert.NoError(t, err)
	assert.Equal(t, uint16(9162), params.Port)
	assert.Equal(t, gosnmp.Version2c, params.Version)
	assert.Equal(t, "udp", params.Transport)
	assert.NotNil(t, params.Logger)
	assert.Equal(t, nil, params.SecurityParameters)
}

func TestDefaultUsers(t *testing.T) {
	config, err := ReadConfig("", makeConfig(t, TrapsConfig{
		CommunityStrings: []string{"public"},
		StopTimeout:      11,
	}))
	assert.NoError(t, err)

	assert.Equal(t, 11, config.StopTimeout)
}

func TestBuildAuthoritativeEngineID(t *testing.T) {
	for hostname, engineID := range expectedEngineIDs {
		config, err := ReadConfig(hostname, makeConfig(t, TrapsConfig{}))
		assert.NoError(t, err)
		assert.Equal(t, engineID, config.authoritativeEngineID)
	}
}

func TestNamespaceIsNormalized(t *testing.T) {
	config, err := ReadConfig("", makeConfig(t, TrapsConfig{
		Namespace: "><\n\r\tfoo",
	}))
	assert.NoError(t, err)

	assert.Equal(t, "--foo", config.Namespace)
}

func TestInvalidNamespace(t *testing.T) {
	_, err := ReadConfig("", makeConfig(t, TrapsConfig{
		Namespace: strings.Repeat("x", 101),
	}))
	assert.Error(t, err)
}

func TestNamespaceSetGlobally(t *testing.T) {
	config, err := ReadConfig("", makeConfigWithGlobalNamespace(t, TrapsConfig{}, "foo"))
	assert.NoError(t, err)

	assert.Equal(t, "foo", config.Namespace)
}

func TestNamespaceSetBothGloballyAndLocally(t *testing.T) {
	config, err := ReadConfig("", makeConfigWithGlobalNamespace(t, TrapsConfig{Namespace: "bar"}, "foo"))
	assert.NoError(t, err)

	assert.Equal(t, "bar", config.Namespace)
}
