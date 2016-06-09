package ship

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"sync"
)

type Crawler struct {
	spiders       []ISpider
	SpiderTimeout time.Duration
	wg sync.WaitGroup
}

func NewCrawler() *Crawler {
	var spiderTimeout = 50 * time.Second
	return &Crawler{
		SpiderTimeout: spiderTimeout,
	}
}

func (c *Crawler) AddSpider(spider ISpider) {
	c.spiders = append(c.spiders, spider)
}

func (c *Crawler) Run() {
	for _, spider := range c.spiders {
		c.wg.Add(1)
		go c.RunSpider(spider)
	}
	c.wg.Wait()
}

func (c *Crawler) Request(method string, url string) (content []byte, err error) {
	log.Printf("HTTP %s %s\n", method, url)
	res, err1 := http.Get(url)
	if err1 != nil {
		return content, err1
	}
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return content, err
	}
	defer res.Body.Close()

	return content, nil

}

func (c *Crawler) RunSpider(spider ISpider) {
	defer c.wg.Done()
	
	log.Println("==================================")
	log.Printf("Work on spider <%s>", spider.GetName())
	log.Println("==================================")

	urls := spider.UrlGenerator()
	items := spider.ItemReceiver()
	errs := spider.GetError()
	timeout := spider.SetTimeout(c.SpiderTimeout)

	go c.BeforeHook(spider)
	//go c.ItemPipeline(spider)

	var method = "GET"
	for {
		select {
		case url := <- urls:
			content, _ := c.Request(method, url)
			spider.ParseItem(content, items)
		case <- errs:
			c.ItemPipeline(spider)
			c.AfterHook(spider)			
			return
		case <-timeout:
			log.Printf("[%s]: Pool of urls is empty...", spider.GetName())
			c.ItemPipeline(spider)
			c.AfterHook(spider)
			return
		}
	}
}


func (c *Crawler) BeforeHook(spider ISpider) {
	err := spider.DoBefore()
	log.Println("Before Hook")
	if err != nil {
		spider.SetError(err)
	}
	log.Printf("[%s]: Finish before work...", spider.GetName())
}

func (c *Crawler) AfterHook(spider ISpider) {
	err := spider.DoAfter()
	if err != nil {
		spider.SetError(err)
	}
	log.Printf("[%s]: Finish after work...", spider.GetName())
}

func (c *Crawler) ItemPipeline(spider ISpider) {
	items := spider.ItemReceiver()
	timeout := spider.SetTimeout(0 * time.Second)

	for {
		select {
		case item := <- items:
			go spider.Pipeline(item)
		case <- timeout:
			return
		}
	}
}
