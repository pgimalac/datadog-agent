// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package agent includes helpers related to the Datadog Agent on Windows
package agent

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner/parameters"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"

	"github.com/stretchr/testify/assert"
)

const (
	// DatadogCodeSignatureThumbprint is the thumbprint of the Datadog Code Signing certificate
	// Valid From: May 2023
	// Valid To:   May 2025
	DatadogCodeSignatureThumbprint = `B03F29CC07566505A718583E9270A6EE17678742`
)

// GetDatadogAgentProductCode returns the product code GUID for the Datadog Agent
func GetDatadogAgentProductCode(client client.VM) (string, error) {
	return windows.GetProductCodeByName(client, "Datadog Agent")
}

// UninstallAgent uninstalls the Datadog Agent
func UninstallAgent(client client.VM, logPath string) error {
	product, err := GetDatadogAgentProductCode(client)
	if err != nil {
		return err
	}
	return windows.UninstallMSI(client, product, logPath)
}

// HasValidDatadogCodeSignature an error if the file at the given path is not validy signed by the Datadog Code Signing certificate
func HasValidDatadogCodeSignature(client client.VM, path string) error {
	sig, err := windows.GetAuthenticodeSignature(client, path)
	if err != nil {
		return err
	}
	if !sig.Valid() {
		return fmt.Errorf("signature status is not valid: %s", sig.StatusMessage)
	}
	if !strings.EqualFold(sig.SignerCertificate.Thumbprint, DatadogCodeSignatureThumbprint) {
		return fmt.Errorf("signature thumbprint is not valid: %s", sig.SignerCertificate.Thumbprint)
	}
	return nil
}

// TestValidDatadogCodeSignatures verifies that the files at the given paths are validly signed by the Datadog Code Signing certificate
// This test is skipped if the verify_code_signature parameter is set to false.
func TestValidDatadogCodeSignatures(t *testing.T, client client.VM, paths []string) bool {
	return t.Run("code signatures", func(t *testing.T) {
		verify, _ := runner.GetProfile().ParamStore().GetBoolWithDefault(parameters.VerifyCodeSignature, true)

		if !verify {
			t.Skip("skipping code signature verification")
		}

		for _, path := range paths {
			err := HasValidDatadogCodeSignature(client, path)
			assert.NoError(t, err, "file %s is not validly signed: %s", path, err)
			// don't break early, check and report on all the files to make it easier to fix
		}
	})
}
