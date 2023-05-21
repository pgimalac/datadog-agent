// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package languagedetection

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/process/procutil"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

type LanguageName string

var (
	python  LanguageName = "python"
	java    LanguageName = "java"
	unknown LanguageName = ""
)

type Language struct {
	Name LanguageName
}

type languageFromCLI struct {
	name      LanguageName
	validator func(exe string) bool
}

// knownPrefixes maps languages names to their prefix
var knownPrefixes = map[string]languageFromCLI{
	"python": {name: python},
	"java": {name: java, validator: func(exe string) bool {
		if exe == "javac" {
			return false
		}
		return true
	}},
}

// exactMatches maps an exact exe name match to a prefix
var exactMatches = map[string]languageFromCLI{
	"py": {name: python},
}

func languageNameFromCommandLine(cmdline []string) (LanguageName, error) {
	exe := getExe(cmdline)

	// First check to see if there is an exact match
	if lang, ok := exactMatches[exe]; ok {
		return lang.name, nil
	}

	for prefix, language := range knownPrefixes {
		if strings.HasPrefix(exe, prefix) {
			if language.validator != nil {
				isValidResult := language.validator(exe)
				if !isValidResult {
					continue
				}
			}
			return language.name, nil
		}
	}

	return unknown, fmt.Errorf("unknown executable: %q", exe)
}

// DetectLanguage uses a combination of commandline parsing and binary analysis to detect a process' language
func DetectLanguage(procs []*procutil.Process) []*Language {
	langs := make([]*Language, len(procs))
	for i, proc := range procs {
		languageName, err := languageNameFromCommandLine(proc.Cmdline)
		if err == nil {
			log.Trace("detected languageName:", languageName)
		}
		langs[i] = &Language{Name: languageName}
	}
	return langs
}
