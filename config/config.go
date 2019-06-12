package config

const BASE_DOWN_PATH = "./images"

const BASE_URL string = "https://www.jdlingyu.mobi/wp-admin/admin-ajax.php"

type RespBody struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
