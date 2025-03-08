package timeutils_test

import (
	"testing"
	"time"

	"github.com/stuckinforloop/fabrik/deps/timeutils"

	"github.com/stretchr/testify/require"
)

func TestTimeNow(t *testing.T) {
	// Mock time.Now() to return a fixed time
	tn := timeutils.TimeNow(func() time.Time {
		return timeutils.FoundingTimeUTC
	})

	// Test that the time returned by tn is the same as the founding time
	require.Equal(t, int64(1740230594), tn().Unix())

	require.Equal(t, int64(1740230654), tn.Tick(time.Minute).Unix())
}
