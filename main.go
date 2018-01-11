package main

import (
    "crypto/tls"
    "encoding/json"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"flag"
    "net/http"
    "net/http/cookiejar"
)

type goods struct {
	RealTitle string
	Price     float64
	GoodsDetailUrl  string
}
type resData struct {
    StateCode int
    Data  struct{
		GoodsList []goods
	}                                           
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}
func main() {
	keywords := flag.String("k", "光方", "查询关键字")
	flag.Parse()
    i:=1
    for{
		f:=getPage(i,*keywords)
		if f{
			i++
		}else{
			return
		}
	}

}

func getPage(page int,keywords string) bool{
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
        if err:=recover();err!=nil{
            fmt.Println(err) 
        }
    }()
	//目标URL，即需要访问的URL
    targetUrl := "https://m.jiaoyimao.com/fe/ajax/goods/?gameId=498&page="+strconv.Itoa(page)
  
    var resp *http.Response
    var err error
    var data []byte
    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    //启用cookie
    client.Jar, _ = cookiejar.New(nil)
    resp, err = client.Get(targetUrl)
    check(err)

    if data, err = ioutil.ReadAll(resp.Body); err != nil {
        check(err)
    }
    var message resData
    err = json.Unmarshal(data, &message)
	check(err)
	if len(message.Data.GoodsList) <1 {
		fmt.Println("========查询完成========")
		return false
	}
	for _,g := range message.Data.GoodsList {
        if strings.Contains(g.RealTitle, keywords){
			fmt.Println(g.RealTitle+"价格："+strconv.FormatFloat(g.Price, 'f', 6, 64)+"   =   "+g.GoodsDetailUrl)
		}
	}
	return true
}
