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
	await.Wait()

	src, depth := "", 0
	for {
		_, err := fmt.Scanln(&src, &depth)
		if err != nil {
			fmt.Println("Проверьте правильность ввода")
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
