package utils

import (
	"fmt"
	"strings"
)

var (
	B       int64  = 1
	KB      int64  = 1024
	MB      int64  = 1024 * 1024
	GB      int64  = 1024 * 1024 * 1024
	TB      int64  = 1024 * 1024 * 1024 * 1024
	formatF string = "%5.2f"
)

func SizeFormat(sizeInBytes int64) string {
	switch {
	case B <= sizeInBytes && sizeInBytes < KB:
		return fmt.Sprintf("%dB", sizeInBytes)
	case KB <= sizeInBytes && sizeInBytes < MB:
		return fmt.Sprintf(formatF+"KB", float64(sizeInBytes)/float64(KB))
	case MB <= sizeInBytes && sizeInBytes < GB:
		return fmt.Sprintf(formatF+"MB", float64(sizeInBytes)/float64(MB))
	case GB <= sizeInBytes && sizeInBytes < TB:
		return fmt.Sprintf(formatF+"GB", float64(sizeInBytes)/float64(GB))
	case TB <= sizeInBytes:
		return fmt.Sprintf(formatF+"TB", float64(sizeInBytes)/float64(TB))
	default:
		return "0"
	}
}

// 场馆营业的时间段处理
func GetStartTime(time string) interface{} {
	if strings.Contains(time, ":") || time == "" {
		return time
	}
	t := []string{
		"00:00", "00:30", "01:00", "01:30", "02:00", "02:30", "03:00", "03:30", "04:00", "04:30", "05:00", "05:30", //0-11
		"06:00", "06:30", "07:00", "07:30", "08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30", //12-23
		"12:00", "12:30", "13:00", "13:30", "14:00", "14:30", "15:00", "15:30", "16:00", "16:30", "17:00", "17:30", //24-35
		"18:00", "18:30", "19:00", "19:30", "20:00", "20:30", "21:00", "21:30", "22:00", "22:30", "23:00", "23:30", //36-47
	}
	return t[ToInt(time)]
}
func GetEndTime(time string) interface{} {
	if strings.Contains(time, ":") || time == "" {
		return time
	}
	t := []string{
		"00:30", "01:00", "01:30", "02:00", "02:30", "03:00", "03:30", "04:00", "04:30", "05:00", "05:30", "06:00", //0-11
		"06:30", "07:00", "07:30", "08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30", "12:00", //12-23
		"12:30", "13:00", "13:30", "14:00", "14:30", "15:00", "15:30", "16:00", "16:30", "17:00", "17:30", "18:00", //24-35
		"18:30", "19:00", "19:30", "20:00", "20:30", "21:00", "21:30", "22:00", "22:30", "23:00", "23:30", "24:00", //36-47
	}
	return t[ToInt(time)]
}
