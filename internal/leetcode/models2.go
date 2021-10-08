package leetcode

import (
	"encoding/json"

	"github.com/pkg/errors"
)

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

var (
	ErrUserNotExist = errors.New("That user does not exist.")
)
