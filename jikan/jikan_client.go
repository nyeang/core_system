package jikan

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 10 * time.Second},
        baseURL:    "https://api.jikan.moe/v4",
    }
}

func (c *Client) FetchAnime(page int) (*JikanResponse, error) {
    url := fmt.Sprintf("%s/anime?page=%d&limit=20&order_by=popularity", c.baseURL, page)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result JikanResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
