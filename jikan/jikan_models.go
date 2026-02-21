package jikan

type JikanResponse struct {
    Pagination Pagination `json:"pagination"`
    Data       []Anime    `json:"data"`
}

type Pagination struct {
    LastVisiblePage int  `json:"last_visible_page"`
    HasNextPage     bool `json:"has_next_page"`
    CurrentPage     int  `json:"current_page"`
    Items           struct {
        Count    int `json:"count"`
        Total    int `json:"total"`
        PerPage  int `json:"per_page"`
    } `json:"items"`
}

type Anime struct {
    MalID   int    `json:"mal_id"`
    Url     string `json:"url"`
    Title   string `json:"title"`
    Type    string `json:"type"`
    Status  string `json:"status"`
    Episodes int   `json:"episodes"`
    Score   float64 `json:"score"`
    // Add more fields as needed, based on the JSON you care about.
}
