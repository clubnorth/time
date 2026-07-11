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

    var fields string
    if req.MediaType == "movie" {
      fields = `"director":"导演（最多两人，用/分隔）","cast":"主演（最多三人，用/分隔）","country":"国家","first_release_date":"首次上映年份(YYYY)","duration":"时长（分钟）","tags":["标签1","标签2","标签3"]`
    } else {
      fields = `"director":"导演（最多两人，用/分隔）","cast":"主演（最多三人，用/分隔）","country":"国家","first_release_date":"首次上映年份(YYYY)","episodes":"集数（如：45min/12集 或 24min/24集）","tags":["标签1","标签2","标签3"]`
    }

    typeLabel := map[string]string{"movie": "电影", "series": "剧集", "anime": "动漫"}[req.MediaType]
    if typeLabel == "" { typeLabel = "影视" }

    prompt := strings.Join([]string{
      `你是影视数据库管理员。请查询` + typeLabel + `《` + req.Name + `》的真实信息，以纯JSON格式返回（不要markdown标记）：`,
      `{` + fields + `}`,
      `tags是3个简短中文标签概括类型或主题。必须填写真实数据，找不到信息就查相关资料合理推测。`,
    }, "\n")

    body := map[string]interface{}{
      "model": "deepseek-chat",
      "messages": []map[string]string{
        {"role": "system", "content": `你是专业的影视数据库管理员。对于任何影视作品，你必须返回具体的JSON数据，包括导演、主演、国家、上映年份、时长/集数和标签。如果数据不完全确定，请根据你的知识合理推测。`},
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

    var info mediaInfoResponse
    if err := json.Unmarshal([]byte(content), &info); err != nil {
      info = mediaInfoResponse{Director: content}
    }

    respond(w, 200, info)
  }
}
