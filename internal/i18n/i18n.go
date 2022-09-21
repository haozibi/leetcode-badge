package i18n

import (
	"encoding/json"
	"strings"

	"github.com/haozibi/leetcode-badge/internal/statics"
)

func InitI18n() error {
	for _, v := range statics.ListI18n() {
		var p Config
		if err := json.Unmarshal(v.Body, &p); err != nil {
			return err
		}
		i18nMap[strings.ReplaceAll(v.Name, ".json", "")] = p
	}
	return nil
}

func Get(name string) *Config {
	if val, ok := i18nMap[name]; ok {
		return &val
	}

	return nil
}

var i18nMap = make(map[string]Config)

type Config struct {
	QuestionProcess QuestionProcess `json:"question_process"`
	ContestRanking  ContestRanking  `json:"contest_ranking"`
}

type QuestionProcess struct {
	KeyDetail string `json:"key_detail"`
	KeyEasy   string `json:"key_easy"`
	KeyMedium string `json:"key_medium"`
	KeyHard   string `json:"key_hard"`
	KeySolved string `json:"key_solved"`
}

type ContestRanking struct {
	KeyDetail        string `json:"key_detail"`
	KeyBeats         string `json:"key_beats"`
	KeyScore         string `json:"key_score"`
	KeyGlobalRanking string `json:"key_global_ranking"`
	KeyLocalRanking  string `json:"key_local_ranking"`
}
