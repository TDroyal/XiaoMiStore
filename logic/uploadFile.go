package logic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 上传图片
func UploadImageFile(c *gin.Context, uploadfileName string) (string, error) { //返回文件保存的目录，第二个参数返回error
	//1.获取上传的文件
	file, upload_err := c.FormFile(uploadfileName)
	if upload_err != nil { //获取上传的文件失败
		fmt.Println(upload_err) //http: no such file
		return "", nil
	}
	//2.获取文件后缀，判断类型是否正确  .jpg .png  .gif  .jpeg
	extName := filepath.Ext(file.Filename)

	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
		".gif":  true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("上传文件类型不合法")
	}

	//3.创建图片保存目录  static/focus/20240628
	date := GetDate()
	// "./static/upload/" + date   //  其中./是相对于根目录而言
	// 在路由中我们配置了静态文件服务r.Static("/static", "./static")，为了更好的存储图片路径，采用static/upload，而不是./static/
	dir := fmt.Sprintf("static/upload/%s", date)
	if err := os.MkdirAll(dir, 0666); err != nil {
		return "", err
	}

	//4.对图片进行压缩（选做）

	//5.生成文件名称和文件保存的目录      xxx.jpeg
	fileName := GetUnixNanoTimestamp() + extName // time.Now().Unix() 将时间转换为 UNIX 时间戳（纳秒级别）
	//6.执行保存
	dst := fmt.Sprintf("%s/%s", dir, fileName)
	c.SaveUploadedFile(file, dst)
	return dst, nil
}
