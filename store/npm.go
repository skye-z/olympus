package store

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/skye-z/olympus/util"
)

const (
	npmRepository = "./repository/npm/"
)

type NpmStore struct {
	RemoteURL string
}

// 获取文件
func (ns NpmStore) GetFile(path string) []byte {
	fileName := ""
	directory := ""
	index := strings.LastIndex(path, "/")
	if index == -1 {
		fileName = path + ".json"
		directory = path
	} else {
		fileName = path[strings.LastIndex(path, "/"):]
		directory = strings.ReplaceAll(path[:strings.LastIndex(path, "/")], "/-", "")
	}

	if util.CheckExist(npmRepository + strings.ReplaceAll(path, "/-", "")) {
		return util.ReadFile(npmRepository + path)
	} else {
		content := ns.getRemoteData(path)
		if index == -1 {
			result := bytes.Replace(content, []byte("https://registry.npmjs.org/"), []byte("http://localhost:27680/npm/"), -1)
			util.SaveFile(npmRepository+directory, fileName, result)
			return result
		} else {
			util.SaveFile(npmRepository+directory, fileName, content)
			return content
		}
	}
}

// 获取远程数据
func (ns NpmStore) getRemoteData(path string) []byte {
	resp, err := http.Get(ns.RemoteURL + path)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return body
}
