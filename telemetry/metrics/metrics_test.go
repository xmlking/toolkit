package metrics

import (
	"context"
	"testing"

	"github.com/xmlking/toolkit/telemetry"
)

func TestInitMetrics(t *testing.T) {
	InitMetrics(context.Background(), &telemetry.MetricsConfig{
		Enabled: true,
		Backend: "prometheus",
	})
}
