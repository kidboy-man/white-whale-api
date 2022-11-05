package repository

import (
	"encoding/json"
	"fetch-app/conf"
	"fetch-app/models"
	"fmt"
	"net/http"
	"time"

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

func NewCurrencyRepository() CurrencyRepository {
	return &currencyRepository{
		client: http.DefaultClient,
		cache:  cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (r *currencyRepository) Convert(from, to string, amount float64) (currency *models.Currency, err error) {
	key := fmt.Sprintf("convert-%s-%s-%v", from, to, amount)
	cached, _ := r.cache.Get(key)
	if cached != nil {
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

	err = r.cache.Add(key, currency, 5*time.Minute)
	return
}

func (r *currencyRepository) GetRate(from, to string) (rateInfo *models.Info, err error) {
	key := fmt.Sprintf("rate-%s-%s", from, to)
	cached, _ := r.cache.Get(key)
	if cached != nil {
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
	err = r.cache.Add(key, rateInfo, 5*time.Minute)
	return
}
