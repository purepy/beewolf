package main

import (
	"beewolf/ship"
	"beewolf/spiders"
)

func main() {
	crawler := ship.NewCrawler()
	crawler.AddSpider(spiders.PandaTV)
	crawler.AddSpider(spiders.QuanminTV)
	crawler.AddSpider(spiders.ZhanqiTV)
	crawler.AddSpider(spiders.HuyaTV)
	crawler.Run()
}
