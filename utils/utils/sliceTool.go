package utils

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

type SliceInt []int

func NewSliceInt(initValue ...int) *SliceInt {
	_slice := new(SliceInt)
	*_slice = append(*_slice, initValue...)
	return _slice
}

// Append ...
func (that *SliceInt) Append(value int, isMerge ...bool) *SliceInt {
	if len(isMerge) > 0 && isMerge[0] == true {
		if !that.InSlice(value) {
			*that = append(*that, int(value))
		}
	} else {
		*that = append(*that, int(value))
	}
	return that
}

// Append ...
func (that *SliceInt) AppendInt64(value int64, isMerge ...bool) *SliceInt {
	_value := int(value)
	*that = *that.Append(_value, isMerge...)
	return that
}

// ToString ...
func (that *SliceInt) ToString() string {
	return strings.Replace(strings.Trim(fmt.Sprint(*that), "[]"), " ", ",", -1)
}

// 判断
func (that *SliceInt) InSlice(value int) bool {
	if len(*that) < 1 {
		return false
	}
	sort.Sort(sort.IntSlice(*that))
	index := sort.Search(len(*that), func(i int) bool {
		return (*that)[i] >= value
	})
	if len(*that) > index {
		if (*that)[index] == value {
			return true
		}
	}
	return false
}

// ToInterface ...
func (that *SliceInt) ToInterface() []interface{} {
	_result := []interface{}{}
	for _, _v := range *that {
		_result = append(_result, _v)
	}
	return _result
}

// ToInt64 ...
func (that *SliceInt) ToInt64() []int64 {
	_result := []int64{}
	for _, _v := range *that {
		_result = append(_result, int64(_v))
	}
	return _result
}

func (that *SliceInt) Reset() *SliceInt {
	*that = SliceInt{}
	return that
}

func (that *SliceInt) Join(sep string) string {
	_slice := []string{}
	for _, _v := range *that {
		_slice = append(_slice, ToString(_v))
	}
	return strings.Join(_slice, sep)
}

func (that *SliceInt) Length() int {
	return len(*that)
}

////////////////////////////////////////////////////////////////

type SliceInt64 []int64

func NewSliceInt64(initValue ...int64) *SliceInt64 {
	_slice := new(SliceInt64)
	*_slice = append(*_slice, initValue...)
	return _slice
}

// ToInterface ...
func (that *SliceInt64) ToInterface() []interface{} {
	_result := []interface{}{}
	for _, _v := range *that {
		_result = append(_result, _v)
	}
	return _result
}

// ToInt64 ...
func (that *SliceInt64) ToInt() []int {
	_result := []int{}
	for _, _v := range *that {
		_result = append(_result, int(_v))
	}
	return _result
}

func (that *SliceInt64) Join(sep string) string {
	_slice := []string{}
	for _, _v := range *that {
		_slice = append(_slice, ToString(_v))
	}
	return strings.Join(_slice, sep)
}

////////////////////////////////////////////////////////////////

type SliceString []string

func NewSliceString(initValue ...string) *SliceString {
	_slice := new(SliceString)
	*_slice = append(*_slice, initValue...)
	return _slice
}

// ToInterface ...
func (that *SliceString) ToInterface() []interface{} {
	_result := []interface{}{}
	for _, _v := range *that {
		_result = append(_result, _v)
	}
	return _result
}

// ToInt64 ...
func (that *SliceString) ToInt64() []int64 {
	_result := []int64{}
	for _, _v := range *that {
		_result = append(_result, ToInt64(_v))
	}
	return _result
}

func (that *SliceString) ToString() []string {
	_result := []string{}
	for _, _v := range *that {
		_result = append(_result, _v)
	}
	return _result
}

func (that *SliceString) ToStr() string {
	_result := strings.Replace(strings.Trim(fmt.Sprint(*that), "[]"), " ", ",", -1)
	return _result
}

func (that *SliceString) In(target string, strArray ...[]string) bool {
	_strArray := make([]string, 0)
	if len(strArray) > 0 {
		_strArray = strArray[0]
	} else {
		_strArray = that.ToString()
	}
	for _, element := range _strArray {
		if target == element {
			return true
		}
	}
	return false
}

// Append ...
func (that *SliceString) Append(value string, isMerge ...bool) *SliceString {
	if len(isMerge) > 0 && isMerge[0] == true {
		if !that.In(value, that.ToString()) {
			*that = append(*that, value)
		}
	} else {
		*that = append(*that, value)
	}
	return that
}

func (that *SliceString) Length() int {
	return len(*that)
}

// 判断slice中是否存在某个item
func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Int()) {
				return true
			}
		}
	}
	return false
}

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// 数组去重
func RemoveIntAndEmpty(list []int64) []int64 {
	x := make([]int64, 0)
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

// 求交集
func IntersectInt(slice1, slice2 []int64) []int64 {
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

////////////////////////////////////////////////////////////////

type SliceFunc []func()

func NewSliceFunc(funcBack ...func()) *SliceFunc {
	_slice := new(SliceFunc)
	*_slice = append(*_slice, funcBack...)
	return _slice
}

// Append ...
func (that *SliceFunc) Append(funcBack ...func()) *SliceFunc {
	*that = append(*that, funcBack...)
	return that
}

// Do ...
func (that *SliceFunc) Do() {
	for _, _f := range *that {
		_f()
	}
}

////////////////////////////////////////////////////////////////

// RandSlice 切片乱序
func RandSlice(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return
	}

	length := rv.Len()
	if length < 2 {
		return
	}

	swap := reflect.Swapper(slice)
	rand.Seed(time.Now().Unix())
	for i := length - 1; i >= 0; i-- {
		j := rand.Intn(length)
		swap(i, j)
	}
	return
}

////////////////////////////////////////////////////////////////

// @Summary 切片分页
// @Param page 当前页
// @Param pageSize 每页显示数量
// @Param nums 数据总数
// @return sliceStart 切片开始
// @return sliceEnd 切片结尾
func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 10
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
