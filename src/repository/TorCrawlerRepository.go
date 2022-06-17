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
)

type TorCrawlerRepository struct {
	client *http.Client
}

func NewTorCrawlerRepository() *TorCrawlerRepository {
	//todo path to main, mb like db
	c, err := newTorClient("D:/Tor Browser/Browser/TorBrowser/Tor/tor.exe")
	if err != nil {
		panic(err)
	}

	return &TorCrawlerRepository{client: c}
}

func newTorClient(torPath string) (*http.Client, error) {
	t, err := tor.Start(nil, &tor.StartConf{ExePath: torPath})
	if err != nil {
		return nil, err
	}
	//todo fmt
	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}, nil
}

func (r *TorCrawlerRepository) DoIndexing(input <-chan src.Text) (map[string]int, error) {
	out := make(chan string)

	go func() {
		wg := sync.WaitGroup{}

		for text := range input {
			wg.Add(1)

			go func(t src.Text, wg *sync.WaitGroup) {
				ws := strings.Split(string(t), " ")

				for _, w := range ws {
					wg.Add(1)

					go func(w string, wg *sync.WaitGroup) {
						w = strings.Trim(w, " \n\t\r") //cutset .?&^%#$;: not cutset \t
						match, err := regexp.MatchString(`.+`, w)
						if !match || err != nil {
							return
						}
						out <- w
						wg.Done()
					}(w, wg)
				}

				wg.Done()
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
	//defer resp.Body.Close() todo

	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	outT := make(chan src.Text)
	outU := make(chan string)

	go func() {
		wg := sync.WaitGroup{}

		wg.Add(1)
		getTextsAndUrls(parsed, outT, outU, &wg)

		wg.Wait()
		close(outT)
		close(outU)
	}()

	return outT, outU, nil
}

func getTextsAndUrls(n *html.Node, text chan<- src.Text, url chan<- string, wg *sync.WaitGroup) {
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

	wg.Done()
}
