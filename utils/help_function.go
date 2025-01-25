package utils

// DateDiff 计算两个日期之前相隔的天数
import (
	"errors"
	"time"
)

func DateDiff(startDate, endDate string) (int64, error) {
	const layout = "2006-01-02" // 假设日期格式为 "YYYY-MM-DD"

	// 解析 startDate
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return 0, errors.New("invalid start date format")
	}

	// 解析 endDate
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return 0, errors.New("invalid end date format")
	}

	// 计算日期差
	diff := end.Sub(start).Hours() / 24

	// 处理边界条件
	if diff < 0 {
		return 0, errors.New("start date is after end date")
	}

	return int64(diff), nil
}
