package config

const BaseDownPath string = "./images"

const BaseURL string = "https://www.jdlingyu.mobi"
const ListURL string = "https://www.jdlingyu.mobi/wp-admin/admin-ajax.php"

type RespBody struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func GetCategory() map[int]map[string]string {
	cat := map[int]map[string]string{
		1: {"title": "cos套图", "url": "https://www.jdlingyu.mobi/tuji/hentai/costt", "cat": "catL1182"},
		2: {"title": "国产套图", "url": "https://www.jdlingyu.mobi/tuji/hentai/gctt", "cat": "catL1183"},
		3: {"title": "日本写真", "url": "https://www.jdlingyu.mobi/tuji/hentai/rbxz", "cat": "catL1184"},
	}
	return cat
}
