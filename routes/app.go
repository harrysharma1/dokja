package routes

import (
	"dokja/db"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func preFill() {
	webNovels := []db.WebNovel{
		{
			Name:          "Omniscient Reader’s Viewpoint",
			AuthorName:    "Sing Shong",
			TotalChapters: 550,
			Info:          "A reader becomes part of the novel he was reading.",
			ImageUrlPath:  "https://static.wikia.nocookie.net/omniscient-readers-viewpoint/images/4/43/Kim_Dokja_Cover.jpg",
			UrlPath:       "/omniscient-readers-viewpoint",
		},
		{
			Name:          "Solo Leveling",
			AuthorName:    "Chugong",
			TotalChapters: 270,
			Info:          "A weak hunter becomes the strongest.",
			ImageUrlPath:  "https://images.immediate.co.uk/production/volatile/sites/3/2025/03/solo-levelling-season-2-finale-993c204.jpg",
			UrlPath:       "/solo-leveling",
		},
	}
	for i := range webNovels {
		db.PutToMongo(webNovels[i])
	}
}

func Listen() {
	engine := django.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		webNovels, err := db.FindAllWebNovels()
		if err != nil {
			log.Fatalf("Error Getting Novels: %s", err)
		}
		return c.Render("index", fiber.Map{
			"Title":     "Dokja",
			"WebNovels": webNovels,
		})
	})

	app.Get("/novels/:name?", func(c *fiber.Ctx) error {
		webNovel, err := db.FindWebNovelBasedOnUrlParam(c.Params("name"))
		if err != nil {
			log.Fatalf("Error Getting Novel: %s", err)
		}
		return c.SendString(fmt.Sprintf("Name: %s,\nAuthor: %s,\nImage: %s,\nPath: %s", webNovel.Name, webNovel.AuthorName, webNovel.ImageUrlPath, webNovel.UrlPath))
	})

	log.Fatal(app.Listen(":6969"))
}
