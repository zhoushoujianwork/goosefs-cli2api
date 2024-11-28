package models

type SparkSubmitRequest struct {
	Args []string `json:"args" binding:"required"`
}

type SparkSubmitResponse struct {
	Results any    `json:"results"`
	Status  string `json:"status"`
}
