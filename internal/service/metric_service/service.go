package metric_service

import "github.com/VadimOcLock/metrics-service/internal/store/some_store"

type Service struct {
	Store some_store.Store
}

func New(s some_store.Store) Service {
	return Service{
		Store: s,
	}
}
