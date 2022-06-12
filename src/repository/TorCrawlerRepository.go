package repository

import (
	"bytes"
	"context"
	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
	"goognion/src"
	"net/http"
	"strings"
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

	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, nil, err
	}
	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

	// Get /
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Grab the <title>
	_, err = html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, err
}

func (r *TorCrawlerRepository) DoIndexing(src []src.Text) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func getTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		var title bytes.Buffer
		if err := html.Render(&title, n.FirstChild); err != nil {
			panic(err)
		}
		return strings.TrimSpace(title.String())
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := getTitle(c); title != "" {
			return title
		}
	}
	return ""
}
