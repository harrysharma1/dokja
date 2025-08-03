package routes

import (
	"dokja/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func Listen() {
	engine := django.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

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

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title":     "Dokja",
			"WebNovels": webNovels,
		})
	})

	log.Fatal(app.Listen(":6969"))
}
