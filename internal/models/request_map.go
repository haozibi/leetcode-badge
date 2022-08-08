package models

import (
	"reflect"
)

type RequestType int

const (
	RequestTypeUserProfile RequestType = iota + 1
	RequestTypeUserQuestionProgress
	RequestTypeUserContestRankingInfo
)

type RequestConfig struct {
	Desc     string
	URI      string
	Method   string
	Query    string
	Response reflect.Type
}

var LeetCodeRequestMap = make(map[RequestType]RequestConfig)

func Register(rt RequestType, c RequestConfig) {
	if _, ok := LeetCodeRequestMap[rt]; ok {
		panic("request type duplicate")
	}
	LeetCodeRequestMap[rt] = c
}
