package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"imgDown/config"
	"imgDown/helper"
	"imgDown/request"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Content struct {
	title string
	img   []string
}

var wg sync.WaitGroup

func main() {
	//获取分类
	getCat()
	//请求地址
	//postForms := request.NewPostForms("catL1182", "zrz_load_more_posts", 1)
	resp, err := http.PostForm(config.ListURL, request.DefaultPostForms())

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var data config.RespBody

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	if data.Status != 200 {
		panic(data.Msg)
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(data.Msg))

	if err != nil {
		panic(err)
	}
	var list []string

	dom.Find(".post-list").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Find(".link-block").Attr("href")
		list = append(list, href)
	})

	downImg(list)
	wg.Wait()

	fmt.Print("job success")
}

//TODO 获取列表页数据
//func getList()  {
//
//}

//获取详情
func getContent(url string) (Content, error) {

	resp, err := http.Get(url)

	if err != nil {
		return Content{}, err
	}
	defer resp.Body.Close()

	var content Content

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Content{}, err
	}
	content.title = dom.Find(".entry-title").Text()

	dom.Find("#entry-content img").Each(func(i int, selection *goquery.Selection) {
		img, exist := selection.Attr("src")
		if !exist {
			return
		}
		content.img = append(content.img, img)
	})
	return content, nil
}

//下载图片
func downImg(list []string) {

	exist, _ := helper.PathExists(config.BaseDownPath)
	if !exist {
		os.Mkdir(config.BaseDownPath, os.ModePerm)
	}
	for _, url := range list {
		content, err := getContent(url)

		if err != nil {
			continue
		}
		for key, img := range content.img {
			wg.Add(1)
			go saveImg(img, content.title, strconv.Itoa(key))
		}
	}

}

//下载图片
func saveImg(url, dir, name string) (n int64, err error) {

	path := config.BaseDownPath + "/" + dir

	exist, _ := helper.PathExists(path)
	if !exist {
		os.Mkdir(path, os.ModePerm)
	}

	downPath := path + "/" + name + ".jpg"
	fmt.Println(downPath)
	out, err := os.Create(downPath)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		wg.Done()
		return
	}
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	wg.Done()
	return
}

func getCat() {
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	var category = config.GetCategory()
	for _, c := range category {
		fmt.Println("1.", c["title"])
	}
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	var c int
	fmt.Print("请选择下载类型:")

	fmt.Scanln(&c)
	cat, ok := category[c]

	if !ok {
		fmt.Println("")
		getCat()
	}

	fmt.Println(cat)
	panic(c)
}
