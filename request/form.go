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

func NewPostForms(postType, action string, page int) url.Values {
	Params := PostParams{}
	Params.Type = append(Params.Type, postType)
	Params.Action = append(Params.Action, action)
	Params.Paged = append(Params.Paged, strconv.Itoa(page))
	return PostForms(Params)
}
