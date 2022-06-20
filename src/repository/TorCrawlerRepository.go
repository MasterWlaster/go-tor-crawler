package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
	"goognion/src"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type TorCrawlerRepository struct {
	client *http.Client
	tor    *tor.Tor
}

func NewTorCrawlerRepository(client *http.Client, tor *tor.Tor) *TorCrawlerRepository {
	return &TorCrawlerRepository{client: client, tor: tor}
}

func (r *TorCrawlerRepository) DoIndexing(input <-chan src.Text) (map[string]int, error) {
	out := make(chan string)

	go func() {
		wg := &sync.WaitGroup{}

		for text := range input {
			wg.Add(1)
			go indexText(text, wg, out)
		}

		wg.Wait()
		close(out)
	}()

	m := map[string]int{}

	for word := range out {
		m[word] += 1
	}

	return m, nil
}

func (r *TorCrawlerRepository) Load(url string) (<-chan src.Text, <-chan string, error) {
	outT := make(chan src.Text)
	outU := make(chan string)

	resp, err := r.client.Get(url)
	if err != nil {
		close(outT)
		close(outU)
		return outT, outU, err
	}
	defer resp.Body.Close()

	parsed, err := html.Parse(resp.Body)
	if err != nil {
		close(outT)
		close(outU)
		return outT, outU, err
	}

	go func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go getTextsAndUrls(parsed, outT, outU, wg)

		wg.Wait()
		close(outT)
		close(outU)
	}()

	return outT, outU, nil
}

func getTextsAndUrls(n *html.Node, text chan<- src.Text, url chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wg.Add(1)
		go getTextsAndUrls(c, text, url, wg)
	}

	switch n.Type {
	case html.ElementNode:
		if n.Data != "a" {
			break
		}
		for _, a := range n.Attr {
			if a.Key == "href" {
				wg.Add(1)
				go func() {
					url <- a.Val
					wg.Done()
				}()
				break
			}
		}
		fallthrough
	case html.TextNode:
		var t bytes.Buffer
		if n == nil || n.FirstChild == nil {
			break
		}
		if err := html.Render(&t, n.FirstChild); err != nil {
			break
		}
		wg.Add(1)
		go func() {
			text <- (src.Text)(t.String())
			wg.Done()
		}()
	}
}

func indexText(t src.Text, wg *sync.WaitGroup, out chan<- string) {
	defer wg.Done()

	ws := strings.Split(string(t), " ")

	for _, w := range ws {
		wg.Add(1)
		go indexWord(w, wg, out)
	}
}

func indexWord(w string, wg *sync.WaitGroup, out chan<- string) {
	defer wg.Done()

	w = strings.TrimFunc(w, func(r rune) bool {
		return unicode.IsPunct(r) ||
			unicode.IsSpace(r) ||
			unicode.IsSymbol(r)
	})
	w = strings.ToLower(w)
	match, err := regexp.MatchString(`^[a-z0-9]+$`, w)
	if !match || err != nil {
		return
	}

	out <- w
}

func NewTorClient(c TorConfig) (*http.Client, *tor.Tor, error) {
	fmt.Println("Запуск Tor...")

	t, err := tor.Start(nil, &tor.StartConf{ExePath: c.ExePath, DataDir: c.DataDir})
	if err != nil {
		return nil, nil, err
	}

	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Tor успешно запущен!")

	return &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}, t, nil
}
