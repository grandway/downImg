package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"imgDown/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Content struct {
	img   []string
	href  []string
	title []string
}

func main() {

	resp, err := http.PostForm(config.BASE_URL, url.Values{"type": {"catL1182"}, "paged": {"1"}, "action": {"zrz_load_more_posts"}})

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

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(data.Msg))

	if err != nil {
		panic(err)
	}

	reg := regexp.MustCompile("http.*?jpg")
	var Content Content

	dom.Find(".post-list").Each(func(i int, selection *goquery.Selection) {
		style, _ := selection.Find(".preview").Attr("style")
		href, _ := selection.Find(".link-block").Attr("href")
		title := selection.Find(".entry-title a").Text()
		img := reg.FindString(style)

		Content.title = append(Content.title, title)
		Content.img = append(Content.img, img)
		Content.href = append(Content.href, href)

	})
	downImg(Content)
}

func downImg(content Content) {

	exist, _ := PathExists(config.BASE_DOWN_PATH)
	if !exist {
		os.Mkdir(config.BASE_DOWN_PATH, os.ModePerm)
	}
	for _, img := range content.img {
		go getImg(img)
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
func getImg(url string) (n int64, err error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	downPath := config.BASE_DOWN_PATH + "/" + name
	fmt.Println(downPath)
	out, err := os.Create(downPath)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	return

}
