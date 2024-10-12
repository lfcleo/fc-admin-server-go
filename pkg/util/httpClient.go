package util

import (
	"io"
	"net/http"
)

// MapWorldGetCode 天地图逆地理编码查询
func MapWorldGetCode(body string) ([]byte, error) {
	url := "http://api.tianditu.gov.cn/geocoder?" + body
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
