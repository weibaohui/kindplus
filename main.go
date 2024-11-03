package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	Clusters []Cluster `yaml:"clusters"`
	Users    []User    `yaml:"users"`
}

type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		CA     string `yaml:"certificate-authority-data"`
		Server string `yaml:"server"`
	} `yaml:"cluster"`
}

type User struct {
	Name string `yaml:"name"`
	User struct {
		CertData string `yaml:"client-certificate-data"`
		KeyData  string `yaml:"client-key-data"`
	} `yaml:"user"`
}

func main() {
	// 指定存放 kubeconfig 文件的目录
	kubeconfigDir := "./kubeconfig" // 替换为实际路径

	// 遍历目录中的所有 kubeconfig 文件
	files, err := ioutil.ReadDir(kubeconfigDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		kubeconfigPath := filepath.Join(kubeconfigDir, file.Name())
		processKubeConfig(kubeconfigPath)
	}
}

func processKubeConfig(path string) {
	// 读取 kubeconfig 文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var kubeConfig KubeConfig
	if err := yaml.Unmarshal(data, &kubeConfig); err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

	for i, user := range kubeConfig.Users {
		if i < len(kubeConfig.Clusters) {
			// 获取对应的 server 域名
			cluster := kubeConfig.Clusters[i]
			serverURL := cluster.Cluster.Server
			parsedURL, err := url.Parse(serverURL)
			if err != nil {
				fmt.Println("Error parsing server URL:", err)
				continue
			}

			// 创建以 certs/ 和 serverName（域名）为名称的文件夹
			serverName := parsedURL.Hostname()
			certsDir := filepath.Join("certs", serverName) // 添加 certs/ 目录前缀
			err = os.MkdirAll(certsDir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				continue
			}

			// 解码并保存证书和密钥
			if err := saveFile(filepath.Join(certsDir, "server.crt"), user.User.CertData); err != nil {
				fmt.Println("Error saving certificate:", err)
			}
			if err := saveFile(filepath.Join(certsDir, "server.key"), user.User.KeyData); err != nil {
				fmt.Println("Error saving key:", err)
			}
			if err := saveFile(filepath.Join(certsDir, "ca.crt"), cluster.Cluster.CA); err != nil {
				fmt.Println("Error saving ca:", err)
			}
		}
	}
}
func saveFile(filename, data string) error {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, decodedData, 0644)
}
