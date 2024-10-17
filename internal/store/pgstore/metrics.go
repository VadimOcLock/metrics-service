package pgstore

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/rs/zerolog/log"
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
    delta = metrics.delta + EXCLUDED.delta,
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

const updateBatch = `
insert into metrics (id, type, delta, value)
values %s
ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, delta = EXCLUDED.delta, value = EXCLUDED.value;
`

func (s *PgStore) UpdateMetricsBatchTx(ctx context.Context, arg metricservice.UpdateMetricsBatchTxParams) error {
	return s.ExecTx(ctx, func(q *Queries) error {
		metrics := *arg.Data
		for _, m := range metrics {
			switch m.MType {
			case enum.GaugeMetricType:
				if m.Value != nil {
					_, err := q.UpsertGaugeMetric(ctx, metricservice.UpsertGaugeMetricParams{
						Name:  m.ID,
						Value: *m.Value,
					})
					if err != nil {
						log.Error().Msgf("Error upserting gauge metric: %v", err)
					}
				}
			case enum.CounterMetricType:
				if m.Delta != nil {
					_, err := q.UpsertCounterMetric(ctx, metricservice.UpsertCounterMetricParams{
						Name:  m.ID,
						Value: *m.Delta,
					})
					if err != nil {
						log.Error().Msgf("Error upserting counter metric: %v", err)
					}
				}
			}
		}
		return nil
	})
}

//func (s *PgStore) UpdateMetricsBatchTx(ctx context.Context, arg metricservice.UpdateMetricsBatchTxParams) error {
//	return s.ExecTx(ctx, func(q *Queries) error {
//		const batchSize = 100
//		metrics := *arg.Data
//		for i := 0; i < len(metrics); i += batchSize {
//			end := i + batchSize
//			if end > len(metrics) {
//				end = len(metrics)
//			}
//			batch := metrics[i:end]
//
//			values := make([]string, 0, len(batch))
//			args := make([]interface{}, 0, len(batch)*4)
//
//			uniqueMetrics := make(map[string]entity.Metrics)
//
//			for _, metric := range batch {
//				if _, exists := uniqueMetrics[metric.ID]; !exists {
//					uniqueMetrics[metric.ID] = metric
//					values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)",
//						len(args)+1, len(args)+2, len(args)+3, len(args)+4))
//					args = append(args, metric.ID, metric.MType, metric.Delta, metric.Value)
//				}
//			}
//			query := fmt.Sprintf(updateBatch, strings.Join(values, ", "))
//
//			if _, err := q.db.Exec(ctx, query, args...); err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//}
