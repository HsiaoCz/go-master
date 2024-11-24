package handlers

import (
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
)

type RecordHandlers struct {
	record mod.RecordModInter
}

func RecordHandlersInit(record mod.RecordModInter) *RecordHandlers {
	return &RecordHandlers{
		record: record,
	}
}

func (rh *RecordHandlers) HandleCreateRecord(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rh *RecordHandlers) HandleGetRecordsByUserID(w http.ResponseWriter, r *http.Request) error {
	return nil
}
