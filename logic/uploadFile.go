package logic

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/models"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

// 这些都应该配置在app.ini中
const (
	OssStatus     = 1
	BackendDomain = "127.0.0.1/" //后端域名
)

// 通过列获取系统设置里面的值  columnName就是结构体的属性名称，例如：OssStatus
func GetSettingFromColumn(columnName string) (string, error) {
	setting := models.Setting{ID: 1}
	if err := dao.DB.First(&setting).Error; err != nil {
		return "", err
	}
	// 反射来获取
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val, nil
}

func GetOssStatus() int {
	return OssStatus
}

func FormatImg(str string) (string, error) { //返回给前端每个图片地址的时候，就应该拼接上前缀域名或者ip
	domain := BackendDomain
	if GetOssStatus() == 1 {
		domain, err := GetSettingFromColumn("OssDomain")
		if err != nil {
			return "", err
		}
		return domain + str, nil
	}
	return domain + str, nil
}

// 上传图片
func UploadImageFile(c *gin.Context, uploadfileName string) (string, error) { //返回文件保存的目录，第二个参数返回error
	// 1.判断是否开启了oss（对象存储服务）
	// 从后端配置app.ini文件中读出是否启用oss服务，不建议设置在数据库中
	OssStatus := GetOssStatus()

	if OssStatus == 1 {
		return OssUploadImageFile(c, uploadfileName)
	}
	return LocalUploadImageFile(c, uploadfileName)
}

// 上传图片到本地服务器
func LocalUploadImageFile(c *gin.Context, uploadfileName string) (string, error) { //返回文件保存的目录，第二个参数返回error
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
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}

// 上传图片到Oss服务器
func OssUploadImageFile(c *gin.Context, uploadfileName string) (string, error) { //返回文件保存的目录，第二个参数返回error
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

	//4.对图片进行压缩（选做）

	//5.生成文件名称和文件保存的目录      xxx.jpeg
	fileName := GetUnixNanoTimestamp() + extName // time.Now().Unix() 将时间转换为 UNIX 时间戳（纳秒级别）
	//6.执行保存 (将文件上传到oss云服务器的dst)
	dst := fmt.Sprintf("%s/%s", dir, fileName)
	return OssUpload(file, dst)
}

// 封装Oss上传的方法
// 官方文档1  https://github.com/aliyun/aliyun-oss-go-sdk
// 官方文档2  https://help.aliyun.com/zh/oss/developer-reference/simple-upload-4

func OssUpload(file *multipart.FileHeader, dst string) (string, error) {

	// 从setting数据库表中读出oss相关的配置
	setting := models.Setting{ID: 1}
	if err := dao.DB.First(&setting).Error; err != nil {
		return "", err
	}

	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	os.Setenv("OSS_ACCESS_KEY_ID", setting.Appid)
	os.Setenv("OSS_ACCESS_KEY_SECRET", setting.AppSecret)

	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		return "", err
	}

	// 创建OSSClient实例。
	client, err := oss.New(setting.EndPoint, "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		return "", err
	}

	// 填写存储空间名称
	bucket, err := client.Bucket(setting.BucketName)
	if err != nil {
		return "", err
	}

	// 读取传前端过来文件
	src, err := file.Open() //参考c.SaveUploadedFile()源码
	if err != nil {
		return "", err
	}
	defer src.Close()
	// 将文件流上传至oss的dst
	err = bucket.PutObject(dst, src) // multipart.File类型和io.Reader类型一致
	if err != nil {
		return "", err
	}
	return dst, nil
}
