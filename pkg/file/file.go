package file

import (
	"io"
	"mime/multipart" //它主要实现了 MIME 的 multipart 解析，主要适用于 HTTP 和常见浏览器生成的 multipart 主体
	"os"
	"path"
	"strings"
)

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)

	return len(content), err
}

// GetExt 获取文件后缀
func GetExt(filename string) string {
	return path.Ext(filename)
}

// CheckExist 检查文件是否存在
/*
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// MKDir 新建文件夹
func MKDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	return err
}

// IsNotExistMkDir 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if exist := CheckExist(src); exist == false {
		if err := MKDir(src); err != nil {
			return err
		}
	}
	return nil
}

/*
Open
调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。如果出现错误，则为*PathError
const (

	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
	O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
	O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
	O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
	O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
	O_SYNC   int = syscall.O_SYNC   // 同步IO
	O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时

)
*/
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, err
}

// GetFolderFilesName 获取文件夹下的文件名称
func GetFolderFilesName(foldPath string) (fs []string, err error) {
	// Given path is a directory.
	dir, err := os.Open(foldPath)
	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if strings.Contains(v.Name(), ".json") {
			fileName := strings.TrimSuffix(v.Name(), ".json")
			fs = append(fs, fileName)
		}
	}
	return fs, err
}
