package snmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFloat64FromString(t *testing.T) {
	snmpValue := &ResultValue{
		SubmissionType: "gauge",
		ResultValue:    "255.745",
	}
	value, err := snmpValue.toFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(255.745), value)
}

func TestToFloat64FromFloat(t *testing.T) {
	snmpValue := &ResultValue{
		SubmissionType: "gauge",
		ResultValue:    float64(255.745),
	}
	value, err := snmpValue.toFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(255.745), value)
}

func TestToFloat64FromInvalidType(t *testing.T) {
	snmpValue := &ResultValue{
		SubmissionType: "gauge",
		ResultValue:    int64(255),
	}
	_, err := snmpValue.toFloat64()
	assert.NotNil(t, err)
}
