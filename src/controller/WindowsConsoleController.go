package controller

import (
	"fmt"
	"goognion/src/service"
	"sync"
)

var await = sync.WaitGroup{}

type WindowsConsoleController struct {
}

func (c *WindowsConsoleController) Run() {
	src, depth := "", 0
	fmt.Println("Ввод:\ndb [глубина индексирования]\nлибо\n[ссылка на страницу] [глубина индексирования]\n-----")
	for {
		await.Wait()
		_, err := fmt.Scanln(&src, &depth)
		if err != nil {
			fmt.Println("Проверьте правильность ввода")
			continue
		}

		switch src {
		case "db":
			urls := service.Crawler.GetUncrawledUrls(depth)

			await.Add(len(urls))
			for _, url := range urls {
				go c.crawl(url, depth, &await)
			}
		default:
			await.Add(1)
			go c.crawl(src, depth, &await)
		}
	}
}

func (c *WindowsConsoleController) crawl(url string, depth int, await *sync.WaitGroup) {
	service.Crawler.Crawl(url, depth)
	await.Done()
}
