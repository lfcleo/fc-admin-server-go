package util

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"path"
	"regexp"
)

// ValidateMobileNumber 正则判断手机号
func ValidateMobileNumber(mobile string) bool {
	// 定义手机号正则表达式
	phonePattern := `^1[34578]\d{9}$`
	re := regexp.MustCompile(phonePattern)
	return re.MatchString(mobile)
}

// ValidatePassword 正则判断密码(包含大小写字母和数字，长度最少6位)
func ValidatePassword(password string) bool {

	// 定义检查各个条件的正则表达式  Go语言的regexp包中是不支持Perl风格的正向先行断言(?=...),所以用以下
	hasLowerCase := regexp.MustCompile(`[a-z]`) //至少包含1个小写字母
	hasUpperCase := regexp.MustCompile(`[A-Z]`) //至少包含1个大写字母
	hasDigit := regexp.MustCompile(`\d`)        //至少包含1个数字
	//hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()\-_=+{};:,<.>]`)	//至少包含1个特殊字符
	minLength := 6

	if len(password) < minLength {
		return false
	}
	if !hasLowerCase.MatchString(password) {
		return false
	}
	if !hasUpperCase.MatchString(password) {
		return false
	}
	if !hasDigit.MatchString(password) {
		return false
	}
	//if !hasSpecialChar.MatchString(password) {
	//	return false
	//}
	return true
}

// ValidateEmail 正则判断邮箱
func ValidateEmail(mobile string) bool {
	// 定义手机号正则表达式
	phonePattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(phonePattern)
	return re.MatchString(mobile)
}

// ValidateArticleImage 正则取出文章(markdown)中第一张图片路径
func ValidateArticleImage(mdContent string) string {
	// 正则表达式匹配Markdown图片语法
	regex := regexp.MustCompile(`!\[.*\]\((.*?)\)`)
	matches := regex.FindAllStringSubmatch(mdContent, -1)

	if len(matches) > 0 {
		firstImagePath := matches[0][1]
		return firstImagePath
	} else {
		return ""
	}
}

// ExtractFilePath url去除域名及端口，保留文件路径
func ExtractFilePath(rawURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// 提取文件路径
	filePath := path.Join(parsedURL.Path)
	return filePath[1:], nil //去除了首字符/
}

// GenCaptchaCode 随机生成6位数验证码
func GenCaptchaCode() (string, error) {
	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		codes[i] = 48 + (codes[i] % 10)
	}
	return string(codes), nil
}

// 字符集，包含大小写字母和数字
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenRandomString 生成随机字符串（用户昵称）
func GenRandomString(length int) (string, error) {
	b := make([]byte, length)
	// 生成随机数
	for i := range b {
		// 从字符集中选取一个字符
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[index.Int64()]
	}

	return string(b), nil
}

// UintArraysEqual 判断两个uint数组是否相等，包括顺序
func UintArraysEqual(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// UintContains 判断数字数组是否含有某个值
func UintContains(arr []uint, target uint) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

//// GetRandomString 生成随机字符串
//func GetRandomString(l int) string {
//	str := "0123456789abcdefghijklmnopqrstuvwxyz"
//	bytes := []byte(str)
//	var result []byte
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	for i := 0; i < l; i++ {
//		result = append(result, bytes[r.Intn(len(bytes))])
//	}
//	return string(result)
//}
//
//// GetRandomIntString 随机生成8位数字字符串,订单号使用
//func GetRandomIntString(t int64) string {
//	rnd := rand.New(rand.NewSource(t))
//	vcode := fmt.Sprintf("%08v", rnd.Int31n(100000000))
//	return vcode
//}
