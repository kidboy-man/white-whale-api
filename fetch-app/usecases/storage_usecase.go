package usecase

import (
	"fetch-app/models"
	repository "fetch-app/repositories"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type StorageUsecase interface {
	GetAgregatedStorages() (agregateds []*models.AgragetedStorage, err error)
	GetStorages() (storages []*models.Storage, err error)
}

type storageUsecase struct {
	storageRepo  repository.StorageRepository
	currencyRepo repository.CurrencyRepository
}

func NewStorageUsecase() StorageUsecase {
	storageRepo := repository.NewStorageRepository()
	currencyRepo := repository.NewCurrencyRepository()
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

func (u *storageUsecase) GetAgregatedStorages() (agregateds []*models.AgragetedStorage, err error) {
	storages, err := u.GetStorages()
	sort.SliceStable(storages, func(i, j int) bool {
		return storages[i].ParsedDate.Before(storages[j].ParsedDate)
	})

	mapPricesIDRToProvince := make(map[string][]float64)
	mapPricesUSDToProvince := make(map[string][]float64)
	mapSizesToProvince := make(map[string][]float64)

	for _, storage := range storages {
		size, _ := strconv.ParseFloat(strings.TrimSpace(storage.Size), 64)
		priceIDR, _ := strconv.ParseFloat(strings.TrimSpace(storage.PriceIDR), 64)
		priceUSD, _ := strconv.ParseFloat(strings.TrimSpace(storage.PriceUSD), 64)

		mapSizesToProvince[storage.Province] = append(mapSizesToProvince[storage.Province], size)
		mapPricesIDRToProvince[storage.Province] = append(mapPricesIDRToProvince[storage.Province], priceIDR)
		mapPricesUSDToProvince[storage.Province] = append(mapPricesUSDToProvince[storage.Province], priceUSD)
	}

	for key, val := range mapSizesToProvince {
		agregateds = append(agregateds, &models.AgragetedStorage{
			Province: key,
			Size: &models.Agregate{
				Min:     val[0],
				Max:     val[0],
				Median:  val[0],
				Average: val[0],
			},
		})
	}
	return
}
