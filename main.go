package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"imgDown/config"
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

//var post_data map[int] string
var wg sync.WaitGroup

func main() {

	//请求地址
	//postForms := request.NewPostForms("catL1182", "zrz_load_more_posts", 1)
	resp, err := http.PostForm(config.BaseURL, request.DefaultPostForms())

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

//todo 获取列表页数据
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

	exist, _ := PathExists(config.BaseDownPath)
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

//文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//下载图片
func saveImg(url, dir, name string) (n int64, err error) {

	path := config.BaseDownPath + "/" + dir

	exist, _ := PathExists(path)
	if !exist {
		os.Mkdir(path, os.ModePerm)
	}

	downPath := path + "/" + name + ".jpg"
	fmt.Println(downPath)
	out, err := os.Create(downPath)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	wg.Done()
	return

}
