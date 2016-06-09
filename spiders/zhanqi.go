package spiders

import (
	"beewolf/ship"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"io/ioutil"
	"net/http"
)

type ZhanqiReturn struct {
	Data struct {
		Total int          `json:"cnt"`
		Rooms []ZhanqiRoom `json:"rooms"`
	}
}

type ZhanqiRoom struct {
	Title        string `json:"title"`
	Owner        string `json:"nickname"`
	CategoryName string `json:"gameName"`
	PersonNum    int    `json:"online,string"`
}

type ZhanqiSpider struct {
	*ship.Spider
	Index string
	TotalPerson int64
}

func (tv *ZhanqiSpider) DoBefore() error {
	res, err1 := http.Get(tv.StartUrl)
	if err1 != nil {
		return err1
	}
	content, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return err2
	}
	defer res.Body.Close()

	var rtn ZhanqiReturn
	err3 := json.Unmarshal(content, &rtn)
	if err3 != nil {
		return err3
	}
	pageCount := rtn.Data.Total / 30
	if rtn.Data.Total%30 > 0 {
		pageCount++
	}
	for i := 1; i < pageCount+1; i++ {
		url := fmt.Sprintf("http://www.zhanqi.tv/api/static/live.hots/30-%d.json", i)
		tv.Urls <- url
	}
	return nil
}

func (tv *ZhanqiSpider) ParseItem(content []byte, items chan interface{}) error {
	var payload ZhanqiReturn
	err := json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}
	for _, room := range payload.Data.Rooms {
		items <- room
	}
	return nil
}

func (tv *ZhanqiSpider) Pipeline(item interface{}) interface{} {
	room := item.(ZhanqiRoom)
	//fmt.Println(room)
	atomic.AddInt64(&tv.TotalPerson, int64(room.PersonNum))
	return room
}

func (tv *ZhanqiSpider) DoAfter() error {
	fmt.Printf("[%s]: 当前在线%d人\n", tv.Name, tv.TotalPerson)
	close(tv.Urls)
	close(tv.Items)
	return nil
}

var ZhanqiTV ship.ISpider = &ZhanqiSpider{
	ship.NewSpider(
		"战旗TV",
		"http://www.zhanqi.tv/api/static/live.hots/30-1.json",
		10000,
		10000),
	"http://www.zhanqi.tv",
	0,
}
