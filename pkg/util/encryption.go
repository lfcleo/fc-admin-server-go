package util

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// AesCBCPkcs7Encrypt aes加密，
func AesCBCPkcs7Encrypt(data []byte, keyStr, iv string) ([]byte, error) {
	// 将密钥和 IV 转换为字节切片
	key := []byte(keyStr)
	ivBytes := []byte(iv)
	// 创建一个新的 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, ivBytes)   // 创建一个新的 CBC 加密器
	paddedPlaintext := pkcs7Padding(data)            // 对明文进行 PKCS#7 填充
	ciphertext := make([]byte, len(paddedPlaintext)) // 创建一个足够大的缓冲区以存储加密后的密文
	mode.CryptBlocks(ciphertext, paddedPlaintext)    // 加密数据
	return ciphertext, nil
}

// pkcs7Padding 使用PKCS7进行填充
func pkcs7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// AesCBCPkcs7Decrypt AES解密,CBC,PKCS7
func AesCBCPkcs7Decrypt(data, key, ivStr string) ([]byte, error) {
	iv := []byte(ivStr)
	if len(iv) != 16 {
		return nil, errors.New("偏移量错误")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	decodeData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	decryptData := make([]byte, len(decodeData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptData, decodeData)
	original, err := pkcs7UnPadding(decryptData)
	return original, err
}

// pkcs7UnPadding 移除 pkcs7 填充
func pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, fmt.Errorf("src length is 0")
	}
	unPadding := int(src[length-1])
	if length < unPadding {
		return src, fmt.Errorf("src length is less than unpadding")
	}
	return src[:(length - unPadding)], nil
}

// EncodeSha1 /*
func EncodeSha1(value string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(value))

	return hex.EncodeToString(sha1.Sum(nil))
}

// EncodeMD5 /*
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// RsaWithSHA256Base64 SHA256withRSA加密 （支付宝授权登录拼接字符串，要求在这里https://opendocs.alipay.com/open/291/106118）
func RsaWithSHA256Base64(signContent string, privateKey string, hash crypto.Hash) (string, error) {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(nil, priKey, hash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// ParsePrivateKey SHA256withRSA加密使用
func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = FormatPrivateKey(privateKey)
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

// FormatPrivateKey SHA256withRSA加密使用
func FormatPrivateKey(privateKey string) string {
	/*
		公钥开头	"-----BEGIN PUBLIC KEY-----"
		公钥结尾	"-----END PUBLIC KEY-----"
		私钥开头	"-----BEGIN RSA PRIVATE KEY-----"
		私钥结尾	"-----END RSA PRIVATE KEY-----"
		别忘了还有换行
	*/
	PemBegin := "-----BEGIN RSA PRIVATE KEY-----\n"
	PemEnd := "\n-----END RSA PRIVATE KEY-----"
	if !strings.HasPrefix(privateKey, PemBegin) {
		privateKey = PemBegin + privateKey
	}
	if !strings.HasSuffix(privateKey, PemEnd) {
		privateKey = privateKey + PemEnd
	}
	return privateKey
}

// CheckSignature 微信公众号签名检查
func CheckSignature(signature, timestamp, nonce, token string) bool {
	arr := []string{timestamp, nonce, token}
	// 字典序排序
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return EncodeSha1(b.String()) == signature
}

// map根据key排序
func sortMapForKeyAnd2String(mp map[string]string) string {
	//1.将map的key放到切片中
	var keyMap = make([]string, 0)
	for k, _ := range mp {
		keyMap = append(keyMap, k)
	}
	//2.对切片排序
	sort.Strings(keyMap)
	//3.遍历切片，然后按key来添加map的值,并且value值做urlencode处理
	var resultStr string
	for _, v := range keyMap {
		resultStr = fmt.Sprintf("%v%v=%v&", resultStr, v, mp[v])
		//resultStr = fmt.Sprintf("%v%v=%v&", resultStr, v, url.QueryEscape(mp[v]))
	}
	return strings.TrimRight(resultStr, "&") //去掉最后一个&字符
}

// map根据key排序
func sortMapForKey(mp map[string]string) map[string]string {
	//1.将map的key放到切片中
	var keyMap = make([]string, 0)
	for k, _ := range mp {
		keyMap = append(keyMap, k)
	}
	//2.对切片排序
	sort.Strings(keyMap)
	//3.遍历切片，然后按key来添加map的值
	var newMap = make(map[string]string)
	for _, v := range keyMap {
		newMap[v] = mp[v]
	}
	return newMap
}

// RemoveDomain 判断字符中是否有域名，去掉域名返回路径（图片等保存在数据库中）
func RemoveDomain(url string) string {
	path := url
	//判断是否以http开头
	if strings.HasPrefix(path, "http") == true {
		a1 := strings.Split(path, "//")[1]         //去除http://，得到剩下的字符串
		a2 := strings.Split(a1, "/")[0]            //获取域名，例如 fanwu-app.oss-cn-hangzhou.aliyuncs.com
		path = strings.Replace(a1, a2+"/", "", -1) //得到去除域名下的文件路径，例如 images/product_images/20201121/1_1605940532_a98832504357e6eef3a4afb862b0c9b5.jpg
	}
	return path
}
