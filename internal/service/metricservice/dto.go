package metricservice

type UpsertGaugeDTO UpsertGaugeMetricParams

func (dto *UpsertGaugeDTO) Valid() error {
	// todo

	return nil
}

type UpsertCounterDTO UpsertCounterMetricParams

func (dto *UpsertCounterDTO) Valid() error {
	// todo

	return nil
}

type FindAllDTO FindAllMetricsParams

func (dto *FindAllDTO) Valid() error {
	// todo

	return nil
}

type FindDTO struct {
	MetricType string
	MetricName string
}

func (dto *FindDTO) Valid() error {
	// todo

	return nil
}
