package store

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/skye-z/olympus/util"
)

const (
	mavenRepository = "./repository/maven/"
)

type MavenStore struct {
	RemoteURL string
}

// 获取文件
func (ms MavenStore) GetFile(path string) []byte {
	fileName := path[strings.LastIndex(path, "/"):]
	directory := path[:strings.LastIndex(path, "/")]
	if util.CheckExist(mavenRepository + path) {
		return util.ReadFile(mavenRepository + path)
	} else {
		content := ms.getRemoteData(path)
		util.SaveFile(mavenRepository+directory, fileName, content)
		return content
	}
}

// 获取远程数据
func (ms MavenStore) getRemoteData(path string) []byte {
	resp, err := http.Get(ms.RemoteURL + path)
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
