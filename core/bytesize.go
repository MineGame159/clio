package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ByteSize uint64

func ParseByteSize(str string) (ByteSize, error) {
	splits := strings.Split(str, " ")
	if len(splits) != 2 {
		return 0, errors.New("invalid byte size format")
	}

	amount, err := strconv.ParseFloat(splits[0], 64)
	if err != nil {
		return 0, err
	}

	var unit uint64

	switch strings.ToLower(splits[1]) {
	case "b":
		unit = 1
	case "kb":
		unit = 1024
	case "mb":
		unit = 1024 * 1024
	case "gb":
		unit = 1024 * 1024 * 1024
	case "tb":
		unit = 1024 * 1024 * 1024 * 1024

	default:
		return 0, errors.New("invalid byte size unit '" + splits[1] + "'")
	}

	return ByteSize(amount * float64(unit)), nil
}

func (b ByteSize) String() string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(b)/(1024))
	}
	if b < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(b)/(1024*1024))
	}
	if b < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(b)/(1024*1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(b)/(1024*1024*1024*1024))
}
