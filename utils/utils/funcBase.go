/**
 * 工具
 * @desc 常用的工具函数
 * ---------------------------------------------------------------------
 * @author		Super <super@papa.com.cn>
 * @date		2019-11-07
 * @copyright	cooper
 * ---------------------------------------------------------------------
 */

package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func TypeOf(t interface{}) string {
	_t := reflect.TypeOf(t)
	if _t != nil {
		if _t.Name() != "" {
			return _t.Name()
		}
		return _t.Elem().Name()
	}
	return "nil"
}

// ToInt ...
func ToInt(d interface{}) int {
	switch n := d.(type) {
	case int:
		return n
	case string:
		if n == "" {
			return 0
		}
		if strings.ToLower(n) == "undefined" {
			return 0
		}
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(n + err.Error())
		}
		return i
	case int64:
		return int(n)
	case float32:
		return int(n)
	case float64:
		return int(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case []byte:
		i, err := strconv.Atoi(string(n))
		if err != nil {
			panic(err.Error())
		}
		return i
	}
	return 0
}

// ToInt64 ...
func ToInt64(d interface{}) int64 {
	switch n := d.(type) {
	case int64:
		return n
	case string:
		if n == "" {
			return 0
		}
		if strings.ToLower(n) == "undefined" {
			return 0
		}
		i, err := strconv.ParseInt(strings.TrimSpace(n), 10, 64)
		if err != nil {
			panic(n + err.Error())
		}
		return i
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case []byte:
		i, err := strconv.Atoi(string(n))
		if err != nil {
			panic(err.Error())
		}
		return int64(i)
	}
	return 0
}

// ToString ...
func ToString(d interface{}, floatDiscard ...bool) string {
	switch n := d.(type) {
	case int:
		return strconv.FormatInt(int64(n), 10)
	case int8:
		return strconv.FormatInt(int64(n), 10)
	case int16:
		return strconv.FormatInt(int64(n), 10)
	case int32:
		return strconv.FormatInt(int64(n), 10)
	case int64:
		return strconv.FormatInt(n, 10)
	case uint8:
		return strconv.FormatInt(int64(n), 10)
	case uint16:
		return strconv.FormatInt(int64(n), 10)
	case uint32:
		return strconv.FormatInt(int64(n), 10)
	case uint64:
		return strconv.FormatInt(int64(n), 10)
	case string:
		_str := n
		if len(floatDiscard) > 0 && floatDiscard[0] {
			_str = strings.TrimSuffix(_str, ".00")
		}
		return _str
	case float32:
		_str := strconv.FormatFloat(float64(n), 'f', 2, 64)
		if len(floatDiscard) > 0 && floatDiscard[0] {
			_str = strings.TrimSuffix(_str, ".00")
		}
		return _str
	case float64:
		_str := strconv.FormatFloat(n, 'f', 2, 64)
		if len(floatDiscard) > 0 && floatDiscard[0] {
			_str = strings.TrimSuffix(_str, ".00")
		}
		return _str
	case []byte:
		_str := string(n)
		if len(floatDiscard) > 0 && floatDiscard[0] {
			_str = strings.TrimSuffix(_str, ".00")
		}
		return _str
	case error:
		return n.Error()
	case bool:
		if n {
			return "1"
		}
		return "0"
	}
	return "0"
}

func ToFloat64(d interface{}) float64 {
	switch n := d.(type) {
	case string:
		float, err := strconv.ParseFloat(n, 64)
		if err != nil {
			panic(n + err.Error())
		}
		return float
	case float64:
		return n
	case float32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	}
	return 0
}

// SliceIntToInt64 ...
func SliceIntToInt64(slice []int) []int64 {
	_result := []int64{}
	for _, _v := range slice {
		_result = append(_result, int64(_v))
	}
	return _result
}

