package pgstore

import (
	"context"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

const upsertGaugeMetric = `
insert into metrics (id, type, delta, value)
values ($1, 'gauge', null, $2)
on conflict (id)
do update set
    type = EXCLUDED.type,
    delta = EXCLUDED.delta,
    value = EXCLUDED.value
returning true as updated;
`

func (q *Queries) UpsertGaugeMetric(ctx context.Context, arg metricservice.UpsertGaugeMetricParams) (bool, error) {
	row := q.db.QueryRow(ctx, upsertGaugeMetric,
		arg.Name,
		arg.Value)
	var updated bool
	err := row.Scan(&updated)

	return updated, err
}

const upsertCounterMetric = `
insert into metrics (id, type, delta, value)
values ($1, 'counter', $2, null)
on conflict (id)
do update set
    type = EXCLUDED.type,
    delta = EXCLUDED.delta,
    value = EXCLUDED.value
returning true as updated;
`

func (q *Queries) UpsertCounterMetric(ctx context.Context, arg metricservice.UpsertCounterMetricParams) (bool, error) {
	row := q.db.QueryRow(ctx, upsertCounterMetric,
		arg.Name,
		arg.Value)
	var updated bool
	err := row.Scan(&updated)

	return updated, err
}

const findAllMetrics = `
select id, type, delta, value
from metrics;
`

func (q *Queries) FindAllMetrics(ctx context.Context, arg metricservice.FindAllMetricsNewParams) ([]entity.Metrics, error) {
	rows, err := q.db.Query(ctx, findAllMetrics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var metrics []entity.Metrics
	for rows.Next() {
		var m entity.Metrics
		if err = rows.Scan(
			&m.ID,
			&m.MType,
			&m.Delta,
			&m.Value,
		); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

const findCounterMetrics = `
select id, type, delta
from metrics
where id = $1;
`

func (q *Queries) FindCounterMetrics(ctx context.Context, arg metricservice.FindCounterMetricParams) (entity.Metrics, error) {
	row := q.db.QueryRow(ctx, findCounterMetrics, arg.MetricName)
	var m entity.Metrics
	err := row.Scan(
		&m.ID,
		&m.MType,
		&m.Delta)

	return m, err
}

const findGaugeMetrics = `
select id, type, value 
from metrics
where id = $1;
`

func (q *Queries) FindGaugeMetrics(ctx context.Context, arg metricservice.FindGaugeMetricParams) (entity.Metrics, error) {
	row := q.db.QueryRow(ctx, findGaugeMetrics, arg.MetricName)
	var m entity.Metrics
	err := row.Scan(
		&m.ID,
		&m.MType,
		&m.Value)

	return m, err
}
