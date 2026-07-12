package handler

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "strings"
  "time"
)

type bookInfoRequest struct {
  BookName string `json:"book_name"`
}

type bookInfoResponse struct {
  Author           string   `json:"author"`
  Nationality      string   `json:"nationality"`
  WordCount        string   `json:"word_count"`
  FirstPublishDate string   `json:"first_publish_date"`
  Tags             []string `json:"tags"`
}

func BookInfo(apiKey string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var req bookInfoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookName == "" {
      respondError(w, 400, "book_name required")
      return
    }

    jsonTemplate := `{"author":"作者","nationality":"国籍","word_count":"字数万字","first_publish_date":"年份","tags":["标签1","标签2","标签3"]}`

    body := map[string]interface{}{
      "model": "deepseek-chat",
      "messages": []map[string]string{
        {"role": "system", "content": `你是图书信息助手。你的训练数据包含大量豆瓣读书、Goodreads、亚马逊的书籍信息。请尽力回忆并返回准确信息。如果记不清，根据类似作品合理推测，但绝不用"未知"作为答案——可以用"可能xxx"。标签要有3个概括类型的简短中文。`},
        {"role": "user", "content": fmt.Sprintf(`书籍《%s》的作者、国籍、字数、首次出版年份、标签，严格按此JSON格式返回（不要markdown）：\n%s`, req.BookName, jsonTemplate)},
      },
      "max_tokens":  300,
      "temperature": 0.6,
    }
    jsonBody, err := json.Marshal(body)
    if err != nil {
      respondError(w, 500, "internal error")
      return
    }

    client := &http.Client{Timeout: 15 * time.Second}
    httpReq, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(jsonBody))
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+apiKey)

    resp, err := client.Do(httpReq)
    if err != nil {
      respondError(w, 500, "API request failed: "+err.Error())
      return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    var result struct {
      Choices []struct {
        Message struct {
          Content string `json:"content"`
        } `json:"message"`
      } `json:"choices"`
    }
    if err := json.Unmarshal(respBody, &result); err != nil || len(result.Choices) == 0 {
      respondError(w, 500, "failed to parse API response")
      return
    }

    content := strings.TrimSpace(result.Choices[0].Message.Content)
    content = strings.TrimPrefix(content, "```json")
    content = strings.TrimPrefix(content, "```")
    content = strings.TrimSuffix(content, "```")
    content = strings.TrimSpace(content)

    var info bookInfoResponse
    if err := json.Unmarshal([]byte(content), &info); err != nil {
      info = bookInfoResponse{Author: content}
    }

    respond(w, 200, info)
  }
}
