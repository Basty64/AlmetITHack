package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var apiKey = os.Getenv("CHATGPT_KEY")

const endpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	var jsonResponse JSONResponse
	var tags []string
	tags = append(tags, "IoT", "security")
	jsonResponse.sendToCrossref(tags)
	//http.HandleFunc("/article", handleArticle)
	//http.HandleFunc("/sentence", handleSentence)
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

type ChatCompletion struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type JSONResponse struct {
	Text  string     `json:"text"`
	Tags  []string   `json:"tags"`
	Links []Articles `json:"links"`
}
type Articles struct {
	Message struct {
		Items []struct {
			Title []string `json:"title"`
			URL   string   `json:"URL"`
		} `json:"items"`
	} `json:"message"`
}

func handleArticle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка чтения тела запроса"))

	}

	messages := &ClientRequest{}
	err = json.Unmarshal(body, &messages)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка десериализации тела запроса"))

	}

	gptBodyReq := GPTBodyReq{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role: "system",
				Content: "Write summarize about this article " +
					"and write five tags as JSON Array about this article. Do it on Russian" +
					"Structure of ur answer is JSON:" +
					"{\"text\": " +
					",\"tags\": " +
					"}",
			},
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

	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(gptBodyReqBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания запроса"))

	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка чтения ответа от API"))
		return
	}

	var chatCompletion ChatCompletion
	if err := json.Unmarshal([]byte(responseBody), &chatCompletion); err != nil {
		fmt.Println("Ошибка разбора JSON:", err)
		return
	}

	var jsonResponse JSONResponse
	if len(chatCompletion.Choices) > 0 {
		content := chatCompletion.Choices[0].Message.Content
		fmt.Println(content)
		json.Unmarshal([]byte(content), &jsonResponse)

	}

	jsonResponse.sendToCrossref(jsonResponse.Tags)

	fmt.Println(jsonResponse)
	b, _ := json.Marshal(jsonResponse)

	w.Write(b)
}

func handleSentence(w http.ResponseWriter, r *http.Request) {

}

func (jsonresp *JSONResponse) sendToCrossref(tags []string) {
	//
	crossrefendpoint := "https://api.crossref.org/works?select=title,URL,type&rows=2&filter=type:journal-article&query=" + tags[0]
	var res string
	for i := 1; i < len(tags); i++ {
		res += "," + tags[i]
	}
	crossrefendpoint += res
	fmt.Println(crossrefendpoint)
	req, _ := http.NewRequest("GET", crossrefendpoint, bytes.NewReader([]byte("")))

	client := &http.Client{}
	resp, _ := client.Do(req)
	var articles Articles
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal([]byte(body), &articles); err != nil {
		fmt.Println("Ошибка разбора JSON:", err)
		return
	}

	jsonresp.Links = append(jsonresp.Links, articles)
	fmt.Println(jsonresp.Links)
}

//func requestBuilder(w http.ResponseWriter, r *http.Request) *http.Response {
//
//}
