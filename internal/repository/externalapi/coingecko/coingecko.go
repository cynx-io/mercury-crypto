package coingecko

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mercury/internal/constant"
	"mercury/internal/dependencies"
	"mercury/internal/helper"
	"mercury/internal/model/entity"
	"mercury/internal/pkg/logger"
	"mercury/internal/repository/database"
	"mercury/internal/repository/externalapi/coingecko/coingeckomodel"
	"sync"
)

type Client struct {
	config   *dependencies.CoinGeckoConfig
	tblCache *database.TblCacheCoinGeckoClient
}

func NewClient(config *dependencies.CoinGeckoConfig, tblCoingeckoCache *database.TblCacheCoinGeckoClient) *Client {
	return &Client{
		config:   config,
		tblCache: tblCoingeckoCache,
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{
		"x-cg-demo-api-key": c.config.ApiKey,
	}
}

func (c *Client) Search(ctx context.Context, query string) (coingeckomodel.SearchResponse, error) {
	url := c.config.BaseUrl + "/api/v3/search?query=" + query
	return helper.HttpRequest[coingeckomodel.SearchResponse](url, nil, c.headers(), helper.GET)
}

func (c *Client) GetCoin(ctx context.Context, tokenId string, useCache bool) (coingeckomodel.GetCoinResponse, error) {

	if useCache {
		// Check for cache
		cachedByte, err := c.tblCache.GetGetCoinCacheByTokenId(ctx, tokenId)
		if err == nil {
			cachedResponse, err := cachedByte.ToGetCoinResponse()
			if err != nil {
				logger.Error("Error on converting cached byte to get coin response: ", err)
				return cachedResponse, err
			}

			return cachedResponse, nil
		}

		if !errors.Is(err, constant.ErrDatabaseNotFound) {
			logger.Warn("Error on getting cache: ", err)
		}
	}

	// No Cache found
	url := c.config.BaseUrl + "/api/v3/coins/" + tokenId
	resp, err := helper.HttpRequest[coingeckomodel.GetCoinResponse](url, nil, c.headers(), helper.GET)
	if err != nil {
		return resp, err
	}

	go func() {
		// Cache response
		stringResponse, err := json.Marshal(resp)
		if err != nil {
			logger.Warn("Error on marshalling cache: ", err)
			return
		}

		cachedResponse := entity.TblCacheCoinGecko{
			TokenId:  tokenId,
			Response: string(stringResponse),
		}

		_, err = c.tblCache.InsertGetCoinCache(ctx, cachedResponse)
		if err != nil {
			logger.Warn("Error on inserting cache: ", err)
			return
		}
	}()

	return resp, nil
}

func (c *Client) GetCoins(ctx context.Context, tokenIds []string) ([]coingeckomodel.GetCoinResponse, error) {
	var results []coingeckomodel.GetCoinResponse

	// Get all cached coins
	cachedCoins, err := c.tblCache.GetMultipleGetCoinCacheByTokenId(ctx, tokenIds)
	if err != nil {
		logger.Warn("Error on getting multiple cached coins: ", err)
	}

	// Prepare a map of tokenId -> response
	cachedMap := make(map[string]coingeckomodel.GetCoinResponse)
	for _, cached := range cachedCoins {
		if cached == nil {
			continue
		}
		resp, err := cached.ToGetCoinResponse()
		if err != nil {
			logger.Warn("Error decoding cache for tokenId : ", cached.TokenId, ": ", err)
			continue
		}
		cachedMap[cached.TokenId] = resp
	}

	// Track tokenIds not found in cache
	var missed []string
	for _, id := range tokenIds {
		if resp, ok := cachedMap[id]; ok {
			results = append(results, resp)
		} else {
			missed = append(missed, id)
		}
	}

	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		wgErrors []error
	)

	for _, tokenId := range missed {
		wg.Add(1)

		go func(tokenId string) {
			defer wg.Done()

			resp, err := helper.HttpRequest[coingeckomodel.GetCoinResponse](
				c.config.BaseUrl+"/api/v3/coins/"+tokenId,
				nil,
				c.headers(),
				helper.GET,
			)
			if err != nil {
				logger.Warn("Error fetching coin from API for tokenId ", tokenId, ": ", err)
				mutex.Lock()
				wgErrors = append(wgErrors, fmt.Errorf("tokenId %s: %w", tokenId, err))
				mutex.Unlock()
				return
			}

			// Append result safely
			mutex.Lock()
			results = append(results, resp)
			mutex.Unlock()

			// Cache in background
			go func(tokenId string, resp coingeckomodel.GetCoinResponse) {
				stringResp, err := json.Marshal(resp)
				if err != nil {
					logger.Warn("Failed to marshal response for caching:", err)
					return
				}
				cacheEntry := entity.TblCacheCoinGecko{
					TokenId:  tokenId,
					Response: string(stringResp),
				}
				if _, err := c.tblCache.InsertGetCoinCache(context.Background(), cacheEntry); err != nil {
					logger.Warn("Failed to insert cache for tokenId ", tokenId, ": ", err)
				}
			}(tokenId, resp)
		}(tokenId)
	}

	wg.Wait()

	// Optionally return aggregate error
	if len(wgErrors) > 0 {
		return results, fmt.Errorf("some requests failed: %v", wgErrors)
	}

	return results, nil
}
