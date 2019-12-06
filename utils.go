package main

import (
	"errors"
	"net/url"
	"strconv"
)

func setPageBounds(limit, page, length int64) (lower int64, upper int64) {
	if limit == 0 {
		return 0, length
	}
	lower = limit * page
	upper = lower + limit
	if length < lower {
		lower, upper = 0, 0
	} else if length < upper {
		upper = length
	}
	return lower, upper
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
