package repository

import (
	"bytes"
	"context"
	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
	"goognion/src"
	"net/http"
	"sync"
	"time"
)

type TorCrawlerRepository struct {
}

func NewTorCrawlerRepository() *TorCrawlerRepository {
	return &TorCrawlerRepository{}
}

func (r *TorCrawlerRepository) Load(url string) ([]src.Text, []string, error) {
	t, err := tor.Start(nil, nil)
	if err != nil {
		return nil, nil, err
	}
	defer t.Close()

	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, nil, err
	}
	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	texts, urls := getTextsAndUrls(parsed, []src.Text{}, []string{})

	return texts, urls, nil
}

func (r *TorCrawlerRepository) DoIndexing(source []src.Text) (map[string]int, error) {
	out := map[string]int{}

	for p := range doIndexing(source) {
		v, ok := out[p.string]

		if !ok {
			out[p.string] = v
			continue
		}
		out[p.string] = v + p.int
	}
	//todo handle errors
	return out, nil
}

func getTextsAndUrls(n *html.Node, texts []src.Text, urls []string) ([]src.Text, []string) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
					break
				}
			}
		case "div":
			fallthrough
		case "p":
			fallthrough
		case "h1":
			fallthrough
		case "h2":
			fallthrough
		case "h3":
			fallthrough
		case "h4":
			fallthrough
		case "h5":
			fallthrough
		case "h6":
			var text bytes.Buffer
			if err := html.Render(&text, n.FirstChild); err != nil {
				panic(err) //todo remove panic
			}
			texts = append(texts, (src.Text)(text.String()))
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getTextsAndUrls(c, texts, urls)
	}

	return texts, urls
}

type pair struct {
	string
	int
}

func doIndexing(source []src.Text) <-chan pair {
	out := make(chan pair, len(source))

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(source))

		for _, v := range source {
			go func(t src.Text) {
				out <- index(t)
				wg.Done()
			}(v)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func index(text src.Text) pair {
	return pair{} //todo index
}
