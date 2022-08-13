package util

import (
	"strconv"

	"github.com/pkg/errors"
)

func ToStrSlice(vals []interface{}) ([]string, error) {
	var resultStrArr = make([]string, 0)
	for _, v := range vals {
		str, ok := v.(string)
		if !ok {
			return resultStrArr, errors.Errorf("Cast to string error")
		}
		resultStrArr = append(resultStrArr, str)
	}

	return resultStrArr, nil
}

func ParseUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}
