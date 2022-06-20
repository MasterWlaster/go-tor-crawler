package controller

import (
	"fmt"
	"goognion/src/service"
	"sync"
)

type ConsoleController struct {
	wg      *sync.WaitGroup
	service *service.Service
}

func NewConsoleController(service *service.Service) *ConsoleController {
	return &ConsoleController{
		wg:      &sync.WaitGroup{},
		service: service}
}

func (c *ConsoleController) Run() {
	src, depth, lim := "", 0, 0
	fmt.Println("\nИспользование:\ndb [глубина индексирования] [лимит страниц, взятых из бд] - начало работы с непроиндексированными страницами\nлибо\n[ссылка на страницу в виде http://..... ] [глубина индексирования]")
	for {
		c.wg.Wait()
		fmt.Println("\nВвод:")

		//_, err := fmt.Scanln(&src, &depth, &lim)
		//if err != nil {
		//	fmt.Println("Проверьте правильность ввода!")
		//	continue
		//}

		_, err := fmt.Scan(&src, &depth)
		if err != nil {
			fmt.Println("Проверьте правильность ввода!")
			continue
		}

		switch src {
		case "db":
			_, err = fmt.Scan(&lim)
			if err != nil {
				fmt.Println("Проверьте правильность ввода!")
				continue
			}

			urls := c.getNotCrawledUrls(lim)
			if urls == nil {
				fmt.Println("В БД не удалось обнаружить непроиндексированных страниц")
				continue
			}

			c.wg.Add(len(urls))
			for _, url := range urls {
				go c.crawl(url, depth)
			}
		default:
			c.wg.Add(1)
			go c.crawl(src, depth)
		}
		fmt.Println("Работаю...")
	}
}

func (c *ConsoleController) crawl(url string, depth int) {
	errs := c.service.Crawler.Crawl(url, depth)
	c.service.Logger.Log(errs)
	c.wg.Done()
}

func (c *ConsoleController) getNotCrawledUrls(limit int) []string {
	var urls []string
	var err error

	if limit <= 0 {
		urls, err = c.service.Crawler.GetNotCrawledUrls()
	} else {
		urls, err = c.service.Crawler.GetNotCrawledUrlsWithLimit(limit)
	}

	c.service.Logger.LogOnce(err)

	return urls
}
