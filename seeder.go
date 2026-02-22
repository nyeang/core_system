package main

import (
    "log"
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

            // ← ADD ERROR CHECK HERE
			insertResult := config.DB.Create(&models.Anime{
				Title:             a.Title,
				Description:       a.Synopsis,
				ReleaseDate:       time.Now(),
				ImageURL:          a.Images.JPG.ImageURL,
				SmallImageURL:     a.Images.JPG.SmallImageURL,
				LargeImageURL:     a.Images.JPG.LargeImageURL,
				ImageURLWebP:      a.Images.WebP.ImageURL,
				SmallImageURLWebP: a.Images.WebP.SmallImageURL,
				LargeImageURLWebP: a.Images.WebP.LargeImageURL,
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
