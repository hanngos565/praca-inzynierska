package structure

type Data struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Body struct {
	ID    string `json:"id"`
	Model string `json:"model"`
	Image string `json:"image"`
}

type Algorithm struct {
	Models map[string][]string `json:"models"`
}

type Images struct {
	Images []string `json:"images"`
}

type Results struct {
	Algorithm string `json:"algorithm"`
	Model     string `json:"model"`
	Image     string `json:"image"`
	Result    string `json:"result"`
	TimeStamp string `json:"timeStamp"`
	Status    string `json:"status"`
}

type DemoResult struct {
	Names  []string  `json:"names"`
	Scores []float64 `json:"scores"`
	BBox   [][]int64 `json:"bbox"`
}
