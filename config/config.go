package config

const BaseDownPath = "./images"

const BaseURL string = "https://www.jdlingyu.mobi/wp-admin/admin-ajax.php"

type RespBody struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
