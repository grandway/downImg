package request

import (
	"net/url"
	"strconv"
)

type PostParams struct {
	Type   []string //分类
	Action []string //方法
	Paged  []string //页码
}

func PostForms(params PostParams) url.Values {
	values := url.Values{}

	values["type"] = params.Type
	values["action"] = params.Action
	values["paged"] = params.Paged
	return values
}

func DefaultPostForms() url.Values {
	return url.Values{"type": {"catL1182"}, "paged": {"1"}, "action": {"zrz_load_more_posts"}}
}

func NewPostForms(kind, action string, paged int) url.Values {
	Params := PostParams{}
	Params.Type = append(Params.Type, kind)
	Params.Action = append(Params.Action, action)
	Params.Paged = append(Params.Paged, strconv.Itoa(paged))
	return PostForms(Params)
}
