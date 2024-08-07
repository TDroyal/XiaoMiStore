package logic

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "github.com/hunterhug/go_image" //加个点，就可以直接调用包里面的方法函数，不需要go_iamge.func()
)

// string转为int类型
func StringToInt(str string) int {
	n, _ := strconv.Atoi(str) //如果出错，n=0
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

// 生成商品缩略图
func ResizeGoodsImage(filename string) {
	// 后缀名
	extname := filepath.Ext(filename)

	//读取数据库中的需要生成的缩略图尺寸
	thumbnailSize, _ := GetSettingFromColumn("ThumbnailSize") //100,200,500
	thumbnailSizeSlice := strings.Split(thumbnailSize, ",")

	// 以前的图片是filename = "static/upload/xxx.png"
	// 生成缩略图的名称为"static/upload/xxx.png_100x100.png"
	for i := 0; i < len(thumbnailSizeSlice); i++ {
		savePath := filename + "_" + thumbnailSizeSlice[i] + "x" + thumbnailSizeSlice[i] + extname
		width := StringToInt(thumbnailSizeSlice[i])
		if err := ThumbnailF2F(filename, savePath, width, width); err != nil { //https://github.com/hunterhug/go_image
			fmt.Println(err) //回头写个日志模块，处理日志
		}
	}

}

// oss不需要生成缩略图，直接在访问图片的地址后面加几个参数即可。
// 图片的oss原地址 http://xxx/111.png
// 缩略图地址  http://xxx/111.png?x-oss-process=image/resize,h_200
