package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ExcelMp = map[int]interface{}

type Excel struct {
	suffix      string    //文件名后缀
	outDir      string    //输出目录
	name        string    //文件名
	columnTitle []string  //列标题
	data        []ExcelMp //数据源
}

// new..
func NewExcel() *Excel {
	return &Excel{}
}

// SetName...
func (that *Excel) SetName(name string) *Excel {
	that.name = name
	return that
}

// SetDir...
func (that *Excel) SetDir(outDir string) *Excel {
	that.outDir = outDir
	return that
}

// SetColumnTitle...
func (that *Excel) SetColumnTitle(columnTitle interface{}) *Excel {
	switch _t := columnTitle.(type) {
	case []string:
		that.columnTitle = _t
	default:
		that.columnTitle = nil
	}
	return that
}

// SetData..
func (that *Excel) SetData(data interface{}) *Excel {
	switch _t := data.(type) {
	case []ExcelMp:
		that.data = _t
	case []interface{}:
		for _, v1 := range _t {
			that.data = append(that.data, v1.(ExcelMp))
		}
	default:
		that.data = nil
	}
	return that
}

// GetDir
func (that *Excel) GetDownLoadUrl() string {
	return that.outDir[1:] + that.name + that.suffix
}

// init..
func (that *Excel) init() {
	if that.name == "" {
		that.name = "demo"
	}
	if that.outDir == "" {
		that.outDir = "./upload/excel/"
	}
	_ = DirCreate(that.outDir)
	if len(that.columnTitle) < 1 {
		that.columnTitle = []string{
			"姓名", "年龄",
		}
	}
	if len(that.data) < 1 {
		that.data = []ExcelMp{
			{0: "jack", 1: 18},
			{0: "mary", 1: 28},
		}
	}
	that.suffix = ".xlsx"
}

// Create...
func (that *Excel) Create() {
	that.init()
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	//
	for clumnNum, v := range that.columnTitle {
		sheetPosition := that.Div(clumnNum+1) + "1"
		//fmt.Print(sheetPosition)
		f.SetCellValue("Sheet1", sheetPosition, v)
	}
	for lineNum, v := range that.data {
		//Set Orderly Reorder the elements..
		var orderly []int
		for k := range v {
			orderly = append(orderly, ToInt(k))
		}
		sort.Ints(orderly)
		// Set value of a cell.
		clumnNum := 0
		for _kk, _ := range orderly {
			clumnNum++
			sheetPosition := that.Div(clumnNum) + strconv.Itoa(lineNum+2)
			switch v[_kk].(type) {
			case string:
				f.SetCellValue("Sheet1", sheetPosition, v[_kk].(string))
				break
			case int:
				f.SetCellValue("Sheet1", sheetPosition, v[_kk].(int))
				break
			case float64:
				f.SetCellValue("Sheet1", sheetPosition, v[_kk].(float64))
				break
			default:
				f.SetCellValue("Sheet1", sheetPosition, v[_kk])
			}
		}
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	that.name = strings.Replace(that.name, ".", "", -1)
	// Save xlsx file by the given path.
	if err := f.SaveAs(that.outDir + that.name + that.suffix); err != nil {
		println(err.Error())
	}
}

// Div 数字转字母
func (that *Excel) Div(Num int) string {
	var (
		Str  string = ""
		k    int
		temp []int //保存转化后每一位数据的值，然后通过索引的方式匹配A-Z
	)
	//用来匹配的字符A-Z
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if Num > 26 { //数据大于26需要进行拆分
		for {
			k = Num % 26 //从个位开始拆分，如果求余为0，说明末尾为26，也就是Z，如果是转化为26进制数，则末尾是可以为0的，这里必须为A-Z中的一个
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			Num = (Num - k) / 26 //减去Num最后一位数的值，因为已经记录在temp中
			if Num <= 26 {       //小于等于26直接进行匹配，不需要进行数据拆分
				temp = append(temp, Num)
				break
			}
		}
	} else {
		return Slice[Num]
	}
	for _, value := range temp {
		Str = Slice[value] + Str //因为数据切分后存储顺序是反的，所以Str要放在后面
	}
	return Str
}

// DirExits 目录是否存在
func DirExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DirCreate 创建目录
func DirCreate(_dir string) error {
	exist, err := DirExits(_dir)
	if err != nil {
		return errors.New("get dir error! " + err.Error())
	}
	if !exist {
		// 创建文件夹
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			return errors.New("mkdir failed! " + err.Error())
		}
	}
	return nil
}

/*
读取在线csv文件

	说明：
	 1、读取csv文件返回的内容为切片类型，可以通过遍历的方式使用或Slicer[0]方式获取具体的值。
	 2、同一个函数或线程内，两次调用Read()方法时，第二次调用时得到的值为每二行数据，依此类推。
	 3、大文件时使用逐行读取，小文件直接读取所有然后遍历，两者应用场景不一样，需要注意。
*/
func ReadCsv(filepath string) ([][]string, error) {
	resp, e := http.Get(filepath)
	if e != nil {
		return nil, e
	}
	if resp == nil {
		return nil, errors.New("resp nil")
	}
	if resp.Body == nil {
		return nil, errors.New("body nil")
	}
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	r := csv.NewReader(bytes.NewReader(buf))
	records, e := r.ReadAll()
	if e != nil {
		return nil, e
	}
	return records, nil
}
