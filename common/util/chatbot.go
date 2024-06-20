package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type ChatbotResponse struct {
	Response string `json:"response"`
}

func GetChatbotResponse(lang string, message string) (string, error) {
	reqUrl, err := url.Parse(os.Getenv("MLAPI_URL"))
	if err != nil {
		return "", err
	}

	query := reqUrl.Query()
	query.Add("language", lang)
	query.Add("input", message)
	reqUrl.RawQuery = query.Encode()

	r, err := http.NewRequest("POST", reqUrl.String(), nil)
	if err != nil {
		return "", err
	}

	fmt.Println("a")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	fmt.Println("b")

	getBody, _ := io.ReadAll(res.Body)
	var getResp ChatbotResponse
	if err := json.Unmarshal(getBody, &getResp); err != nil {
		return "", err
	}

	fmt.Println("c")

	if res.StatusCode != http.StatusOK {
		return "", errors.New("request failed")
	}
	fmt.Println(getResp.Response)
	return getResp.Response, nil
}
