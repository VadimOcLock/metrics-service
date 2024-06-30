package metricservice

import "github.com/VadimOcLock/metrics-service/internal/store/somestore"

type Service struct {
	Store somestore.Store
}

func New(s somestore.Store) Service {
	return Service{
		Store: s,
	}
}
