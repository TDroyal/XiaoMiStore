package logic

import "strconv"

//string转为int类型
func StringToInt(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

//int转为string类型
func IntToString(n int) string {
	str := strconv.Itoa(n)
	return str
}
