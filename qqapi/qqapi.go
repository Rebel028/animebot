package qqapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func RequestImage(base64 string) Results {
	url := "https://ai.tu.qq.com/trpc.shadow_cv.ai_processor_cgi.AIProcessorCgi/Process"

	str := fmt.Sprintf("{\n\t\"busiId\": \"ai_painting_anime_entry\",\n\t\"images\": [\"%s\"]\n}", base64)

	payload := strings.NewReader(str)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", "https://h5.tu.qq.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)

	fullResp, err := UnmarshalQqResponse(body)
	if err != nil {
		log.Panic(err)
	}

	log.Println(fullResp)

	results, _ := UnmarshalResults([]byte(fullResp.Extra))

	log.Println(results)

	return results
}

func UnmarshalQqResponse(data []byte) (QqResponse, error) {
	var r QqResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *QqResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type QqResponse struct {
	Code   int64         `json:"code"`
	Msg    string        `json:"msg"`
	Images []interface{} `json:"images"`
	Faces  []interface{} `json:"faces"`
	Extra  string        `json:"extra"`
	Videos []interface{} `json:"videos"`
}

func UnmarshalResults(data []byte) (Results, error) {
	var r Results
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Results) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Results struct {
	VideoUrls []string `json:"video_urls"`
	ImgUrls   []string `json:"img_urls"`
}
