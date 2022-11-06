package repository

import (
	"encoding/json"
	"fetch-app/conf"
	"fetch-app/models"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/patrickmn/go-cache"
)

type CurrencyRepository interface {
	Convert(from, to string, amount float64) (currency *models.Currency, err error)
	GetRate(from, to string) (rateInfo *models.Info, err error)
}

type currencyRepository struct {
	client *http.Client
	cache  *cache.Cache
}

func NewCurrencyRepository(cache *cache.Cache) CurrencyRepository {
	return &currencyRepository{
		client: http.DefaultClient,
		cache:  cache,
	}
}

func (r *currencyRepository) Convert(from, to string, amount float64) (currency *models.Currency, err error) {
	key := fmt.Sprintf("convert-%s-%s-%v", from, to, amount)
	cached, isExist := r.cache.Get(key)
	logs.Debug("is cache exist: %v", isExist)
	logs.Debug("cached value: %v", cached)
	if cached != nil {
		logs.Info("using cached %s", key)
		return cached.(*models.Currency), nil
	}

	url := fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/convert?to=%s&from=%s&amount=%v",
		to,
		from,
		amount,
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("apikey", conf.AppConfig.ApilayerKey)

	response, err := r.client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&currency)
	if err != nil {
		return
	}

	logs.Info("set cache for %s value %v", key, currency)
	err = r.cache.Add(key, currency, 5*time.Minute)
	if err != nil {
		logs.Error("error setting cache %s value %v", key, currency)
	}
	return
}

func (r *currencyRepository) GetRate(from, to string) (rateInfo *models.Info, err error) {
	key := fmt.Sprintf("rate-%s-%s", from, to)
	cached, isExist := r.cache.Get(key)
	logs.Debug("is cache exist: %v", isExist)
	logs.Debug("cached value: %v", cached)
	if cached != nil {
		logs.Info("using cached %s", key)
		return cached.(*models.Info), nil
	}

	url := fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/convert?to=%s&from=%s&amount=1000",
		to,
		from,
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("apikey", conf.AppConfig.ApilayerKey)

	response, err := r.client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	var currency *models.Currency
	err = json.NewDecoder(response.Body).Decode(&currency)
	if err != nil {
		return
	}

	rateInfo = currency.Info
	logs.Info("set cache for %s value %v", key, rateInfo)
	err = r.cache.Add(key, rateInfo, 5*time.Minute)
	if err != nil {
		logs.Error("error setting cache %s value %v", key, rateInfo)
	}
	return
}
