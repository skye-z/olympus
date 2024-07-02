package util

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// 获取远程数据
func GetData(remote, path string, check bool) []byte {
	resp := GetResp(remote, path, nil, check)
	if resp == nil {
		return nil
	}

	defer resp.Body.Close()

	if check && resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Req] error:", err)
		return nil
	}

	return body
}

// 获取远程响应
func GetResp(remote, path string, header http.Header, check bool) *http.Response {
	var client *http.Client
	proxy := GetString("basic.proxy")
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Println("[Req] error:", err)
			return nil
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: transport,
		}
		log.Println("[Req] proxy get " + remote + path)
	} else {
		client = &http.Client{}
		log.Println("[Req] direct get " + remote + path)
	}

	req, err := http.NewRequest(http.MethodGet, remote+path, nil)
	if err != nil {
		log.Println("[Req] error:", err)
		return nil
	}
	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[Req] error:", err)
		return nil
	}
	return resp
}
