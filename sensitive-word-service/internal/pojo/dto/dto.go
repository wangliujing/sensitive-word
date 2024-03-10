package dto

type DetectResult struct {
	Text              string `json:"text"`
	ContainsSensitive bool   `json:"containsSensitive"`
	Keyword           []any  `json:"keyword"`
}