// MapToSlice map 转 slice
func MapToSlice(m map[int64]int64) []int64 {
	s := make([]int64, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// SliceToMap slice 转 map
func SliceToMap(m interface{}) interface{} {
	switch _slice := m.(type) {
	case []int:
		_map := make(map[int]int, 0)
		for _, _v := range _slice {
			_map[_v] = _v
		}
		return _map
	case []int64:
		_map := make(map[int64]int64, 0)
		for _, _v := range _slice {
			_map[_v] = _v
		}
		return _map
	case []string:
		_map := make(map[string]string, 0)
		for _, _v := range _slice {
			_map[_v] = _v
		}
		return _map
	default:
		return nil
	}
}

// ForMatToStr 格式化时间(当前格式，目标格式)
func ForMatToStr(times interface{}, baseFormat ...string) string {
	_forMat := "2006-01-02 15:04:05"
	_toMat := "2006-01-02 15:04:05"
	if len(baseFormat) > 0 {
		_forMat = baseFormat[0]
	}
	if len(baseFormat) > 1 {
		_toMat = baseFormat[1]
	}
	parseStrTime, _ := time.Parse(_forMat, ToString(times))
	return parseStrTime.Format(_toMat)
}

// TimeStr 格式化时间
func TimeStr(format string) string {
	format = strings.Replace(format, "Y", "2006", 1)
	format = strings.Replace(format, "m", "01", 1)
	format = strings.Replace(format, "d", "02", 1)
	format = strings.Replace(format, "h", "15", 1)
	format = strings.Replace(format, "i", "04", 1)
	format = strings.Replace(format, "s", "05", 1)
	return time.Now().Format(format)
}

// DateStrEncode 格式化时间
func DateStrEncode(dateStr interface{}) string {
	_dateSlice := strings.Split(ToString(dateStr), "-")
	return strings.Join(_dateSlice, "")
}

// DateStrDecode 格式化时间
func DateStrDecode(dateStr interface{}) string {
	_dateSlice := strings.Split(ToString(dateStr), "")
	if len(_dateSlice) < 8 {
		return ""
	}
	_slice := []string{
		strings.Join(_dateSlice[0:4], ""),
		strings.Join(_dateSlice[4:6], ""),
		strings.Join(_dateSlice[6:8], ""),
	}
	return strings.Join(_slice, "-")
}

// UnixToStr 时间戳转字符串日期
func UnixToStr(timestamp int64, params ...string) string {
	if timestamp == 0 {
		return ""
	}
	timeNow := time.Unix(timestamp, 0) //2017-08-30 16:19:19 +0800 CST
	format := "2006-01-02 15:04:05"
	if len(params) > 0 {
		format = strings.Replace(params[0], "Y", "2006", 1)
		format = strings.Replace(format, "m", "01", 1)
		format = strings.Replace(format, "d", "02", 1)
		format = strings.Replace(format, "h", "15", 1)
		format = strings.Replace(format, "i", "04", 1)
		format = strings.Replace(format, "s", "05", 1)
	}
	return timeNow.Format(format) //2006-01-02 15:04:05
}

// 时间字符串 to 时间戳，例如：2020-05-07，2020-05-07
func DateToTimeStamp(ApplyStart, ApplyEnd string) (Start, End int64) {
	var UseFormal = "2006-01-02 15:04:05"
	var his = " 00:00:00"
	var his2 = " 23:59:59"
	if start := strings.Split(ApplyStart, "-"); len(start) == 1 {
		ApplyStart += "-01-01"
	}
	if end := strings.Split(ApplyEnd, "-"); len(end) == 1 {
		ApplyEnd += "-12-31"
	}
	loc, _ := time.LoadLocation("Local") //设置时区
	TStart, _ := time.ParseInLocation(UseFormal, ApplyStart+his, loc)
	TEnd, _ := time.ParseInLocation(UseFormal, ApplyEnd+his2, loc)
	return TStart.Unix(), TEnd.Unix()
}

// 字符串转时间戳
func StrToTimeStamp(str string, baseFormat ...string) int64 {
	_format := "2006-01-02 15:04:05"
	if len(baseFormat) > 0 {
		_format = baseFormat[0]
	}
	loc, _ := time.LoadLocation("Local")                     //重要：获取时区
	theTime, _err := time.ParseInLocation(_format, str, loc) //使用模板在对应时区转化为time.time类型
	if _err != nil {
		return 0
	}
	return theTime.Unix()
}

// GetDaysBetween2Date 两个日期相差天数
// GetDaysBetween2Date("20060102", "20220613", "20220601") //结果为12
// GetDaysBetween2Date("2006-01-02", "2022-06-13", "2022-06-01")  //结果为12
func GetDaysBetween2Date(format, date1Str, date2Str string) (int, error) {
	// 将字符串转化为Time格式
	date1, err := time.ParseInLocation(format, date1Str, time.Local)
	if err != nil {
		return 0, err
	}
	// 将字符串转化为Time格式
	date2, err := time.ParseInLocation(format, date2Str, time.Local)
	if err != nil {
		return 0, err
	}
	//计算相差天数
	return int(date1.Sub(date2).Hours() / 24), nil
}

// UID ...
func UID(preID, randID interface{}) string {
	_idSlice := []string{ToString(preID)}
	_idSlice = append(_idSlice, strings.Split(TimeStr("Ymdhis"), "")[2:]...)
	_idSlice = append(_idSlice, ToString(randID))
	return strings.Join(_idSlice, "")
}

// 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
	//endregion
}

// Md5Encode ...
func Md5Encode(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5加密
func MD5(str string) string {
	data := []byte(str)
	sum := fmt.Sprintf("%x\n", md5.Sum(data))
	return sum
}

// 判断
func InSlice(slice []int, x int) bool {
	if len(slice) < 1 {
		return false
	}
	sort.Sort(sort.IntSlice(slice))
	index := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= x
	})
	if len(slice) > index {
		if slice[index] == x {
			return true
		}
	}
	return false
}

