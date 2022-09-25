package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestMetricsMarshal(t *testing.T) {
	hc := HealthCheck{
		Enabled:  true,
		Interval: 5 * time.Second,
	}
	data, err := yaml.Marshal(hc)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

	var hc2 HealthCheck
	err = yaml.Unmarshal(data, &hc2)
	if err != nil {
		t.Error(err)
	}
	assert.Exactly(t, hc, hc2)
}
