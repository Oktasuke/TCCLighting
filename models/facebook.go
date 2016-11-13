package models

import "fmt"

func NewfacebookReq() facebookReq {
	return facebookReq{}
}

type facebookReq struct {
	Entry  []Entry `json:"entry"`
	Object string  `json:"object"`
}

type Entry struct {
	Changes []Change `json:"changes"`
	Id      string   `json:"id"`
	Time    int      `json:"time"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	Item    string `json:"item"`
	Verb    string `json:"verb"`
	User_id int    `json:"user_id"`
}

func Assert(data interface{}) {
	switch data.(type) {
	case string:
		fmt.Print(data.(string))
	case float64:
		fmt.Print(data.(float64))
	case bool:
		fmt.Print(data.(bool))
	case nil:
		fmt.Print("null")
	case []interface{}:
		fmt.Print("[")
		for _, v := range data.([]interface{}) {
			Assert(v)
			fmt.Print(" ")
		}
		fmt.Print("]")
	case map[string]interface{}:
		fmt.Print("{")
		for k, v := range data.(map[string]interface{}) {
			fmt.Print(k + ":")
			Assert(v)
			fmt.Print(" ")
		}
		fmt.Print("}")
	default:
	}
}
