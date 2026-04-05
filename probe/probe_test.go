package probe_test

import (
	"testing"

	"github.com/jefflunt/fcp-cli/probe"
)

func TestSecondsToTicks(t *testing.T) {
	cases := []struct {
		seconds float64
		fps     int
		want    int64
	}{
		{seconds: 10.0, fps: 30, want: 300},
		{seconds: 0.5, fps: 30, want: 15},
		{seconds: 1.0, fps: 24, want: 24},
		{seconds: 60.0, fps: 30, want: 1800},
		{seconds: 0.0, fps: 30, want: 0},
	}

	for _, tc := range cases {
		got := probe.SecondsToTicks(tc.seconds, tc.fps)
		if got != tc.want {
			t.Errorf("SecondsToTicks(%.2f, %d) = %d, want %d", tc.seconds, tc.fps, got, tc.want)
		}
	}
}
