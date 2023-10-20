// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package profile

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cihub/seelog"
	assert "github.com/stretchr/testify/require"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/networkdevice/profile/profiledefinition"
)

func getMetricFromProfile(p profiledefinition.ProfileDefinition, metricName string) *profiledefinition.MetricsConfig {
	for _, m := range p.Metrics {
		if m.Symbol.Name == metricName {
			return &m
		}
	}
	return nil
}

func Test_resolveProfileDefinitionPath(t *testing.T) {
	defaultTestConfdPath, _ := filepath.Abs(filepath.Join("..", "test", "user_profiles.d"))
	config.Datadog.Set("confd_path", defaultTestConfdPath)

	absPath, _ := filepath.Abs(filepath.Join("tmp", "myfile.yaml"))
	tests := []struct {
		name               string
		definitionFilePath string
		expectedPath       string
	}{
		{
			name:               "abs path",
			definitionFilePath: absPath,
			expectedPath:       absPath,
		},
		{
			name:               "relative path with default profile",
			definitionFilePath: "p2.yaml",
			expectedPath:       filepath.Join(config.Datadog.Get("confd_path").(string), "snmp.d", "default_profiles", "p2.yaml"),
		},
		{
			name:               "relative path with user profile",
			definitionFilePath: "p3.yaml",
			expectedPath:       filepath.Join(config.Datadog.Get("confd_path").(string), "snmp.d", "profiles", "p3.yaml"),
		},
		{
			name:               "relative path with user profile precedence",
			definitionFilePath: "p1.yaml",
			expectedPath:       filepath.Join(config.Datadog.Get("confd_path").(string), "snmp.d", "profiles", "p1.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := resolveProfileDefinitionPath(tt.definitionFilePath)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}

func Test_loadDefaultProfiles(t *testing.T) {
	SetConfdPathAndCleanProfiles()
	SetGlobalProfileConfigMap(nil)
	defaultProfiles, err := loadYamlProfiles()
	assert.Nil(t, err)
	defaultProfiles2, err := loadYamlProfiles()
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("%p", defaultProfiles), fmt.Sprintf("%p", defaultProfiles2))
}

func Test_loadDefaultProfiles_withUserProfiles(t *testing.T) {
	defaultTestConfdPath, _ := filepath.Abs(filepath.Join("..", "test", "user_profiles.d"))
	SetGlobalProfileConfigMap(nil)
	config.Datadog.Set("confd_path", defaultTestConfdPath)

	defaultProfiles, err := loadYamlProfiles()
	assert.Nil(t, err)

	assert.Len(t, defaultProfiles, 4)
	assert.NotNil(t, defaultProfiles)

	p1 := defaultProfiles["p1"].Definition // user p1 overrides datadog p1
	p2 := defaultProfiles["p2"].Definition // datadog p2
	p3 := defaultProfiles["p3"].Definition // user p3
	p4 := defaultProfiles["p4"].Definition // user p3

	assert.Equal(t, "p1_user", p1.Device.Vendor) // overrides datadog p1 profile
	assert.NotNil(t, getMetricFromProfile(p1, "p1_metric_override"))

	assert.Equal(t, "p2_datadog", p2.Device.Vendor)
	assert.NotNil(t, getMetricFromProfile(p2, "p2_metric"))

	assert.Equal(t, "p3_user", p3.Device.Vendor)
	assert.NotNil(t, getMetricFromProfile(p3, "p3_metric"))

	assert.Equal(t, "p4_user", p4.Device.Vendor)
	assert.NotNil(t, getMetricFromProfile(p4, "p4_user_metric"))
	assert.NotNil(t, getMetricFromProfile(p4, "p4_default_metric"))
}

func Test_loadDefaultProfiles_invalidDir(t *testing.T) {
	invalidPath, _ := filepath.Abs(filepath.Join(".", "tmp", "invalidPath"))
	config.Datadog.Set("confd_path", invalidPath)
	SetGlobalProfileConfigMap(nil)

	defaultProfiles, err := loadYamlProfiles()
	assert.Nil(t, err)
	assert.Len(t, defaultProfiles, 0)
}

func Test_loadDefaultProfiles_invalidExtendProfile(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	l, err := seelog.LoggerFromWriterWithMinLevelAndFormat(w, seelog.DebugLvl, "[%LEVEL] %FuncShort: %Msg")
	assert.Nil(t, err)
	log.SetupLogger(l, "debug")

	profilesWithInvalidExtendConfdPath, _ := filepath.Abs(filepath.Join("..", "test", "invalid_ext.d"))
	config.Datadog.Set("confd_path", profilesWithInvalidExtendConfdPath)
	SetGlobalProfileConfigMap(nil)

	defaultProfiles, err := loadYamlProfiles()

	w.Flush()
	logs := b.String()
	assert.Nil(t, err)

	assert.Equal(t, 1, strings.Count(logs, "[WARN] loadResolveProfiles: failed to expand profile `f5-big-ip"), logs)
	assert.Equal(t, ProfileConfigMap{}, defaultProfiles)
}

func Test_loadDefaultProfiles_validAndInvalidProfiles(t *testing.T) {
	// Valid profiles should be returned even if some profiles are invalid
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	l, err := seelog.LoggerFromWriterWithMinLevelAndFormat(w, seelog.DebugLvl, "[%LEVEL] %FuncShort: %Msg")
	assert.Nil(t, err)
	log.SetupLogger(l, "debug")

	profilesWithInvalidExtendConfdPath, _ := filepath.Abs(filepath.Join("..", "test", "valid_invalid.d"))
	config.Datadog.Set("confd_path", profilesWithInvalidExtendConfdPath)
	SetGlobalProfileConfigMap(nil)

	defaultProfiles, err := loadYamlProfiles()

	for _, profile := range defaultProfiles {
		profiledefinition.NormalizeMetrics(profile.Definition.Metrics)
	}

	w.Flush()
	logs := b.String()
	assert.Nil(t, err)

	assert.Equal(t, 1, strings.Count(logs, "unmarshal errors"), logs)
	assert.Contains(t, defaultProfiles, "f5-big-ip")
	assert.NotContains(t, defaultProfiles, "f5-invalid")
}

func Test_loadProfiles(t *testing.T) {
	tests := []struct {
		name               string
		confdDir           string
		initConfigProfiles ProfileConfigMap
		userProfiles       ProfileConfigMap
		expectedProfiles   []string
		expectedErr        string
		expectedLogs       []string
	}{
		{
			name:               "OK only default profiles",
			confdDir:           "default_profiles.d",
			initConfigProfiles: nil,
			userProfiles:       nil,
			expectedProfiles: []string{
				"p1",
				"p2",
			},
		},
		{
			name:     "OK init config profiles",
			confdDir: "default_profiles.d",
			initConfigProfiles: ProfileConfigMap{
				"init_config_p1": ProfileConfig{
					Definition: profiledefinition.ProfileDefinition{
						Name: "init_config_p1",
					},
				},
			},
			userProfiles: nil,
			expectedProfiles: []string{
				"init_config_p1",
			},
		},
		{
			name:               "OK user profiles",
			confdDir:           "default_profiles.d",
			initConfigProfiles: nil,
			userProfiles: ProfileConfigMap{
				"user_p1": ProfileConfig{
					Definition: profiledefinition.ProfileDefinition{
						Name: "user_p1",
					},
				},
			},
			expectedProfiles: []string{
				"user_p1",
				"p1",
				"p2",
			},
		},
		{
			name:               "ERROR invalid missing confd user and default profiles",
			confdDir:           "does-not-exist.d",
			initConfigProfiles: nil,
			expectedProfiles:   []string{},
			expectedLogs: []string{
				"failed to get user profile definitions",
				"failed to get default profile definitions",
			},
		},
		{
			name:     "ERROR both init config and user profiles",
			confdDir: "default_profiles.d",
			initConfigProfiles: ProfileConfigMap{
				"init_config_p1": ProfileConfig{
					Definition: profiledefinition.ProfileDefinition{
						Name: "init_config_p1",
					},
				},
			},
			userProfiles: ProfileConfigMap{
				"user_p1": ProfileConfig{
					Definition: profiledefinition.ProfileDefinition{
						Name: "user_p1",
					},
				},
			},
			expectedErr: "passing both initConfigProfiles and userProfiles is not expected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			l, err := seelog.LoggerFromWriterWithMinLevelAndFormat(w, seelog.DebugLvl, "[%LEVEL] %FuncShort: %Msg")
			assert.Nil(t, err)
			log.SetupLogger(l, "debug")

			profilesWithInvalidExtendConfdPath, _ := filepath.Abs(filepath.Join("..", "test", tt.confdDir))
			config.Datadog.Set("confd_path", profilesWithInvalidExtendConfdPath)
			SetGlobalProfileConfigMap(nil)

			actualProfiles, err := loadProfiles(tt.initConfigProfiles, tt.userProfiles)
			assert.ElementsMatch(t, tt.expectedProfiles, convertToProfileNames(actualProfiles))
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			}

			w.Flush()
			logs := b.String()
			for _, expectedLog := range tt.expectedLogs {
				assert.Contains(t, logs, expectedLog)
			}
		})
	}
}
