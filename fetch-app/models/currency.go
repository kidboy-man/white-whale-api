package models

// {
// 	"date": "2022-11-05",
// 	"info": {
// 	  "rate": 6.4042627e-05,
// 	  "timestamp": 1667658363
// 	},
// 	"query": {
// 	  "amount": 20000,
// 	  "from": "IDR",
// 	  "to": "USD"
// 	},
// 	"result": 1.280853,
// 	"success": true
//   }

type Info struct {
	Rate      float64 `json:"rate"`
	Timestamp int     `json:"timestamp"`
}

type Query struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type Currency struct {
	Date    string  `json:"string"`
	Result  float64 `json:"result"`
	Success bool    `json:"success"`
	Info    *Info   `json:"info"`
	Query   *Query  `json:"query"`
}
