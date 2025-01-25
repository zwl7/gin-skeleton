/**
 * 工具
 * @desc 接口返回的通用类型：方便控制器json格式化
 * ---------------------------------------------------------------------
 * @author		Super <super@papa.com.cn>
 * @date		2019-11-07
 * @copyright	cooper
 * ---------------------------------------------------------------------
 */

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// NewMP ...
func NewMP() Mp {
	return make(map[string]interface{}, 0)
}

// StatisticListJoin 报表类分页列表
func StatisticListJoin(time_option, filter, slice, total, page, size interface{}) Mp {
	return NewMP().Set("time_option", time_option).Set("filter", filter).Set("list", slice).Set("count", total).Set("page", page).Set("size", size)
}

// ListJoin 分页列表
func ListJoin(slice, total, page, size interface{}) Mp {
	return NewMP().Set("list", slice).Set("count", total).Set("page", page).Set("size", size)
}

// 线路列表
func LineListJoin(slice, total, page, size, identifier interface{}) Mp {
	return NewMP().Set("list", slice).Set("count", total).Set("page", page).Set("size", size).Set("identifier", identifier)
}

// Mp ...
type Mp map[string]interface{}

func (that Mp) Init(mode interface{}) (Mp, error) {
	_json, _err := json.Marshal(mode)
	if _err != nil {
		return nil, _err
	}
	if _err := json.Unmarshal(_json, &that); _err != nil {
		return nil, _err
	}
	return that, nil
}

// Set ...
func (that Mp) Set(key string, value interface{}) Mp {
	that[key] = value
	return that
}

func (that Mp) SetString(key string, value interface{}) Mp {
	that[key] = ToString(value)
	return that
}

// Has ...
func (that Mp) Has(keys ...string) bool {
	for _, _key := range keys {
		if _, _ok := that[_key]; !_ok {
			return false
		}
	}
	return true
}

// VerifyKey ...
func (that Mp) VerifyKey(keys ...string) string {
	for _, _v := range keys {
		if _, _ok := that[_v]; !_ok {
			return _v
		}
	}
	return ""
}

func (that Mp) Copy() Mp {
	_mp := NewMP()
	for _k, _v := range that {
		_mp.Set(_k, _v)
	}
	return _mp
}

// Get ...
func (that Mp) Get(key string) interface{} {
	return that[key]
}

// Del ...
func (that Mp) Del(key string) Mp {
	if that.Has(key) {
		delete(that, key)
	}
	return that
}

func (that Mp) String(key string) (string, error) {
	if that.Has(key) {
		switch n := that.Get(key).(type) {
		case int64:
			return strconv.FormatInt(n, 10), nil
		case string:
			return n, nil
		case int:
			return strconv.Itoa(n), nil
		case float32:
			return strconv.Itoa(int(n)), nil
		case float64:
			return strconv.Itoa(int(n)), nil
		case []byte:
			return string(n), nil
		case bool:
			if n {
				return "1", nil
			}
			return "0", nil
		}
		return "", errors.New("the type of data in map , can not to string")
	}
	return "0", errors.New("map not have this key")
}

// DefaultString ...
func (that Mp) DefaultString(key string, defaultValue string) string {
	if that.Has(key) {
		switch n := that.Get(key).(type) {
		case int64:
			return strconv.FormatInt(n, 10)
		case string:
			return n
		case int:
			return strconv.Itoa(n)
		case float32:
			return strconv.Itoa(int(n))
		case float64:
			return strconv.Itoa(int(n))
		case []byte:
			return string(n)
		case bool:
			if n {
				return "1"
			}
			return "0"
		}
		return "0"
	}
	return defaultValue
}

