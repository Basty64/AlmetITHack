package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const apiKey = "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"
const endpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	http.HandleFunc("/chat", handleChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type GPTBodyReq struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ClientRequest struct {
	Message string `json:"message"`
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка чтения тела запроса"))
		return
	}

	messages := &ClientRequest{}
	err = json.Unmarshal(body, &messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка десериализации тела запроса"))
		return
	}

	fmt.Println(messages)

	gptBodyReq := GPTBodyReq{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: messages.Message,
			},
		},
		Temperature: 0.2,
	}

	gptBodyReqBytes, err := json.Marshal(gptBodyReq)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка сериализации тела запроса"))
		return
	}

	fmt.Println(gptBodyReq)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(gptBodyReqBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания запроса"))
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка отправки запроса к API"))
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка чтения ответа от API"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
