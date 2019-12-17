package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"os/exec"
	"regexp"
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

func stringIDToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return [12]byte{}, err
	}
	return objectID, nil
}

// reads through the temperatures provided by lm-sensors
// giving every device a different reading
// on my laptop it provides 3 readings thus 3 devices
func valueService(n int) float64 {
	out, _ := exec.Command("sensors").Output()
	regexResult := regexp.MustCompile(".{4}.C\\s").FindAll(out, 6)
	result, _ := strconv.ParseFloat(string(regexResult[n]), 64)
	return result
}
