package spiders

import (
	"beewolf/ship"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

type PandaReturn struct {
	Data struct {
		Total int         `json:"total"`
		Items []PandaRoom `json:"items"`
	} `json:"data"`
}

type PandaRoom struct {
	Title     string `json:"name"`
	PersonNum int    `json:"person_num,string"`
}

type PandaSpider struct {
	*ship.Spider
	Index       string
	TotalPerson int64
}

func (tv *PandaSpider) DoBefore() error {
	res, err1 := http.Get(tv.StartUrl)
	if err1 != nil {
		return err1
	}
	content, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return err2
	}
	defer res.Body.Close()

	var rtn PandaReturn
	err3 := json.Unmarshal(content, &rtn)
	if err3 != nil {
		return err3
	}

	pageCount := rtn.Data.Total / 120
	if rtn.Data.Total%120 > 0 {
		pageCount++
	}
	fmt.Printf("PageCount: %d\n", pageCount)
	for i := 1; i < pageCount+1; i++ {
		url := fmt.Sprintf("http://www.panda.tv/live_lists?status=2&order=person_num&pageno=%d&pagenum=120", i)
		tv.Urls <- url
	}
	return nil
}

func (tv *PandaSpider) ParseItem(content []byte, items chan interface{}) error {
	var payload PandaReturn
	err := json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	for _, room := range payload.Data.Items {
		items <- room
	}
	return nil
}

func (tv *PandaSpider) Pipeline(item interface{}) interface{} {
	room := item.(PandaRoom)
	//fmt.Println(room)
	atomic.AddInt64(&tv.TotalPerson, int64(room.PersonNum))
	return room
}

func (tv *PandaSpider) DoAfter() error {
	fmt.Printf("[%s]: 当前在线%d人\n", tv.Name, tv.TotalPerson)
	close(tv.Urls)
	close(tv.Items)
	return nil
}

var PandaTV ship.ISpider = &PandaSpider{
	ship.NewSpider(
		"熊猫TV",
		"http://www.panda.tv/live_lists?status=2&order=person_num&pageno=1&pagenum=120",
		10000,
		10000),
	"http://www.panda.tv",
	0,
}
