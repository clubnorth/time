package handler

import (
  "bytes"
  "encoding/json"
  "io"
  "net/http"
  "strings"
  "time"
)

type mediaInfoRequest struct {
  Name      string `json:"name"`
  MediaType string `json:"type"`
}

type mediaInfoResponse struct {
  Director         string   `json:"director"`
  Cast             string   `json:"cast"`
  Country          string   `json:"country"`
  FirstReleaseDate string   `json:"first_release_date"`
  Duration         string   `json:"duration"`
  Episodes         string   `json:"episodes"`
  Tags             []string `json:"tags"`
}

func MediaInfo(apiKey string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var req mediaInfoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
      respondError(w, 400, "name required")
      return
    }

    var jsonTemplate string
    if req.MediaType == "movie" {
      jsonTemplate = `{"director":"导演","cast":"主演","country":"国家","first_release_date":"年份","duration":"分钟","tags":["标签1","标签2","标签3"]}`
    } else {
      jsonTemplate = `{"director":"导演","cast":"主演","country":"国家","first_release_date":"年份","episodes":"集数","tags":["标签1","标签2","标签3"]}`
    }

    typeLabel := map[string]string{"movie": "电影", "series": "剧集", "anime": "动漫"}[req.MediaType]
    if typeLabel == "" { typeLabel = "影视" }

    body := map[string]interface{}{
      "model": "deepseek-chat",
      "messages": []map[string]string{
        {"role": "system", "content": `你是影视信息助手。你的训练数据包含大量豆瓣、IMDb、TMDB、维基百科的电影电视剧动漫信息。请尽力回忆并返回准确信息。如果记不清，根据类似作品合理推测，但绝不用"未知"作为答案——可以用"可能xxx"。标签要有3个概括类型的简短中文。`},
        {"role": "user", "content": typeLabel + "《" + req.Name + "》的导演、主演、国家、年份、时长等，严格按此JSON格式返回（不要markdown）：\n" + jsonTemplate},
      },
      "max_tokens":  350,
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

    var info mediaInfoResponse
    if err := json.Unmarshal([]byte(content), &info); err != nil {
      info = mediaInfoResponse{Director: content}
    }

    respond(w, 200, info)
  }
}
