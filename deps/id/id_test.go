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
	testIDOne = "0194fdc2-fa2f-4cc0-81d3-ff12045b73c8"
	testIDTwo = "6e4ff95f-f662-45ee-a82a-bdf44a2d0b75"
)

func TestID(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, true)

	idOne, err := source.Generate()
	require.NoError(t, err)
	require.Equal(t, testIDOne, idOne)

	idTwo, err := source.Generate()
	require.NoError(t, err)
	require.Equal(t, testIDTwo, idTwo)
}

func TestIDMustGenerate(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, true)

	idOne := source.MustGenerate()
	require.Equal(t, testIDOne, idOne)

	idTwo := source.MustGenerate()
	require.Equal(t, testIDTwo, idTwo)
}

func TestIDNonDeterministic(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var timeNow timeutils.TimeNow = func() time.Time {
		return timeutils.FoundingTimeUTC
	}
	source := id.New(rnd, timeNow, false)

	idOne, err := source.Generate()
	require.NoError(t, err)
	require.NotEqual(t, testIDOne, idOne)

	idTwo, err := source.Generate()
	require.NoError(t, err)
	require.NotEqual(t, testIDTwo, idTwo)
}
