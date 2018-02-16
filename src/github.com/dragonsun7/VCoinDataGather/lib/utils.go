package lib

import (
	"github.com/buger/jsonparser"
	"strconv"
)

func JsonStringValueToInt64(value []byte, key string) (int64, error) {
	s, err := jsonparser.GetString(value, key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}

func JsonStringValueToFloat64(value []byte, key string) (float64, error) {
	s, err := jsonparser.GetString(value, key)
	if err != nil {
		return 0.0, err
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}

	return f, nil
}


