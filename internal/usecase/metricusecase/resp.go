package metricusecase

type UpdateResp struct {
	Message string `json:"message"`
}

type FindAllResp struct {
	HTML string
}

type FindResp struct {
	MetricValue string
}
