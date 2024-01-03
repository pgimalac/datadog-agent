// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package agent includes helpers related to the Datadog Agent on Windows
package agent

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows/agent/installers/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	defaultMajorVersion           = "7"
	defaultArch                   = "x86_64"
	agentInstallerListProductName = "datadog-agent"
	agentS3BucketRelease          = "ddagent-windows-stable"
	agentS3BucketTesting          = "dd-agent-mstesting"
	betaChannel                   = "beta"
	betaURL                       = "https://s3.amazonaws.com/dd-agent-mstesting/builds/beta/installers_v2.json"
	stableChannel                 = "stable"
	stableURL                     = "https://s3.amazonaws.com/dd-agent-mstesting/builds/stable/installers_v2.json"
)

// GetBetaMSIURL returns the URL for the beta agent MSI
// majorVersion: 6, 7
// arch: x86_64
func GetBetaMSIURL(version string, arch string) (string, error) {
	return GetMSIURL(betaChannel, version, arch)
}

// GetStableMSIURL returns the URL for the stable agent MSI
// majorVersion: 6, 7
// arch: x86_64
func GetStableMSIURL(version string, arch string) (string, error) {
	return GetMSIURL(stableChannel, version, arch)
}

// GetMSIURL returns the URL for the agent MSI
// channel: beta, stable
// majorVersion: 6, 7
// arch: x86_64
func GetMSIURL(channel string, version string, arch string) (string, error) {
	channelURL, err := GetChannelURL(channel)
	if err != nil {
		return "", err
	}

	return installers.GetProductURL(channelURL, agentInstallerListProductName, version, arch)
}

// GetChannelURL returns the URL for the channel name
// channel: beta, stable
func GetChannelURL(channel string) (string, error) {
	if strings.EqualFold(channel, betaChannel) {
		return betaURL, nil
	} else if strings.EqualFold(channel, stableChannel) {
		return stableURL, nil
	}

	return "", fmt.Errorf("unknown channel %v", channel)
}

// GetLatestMSIURL returns the URL for the latest agent MSI
// majorVersion: 6, 7
// arch: x86_64
func GetLatestMSIURL(majorVersion string, arch string) string {
	// why do we use amd64 for the latest URL and x86_64 everywhere else?
	if arch == "x86_64" {
		arch = "amd64"
	}
	return fmt.Sprintf(`https://s3.amazonaws.com/`+agentS3BucketRelease+`/datadog-agent-%s-latest.%s.msi`,
		majorVersion, arch)
}

// GetPipelineMSIURL returns the URL for the agent MSI built by the pipeline
// majorVersion: 6, 7
// arch: x86_64
func GetPipelineMSIURL(pipelineID string, majorVersion string, arch string) (string, error) {
	// dd-agent-mstesting is a public bucket so we can use anonymous credentials
	config, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithCredentialsProvider(aws.AnonymousCredentials{}))
	if err != nil {
		return "", err
	}

	s3Client := s3.NewFromConfig(config)

	// Manual URL example: https://s3.amazonaws.com/dd-agent-mstesting?prefix=pipelines/A7/25309493
	result, err := s3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(agentS3BucketTesting),
		Prefix: aws.String(fmt.Sprintf("pipelines/A%s/%s", majorVersion, pipelineID)),
	})

	if err != nil {
		return "", err
	}

	if len(result.Contents) <= 0 {
		return "", fmt.Errorf("no agent MSI found for pipeline %v", pipelineID)
	}

	// match the arch
	for _, obj := range result.Contents {
		if strings.Contains(*obj.Key, arch) {
			return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", agentS3BucketTesting, *obj.Key), nil
		}
	}

	return "", fmt.Errorf("no agent MSI found for pipeline %v and arch %v", pipelineID, arch)
}

// LookupChannelFromEnv looks at environment variabes to select the agent channel, if the value
// is found it is returned along with true, otherwise a default value and false are returned.
//
// WINDOWS_AGENT_CHANNEL: beta, stable
//
// default is stable channel
func LookupChannelFromEnv() (string, bool) {
	channel := os.Getenv("WINDOWS_AGENT_CHANNEL")
	if channel != "" {
		return channel, true
	}
	return stableChannel, false
}

// LookupVersionFromEnv looks at environment variabes to select the agent version, if the value
// is found it is returned along with true, otherwise a default value and false are returned.
//
// In order of priority:
//
// WINDOWS_AGENT_VERSION: The complete version, e.g. 7.49.0-1, 7.49.0-rc.3-1, or a major version, e.g. 7
//
// AGENT_MAJOR_VERSION: The major version of the agent, 6 or 7
//
// If only a major version is provided, the latest version of that major version is used.
//
// Default version: 7
func LookupVersionFromEnv() (string, bool) {
	version := os.Getenv("WINDOWS_AGENT_VERSION")
	if version != "" {
		return version, true
	}

	// Currently commonly used in CI, should we keep it or transition to WINDOWS_AGENT_VERSION?
	version = os.Getenv("AGENT_MAJOR_VERSION")
	if version != "" {
		return version, true
	}

	return defaultMajorVersion, false
}

