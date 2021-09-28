package leetcode

import "encoding/json"

type Data struct {
	Data   json.RawMessage `json:"data"`
	Errors []Errors        `json:"errors"`
}

type AllNum struct {
	AllNum int `json:"allNum"`
}

type DataFollow struct {
	Followers         AllNum `json:"followers"`
	FollowingEntities AllNum `json:"followingEntities"`
}
