// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux

package packagesigningimpl

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/comp/core/log"
	"github.com/DataDog/datadog-agent/comp/core/log/logimpl"
	pkgUtils "github.com/DataDog/datadog-agent/comp/metadata/packagesigning/utils"
	"github.com/DataDog/datadog-agent/pkg/serializer"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

const (
	publicKeyWithoutExpiration = `The following public key can be used to verify RPM packages built and
signed by Red Hat, Inc.  This key is used for packages in Red Hat
products shipped after November 2009, and for all updates to those
products.

Questions about this key should be sent to security@redhat.com.

pub  4096R/FD431D51 2009-10-22 Red Hat, Inc. (release key 2) <security@redhat.com>

-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBErgSTsBEACh2A4b0O9t+vzC9VrVtL1AKvUWi9OPCjkvR7Xd8DtJxeeMZ5eF
0HtzIG58qDRybwUe89FZprB1ffuUKzdE+HcL3FbNWSSOXVjZIersdXyH3NvnLLLF
0DNRB2ix3bXG9Rh/RXpFsNxDp2CEMdUvbYCzE79K1EnUTVh1L0Of023FtPSZXX0c
u7Pb5DI5lX5YeoXO6RoodrIGYJsVBQWnrWw4xNTconUfNPk0EGZtEnzvH2zyPoJh
XGF+Ncu9XwbalnYde10OCvSWAZ5zTCpoLMTvQjWpbCdWXJzCm6G+/hx9upke546H
5IjtYm4dTIVTnc3wvDiODgBKRzOl9rEOCIgOuGtDxRxcQkjrC+xvg5Vkqn7vBUyW
9pHedOU+PoF3DGOM+dqv+eNKBvh9YF9ugFAQBkcG7viZgvGEMGGUpzNgN7XnS1gj
/DPo9mZESOYnKceve2tIC87p2hqjrxOHuI7fkZYeNIcAoa83rBltFXaBDYhWAKS1
PcXS1/7JzP0ky7d0L6Xbu/If5kqWQpKwUInXtySRkuraVfuK3Bpa+X1XecWi24JY
HVtlNX025xx1ewVzGNCTlWn1skQN2OOoQTV4C8/qFpTW6DTWYurd4+fE0OJFJZQF
buhfXYwmRlVOgN5i77NTIJZJQfYFj38c/Iv5vZBPokO6mffrOTv3MHWVgQARAQAB
tDNSZWQgSGF0LCBJbmMuIChyZWxlYXNlIGtleSAyKSA8c2VjdXJpdHlAcmVkaGF0
LmNvbT6JAjYEEwECACAFAkrgSTsCGwMGCwkIBwMCBBUCCAMEFgIDAQIeAQIXgAAK
CRAZni+R/UMdUWzpD/9s5SFR/ZF3yjY5VLUFLMXIKUztNN3oc45fyLdTI3+UClKC
2tEruzYjqNHhqAEXa2sN1fMrsuKec61Ll2NfvJjkLKDvgVIh7kM7aslNYVOP6BTf
C/JJ7/ufz3UZmyViH/WDl+AYdgk3JqCIO5w5ryrC9IyBzYv2m0HqYbWfphY3uHw5
un3ndLJcu8+BGP5F+ONQEGl+DRH58Il9Jp3HwbRa7dvkPgEhfFR+1hI+Btta2C7E
0/2NKzCxZw7Lx3PBRcU92YKyaEihfy/aQKZCAuyfKiMvsmzs+4poIX7I9NQCJpyE
IGfINoZ7VxqHwRn/d5mw2MZTJjbzSf+Um9YJyA0iEEyD6qjriWQRbuxpQXmlAJbh
8okZ4gbVFv1F8MzK+4R8VvWJ0XxgtikSo72fHjwha7MAjqFnOq6eo6fEC/75g3NL
Ght5VdpGuHk0vbdENHMC8wS99e5qXGNDued3hlTavDMlEAHl34q2H9nakTGRF5Ki
JUfNh3DVRGhg8cMIti21njiRh7gyFI2OccATY7bBSr79JhuNwelHuxLrCFpY7V25
OFktl15jZJaMxuQBqYdBgSay2G0U6D1+7VsWufpzd/Abx1/c3oi9ZaJvW22kAggq
dzdA27UUYjWvx42w9menJwh/0jeQcTecIUd0d0rFcw/c1pvgMMl/Q73yzKgKYw==
=zbHE
-----END PGP PUBLIC KEY BLOCK-----
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGIpIp4BEAC/o5e1WzLIsS6/JOQCs4XYATYTcf6B6ALzcP05G0W3uRpUQSrL
FRKNrU8ZCelm/B+XSh2ljJNeklp2WLxYENDOsftDXGoyLr2hEkI5OyK267IHhFNJ
g+BN+T5Cjh4ZiiWij6o9F7x2ZpxISE9M4iI80rwSv1KOnGSw5j2zD2EwoMjTVyVE
/t3s5XJxnDclB7ZqL+cgjv0mWUY/4+b/OoRTkhq7b8QILuZp75Y64pkrndgakm1T
8mAGXV02mEzpNj9DyAJdUqa11PIhMJMxxHOGHJ8CcHZ2NJL2e7yJf4orTj+cMhP5
LzJcVlaXnQYu8Zkqa0V6J1Qdj8ZXL72QsmyicRYXAtK9Jm5pvBHuYU2m6Ja7dBEB
Vkhe7lTKhAjkZC5ErPmANNS9kPdtXCOpwN1lOnmD2m04hks3kpH9OTX7RkTFUSws
eARAfRID6RLfi59B9lmAbekecnsMIFMx7qR7ZKyQb3GOuZwNYOaYFevuxusSwCHv
4FtLDIhk+Fge+EbPdEva+VLJeMOb02gC4V/cX/oFoPkxM1A5LHjkuAM+aFLAiIRd
Np/tAPWk1k6yc+FqkcDqOttbP4ciiXb9JPtmzTCbJD8lgH0rGp8ufyMXC9x7/dqX
TjsiGzyvlMnrkKB4GL4DqRFl8LAR02A3846DD8CAcaxoXggL2bJCU2rgUQARAQAB
tDVSZWQgSGF0LCBJbmMuIChhdXhpbGlhcnkga2V5IDMpIDxzZWN1cml0eUByZWRo
YXQuY29tPokCUgQTAQgAPBYhBH5GJCWMQGU11W1vE1BU5KRaY0CzBQJiKSKeAhsD
BQsJCAcCAyICAQYVCgkICwIEFgIDAQIeBwIXgAAKCRBQVOSkWmNAsyBfEACuTN/X
YR+QyzeRw0pXcTvMqzNE4DKKr97hSQEwZH1/v1PEPs5O3psuVUm2iam7bqYwG+ry
EskAgMHi8AJmY0lioQD5/LTSLTrM8UyQnU3g17DHau1NHIFTGyaW4a7xviU4C2+k
c6X0u1CPHI1U4Q8prpNcfLsldaNYlsVZtUtYSHKPAUcswXWliW7QYjZ5tMSbu8jR
OMOc3mZuf0fcVFNu8+XSpN7qLhRNcPv+FCNmk/wkaQfH4Pv+jVsOgHqkV3aLqJeN
kNUnpyEKYkNqo7mNfNVWOcl+Z1KKKwSkIi3vg8maC7rODsy6IX+Y96M93sqYDQom
aaWue2gvw6thEoH4SaCrCL78mj2YFpeg1Oew4QwVcBnt68KOPfL9YyoOicNs4Vuu
fb/vjU2ONPZAeepIKA8QxCETiryCcP43daqThvIgdbUIiWne3gae6eSj0EuUPoYe
H5g2Lw0qdwbHIOxqp2kvN96Ii7s1DK3VyhMt/GSPCxRnDRJ8oQKJ2W/I1IT5VtiU
zMjjq5JcYzRPzHDxfVzT9CLeU/0XQ+2OOUAiZKZ0dzSyyVn8xbpviT7iadvjlQX3
CINaPB+d2Kxa6uFWh+ZYOLLAgZ9B8NKutUHpXN66YSfe79xFBSFWKkJ8cSIMk13/
Ifs7ApKlKCCRDpwoDqx/sjIaj1cpOfLHYjnefg==
=UZd/
-----END PGP PUBLIC KEY BLOCK-----`
	datadogPublicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1.4.11 (GNU/Linux)

mQINBFd0Tp4BEACr6XBAdWWDQ0tc5ykBRPyI9lbUuzDFOKFLjAxblAcfJ66rzGL8
Sv+XROCDY5eFm7qAGYhUYFNbbrch0A3/W4SklJVXUFdUEA9EZu6+typRHYs058TR
vk7gyRvBLdfSrqxzwa4i9RNB86o45Kgm0bPDF/v1vnZJeKXNqRHI+HgxhgIE3XaN
MaT8TGIDHjSoqk5Do9ZtnBxRKdW5dcFkpYPphJ6kxfmF8LiNSIyr7ve6ky9z9JVT
lYPWliakK5mk4eyltKH+crVu/Ls3lKkWOvYusGcQqj5OiPz2mnyapUG2x253dlmi
4GDMgq3LwflTQvMlGm7jOL/gwj1zlisFOuasB8lHakPptU784Mpp5h5F9FhGMaPh
ynsB/i3pXMfUv/s5ylZjkT/Csko2hSK+WEXBxoQbAqc+nuq2f2n80eTK3mfDtn9d
AG1oHsSeMUWPax6VGNQIp79gwzn6gMwDRrfykoPVZw1IKHB6rIxT3WfDbBiGnSPW
o8eCC98dW3cRUr6hFY6XB//IzgUScU+QHwh5byjk0vMkMVjGSEeBwns6FWkwEVxq
Lr0yhihHYZif8ETt4TEfnIrGkgFFkTtSbGBG8Hf22RQqyMoJBHwxDHD/2Rorp5DI
to/8bC9BVHRBDUMIc7ahmSVq3ozvAISwzdvRCv1AKlAMFaPrAPUzDxSrzwARAQAB
tCREYXRhZG9nLCBJbmMgPHBhY2thZ2VAZGF0YWRvZ2hxLmNvbT6JAj4EEwECACgF
Ald0Tp4CGwMFCQtHNQAGCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEPEGjhTg
lCKzw44P/R9AgSzXUwTFegdZjpR3ocFhOnV5xliZHlsCeX4wqYWyO5g0ZqSaSUHc
cExnHFP0Imu1iAkCm+OPs9VicqyT+a7dllfCvv3lgqKYxk29lzRsUsqPe0iWzPLv
ZI+p4vOb1oXzJH/+rYMVuHcHmtQLe2AVKFH4UtaFIJ+uBeXyC8uWKamJPykpq5Fs
l08TWMjsir8vPNVx1pitLHs/XgrshNatsFg+7c//nqpckJbNYG240vHdbn1VJN7v
yHImZywJS7N7y607adj6mRoLrF+kIl040hheZbU9A67qRed8XY7gAj5MgqcipXSl
FvqepnUvzvkv+v2xGo1ohPCNPhW5nZu1uVU566wDzJaoUQILxmAX8wIxqFNewEp7
N7DCcI8gUSgX0hKZz2VPXlI2/X3ODu8Eu5+3JNqzXFerj4o209uv7sKaIkxlVPno
cCqmth1e8v8JKQw5+2fvtVa7TavXoG33+KKCC/MgP9uOppG/GLns/2P45K+QA5D7
kiHM/eaABqnwdyg1XNiPbf1GqyT4bLa5CfWJ+RrgoUiZ4zfYrA1NbFFafcRkksSy
nHqwA45eQzOye+95moLpip69fGWCX12OxjAH2jFieUX4yxHczYdc/CqyH9eyoKxL
3wpMulnDuv4BQOfxyWTRoZBPmbOAMWBCiVSiwqiTTPjnsTfDXXyQ
=5JfY
-----END PGP PUBLIC KEY BLOCK-----`
)

func TestDecryptGPGReader(t *testing.T) {
	// Example data for testing
	keys := make(map[string]SigningKey)

	testCases := []struct {
		name    string
		reader  io.Reader
		armored bool
		keyType string
		output  SigningKey
	}{
		{
			name:    "Key without expiration",
			reader:  strings.NewReader(publicKeyWithoutExpiration),
			keyType: "RSA",
			armored: true,
			output: SigningKey{
				Fingerprint:    "199E2F91FD431D51",
				ExpirationDate: "9999-12-31",
				KeyType:        "RSA",
			},
		},
		{
			name:    "Datadog key with expiration date",
			reader:  strings.NewReader(datadogPublicKey),
			keyType: "RSA",
			armored: true,
			output: SigningKey{
				Fingerprint:    "F1068E14E09422B3",
				ExpirationDate: "2022-06-28",
				KeyType:        "RSA",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			decryptGPGReader(keys, testCase.reader, testCase.armored, testCase.keyType, nil)

			retrieved, ok := keys[testCase.output.Fingerprint]
			if !ok || !compareKeys(retrieved, testCase.output) {
				t.Errorf("Expected key %s|%s to be present in the map", testCase.output.Fingerprint, testCase.output.ExpirationDate)
				t.Logf("Key %s|%s", retrieved.Fingerprint, retrieved.ExpirationDate)
			}
		})

	}
}

func compareKeys(a, b SigningKey) bool {
	if a.Fingerprint != b.Fingerprint {
		return false
	}
	if a.ExpirationDate != b.ExpirationDate {
		return false
	}
	if a.KeyType != b.KeyType {
		return false
	}
	if a.Repositories == nil && b.Repositories == nil {
		for idx, repo := range a.Repositories {
			if repo.RepoName != b.Repositories[idx].RepoName {
				return false
			}
		}
	}
	return true
}

func TestUpdateWithRepoFile(t *testing.T) {
	t.Cleanup(func() {
		pkgUtils.YumConf = "/etc/yum.conf"
		pkgUtils.YumRepo = "/etc/yum.repos.d/"
	})
	pkgUtils.YumConf = "testdata/yum.conf"
	pkgUtils.YumRepo = "testdata/yum.repos.d/"
	testCases := []struct {
		name    string
		allKeys map[string]SigningKey
	}{
		{
			name: "Update with repo files",
			allKeys: map[string]SigningKey{"32637D44F14F620E": {Fingerprint: "32637D44F14F620E", ExpirationDate: "2032-09-05", KeyType: "repo", Repositories: []pkgUtils.Repositories{{RepoName: "https://versailles.com"}}},
				"E6266D4AC0962C7D": {Fingerprint: "E6266D4AC0962C7D", ExpirationDate: "2028-04-18", KeyType: "repo", Repositories: []pkgUtils.Repositories{{RepoName: "https://versailles.com"}}},
				"D3A80E30382E94DE": {Fingerprint: "D3A80E30382E94DE", ExpirationDate: "2022-06-28", KeyType: "repo", Repositories: []pkgUtils.Repositories{{RepoName: "https://versailles.com"}}},
			},
		},
	}
	allKeys := make(map[string]SigningKey)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			updateWithRepoFiles(allKeys, "yum", nil, nil)
			for key := range allKeys {
				if _, ok := testCase.allKeys[key]; !ok {
					t.Errorf("Unexpected key %s", key)
				}
			}
		})
	}
}

func TestGetDebsigPath(t *testing.T) {
	t.Cleanup(func() {
		debsigPolicies = "/etc/debsig/policies/"
		debsigKeyring = "/usr/share/debsig/keyrings/"
	})

	debsigPolicies = "testdata/debsig/policies"
	debsigKeyring = "testdata/debsig/keyrings"
	testCases := []struct {
		name  string
		files []string
	}{
		{
			name:  "Find debsigfiles",
			files: []string{"testdata/debsig/keyrings/F1E2D3C4B5/debsig.gpg"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			debsigFiles := getDebsigKeyPaths()
			for idx, file := range debsigFiles {
				if file != testCase.files[idx] {
					t.Errorf("Expected file %s, got %s", testCase.files[idx], file)
				}
			}
		})

	}

}

func TestParseSourceListFile(t *testing.T) {
	testCases := []struct {
		name        string
		fileName    string
		reposPerKey map[string][]pkgUtils.Repositories
	}{
		{
			name:     "Source list file with several repo config",
			fileName: "testdata/datadog.list",
			reposPerKey: map[string][]pkgUtils.Repositories{"/usr/share/keyrings/datadog-archive-keyring.gpg": {{RepoName: "https://apt.datadoghq.com//stable/7"}, {RepoName: "https://apt.datadoghq.com//stable/6"}, {RepoName: "https://apt.datadoghq.com//beta/7"}},
				"/usr/vinz/clortho/keyring.gpg": {{RepoName: "https://apt.ghostbusters.com//stable/84"}},
				"/don/rosa/carl/barks":          {{RepoName: "https://duck.family.com/scrooge/donald/huey/dewey/louie"}},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reposPerKey := parseSourceListFile(testCase.fileName)
			errorData := pkgUtils.CompareRepoPerKeys(reposPerKey, testCase.reposPerKey)
			if len(errorData) > 0 {
				for _, key := range errorData {
					if _, ok := testCase.reposPerKey[key]; !ok {
						t.Errorf("Unexpected key %s", key)
					} else {
						if _, ok := reposPerKey[key]; !ok {
							t.Errorf("Missing key %s", key)
						} else {
							t.Errorf("Wrong key %s expected %v got %v", key, testCase.reposPerKey[key], reposPerKey[key])
						}
					}
				}
			}
		})
	}
}

func TestGetAPTPayload(t *testing.T) {
	setupAPTSigningMock(t)

	expectedMetadata := &signingMetadata{
		SigningKeys: []SigningKey{
			{Fingerprint: "F1068E14E09422B3", ExpirationDate: "2022-06-28", KeyType: "signed-by", Repositories: []pkgUtils.Repositories{{RepoName: "https://apt.datadoghq.com//stable/7"}}},
			{Fingerprint: "FD4BF915", ExpirationDate: "9999-12-31", KeyType: "trusted"},
		},
	}

	ih := getTestPackageSigning(t)

	p := ih.getPayload().(*Payload)
	assert.Equal(t, expectedMetadata, p.Metadata)
}

func TestGetYUMPayload(t *testing.T) {
	setupYUMSigningMock(t)

	expectedMetadata := &signingMetadata{
		SigningKeys: []SigningKey{
			{Fingerprint: "AL1C1AK3YS", ExpirationDate: "9999-12-31", KeyType: "repo", Repositories: []pkgUtils.Repositories{{RepoName: "https://yum.datadoghq.com/stable/7/x86_64/"}}},
			{Fingerprint: "733142A241337", ExpirationDate: "2030-03-02", KeyType: "rpm"},
		},
	}

	ih := getTestPackageSigning(t)

	p := ih.getPayload().(*Payload)
	assert.Equal(t, expectedMetadata, p.Metadata)
}

func setupAPTSigningMock(t *testing.T) {
	t.Cleanup(func() {
		getPkgManager = pkgUtils.GetPackageManager
		getAPTKeys = getAPTSignatureKeys
		getYUMKeys = getYUMSignatureKeys
	})

	getPkgManager = getPackageAPTMock
	getAPTKeys = getAPTKeysMock
}
func setupYUMSigningMock(t *testing.T) {
	setupAPTSigningMock(t)

	getPkgManager = getPackageYUMMock
	getYUMKeys = getYUMKeysMock
}
func getPackageAPTMock() string { return "apt" }
func getPackageYUMMock() string { return "yum" }
func getAPTKeysMock(_ *http.Client, _ log.Component) []SigningKey {
	return []SigningKey{
		{Fingerprint: "F1068E14E09422B3", ExpirationDate: "2022-06-28", KeyType: "signed-by", Repositories: []pkgUtils.Repositories{{RepoName: "https://apt.datadoghq.com//stable/7"}}},
		{Fingerprint: "FD4BF915", ExpirationDate: "9999-12-31", KeyType: "trusted"},
	}
}
func getYUMKeysMock(_ string, _ *http.Client, _ log.Component) []SigningKey {
	return []SigningKey{
		{Fingerprint: "AL1C1AK3YS", ExpirationDate: "9999-12-31", KeyType: "repo", Repositories: []pkgUtils.Repositories{{RepoName: "https://yum.datadoghq.com/stable/7/x86_64/"}}},
		{Fingerprint: "733142A241337", ExpirationDate: "2030-03-02", KeyType: "rpm"},
	}
}

func getTestPackageSigning(t *testing.T) *pkgSigning {
	p := newPackageSigningProvider(
		fxutil.Test[dependencies](
			t,
			logimpl.MockModule(),
			config.MockModule(),
			fx.Provide(func() serializer.MetricSerializer { return &serializer.MockSerializer{} }),
		),
	)
	return p.Comp.(*pkgSigning)
}