package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
 *发送get请求   Author:tang
 *@param url api地址
 *@param string 返回的数据string， 返回的数据byte
 */
func GetHttpContent(url string) (string, []byte) {
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("User-Agent", "chrome 100")
	response, _ := client.Do(reqest)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		fmt.Println("---请求后的数据---", bodystr)
		return bodystr, body
	}
	return "", nil
}

// CheckUrl 检测URL前缀
func CheckUrl(str string) string {
	a := strings.HasPrefix(str, "http://")
	b := strings.HasPrefix(str, "https://")
	if a || b {
		return str
	}
	return "http://" + str
}
