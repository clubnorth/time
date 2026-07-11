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

    prompt := fmt.Sprintf(`你是图书管理员。请查询书籍《%s》的真实信息，以纯JSON格式返回（不要markdown标记）：
{"author":"作者姓名","nationality":"作者国籍","word_count":"字数（万字，如90）","first_publish_date":"首次出版年份(YYYY)","tags":["标签1","标签2","标签3"]}
tags是3个简短中文标签概括这本书的类型或主题。必须填写真实数据，找不到信息就查相关资料合理推测。`, req.BookName)

    body := map[string]interface{}{
      "model": "deepseek-chat",
      "messages": []map[string]string{
        {"role": "system", "content": `你是专业的图书管理员。对于任何书籍，你必须返回具体的JSON数据，包括作者、国籍、字数和首次出版日期。如果数据不完全确定，请根据你的知识合理推测。`},
        {"role": "user", "content": prompt},
      },
      "max_tokens": 300,
      "temperature": 0.3,
    }
    jsonBody, _ := json.Marshal(body)

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
      info = bookInfoResponse{Author: content, Nationality: "", WordCount: "", FirstPublishDate: ""}
    }

    respond(w, 200, info)
  }
}
