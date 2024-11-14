package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
)

type RecordModInter interface {
	CreateRecord(context.Context, *types.Records) (*types.Records, error)
}

