package metricservice

type UpdateGaugeDTO UpsertGaugeMetricParams

func (dto *UpdateGaugeDTO) Valid() error {
	// todo

	return nil
}

type UpdateCounterDTO UpsertCounterMetricParams

func (dto *UpdateCounterDTO) Valid() error {
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

type UpdateBatchDTO UpdateMetricsBatchTxParams

func (m *UpdateBatchDTO) Valid() error { return nil }
