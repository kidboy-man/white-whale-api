package helpers

import (
	"fetch-app/utils"
	"sort"
)

func CalculateAggregate(data []float64) (min, max, median, average float64) {
	sort.Float64s(data)
	min = data[0]
	max = data[len(data)-1]
	median = utils.Median(data)
	average = utils.Average(data)
	return
}
