package utils

import (
	"github.com/yanyiwu/gojieba"
)

func WordCount(rawText string) []string {
	var words []string
	x := gojieba.NewJieba()
	defer x.Free()
	words = x.Cut(rawText, true)
	return words
}

