package ship

import (
	"time"
	"log"
)

type Spider struct {
	Name     string
	StartUrl string
	Urls     chan string
	Items    chan interface{}

	errs chan error
}
/*
包含了一个爬虫的基本组成
StartUrl为起始url，Urls用于接收url爬取器的输出，Items用于接收解析器的输出结果
*/

func NewSpider(name string, startUrl string, urlMax int, itemMax int) *Spider {
	urls := make(chan string, urlMax)
	items := make(chan interface{}, itemMax)
	errs := make(chan error, 10)
	return &Spider{
		Name:     name,
		StartUrl: startUrl,
		Urls:     urls,
		Items:    items,
		errs: errs,
	}
}

func (s *Spider) GetName() string {
	return s.Name
}

func (s *Spider) UrlGenerator() chan string {
	return s.Urls
}

func (s *Spider) ItemReceiver() chan interface{} {
	return s.Items
}

func (s *Spider) SetTimeout(d time.Duration) chan bool {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(d)
		timeout <- true
	}()
	return timeout
}

func (s *Spider) SetError(err error) {
	s.errs <- err
	log.Printf("[%s]: %s", s.GetName(), err)
}

func (s *Spider) GetError() chan error {
	return s.errs
}
