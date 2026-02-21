package jikan

type JikanResponse struct {
    Data []AnimeData `json:"data"`
}

type AnimeData struct {
    MalID    int     `json:"mal_id"`
    Title    string  `json:"title"`
    Synopsis string  `json:"synopsis"`
    Score    float64 `json:"score"`
    Episodes int     `json:"episodes"`
    Images   Images  `json:"images"`
}

type Images struct {
    JPG  ImageFormat `json:"jpg"`
    WebP ImageFormat `json:"webp"`
}

type ImageFormat struct {
    ImageURL      string `json:"image_url"`
    SmallImageURL string `json:"small_image_url"`
    LargeImageURL string `json:"large_image_url"`
}
