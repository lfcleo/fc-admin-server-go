package upload

import (
	"bytes"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/pkg/file"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// GetImagePath 获取图片的保存路径，就是配置文件设置的   documents/images/
func GetImagePath() string {
	return config.Data.Server.ImageSavePath
}

// GetImageDateName 图片的日期文件夹      20190730/
func GetImageDateName() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d/", t.Year(), t.Month(), t.Day())
}

// GetImageFullUrl 获取图片完整访问URL       http://127.0.0.1:8000/documents/images/20190730/******.png
func GetImageFullUrl(name string) string {
	return config.Data.Server.DomainName + GetImagePath() + GetImageDateName() + name
}

// SetImageName 为图片设置新名称，名称格式为   "md5(当前时间戳+用户uuid).图片格式"
func SetImageName(uid, name string) string {
	ext := path.Ext(name)      //ext返回路径使用的文件扩展名
	timeN := time.Now().Unix() //时间戳
	timeStr := strconv.FormatInt(timeN, 10)
	//生成随机的文件名
	var bb bytes.Buffer
	bb.WriteString(timeStr)
	bb.WriteString(uid)
	//bb.WriteString(util.GetRandomString(9))
	newFileName := fmt.Sprintf("%s_%s", timeStr, util.EncodeMD5(bb.String()))

	return newFileName + ext
}

// GetImageFullPath 获取图片在项目中的目录	runtime/documents/images/(头像/商品图片的路径)
func GetImageFullPath(path string) string {
	return config.Data.Server.RuntimeRootPath + GetImagePath() + path + "/"
}

// CheckImageExt 检查图片后缀,是否是属于配置中允许的后缀名
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allExt := range config.Data.Server.ImageAllowExts {
		if strings.ToLower(allExt) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

// CheckImageSizeByNum 检查图片大小
func CheckImageSizeByNum(size int64) bool {
	// 1M = 1024KB = 1048576B = 1024 * 1024B
	return size <= int64(config.Data.Server.ImageMaxSize*1024*1024)
}

// CheckImage 检查图片
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src) //如果不存在则新建文件夹
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src) //检查文件权限
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

// httpContentType 类型
var httpContentType = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// CheckHttpTypeImage 检查http上传的文件类型
func CheckHttpTypeImage(httType string) bool {
	if !httpContentType[httType] {
		return false
	}
	return true
}
