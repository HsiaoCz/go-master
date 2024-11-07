package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
)

type UserModInter interface {
	CreateUser(context.Context, *types.Users) (*types.Users, error)
}
