package database

import (
	"context"
	"gorm.io/gorm"
	"mercury/internal/model/entity"
)

const (
	endpointGetCoin = "GET_COIN"
)

type TblCacheCoinGeckoClient struct {
	DB *gorm.DB
}

func NewCacheCoinGeckoClient(DB *gorm.DB) *TblCacheCoinGeckoClient {
	return &TblCacheCoinGeckoClient{
		DB: DB,
	}
}

func (db *TblCacheCoinGeckoClient) UpsertGetCoinCacheByTokenId(ctx context.Context, coingeckoCache entity.TblCacheCoinGecko) error {

	coingeckoCache.Endpoint = endpointGetCoin

	tx := db.DB.WithContext(ctx).Where("token_id = ?", coingeckoCache.TokenId).Updates(&coingeckoCache)
	if tx.Error != nil {

		return tx.Error
	}

	if tx.RowsAffected == 0 {
		_, err := db.InsertGetCoinCache(ctx, coingeckoCache)
		return err
	}

	return nil
}

func (db *TblCacheCoinGeckoClient) InsertGetCoinCache(ctx context.Context, coingeckoCache entity.TblCacheCoinGecko) (int, error) {

	coingeckoCache.Endpoint = endpointGetCoin

	err := db.DB.WithContext(ctx).Create(&coingeckoCache).Error
	if err != nil {
		return 0, err
	}

	return coingeckoCache.Id, nil
}

func (db *TblCacheCoinGeckoClient) GetGetCoinCacheByTokenId(ctx context.Context, tokenId string) (*entity.TblCacheCoinGecko, error) {
	var coingeckoCache entity.TblCacheCoinGecko
	err := db.DB.WithContext(ctx).Where("endpoint = ? AND token_id = ?", endpointGetCoin, tokenId).First(&coingeckoCache).Error
	if err != nil {
		return &coingeckoCache, err
	}

	return &coingeckoCache, nil
}

func (db *TblCacheCoinGeckoClient) GetMultipleGetCoinCacheByTokenId(ctx context.Context, tokenIds []string) ([]*entity.TblCacheCoinGecko, error) {
	var coingeckoCaches []*entity.TblCacheCoinGecko
	err := db.DB.WithContext(ctx).Where("endpoint = ? AND token_id IN (?)", endpointGetCoin, tokenIds).Find(&coingeckoCaches).Error
	if err != nil {
		return coingeckoCaches, err
	}

	return coingeckoCaches, nil
}
