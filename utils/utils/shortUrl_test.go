package utils

import (
	"fmt"
	"net/url"
	"testing"
)

//func TestEncodeURL(t *testing.T) {
//	_encoder := NewURLEncoder(&URLEncoderConfig{alphabet: "mn6j2c4rv8bpygw95z7hsdaetxuk3fq"})
//	fmt.Println(_encoder.DecodeURL("pvv8mqy"))
//	_i := 1
//	for {
//		if _i > 10000 {
//			break
//		}
//		fmt.Println(_encoder.EncodeURL(uint64(_i)))
//		_i++
//	}
//}

func TestSlice(t *testing.T) {
	_str := "abcdeafgijklmnopqrstuzvwxyz"
	fmt.Println(_str+"最长连续字符串为：", getMaxStr(_str))
}

func TestMD5(t *testing.T) {
	_str := SignMD5("9ZwFiE9cbyvfoHnXuib2UUaEseoCnAnW", url.Values{
		"aaa": []string{"1.2356"},
		"bbb": []string{"1.0001"},
		"sss": []string{"1"},
		"ccc": []string{"safasfsafsafas"},
	})
	fmt.Println(_str)
}
