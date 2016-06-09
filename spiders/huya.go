package spiders

import (
	"beewolf/ship"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"io/ioutil"
	"net/http"
)

type HuyaReturn struct {
	Data struct {
		List  []HuyaRoom `json:"list"`
		Total int        `json:"total,string"`
	} `json:"data"`
}

type HuyaRoom struct {
	Title        string `json:"roomName"`
	Owner        string `json:"nick"`
	CategoryName string `json:"gameFullName"`
	PersonNum    int    `json:"totalCount,string"`
}

type HuyaSpider struct {
	*ship.Spider
	Index string
	TotalPerson int64
}

func (tv *HuyaSpider) DoBefore() error {
	res, err1 := http.Get(tv.StartUrl)
	if err1 != nil {
		return err1
	}
	content, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return err2
	}
	defer res.Body.Close()

	var rtn HuyaReturn
	err3 := json.Unmarshal(content, &rtn)
	if err3 != nil {
		return err3
	}

	pageCount := rtn.Data.Total / 20
	fmt.Println(rtn.Data.Total)
	fmt.Println(pageCount)
	if rtn.Data.Total%20 > 0 {
		pageCount++
	}

	var i int
	for i = 1; i < pageCount+1; i++ {
		url := fmt.Sprintf("http://www.huya.com/index.php?m=Live&do=ajaxAllLiveByPage&page=%d&pageNum=1", i)
		tv.Urls <- url
	}
	fmt.Printf("共%d页，%d个链接\n", pageCount, i)
	return nil
}

func (tv *HuyaSpider) ParseItem(content []byte, items chan interface{}) error {
	var payload HuyaReturn
	err := json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	for _, room := range payload.Data.List {
		items <- room
	}
	return nil
}

func (tv *HuyaSpider) Pipeline(item interface{}) interface{} {
	room := item.(HuyaRoom)
	//fmt.Println(room)
	atomic.AddInt64(&tv.TotalPerson, int64(room.PersonNum))
	return room
}

func (tv *HuyaSpider) DoAfter() error {
	fmt.Printf("[%s]: 当前在线%d人\n",tv.Name, tv.TotalPerson)
	close(tv.Urls)
	close(tv.Items)
	return nil
}

var HuyaTV ship.ISpider = &HuyaSpider{
	ship.NewSpider(
		"虎牙TV",
		"http://www.huya.com/index.php?m=Live&do=ajaxAllLiveByPage&page=1&pageNum=1",
		10000,
		10000),
	"http://www.huya.com",
	0,
}
