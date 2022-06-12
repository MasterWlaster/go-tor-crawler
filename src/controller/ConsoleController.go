package controller

import (
	"fmt"
	"goognion/src/service"
	"sync"
)

type ConsoleController struct {
	await   *sync.WaitGroup
	service *service.Service
}

func NewConsoleController(service *service.Service) *ConsoleController {
	return &ConsoleController{
		await:   &sync.WaitGroup{},
		service: service}
}

func (c *ConsoleController) Run() {
	src, depth := "", 0
	fmt.Println("\nВвод:\ndb [глубина индексирования]\nлибо\n[ссылка на страницу] [глубина индексирования]\n-----")
	for {
		c.await.Wait()
		_, err := fmt.Scanln(&src, &depth)
		if err != nil {
			fmt.Println("Проверьте правильность ввода!")
			continue
		}

		switch src {
		case "db":
			urls := c.getNotCrawledUrls(depth)
			if urls == nil {
				fmt.Println("В БД не удалось обнаружить непроиндексированных страниц")
			}

			c.await.Add(len(urls))
			for _, url := range urls {
				go c.crawl(url, depth, c.await)
			}
		default:
			c.await.Add(1)
			go c.crawl(src, depth, c.await)
		}
	}
}

func (c *ConsoleController) crawl(url string, depth int, await *sync.WaitGroup) {
	c.service.Crawler.Crawl(url, depth)
	await.Done()
}

func (c *ConsoleController) getNotCrawledUrls(maxDepth int) []string {
	//todo
	return nil
}
