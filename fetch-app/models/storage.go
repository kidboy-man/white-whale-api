package models

import "time"

// to fetch from json:
//
//	{
//	    "uuid": "2c75f19f-ac7a-4c28-9459-8fbb64c8cb00",
//	    "komoditas": "GURAME",
//	    "area_provinsi": "JAWA TENGAH",
//	    "area_kota": " PURWOREJOL",
//	    "size": "40",
//	    "price": "55000",
//	    "tgl_parsed": "2022-01-01T13:11:46Z",
//	    "timestamp": "1641042706799"
//	  },
type Storage struct {
	UUID       string    `json:"uuid"`
	Comodity   string    `json:"komoditas"`
	Province   string    `json:"area_provinsi"`
	City       string    `json:"area_kota"`
	Size       string    `json:"size"`
	PriceIDR   string    `json:"price"`
	PriceUSD   string    `json:"price_usd"`
	ParsedDate time.Time `json:"tgl_parsed"`
	Timestamp  string    `json:"timestamp"`
}

type AgragetedStorage struct {
	Province string               `json:"area_provinsi"`
	Size     *Agregate            `json:"size"`
	Price    map[string]*Agregate `json:"price"`
}

type Agregate struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Median  float64 `json:"median"`
	Average float64 `json:"average"`
}
