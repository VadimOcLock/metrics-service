package metricusecase

import (
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
)

type MetricUseCase struct {
	metricService MetricService
	htmlBuilder   HTMLBuilder
}

var _ MetricService = (*metricservice.Service)(nil)

type Options struct {
	htmlBuilder HTMLBuilder
}

type OptFunc func(*Options)

func New(metricService MetricService, opts ...OptFunc) MetricUseCase {
	options := &Options{
		htmlBuilder: NewHTMLBuilder(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return MetricUseCase{
		metricService: metricService,
		htmlBuilder:   options.htmlBuilder,
	}
}

func WithHTMLBuilder(builder HTMLBuilder) OptFunc {
	return func(o *Options) {
		o.htmlBuilder = builder
	}
}
