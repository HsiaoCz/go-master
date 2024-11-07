package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
)

type BookModInter interface {
	CreateBook(context.Context, *types.Books) (*types.Books, error)
}
