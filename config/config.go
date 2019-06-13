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
		1: {"title": "测试1", "url": "https://www.jdlingyu.mobi/tuji/hentai/costt", "cat": "catL1182"},
	}
	return cat
}
