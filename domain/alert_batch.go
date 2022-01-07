package domain

type AlertBatch struct {
	Scanner     string `json:"scanner"`
	ChainID     int64  `json:"chainId"`
	BlockStart  int64  `json:"blockStart"`
	BlockEnd    int64  `json:"blockEnd"`
	AlertCount  int64  `json:"alertCount"`
	MaxSeverity int64  `json:"maxSeverity"`
	Ref         string `json:"ref"`
}