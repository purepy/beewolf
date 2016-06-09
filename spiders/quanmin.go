package spiders

import (
	"beewolf/ship"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"io/ioutil"
	"net/http"
)

type QuanminReturn struct {
	Total     int           `json:"total"`
	PageCount int           `json:"pageCount"`
	PageNo    int           `json:"page"`
	Size      int           `json:"size"`
	Data      []QuanminRoom `json:"data"`
}

type QuanminRoom struct {
	Title        string `json:"title"`
	Owner        string `json:"nick"`
	CategoryName string `json:"category_name"`
	PersonNum    int    `json:"view,string"`
}

type QuanminSpider struct {
	*ship.Spider
	Index string
	TotalPerson int64
}

/*
实现一个简单的蜘蛛需要实现ISpider的如下4个方法:
DoBefore() error
ParseItem(content []byte, items interface{}) error
Pipeline(item interface{}) interface{}
DoAfter() error
*/

func (tv *QuanminSpider) DoBefore() error {
	res, err := http.Get(tv.StartUrl)
	if err != nil {
		return err
	}

	content, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return err2
	}
	defer res.Body.Close()

	var rtn QuanminReturn
	err3 := json.Unmarshal(content, &rtn)
	if err3 != nil {
		return err3
	}

	fmt.Printf("共有%d页\n", rtn.PageCount)
	tv.Urls <- tv.StartUrl
	for i := 1; i < rtn.PageCount+1; i++ {
		url := fmt.Sprintf("http://www.quanmin.tv/json/play/list_%d.json", i)
		tv.Urls <- url
	}

	return nil
}

func (tv *QuanminSpider) ParseItem(content []byte, items chan interface{}) error {
	var payload QuanminReturn
	err := json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}
	for _, room := range payload.Data {
		items <- room
	}

	return nil
}

func (tv *QuanminSpider) Pipeline(item interface{}) interface{} {
	room := item.(QuanminRoom)
	atomic.AddInt64(&tv.TotalPerson, int64(room.PersonNum))
	//fmt.Println(room)
	return room
}

func (tv *QuanminSpider) DoAfter() error {
	fmt.Printf("[%s]: 当前在线%d人\n", tv.Name, tv.TotalPerson)
	close(tv.Urls)
	close(tv.Items)
	return nil
}

var QuanminTV ship.ISpider = &QuanminSpider{
	ship.NewSpider(
		"全民TV",
		"http://www.quanmin.tv/json/play/list.json",
		1000,
		1000),
	"http://www.quanmin.tv",
	0,
}
