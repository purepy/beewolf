package ship

import (
	"time"
)

type ISpider interface {
	GetName() string
	UrlGenerator() chan string
	DoBefore() error
	ItemReceiver() chan interface{}
	ParseItem(content []byte, items chan interface{}) error
	Pipeline(item interface{}) interface{}
	DoAfter() error
	SetTimeout(d time.Duration) chan bool
	SetError(err error)
	GetError() chan error
}

/*
ISpiders的接口方法

GetName() string 需要返回Spider的Name
UrlGenerator() chan string 需要返回一个url存储器
ItemReceiver() chan interface{} 需要返回一个item存储器
ParseItem(content []byte, items chan interface{}) error 需要实现对于content的解析，并将解析结果放入信道items中
Pipeline(item interface{}) interface{} 对传入的单个item进行自定义的数据操作，完成之后需要返回一个新item
SetTimeout(d time.Duration) chan bool 设置并返回超时对象，用于chan的相关超时操作
DoBefore() error 运行自定义的初始化工作
DoAfter() error 运行自定义的清洁工作
*/
