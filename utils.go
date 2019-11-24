package main

import (
	"errors"
	"net/url"
	"strconv"
)

func setPageBounds(limit int, page int, len int) (int, int) {
	if limit == 0 {
		return 0, len
	}
	lower := limit * page
	upper := lower + limit
	if len < lower {
		lower, upper = 0, 0
	} else if len < upper {
		upper = len
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
