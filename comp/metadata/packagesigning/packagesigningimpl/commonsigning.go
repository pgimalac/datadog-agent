// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package packagesigningimpl

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-agent/comp/core/log"
	pgp "github.com/ProtonMail/go-crypto/openpgp"
)

// SigningKey represents relevant fields for a package signature key
type SigningKey struct {
	Fingerprint    string         `json:"signing_key_fingerprint"`
	ExpirationDate string         `json:"signing_key_expiration_date"`
	KeyType        string         `json:"signing_key_type"`
	Repositories   []repositories `json:"repositories"`
}

type repositories struct {
	RepoName string `json:"repo_name"`
}

type repoFile struct {
	filename     string
	repositories []repositories
}

const (
	aptPath    = "/etc/apt"
	yumPath    = "/etc/yum"
	dnfPath    = "/etc/dnf"
	zyppPath   = "/etc/zypp"
	noExpDate  = "9999-12-31"
	formatDate = "2006-01-02"
)

// getPackageManager is a lazy implementation to detect if we use APT or YUM (RH or SUSE)
func getPackageManager() string {
	if _, err := os.Stat(aptPath); !os.IsNotExist(err) {
		return "apt"
	} else if _, err := os.Stat(yumPath); !os.IsNotExist(err) {
		return "yum"
	} else if _, err := os.Stat(dnfPath); !os.IsNotExist(err) {
		return "dnf"
	} else if _, err := os.Stat(zyppPath); !os.IsNotExist(err) {
		return "zypper"
	}
	return ""
}

// decryptGPGFile parse a gpg file (local or http) and extract signing keys information
// Some files can contain a list of repositories.
func decryptGPGFile(allKeys map[string]SigningKey, gpgFile repoFile, keyType string, client *http.Client, logger log.Component) {
	var reader io.Reader
	if strings.HasPrefix(gpgFile.filename, "http") {
		response, err := client.Get(gpgFile.filename)
		if err != nil {
			return
		}
		defer response.Body.Close()
		reader = response.Body
	} else {
		file, err := os.Open(strings.Replace(gpgFile.filename, "file://", "", 1))
		if err != nil {
			return
		}
		defer file.Close()
		reader = file
	}
	err := decryptGPGReader(allKeys, reader, keyType, gpgFile.repositories)
	if err != nil {
		logger.Infof("Error while parsing gpg file %s: %s", gpgFile.filename, err)
	}
}

// decryptGPGReader extract keys from a reader, useful for rpm extraction
func decryptGPGReader(allKeys map[string]SigningKey, reader io.Reader, keyType string, repositories []repositories) error {
	keyList, err := pgp.ReadArmoredKeyRing(reader)
	if err != nil {
		// Try a non armored keyring
		keyList, err = pgp.ReadKeyRing(reader)
		if err != nil {
			return err
		}
	}
	for _, key := range keyList {
		fingerprint := key.PrimaryKey.KeyIdString()
		expDate := noExpDate
		i := key.PrimaryIdentity()
		keyLifetime := i.SelfSignature.KeyLifetimeSecs
		if keyLifetime != nil {
			expiry := key.PrimaryKey.CreationTime.Add(time.Duration(*i.SelfSignature.KeyLifetimeSecs) * time.Second)
			expDate = expiry.Format(formatDate)
		}
		allKeys[fingerprint] = SigningKey{
			Fingerprint:    fingerprint,
			ExpirationDate: expDate,
			KeyType:        keyType,
			Repositories:   repositories,
		}
	}
	return nil
}
