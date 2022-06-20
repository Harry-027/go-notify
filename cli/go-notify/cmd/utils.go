package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	ussr "os/user"
	"path"

	"github.com/Harry-027/go-notify/api-server/config"
)

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Config struct {
	Token string `json:"token"`
}

func fetchBasePath() string {
	basePath := config.GetConfig(config.SERVER_URL)
	if basePath != "" {
		return basePath
	}
	return "http://localhost:3001" // default server url for local setup
}

func makeCallToServer(methodType, callType, url string, token string, data []byte) (string, []byte) {

	baseApi := fetchBasePath()
	var payload []byte
	uri := fmt.Sprintf("%s%s", baseApi, url)
	if data != nil {
		payload = data
	}
	return makeAuthCall(methodType, callType, uri, token, bytes.NewBuffer(payload))
}

func makeAuthCall(methodType, callType, uri, token string, data *bytes.Buffer) (string, []byte) {
	client := http.Client{}
	req, err := http.NewRequest(methodType, uri, data)
	req.Header.Set("Content-type", "application/json")

	if callType == "USER" {
		req.Header.Set("Authorization", token)
	}

	if err != nil {
		log.Fatal("An error occurred: ", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("An error occurred: ", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("An error occurred while reading response: ", err.Error())
	}
	return resp.Status, body
}

func configPath() string {
	cfgFile := ".go-notify-config"
	usr, _ := ussr.Current()
	return path.Join(usr.HomeDir, cfgFile)
}

func saveConfig(c Config) (err error) {
	jsonC, _ := json.Marshal(c)
	_, err = os.Create(configPath())
	if err != nil {
		return
	}

	err = os.WriteFile(configPath(), jsonC, os.ModeAppend)
	return
}

func readConfig() Config {
	data, _ := ioutil.ReadFile(configPath())
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

func dataConversion(body []byte, keys []string) [][]string {
	var mapSlice []map[string]interface{}
	if err := json.Unmarshal(body, &mapSlice); err != nil {
		panic(err)
	}
	var stringSlice [][]string
	for _, m := range mapSlice {
		var resultMap []string
		for _, v := range keys {
			resultMap = append(resultMap, fmt.Sprintf("%v", m[v]))
		}
		stringSlice = append(stringSlice, resultMap)
	}
	return stringSlice
}
