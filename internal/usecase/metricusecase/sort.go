package metricusecase

import (
	"sort"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

func SortMetrics(metrics *[]entity.Metric) {
	sort.Slice(*metrics, func(i, j int) bool {
		if (*metrics)[i].Type == (*metrics)[j].Type {
			return (*metrics)[i].Name < (*metrics)[j].Name
		}

		return (*metrics)[i].Type < (*metrics)[j].Type
	})
}
