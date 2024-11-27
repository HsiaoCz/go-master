package service

import (
	"context"

	v1 "github.com/HsiaoCz/go-master/g-record/pb/v1"
)

type RecordService struct {
	v1.UnimplementedRecordServiceServer
}

func (r *RecordService) CreateRecord(ctx context.Context, create_record_request *v1.CreateRecordRequest) (*v1.CreateRecordResponse, error) {
	return nil, nil
}
