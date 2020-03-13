package model

type Paging struct {
	IsStart  bool   `json:"is_start"`
	IsEnd    bool   `json:"is_end"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Totals   int    `json:"totals"`
}