// LookupArchFromEnv looks at environment variabes to select the agent arch, if the value
// is found it is returned along with true, otherwise a default value and false are returned.
//
// WINDOWS_AGENT_ARCH: The arch of the agent, x86_64
//
// Default arch: x86_64
func LookupArchFromEnv() (string, bool) {
	arch := os.Getenv("WINDOWS_AGENT_ARCH")
	if arch != "" {
		return arch, true
	}
	return defaultArch, false
}

// LookupChannelURLFromEnv looks at environment variabes to select the agent channel URL, if the value
// is found it is returned along with true, otherwise a default value and false are returned.
//
// WINDOWS_AGENT_CHANNEL_URL: URL to installers_v2.json
//
// See also LookupChannelFromEnv()
//
// default is stable channel
func LookupChannelURLFromEnv() (string, bool) {
	channelURL := os.Getenv("WINDOWS_AGENT_CHANNEL_URL")
	if channelURL != "" {
		return channelURL, true
	}

	channel, _ := LookupChannelFromEnv()
	channelURL, err := GetChannelURL(channel)
	if err != nil {
		return channelURL, true
	}

	return stableURL, false
}

// GetMSIURLFromEnv looks at environment variabes to select the agent MSI URL.
//
// The channel, version, and arch parameters are optional, if not provided they are
// read from the environment.
//
// Primary environment variables in order of priority:
//
// WINDOWS_AGENT_MSI_URL: manually provided URL (all other parameters are ignored)
//
// CI_PIPELINE_ID: use the URL from a specific CI pipeline, major version and arch are used, channel is ignored
//
// WINDOWS_AGENT_VERSION: The complete version, e.g. 7.49.0-1, 7.49.0-rc.3-1, or a major version, e.g. 7, arch and channel are used
//
// Other environment variables:
//
// WINDOWS_AGENT_CHANNEL: beta or stable
//
// WINDOWS_AGENT_ARCH: The arch of the agent, x86_64
//
// Since not all versions in the beta channel contain `rc`, the version is not used to assume
// the channel. To use a beta version, set the channel option.
//
// See other Get*FromEnv functions for more options and details.
//
// If none of the above are set, the latest stable version is used.
func GetMSIURLFromEnv(channel string, version string, arch string) (string, error) {
	// check for manually provided MSI URL
	url := os.Getenv("WINDOWS_AGENT_MSI_URL")
	if url != "" {
		return url, nil
	}

	// If direct URL is not provided, see if a version was provided
	if version == "" {
		version, _ = LookupVersionFromEnv()
	}

	if arch == "" {
		arch, _ = LookupArchFromEnv()
	}

	majorVersion := strings.Split(version, ".")[0]

	// check if we should use the URL from a specific CI pipeline
	pipelineID := os.Getenv("CI_PIPELINE_ID")
	if pipelineID != "" {
		url, err := GetPipelineMSIURL(pipelineID, majorVersion, arch)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	// if version is a complete version, e.g. 7.49.0-1, use it as is
	if strings.Contains(version, ".") {
		channelURL, err := selectChannelWithVersion(channel, version)
		if err != nil {
			return "", err
		}
		return installers.GetProductURL(channelURL, agentInstallerListProductName, version, arch)
	}

	// Default to latest stable
	return GetLatestMSIURL(majorVersion, arch), nil
}

// selectChannelWithVersion returns the channel URL based on the provided channel and version. If
// a channel is provided, it is used. If a channel is not provided, the version is used to determine
// the channel. If the version contains `-rc.`, the beta channel is used, otherwise the stable channel is used.
func selectChannelWithVersion(channel string, version string) (string, error) {
	if channel != "" {
		// if channel name is provided, lookup its URL
		channelURL, err := GetChannelURL(channel)
		if err != nil {
			return "", err
		}
		return channelURL, nil
	}

	channelURL, found := LookupChannelURLFromEnv()
	if found {
		return channelURL, nil
	}

	// if channel is not provided, check if we can infer it from the version,
	// If version contains `-rc.`, assume it is a beta version
	if strings.Contains(strings.ToLower(version), `-rc.`) {
		channelURL, err := GetChannelURL(betaChannel)
		if err != nil {
			return "", err
		}
		return channelURL, nil
	}

	// if not then the returned default is used.
	return channelURL, nil
}
