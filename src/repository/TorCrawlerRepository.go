package repository

import (
	"bytes"
	"context"
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

func NewTorClient(torPath string, torDataDir string) (*http.Client, *tor.Tor, error) {
	t, err := tor.Start(nil, &tor.StartConf{ExePath: torPath})
	if err != nil {
		return nil, nil, err
	}
	//todo fmt
	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}, t, nil
}

func (r *TorCrawlerRepository) DoIndexing(input <-chan src.Text) (map[string]int, error) {
	out := make(chan string)

	go func() {
		wg := sync.WaitGroup{}

		for text := range input {
			wg.Add(1)

			go func(t src.Text, wg *sync.WaitGroup) {
				defer wg.Done()

				ws := strings.Split(string(t), " ")

				for _, w := range ws {
					wg.Add(1)

					go func(w string, wg *sync.WaitGroup) {
						defer wg.Done()

						w = strings.TrimFunc(w, func(r rune) bool {
							return unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsSymbol(r)
						})
						w = strings.ToLower(w)
						match, err := regexp.MatchString(`.+`, w)
						if !match || err != nil {
							return
						}

						out <- w
					}(w, wg)
				}
			}(text, &wg)
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
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	outT := make(chan src.Text)
	outU := make(chan string)

	go func() {
		wg := sync.WaitGroup{}

		wg.Add(1)
		go getTextsAndUrls(parsed, outT, outU, &wg)

		wg.Wait()
		//get(parsed, outT, outU)
		close(outT)
		close(outU)
	}()

	return outT, outU, nil
}

func getTextsAndUrls(n *html.Node, text chan<- src.Text, url chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wg.Add(1)
		go func(n *html.Node, wg *sync.WaitGroup) {
			getTextsAndUrls(n, text, url, wg)
		}(c, wg)
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					go func() {
						url <- a.Val
						wg.Done()
					}()
					break
				}
			}
			wg.Add(1)
			//fallthrough
		case "div":
			fallthrough
		case "label":
			fallthrough
		case "button":
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
			var t bytes.Buffer
			if n == nil || n.FirstChild == nil {
				break
			}
			if err := html.Render(&t, n.FirstChild); err != nil {
				break
			}
			go func() {
				text <- (src.Text)(t.String())
				wg.Done()
			}()
			wg.Add(1)
		}
	}
}

func get(n *html.Node, text chan<- src.Text, url chan<- string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		get(c, text, url)
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					url <- a.Val
					break
				}
			}
			fallthrough
		case "div":
			fallthrough
		case "label":
			fallthrough
		case "button":
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
			var t bytes.Buffer
			if n == nil || n.FirstChild == nil {
				break
			}
			if err := html.Render(&t, n.FirstChild); err != nil {
				break
			}
			text <- (src.Text)(t.String())
		}
	}
}
