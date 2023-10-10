package package1

import (
	"context"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/util/hostname"
	"github.com/stretchr/testify/require"
)

func TestHostname(t *testing.T) {
	_, err := hostname.Get(context.Background())
	require.NoError(t, err)
}
