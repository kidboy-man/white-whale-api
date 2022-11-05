package utils

func CalculateOffset(limit, page int) (offset int) {
	offset = (page - 1) * limit
	return
}
