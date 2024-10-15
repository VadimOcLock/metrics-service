package somestore

type Impl struct {
	s MemStorage
}

func New() Impl {
	return Impl{
		s: MemStorage{
			gauges:   make(map[string]float64),
			counters: make(map[string]int64),
		},
	}
}
