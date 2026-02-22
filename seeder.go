package main

import (
    "log"
    "strings"
    "time"
    "core-anime/config"
    "core-anime/jikan"
    "core-anime/models"
)

func SeedAnimeFromJikan() {
    client := jikan.NewClient()

    for page := 1; page <= 3; page++ {
        result, err := client.FetchAnime(page)
        if err != nil {
            log.Println("Jikan error:", err)
            continue
        }

        for _, a := range result.Data {
            // Skip if already exists
            var existing models.Anime
            if config.DB.Where("title = ?", a.Title).First(&existing).Error == nil {
                continue
            }

            // Build genre string
            genreNames := []string{}
            for _, g := range a.Genres {
                genreNames = append(genreNames, g.Name)
            }

            // Extract YouTube ID from embed URL
            youtubeID := extractYoutubeID(a.Trailer.EmbedURL)

            // Build proper trailer URL
            trailerURL := a.Trailer.URL
            if trailerURL == "" && youtubeID != "" {
                trailerURL = "https://www.youtube.com/watch?v=" + youtubeID
            }

            insertResult := config.DB.Create(&models.Anime{
                Title:             a.Title,
                Description:       a.Synopsis,
                Genres:            strings.Join(genreNames, ", "),
                ReleaseDate:       time.Now(),
                ImageURL:          a.Images.JPG.ImageURL,
                SmallImageURL:     a.Images.JPG.SmallImageURL,
                LargeImageURL:     a.Images.JPG.LargeImageURL,
                ImageURLWebP:      a.Images.WebP.ImageURL,
                SmallImageURLWebP: a.Images.WebP.SmallImageURL,
                LargeImageURLWebP: a.Images.WebP.LargeImageURL,
                TrailerURL:        trailerURL,  
                TrailerEmbedURL:   a.Trailer.EmbedURL,
                TrailerYoutubeID:  youtubeID,  
            })
            if insertResult.Error != nil {
                log.Println("❌ Insert error:", insertResult.Error)
            } else {
                log.Println("✓ Inserted:", a.Title)
            }
        }

        log.Printf("✓ Page %d seeded", page)
        time.Sleep(1 * time.Second)
    }

    log.Println("✓ Anime seeding complete!")
}

// Extract YouTube ID from embed URL
// e.g. https://www.youtube-nocookie.com/embed/qig4KOK2R2g?... → qig4KOK2R2g
func extractYoutubeID(embedURL string) string {
    if embedURL == "" {
        return ""
    }
    parts := strings.Split(embedURL, "/embed/")
    if len(parts) < 2 {
        return ""
    }
    return strings.Split(parts[1], "?")[0]
}
