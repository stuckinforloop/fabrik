package id

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

func (s *Source) ULID() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.modeTest {
		return s.testULID()
	}

	ts := s.nowFunc().UnixMilli()
	if ts < 0 {
		return "", fmt.Errorf("invalid timestamp: %d", ts)
	}
	ms := uint64(ts)

	entropy := ulid.Monotonic(s.rnd, 0)
	id, err := ulid.New(ms, entropy)
	if err != nil {
		return "", fmt.Errorf("generate ulid: %w", err)
	}

	return id.String(), nil
}

func (s *Source) MustULID() string {
	id, err := s.ULID()
	if err != nil {
		panic(err)
	}
	return id
}

func (s *Source) testULID() (string, error) {
	ts := s.nowFunc().UnixMilli()
	if ts < 0 {
		return "", fmt.Errorf("invalid timestamp: %d", ts)
	}
	ms := uint64(ts)

	id := ulid.MustNew(ms, s.rnd)
	return id.String(), nil
}
