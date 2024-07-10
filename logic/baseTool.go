package logic

import (
	"strconv"
	"time"
)

// string转为int类型
func StringToInt(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

// string转为float64类型
func StringToFloat(str string) float64 {
	n, _ := strconv.ParseFloat(str, 64)
	return n
}

// int转为string类型
func IntToString(n int) string {
	str := strconv.Itoa(n)
	return str
}

// 获取当前的日期 年月日
func GetDate() string {
	time := time.Now()
	return IntToString(time.Year()) + IntToString(int(time.Month())) + IntToString(time.Day())
}

// 获得当前的Unix时间戳(毫秒)
func GetUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// 获得当前的Unix时间戳(纳秒)
func GetUnixNanoTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
