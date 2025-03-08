package id_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stuckinforloop/fabrik/deps/id"
	"github.com/stuckinforloop/fabrik/deps/timeutils"
)

var (
	testULIDOne = "01JMPX73EG06AFVGQT5ZYC0GEK"
	testULIDTwo = "01JMPX73EGZW908PVKS1Q4ZYAZ"
)

func TestULID(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, true)

	idOne, err := source.ULID()
	require.NoError(t, err)
	require.Equal(t, testULIDOne, idOne)

	idTwo, err := source.ULID()
	require.NoError(t, err)
	require.Equal(t, testULIDTwo, idTwo)
}

func TestULIDMustGenerate(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, true)

	idOne := source.MustULID()
	require.Equal(t, testULIDOne, idOne)

	idTwo := source.MustULID()
	require.Equal(t, testULIDTwo, idTwo)
}

func TestULIDNonDeterministic(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, false)

	idOne, err := source.ULID()
	require.NoError(t, err)
	require.NotEqual(t, testULIDOne, idOne)

	idTwo, err := source.ULID()
	require.NoError(t, err)
	require.NotEqual(t, testULIDTwo, idTwo)
}

func TestULIDNegativeTimestamp(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	var timeNow timeutils.TimeNow = func() time.Time {
		return time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	source := id.New(rnd, timeNow, true)

	_, err := source.ULID()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid timestamp")
}