// 除法保留位数(默认保留2位小数)。例如：1/100 = 0.01
func Float64Decimal(value float64, params ...interface{}) float64 {
	_point := "2"
	if len(params) > 0 {
		_point = ToString(params[0])
	}
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+_point+"f", value/100.0), 64)
	return value
}

// StringToFloat64...
func StringToFloat64(value interface{}, params ...string) float64 {
	_point := "2"
	if len(params) > 0 {
		_point = params[0]
	}
	_price, _ := strconv.ParseFloat(ToString(value), 64)
	_folat, _ := strconv.ParseFloat(fmt.Sprintf("%."+_point+"f", _price), 64)
	return _folat
}

// 截取字符串（支持中文）
func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// GetAge 获取年龄
func GetAge(birthday string) int64 {
	var Age int64
	var pslTime string
	if strings.Index(birthday, ".") != -1 {
		pslTime = "2006.01.02"
	} else if strings.Index(birthday, "-") != -1 {
		pslTime = "2006-01-02"
	} else {
		pslTime = "2006/01/02"
	}
	t1, err := time.ParseInLocation(pslTime, birthday, time.Local)
	if err != nil {
		return Age
	}
	diff := time.Now().Local().Unix() - t1.Local().Unix()
	if diff > 0 {
		Age = diff / (3600 * 365 * 24)
	}
	return Age
}

// SignMD5 签名算法（md5加密）
func SignMD5(key string, param interface{}) (sign string) {
	var pList = make([]string, 0, 0)
	switch _param := param.(type) {
	case url.Values:
		for _key := range _param {
			var value = _param.Get(_key)
			if len(value) > 0 {
				if _t, _err := strconv.ParseFloat(value, 10); _err == nil && strings.Index(value, ".") >= 1 {
					//兼容内部调用验签：对数值浮点型的取值进行处理
					pList = append(pList, _key+"="+ToString(_t, true))
				} else {
					pList = append(pList, _key+"="+value)
				}
			}
		}
	case Mp:
		for _key, _value := range _param {
			pList = append(pList, _key+"="+ToString(_value, true))
		}
	}
	sort.Strings(pList)
	pList = append(pList, "key="+key)
	var src = strings.Join(pList, "&")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := hex.EncodeToString(md5Ctx.Sum(nil))
	sign = strings.ToUpper(strings.Join(strings.Split(cipherStr, "")[5:20], ""))
	fmt.Println("--------------生成签名前的参数:----------" + src + ",签名值为:" + sign)
	return sign
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomInteger(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 16进制转10进制
func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

// RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandInt(length int) int {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	return r.Intn(length)
}

// 只允许英文和数字且不能为纯数字
func IsAlphaNumeric(s string) bool {
	if IsNumeric(s) {
		return false
	}
	isAlpha := regexp.MustCompile(`^[A-Za-z1-9]+$`).MatchString
	return isAlpha(s)
}

// 是否为英文
func IsEnglish(s string) bool {
	alpha := "abcdefghijklmnopqrstuvwxyz"
	for _, char := range s {
		if !strings.Contains(alpha, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

// 是否为数值类型
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// 提取中文
func UnicodeHan(str string) string {
	s := ""
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			s += string(r)
		}
	}
	return s
}

func CalculateAge(birthDay string) int64 {
	birthDateTime, _ := time.Parse("2006-01-02", birthDay)
	now := time.Now()
	//age := ToInt64(now.Year()) - ToInt64(birthDateTime.Year())
	age := now.Year() - birthDateTime.Year()
	if now.Month() < birthDateTime.Month() || (now.Month() == birthDateTime.Month() && now.Day() < birthDateTime.Day()) {
		age--
	}
	return ToInt64(age)
}

// Get information from ID card number. Birthday, age, gender
func GetIdCardNoInfo(idNumber string) (string, int64, int8) {
	birthday := idNumber[6:10] + "-" + idNumber[10:12] + "-" + idNumber[12:14]
	age := CalculateAge(birthday)
	sex := int8(1)
	genderMask, _ := strconv.Atoi(string(idNumber[16]))
	if genderMask%2 == 0 {
		sex = int8(2)
	}
	return birthday, age, sex
}

// 打印堆栈信息
func PrintStackTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}

// 过滤表情符
func FilterEmojis(value string) string {
	_regexp := regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{2600}-\x{26FF}]`)
	return _regexp.ReplaceAllString(value, "")
}

// 求交集
func IntersectSlice(slice1, slice2 []int64) []int64 {
	m := make(map[int64]int)
	nn := make([]int64, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

func GetIntersection(arr1, arr2 []int64) []int64 {
	m := make(map[int64]bool)
	var result []int64
	for _, num := range arr1 {
		m[num] = true
	}
	for _, num := range arr2 {
		if m[num] {
			result = append(result, num)
		}
	}
	return result
}

// 字符串数组是否存在
func IsStringInArray(str string, arr []string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

// CheckIdCard 检验身份证
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}
