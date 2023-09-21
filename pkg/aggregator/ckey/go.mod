module github.com/DataDog/datadog-agent/pkg/aggregator/ckey

go 1.20

replace (
	github.com/DataDog/datadog-agent/pkg/tagset => ../../../pkg/tagset/
	github.com/DataDog/datadog-agent/pkg/util/util_sort => ../../../pkg/util/util_sort
)

require (
	github.com/DataDog/datadog-agent/pkg/tagset v0.0.0-00010101000000-000000000000
	github.com/DataDog/datadog-agent/pkg/util/util_sort v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	github.com/twmb/murmur3 v1.1.8
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)