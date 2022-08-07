package utils

import "strconv"

func Uint64(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}

func Int64(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func Int(id string) (int, error) {
	return strconv.Atoi(id)
}

func Uint(id string) (uint, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return uint(0), err
	}
	return uint(ID), nil
}