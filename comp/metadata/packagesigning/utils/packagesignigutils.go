// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package utils defines shared methods in package signing component
package utils

import (
	"os"
	"runtime"
)

// GetLinuxPackageSigningPolicy returns the global GPG signing policy of the host
func GetLinuxPackageSigningPolicy() (bool, bool) {
	if runtime.GOOS == "linux" {
		pkgManager := GetPackageManager()
		switch pkgManager {
		case "apt":
			return getNoDebsig(), false
		case "yum", "dnf", "zypper":
			return getMainGPGCheck(pkgManager)
		default: // should not happen, tested above
			return false, false
		}
	}
	return false, false
}

const (
	aptPath  = "/etc/apt"
	yumPath  = "/etc/yum"
	dnfPath  = "/etc/dnf"
	zyppPath = "/etc/zypp"
)

// GetPackageManager is a lazy implementation to detect if we use APT or YUM (RH or SUSE)
func GetPackageManager() string {
	if _, err := os.Stat(aptPath); err == nil {
		return "apt"
	} else if _, err := os.Stat(yumPath); err == nil {
		return "yum"
	} else if _, err := os.Stat(dnfPath); err == nil {
		return "dnf"
	} else if _, err := os.Stat(zyppPath); err == nil {
		return "zypper"
	}
	return ""
}

// CompareRepoPerKeys is a method used on tests
func CompareRepoPerKeys(a, b map[string][]Repositories) []string {
	errorKeys := make([]string, 0)
	if len(a) < len(b) {
		for key := range b {
			if _, ok := a[key]; !ok {
				errorKeys = append(errorKeys, key)
			}
		}
	} else if len(a) > len(b) {
		for key := range a {
			if _, ok := b[key]; !ok {
				errorKeys = append(errorKeys, key)
			}
		}
	} else {
		errorKeys = append(errorKeys, compareKey(a, b)...)
		errorKeys = append(errorKeys, compareKey(b, a)...)
	}
	return errorKeys
}
func compareKey(a, b map[string][]Repositories) []string {
	errorKeys := make([]string, 0)
	for key := range a {
		if _, ok := b[key]; !ok {
			errorKeys = append(errorKeys, key)
		} else {
			if len(a[key]) == len(b[key]) {
				if anyMissingRepository(a[key], b[key]) {
					errorKeys = append(errorKeys, key)
				}
				if anyMissingRepository(b[key], a[key]) {
					errorKeys = append(errorKeys, key)
				}
			} else {
				errorKeys = append(errorKeys, key)
			}
		}
	}
	return errorKeys
}
func anyMissingRepository(r, s []Repositories) bool {
	for _, src := range r {
		found := false
		for _, dest := range s {
			if src.RepoName == dest.RepoName {
				found = true
				break
			}

		}
		if !found {
			return true
		}
	}
	return false
}
