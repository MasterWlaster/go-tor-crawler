package src

type Text string

type Index struct {
	Word  string
	Count int
}

type Page struct {
	Url     string
	Indexes map[string]int
}
