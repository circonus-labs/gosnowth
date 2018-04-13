package gosnowth

import (
	"encoding/json"
	"testing"
)

func TestTextValue(t *testing.T) {
	var (
		data = `[[1380000000,"hello"],[1380000300,"world"]]`
		tvr  = TextValueResponse{}
	)
	if err := json.Unmarshal([]byte(data), &tvr); err != nil {
		t.Error("error unmarshalling: ", err)
	}
}
