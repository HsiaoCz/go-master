package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/pkg"
	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type UserModInter interface {
	CreateUser(context.Context, *types.Users) (*types.Users, error)
	GetUserByID(context.Context, string) (*types.Users, error)
	GetUserByPhoneAndPassword(context.Context, *types.Login) error
	DeleteUserByID(context.Context, string) error
}

type UserMod struct {
	db *gorm.DB
}

func UserModInit(db *gorm.DB) *UserMod {
	return &UserMod{
		db: db,
	}
}

func (u *UserMod) CreateUser(ctx context.Context, user *types.Users) (*types.Users, error) {
	tx := u.db.Debug().WithContext(ctx).Model(&types.Users{}).Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (u *UserMod) GetUserByID(ctx context.Context, user_id string) (*types.Users, error) {
	var user types.Users
	tx := u.db.Debug().Model(&types.Users{}).Where("user_id = ?", user_id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (u *UserMod) GetUserByPhoneAndPassword(ctx context.Context, login *types.Login) error {
	var user types.Users
	tx := u.db.Debug().Model(&types.Users{}).Where("phone = ? AND hash_password = ?", login.Phone, pkg.EncryPassword(login.Password)).First(&user)
	return tx.Error
}

func (u *UserMod) DeleteUserByID(ctx context.Context, user_id string) error {
	var user types.Users
	tx := u.db.Debug().WithContext(ctx).Model(&types.Users{}).Where("user_id = ?", user_id).Delete(&user)
	return tx.Error
}
