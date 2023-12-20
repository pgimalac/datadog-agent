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
	agentInstallerListProductName = "datadog-agent"
	agentS3BucketRelease          = "ddagent-windows-stable"
	agentS3BucketTesting          = "dd-agent-mstesting"
	betaURL                       = "https://s3.amazonaws.com/dd-agent-mstesting/builds/beta/installers_v2.json"
	stableURL                     = "https://s3.amazonaws.com/dd-agent-mstesting/builds/stable/installers_v2.json"
)

// GetBetaMSIURL returns the URL for the beta agent MSI
// majorVersion: 6, 7
// arch: x86_64
func GetBetaMSIURL(version string, arch string) (string, error) {
	return installers.GetProductURL(betaURL, agentInstallerListProductName, version, arch)
}

// GetStableMSIURL returns the URL for the stable agent MSI
// majorVersion: 6, 7
// arch: x86_64
func GetStableMSIURL(version string, arch string) (string, error) {
	return installers.GetProductURL(stableURL, agentInstallerListProductName, version, arch)
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

// GetMajorVersionFromEnv looks at environment variabes to select the agent major version.
//
// WINDOWS_AGENT_MAJOR_VERSION: The major version of the agent, 6 or 7
//
// Default major version: 7
func GetMajorVersionFromEnv() string {
	majorVersion := os.Getenv("WINDOWS_AGENT_MAJOR_VERSION")
	if majorVersion == "" {
		majorVersion = "7"
	}
	return majorVersion
}

// GetMSIURLFromEnv looks at environment variabes to select the agent MSI URL.
//
// The following environment variables select the agent version:
//
//   - WINDOWS_AGENT_MAJOR_VERSION: The major version of the agent, 6 or 7
//
//   - WINDOWS_AGENT_ARCH: The arch of the agent, x86_64
//
// The following environment variables select the package:
//
//   - WINDOWS_AGENT_MSI_URL: manually provided URL (version and arch are ignored)
//
//   - CI_PIPELINE_ID: use the URL from a specific CI pipeline
//
// If none of the above are set, the latest stable version is used.
func GetMSIURLFromEnv() (string, error) {
	// check for manually provided URL
	url := os.Getenv("WINDOWS_AGENT_MSI_URL")
	if url != "" {
		return url, nil
	}

	majorVersion := GetMajorVersionFromEnv()

	arch := os.Getenv("WINDOWS_AGENT_ARCH")
	if arch == "" {
		arch = "x86_64"
	}

	// check if we should use the URL from a specific CI pipeline
	pipelineID := os.Getenv("CI_PIPELINE_ID")
	if pipelineID != "" {
		url, err := GetPipelineMSIURL(pipelineID, majorVersion, arch)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	// Default to latest stable
	return GetLatestMSIURL(majorVersion, arch), nil
}
