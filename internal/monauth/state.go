package monauth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"github.com/pkg/errors"
)

type state struct {
	state string
	sync.Mutex
}

func (s *state) Get() string {
	s.Lock()
	defer s.Unlock()
	return s.state
}

func (s *state) Set(state string) {
	s.Lock()
	defer s.Unlock()
	s.state = state
}

func generateStateValue() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "reading random")
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
