package handlers

import "github.com/HsiaoCz/go-master/recommend/mod"

type RecodeHandlers struct {
	record mod.RecordModInter
}

func RecodeHandlersInit(record mod.RecordModInter) *RecodeHandlers {
	return &RecodeHandlers{
		record: record,
	}
}
