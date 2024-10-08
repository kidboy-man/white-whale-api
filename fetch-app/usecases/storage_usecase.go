package usecase

import (
	"encoding/json"
	"fetch-app/conf"
	"fetch-app/helpers"
	"fetch-app/models"
	repository "fetch-app/repositories"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type StorageUsecase interface {
	GetAggregatedStorages() (aggregateds []*models.AggregatedStorage, err error)
	GetStorages() (storages []*models.Storage, err error)
}

type storageUsecase struct {
	storageRepo  repository.StorageRepository
	currencyRepo repository.CurrencyRepository
}

func NewStorageUsecase() StorageUsecase {
	storageRepo := repository.NewStorageRepository()
	currencyRepo := repository.NewCurrencyRepository(conf.AppConfig.Cache)
	return &storageUsecase{
		storageRepo:  storageRepo,
		currencyRepo: currencyRepo,
	}
}

func (u *storageUsecase) GetStorages() (storages []*models.Storage, err error) {
	storages, err = u.storageRepo.FetchStorages()
	if err != nil || len(storages) == 0 {
		return
	}

	rateInfo, err := u.currencyRepo.GetRate("IDR", "USD")
	if err != nil {
		return
	}

	rate := rateInfo.Rate
	for _, storage := range storages {
		priceIDR, _ := strconv.ParseFloat(strings.TrimSpace(storage.PriceIDR), 64)
		storage.PriceUSD = fmt.Sprintf("%v", priceIDR*rate)
	}
	return
}

func (u *storageUsecase) GetAggregatedStorages() (aggregateds []*models.AggregatedStorage, err error) {
	storages, err := u.GetStorages()
	if err != nil {
		return
	}

	mapPricesIDRToProvince := make(map[string][]float64)
	mapPricesUSDToProvince := make(map[string][]float64)
	mapSizesToProvince := make(map[string][]float64)
	mapTotalSizeToProvince := make(map[string]float64)
	mapTotalPriceIDRToProvince := make(map[string]float64)
	mapTotalPriceUSDToProvince := make(map[string]float64)

	for _, storage := range storages {
		if storage.ParsedDate == "" && storage.Timestamp == "" {
			continue
		}

		date, err := time.Parse(time.RFC3339, storage.ParsedDate)
		if err != nil {
			jsonBytes, _ := json.Marshal(storage)
			logs.Debug("storage parsed", string(jsonBytes))
			logs.Error("error parsing date %s: %v", storage.ParsedDate, err)
			number, err := strconv.ParseFloat(storage.Timestamp, 64)
			if err != nil {
				logs.Error(err)
				return nil, err
			}

			date = time.Unix(int64(number), 0)
		}
		year, week := date.ISOWeek()
		size, _ := strconv.ParseFloat(strings.TrimSpace(storage.Size), 64)
		priceIDR, _ := strconv.ParseFloat(strings.TrimSpace(storage.PriceIDR), 64)
		priceUSD, _ := strconv.ParseFloat(strings.TrimSpace(storage.PriceUSD), 64)

		key := fmt.Sprintf("%s-%v-%v", storage.Province, year, week)
		mapSizesToProvince[key] = append(mapSizesToProvince[key], size)
		mapPricesIDRToProvince[key] = append(mapPricesIDRToProvince[key], priceIDR)
		mapPricesUSDToProvince[key] = append(mapPricesUSDToProvince[key], priceUSD)
		mapTotalSizeToProvince[key] += size
		mapTotalPriceIDRToProvince[key] += priceIDR
		mapTotalPriceUSDToProvince[key] += priceUSD
	}

	mapSummaries := make(map[string][]*models.Summary)
	for key, val := range mapSizesToProvince {

		splitKey := strings.Split(key, "-")
		province := splitKey[0]
		weekGroup := fmt.Sprintf("%s-%s", splitKey[1], splitKey[2])

		sizeMin, sizeMax, sizeMedian, sizeAvg := helpers.CalculateAggregate(val)
		priceIDRMin, priceIDRMax, priceIDRMedian, priceIDRAvg := helpers.CalculateAggregate(mapPricesIDRToProvince[key])
		priceUSDMin, priceUSDMax, priceUSDMedian, priceUSDAvg := helpers.CalculateAggregate(mapPricesUSDToProvince[key])

		mapSummaries[province] = append(mapSummaries[province], &models.Summary{
			Week: weekGroup,
			Size: &models.Aggregate{
				Min:     sizeMin,
				Max:     sizeMax,
				Median:  sizeMedian,
				Average: sizeAvg,
			},
			PriceIDR: &models.Aggregate{
				Min:     priceIDRMin,
				Max:     priceIDRMax,
				Median:  priceIDRMedian,
				Average: priceIDRAvg,
			},
			PriceUSD: &models.Aggregate{
				Min:     priceUSDMin,
				Max:     priceUSDMax,
				Median:  priceUSDMedian,
				Average: priceUSDAvg,
			},
		})
	}

	for key, val := range mapSummaries {
		aggregateds = append(aggregateds, &models.AggregatedStorage{
			Province:  key,
			Summaries: val,
		})
	}
	return
}
