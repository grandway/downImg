package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/go.uuid"
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
type Category struct {
	title string
	url   string
	cat   string
}

var wg sync.WaitGroup

var category Category

func main() {
	var Category Category
	var flag = false
	for !flag {
		cate, f := getCat(Category)
		if f {
			flag = true
		}
		category = cate
	}

	totalPage := getTotalPage(category)

	fmt.Println(category.title, "开始下载...")

	for i := totalPage; i > 0; i-- {
		getList(category.cat, i)
	}
	wg.Wait()

	fmt.Println("job success")
}

//获取列表
func getList(category string, page int) {
	//请求地址
	postForms := request.NewPostForms(category, "zrz_load_more_posts", page)

	resp, err := http.PostForm(config.ListURL, postForms)

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
	if list[0] == "" {
		return
	}
	downImg(list)
}

//获取总页数
func getTotalPage(cate Category) (totalPage int) {
	resp, err := http.Get(cate.url)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)

	}
	page, exists := dom.Find("page-nav").Attr(":pages")
	if !exists {
		panic(err)
	}
	totalPage, _ = strconv.Atoi(page)
	fmt.Println("总页数为：", page)
	return totalPage
}

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
		if !exist || (len(img) < len(config.BaseURL)) {
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

		wg.Add(1)
		go func(url string) {
			content, err := getContent(url)
			if err == nil {
				path := config.BaseDownPath + "/" + content.title
				exist, _ := helper.PathExists(path)
				if !exist {
					os.Mkdir(path, os.ModePerm)
				}
				for _, img := range content.img {
					wg.Add(1)
					go saveImg(img, path, uuid.NewV4().String())
				}
			}
			wg.Done()
		}(url)

	}

}

//下载图片
func saveImg(url, dir, name string) (n int64, err error) {

	downPath := dir + "/" + name + ".jpg"
	fmt.Println(downPath)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		out, _ := os.Create(downPath)
		defer out.Close()
		pix, _ := ioutil.ReadAll(resp.Body)
		n, err = io.Copy(out, bytes.NewReader(pix))
	}
	wg.Done()
	return
}

func getCat(cate Category) (Category, bool) {
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	var category = config.GetCategory()
	for i, c := range category {
		fmt.Println(i, ".", c["title"])
	}
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	var c int
	fmt.Print("请选择下载类型:")

	fmt.Scanln(&c)

	cat, ok := category[c]

	if !ok {
		return cate, false
	}

	cate.title = cat["title"]
	cate.url = cat["url"]
	cate.cat = cat["cat"]

	return cate, true
}
