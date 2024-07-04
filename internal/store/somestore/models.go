package somestore

import "github.com/VadimOcLock/metrics-service/internal/entity"

type Metric entity.Metric

func (m Metric) Entity() entity.Metric {
	return entity.Metric{
		Type:  m.Type,
		Name:  m.Name,
		Value: m.Value,
	}
}
