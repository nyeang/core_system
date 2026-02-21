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

func (c *Client) SearchAction(page, limit int) (*JikanResponse, error) {
    url := fmt.Sprintf("%s/anime?q=action&page=%d&limit=%d", c.baseURL, page, limit)

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %s", resp.Status)
    }

    var result JikanResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}
