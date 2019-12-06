package main

import (
	"errors"
	"net/url"
	"strconv"
)

func setPageBoundsToInt64(limit, page int) (lower int64, upper int64) {
	return int64(page * limit), int64(page*limit + limit)
}

func convertToPositiveInteger(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, errors.New("input is a negative number")
	}
	return id, nil
}

func readIntFromQueryParameter(url *url.URL, param string, defaultValue int) (int, error) {
	valueStr := url.Query().Get(param)
	if valueStr == "" {
		return defaultValue, nil
	}
	return convertToPositiveInteger(valueStr)
}
