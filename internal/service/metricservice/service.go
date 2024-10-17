package metricservice

type Service struct {
	Store Store
}

func New(s Store) Service {
	return Service{
		Store: s,
	}
}
