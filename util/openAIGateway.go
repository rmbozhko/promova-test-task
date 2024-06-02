package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type moderationRequest struct {
	Input string `json:"input"`
}

const OpenAIToken string = "OPENAI_API_KEY"

func InputDataIsSafe(inputData string) (bool, error) {
	apiKey := os.Getenv(OpenAIToken)
	url := "https://api.openai.com/v1/moderations"

	data := moderationRequest{
		Input: inputData,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return false, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}

	ch := make(chan *http.Response)

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		ch <- resp
	}()

	resp := <-ch
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error fetching response:", err)
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false, err
	}

	var moderationResult map[string]interface{}
	err = json.Unmarshal(body, &moderationResult)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return false, err
	}

	if v, ok := moderationResult["error"]; ok {
		errorResponseData := v.(map[string]interface{})
		return false, fmt.Errorf(errorResponseData["message"].(string))
	} else {
		results := moderationResult["results"].([]interface{})
		firstResult := results[0].(map[string]interface{})
		flagged := firstResult["flagged"].(bool)
		return flagged, nil
	}
}
