package resumeSdk

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func ParseByUrl(url string, fileName string, appCode string) map[string]any {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	content, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return parse(content, fileName, appCode)
}

func ParseByFilePath(path string, fileName string, appCode string) map[string]any {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return parse(content, fileName, appCode)
}

func parse(content []byte, fileName string, appCode string) map[string]any {
	client := &http.Client{}
	payload := struct {
		FileName   string `json:"file_name,omitempty"`
		FileCont   string `json:"file_cont,omitempty"`
		NeedAvatar int    `json:"need_avatar,omitempty"`
		OcrType    int    `json:"ocr_type,omitempty"`
	}{
		fileName,
		base64.StdEncoding.EncodeToString(content),
		1,
		1,
	}
	value, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://resumesdk.market.alicloudapi.com/ResumeParser", bytes.NewReader(value))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "APPCODE "+appCode)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Request fail")
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	data := result["result"].(map[string]any)
	return data
}
