package utils

import (
	"sort"
	"strings"
)

// 取字符串中{} [] 中的数据返回
func GetValue(data string, params ...string) string {
	lstr, rstr, gstr := "{", "}", ""
	isL, isR, str := false, false, ""
	if len(params) > 0 {
		lstr = params[0]
	}
	if len(params) > 1 {
		rstr = params[1]
	}
	if len(params) > 2 {
		gstr = params[2]
	}
	if ok := strings.Index(data, lstr); ok < 0 {
		return str
	}
	// 截取出包含在{}里面的内容
	ss := strings.Split(data, "")
	for key, val := range ss {
		if val == rstr {
			if isR == false {
				str += gstr
			}
			isL = false
			continue
		}
		if val == lstr && val != rstr {
			if key+1 < len(ss) && ss[key+1] != rstr {
				str += ss[key+1]
			}
			isL, isR = true, false
			continue
		}
		if isL && val != rstr && isR == false {
			if val == rstr {
				isR = true
				continue
			}
			if ss[key+1] != rstr && key+1 < len(ss) {
				str += ss[key+1]
			}
			isL = true
		}
	}
	return str[0 : len(str)-1]
}

// SnakeString 将字符串"驼峰格式"转化成"下划线格式"
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func SnakeString2(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	if num > 1 {
		_slice := strings.Split(s, "")
		if _slice[num-1] == "D" && _slice[num-2] == "I" {
			_slice[num-1] = "d"
		}
		s = strings.Join(_slice, "")
	}
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func StrCount(s string) int {
	_runeSlice := []rune(s)
	count := 0
	for _, _v := range _runeSlice {
		_b := []byte(string(_v))
		switch len(_b) {
		case 3: //一般是中文，宽度占2个英文字符位
			count = count + 2
		default:
			count = count + len(_b)
		}
	}
	return count
}

func StrSubByCount(s string, count int) (string, string) {
	_runeSlice := []rune(s)
	_byte := []byte{}
	_i := 0
	for _, _v := range _runeSlice {
		_b := []byte(string(_v))
		switch len(_b) {
		case 3: //一般是中文，宽度占2个英文字符位
			count = count - 2
		default:
			count = count - len(_b)
		}
		if count < 0 {
			break
		}
		_byte = append(_byte, _b...)
		_i++
	}
	return string(_byte), string(_runeSlice[_i:])
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// 字符串数组去重、去空
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	sort.Strings(a)
	length := len(a)
	for i := 0; i < length; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
