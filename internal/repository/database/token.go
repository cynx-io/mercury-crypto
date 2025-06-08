package database

import (
	"context"
	"gorm.io/gorm"
	"mercury/internal/model/entity"
)

type TblToken struct {
	DB *gorm.DB
}

func NewTblToken(DB *gorm.DB) *TblToken {
	return &TblToken{
		DB: DB,
	}
}

func (db *TblToken) InsertToken(ctx context.Context, token entity.TblToken) (int, error) {
	err := db.DB.WithContext(ctx).Create(&token).Error
	if err != nil {
		return 0, err
	}

	return token.Id, nil
}

func (db *TblToken) GetTokenByTokenId(ctx context.Context, value string) (*entity.TblToken, error) {
	var token entity.TblToken
	err := db.DB.WithContext(ctx).Where("token_id = ?", value).First(&token).Error
	if err != nil {
		return &token, err
	}

	return &token, nil
}
