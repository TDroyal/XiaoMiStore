package logic

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
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

// 获取当前的日期 年月日  格式：YYYYMMDD
func GetDate() string {
	time := time.Now()

	m := int(time.Month())
	d := time.Day()

	parseFunc := func(n int) string {
		if n < 10 {
			return "0" + IntToString(n)
		}
		return IntToString(n)
	}

	return IntToString(time.Year()) + parseFunc(m) + parseFunc(d)
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

// 将markdown文本转为html文本

/*
如果str是

### 嘻嘻
**cnm**

那么解析出来的html文本是

<h3>嘻嘻</h3>
<strong>cnm</strong>


*/

func FormatAttr(str string) string {

	tempSlice := strings.Split(str, "\n")

	var htmlStr string
	for _, v := range tempSlice {
		md := []byte(v)
		htmlStr += string(markdown.ToHTML(md, nil, nil))
	}

	return htmlStr
}

// 生成一串随机数
func GetRandomNum(n int) string { // 传入生成随机数的长度
	var randomstr string
	for i := 0; i < n; i++ {
		randomstr += IntToString(rand.Intn(10))
	}
	return randomstr
}

// 获取订单ID
func GetOrderID() string {
	order_id := GetDate() + GetRandomNum(12)
	return order_id
}