// removeDuplicateElement 删除数组重复元素...
func (that Mp) RemoveDuplicateElement(originals interface{}) (interface{}, error) {
	temp := map[string]struct{}{}
	switch slice := originals.(type) {
	case []string:
		result := make([]string, 0, len(originals.([]string)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result, nil
	case []int64:
		result := make([]int64, 0, len(originals.([]int64)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result, nil
	default:
		return nil, errors.New("Unknown type...")
	}
}

// ToJson ...
func (that Mp) ToJson() ([]byte, error) {
	return json.Marshal(that)
}

func (that Mp) Val(key string) InterfaceVal {
	return InterfaceVal{
		value: that[key],
	}
}

func (that Mp) List(key string) ([]interface{}, bool) {
	_data, _ok := that[key]
	if !_ok {
		return nil, false
	}
	_list, _ok := _data.([]interface{})
	if !_ok {
		return nil, false
	}
	return _list, true
}

type InterfaceVal struct {
	value interface{}
}

func (that InterfaceVal) String() string {
	return ToString(that.value)
}

func (that InterfaceVal) Int() int {
	return ToInt(that.value)
}

func (that InterfaceVal) Int64() int64 {
	return ToInt64(that.value)
}

// JsonMerge....
func (that Mp) JsonMerge(dst, src Mp) Mp {
	return jsMerge(dst, src, 0)
}

var jsonMergeDepth = 32

func jsMerge(dst, src Mp, depth int) Mp {
	if depth > jsonMergeDepth {
		return dst
		// panic("too deep!")
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := jsMapify(srcVal)
			dstMap, dstMapOk := jsMapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = jsMerge(dstMap, srcMap, depth+1)
			}
		}
		dst[key] = srcVal
	}
	return dst
}
func jsMapify(i interface{}) (Mp, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := Mp{}
		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return Mp{}, false
}

// 将十六进制的 数值 1-16 转换为 00 01 02 ... 0f
func GetHexToHexStr(i int8) string {

	hex_value_map_str := map[int8]string{
		1: "00", 2: "01", 3: "02", 4: "03", 5: "04", 6: "05", 7: "06", 8: "07", 9: "08", 10: "09", 11: "0a", 12: "0b", 13: "0c", 14: "0d", 15: "0e", 16: "0f",
		17: "10", 18: "11", 19: "12", 20: "13", 21: "14", 22: "15", 23: "16", 24: "17", 25: "18", 26: "19", 27: "1a", 28: "1b", 29: "1c", 30: "1d", 31: "1e", 32: "1f",
		33: "20", 34: "21", 35: "22", 36: "23", 37: "24", 38: "25", 39: "26", 40: "27", 41: "28", 42: "29", 43: "2a", 44: "2b", 45: "2c", 46: "2d", 47: "2e", 48: "2f",
	}

	return hex_value_map_str[i]
}

// 获取 开/关灯完整指令
func GetCRCCode(arr map[int64]string) string {

	last_crc_code := ""
	var i int64
	i = 0
	for key, _ := range arr {
		key++
		if i == 1 {
			last_crc_code = arr[i]
		} else if i > 1 {
			last_crc_code = DecHex(HexDec(last_crc_code) ^ HexDec(arr[i]))
		}
		i++
	}
	return strings.ToLower(MapToString(arr) + last_crc_code)
}

// 把十进制转换为十六进制
func DecHex(n int64) string {
	if n < 0 {
		log.Println("Decimal to hexadecimal error: the argument must be greater than zero.")
		return ""
	}
	if n == 0 {
		return "0"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

// 把十六进制转换为十进制
func HexDec(h string) (n int64) {
	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := float64(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			log.Println("Hexadecimal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(16, float64(l-i-1))
	}
	return int64(d)
}

// 将map转成string
func MapToString(arr map[int64]string) string {

	a := ""
	var i int64
	i = 0
	for key, _ := range arr {
		key++
		a += arr[i]
		i++
	}

	return a
}

// 将字符串进行分割为 map[int]string
func DivisionMap(a string) map[int64]string {

	FF := make(map[int64]string) //先声明一个切片来储存分割后的内容
	DD := []rune(a)              //需要分割的字符串内容，将它转为字符，然后取长度。
	var j int64
	j = 0
	for i := 0; i < len(DD); i = i + 2 {
		a := ""
		a = string(DD[i]) + string(DD[i+1])
		FF[j] = a
		j++
	}

	return FF
}

// 将json参数-jsonContent解析到map
func UnmarJson(params string) Mp {
	_jsonContent := NewMP()
	err := json.Unmarshal([]byte(params), &_jsonContent)
	if err != nil {
		fmt.Println("err = ", err)
	}
	return _jsonContent
}

// 将map参数-jsonContent编码到json
func MarJson(params interface{}) string {
	result, err := json.MarshalIndent(params, "", "    ")
	if err != nil {
		fmt.Println("err = ", err)
	}
	return string(result)
}

// 总金额均摊到sku子集
func (that Mp) AverageItemMoney(averageTotal int64) {
	//当销售价为0／小于0时，item均为0
	if averageTotal <= 0 {
		for _k, _ := range that {
			that[_k] = 0
		}
		return
	}
	//开始均摊
	total := 0.00
	for _, _v := range that {
		total = (total*1000 + ToFloat64(_v)*1000) / 1000
	}
	count := len(that)
	add := 0.00
	i := 1
	for _k, _v := range that {
		var _proportionValue = 0.0
		if total == 0 {
			_proportionValue = ToFloat64(ChangeNumber(float64(averageTotal / ToInt64(count))))
		} else {
			_proportionValue = ToFloat64(ChangeNumber(float64(averageTotal) * ToFloat64(_v) / total))
		}
		if i >= count {
			that[_k] = (float64(averageTotal*1000) - (add * 1000)) / 1000
		} else {
			that[_k] = _proportionValue
		}
		i++
		add = (add*1000 + ToFloat64(ToString(that[_k]))*1000) / 1000.0
	}
}

// 截取小数位数,默认保留2位小数
func ChangeNumber(f float64, m ...int) interface{} {
	_point := 2
	if len(m) > 0 {
		_point = m[0]
	}
	n := strconv.FormatFloat(f, 'f', -1, 32)
	if n == "" {
		return ""
	}
	if _point >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if len(newn) < 2 || _point >= len(newn[1]) {
		return n
	}
	_sum := ToFloat64(newn[0]+"."+newn[1][:_point]) * 100 / 100
	return _sum
}
