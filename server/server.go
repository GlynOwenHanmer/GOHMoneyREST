package server

import (
	"net/http"

	"github.com/glynternet/accounting-rest/pkg/storage"
	"github.com/pkg/errors"
)

func New(store storage.Storage) (*server, error) {
	if store == nil {
		return nil, errors.New("nil store")
	}
	return &server{storage: store}, nil
}

type server struct {
	storage storage.Storage
}

func (s *server) ListenAndServe(addr string) error {
	//logDBState()  //TODO: logDBState?
	router, err := s.newRouter()
	if err != nil {
		return errors.Wrap(err, "creating new Router")
	}
	return http.ListenAndServe(addr, router)
}
