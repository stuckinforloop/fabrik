package id

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"

	"github.com/stuckinforloop/fabrik/deps/timeutils"
)

type Source struct {
	rnd      *rand.Rand
	nowFunc  timeutils.TimeNow
	mutex    sync.Mutex
	modeTest bool
}

func New(rnd *rand.Rand, nowFunc timeutils.TimeNow, modeTest bool) *Source {
	return &Source{
		rnd:      rnd,
		nowFunc:  nowFunc,
		modeTest: modeTest,
	}
}

func (s *Source) Generate() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.modeTest {
		return s.test()
	}

	id, err := uuid.NewV7FromReader(s.rnd)
	if err != nil {
		return "", fmt.Errorf("new uuid: %w", err)
	}

	return id.String(), nil
}

func (s *Source) MustGenerate() string {
	id, err := s.Generate()
	if err != nil {
		panic(err)
	}

	return id
}

func (s *Source) test() (string, error) {
	id, err := uuid.NewRandomFromReader(s.rnd)
	if err != nil {
		return "", fmt.Errorf("new uuid: %w", err)
	}

	return id.String(), nil
}
