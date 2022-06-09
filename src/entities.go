package src

type Text string

type Page struct {
	Url     string
	Indexes map[string]int
}
