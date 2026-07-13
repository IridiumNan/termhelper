package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

const defaultOllamaURL = "http://100.120.81.67:11434/api/generate"

func main() {
	start := time.Now()
	msg := Message{
		Role: "user",
		Content: `❯ python test.py
  File "/tmp/cai/go-ollama-demo/test.py", line 5
    else:
    ^^^^
SyntaxError: invalid syntax
		我是一个刚学习python的新手， 这个报错的中文意思是什么,你只需要解释其中的中文意思，不需要帮我解决报错?`,
	}
	req := Request{
		Model:    "qwen3.5:4b",
		Stream:   false,
		Messages: []Message{msg},
	}
	resp, err := talkToOllama(defaultOllamaURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.Message.Content)
	fmt.Printf("Completed in %v", time.Since(start))
}

func talkToOllama(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}
