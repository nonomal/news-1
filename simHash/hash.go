package simHash

import (
	"fmt"
	"strconv"

	"github.com/yanyiwu/gojieba"
)

func hashWord(key []byte, weight float64, result *[]float64) {
	var hash uint32

	for _, b := range key {
		hash += uint32(b)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)

	binary := []byte(strconv.FormatUint(uint64(hash), 2))

	for index, byte := range binary {
		byteNum, _ := strconv.Atoi(string(byte))
		if byteNum == 1 {
			(*result)[index] = (*result)[index] + weight
		} else {
			(*result)[index] = (*result)[index] - weight
		}
	}
}

func Calc(str string) []byte {
	x := gojieba.NewJieba()
	defer x.Free()
	keywords := x.ExtractWithWeight(str, 5)
	fmt.Println("keywords", keywords)
	result := make([]float64, 64)
	for _, keyword := range keywords {
		hashWord([]byte(keyword.Word), keyword.Weight, &result)
	}
	hash := make([]byte, 64)
	for index, num := range result {
		if num > 0 {
			hash[index] = '1'
		} else {
			hash[index] = '0'
		}
	}
	return hash
}

func Distance(hash []byte) uint64 {
	var res uint64
	for _, hasByte := range hash {
		res <<= 1
		if hasByte == '1' {
			res += 1
		}
	}
	return res
}

func IsEqual(lDistance uint64, rDistance uint64, n int) bool {
	d := int(popcnt64Go(lDistance ^ rDistance))
	fmt.Println("d", lDistance, rDistance, d)
	return d <= n
}

func popcnt64Go(x uint64) uint64 {
	x = (x & 0x5555555555555555) + ((x & 0xAAAAAAAAAAAAAAAA) >> 1)
	x = (x & 0x3333333333333333) + ((x & 0xCCCCCCCCCCCCCCCC) >> 2)
	x = (x & 0x0F0F0F0F0F0F0F0F) + ((x & 0xF0F0F0F0F0F0F0F0) >> 4)
	x *= 0x0101010101010101
	return ((x >> 56) & 0xFF)
}
