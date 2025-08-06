package routes

import (
	"dokja/db"
	"dokja/util"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func preFill() {
	webNovels := []db.WebNovel{
		{
			Name:          "Omniscient Reader’s Viewpoint",
			AuthorName:    "Sing Shong",
			TotalChapters: 550,
			Info: `Kim Dokja does not consider himself the protagonist of his own life. Befitting the name his parents gave him, he is a solitary person whose sole hobby is reading web novels. For over a decade, he has lived vicariously through Yu Junghyeok, the main character of the web novel Three Ways to Survive the Apocalypse (TWSA).
Through Junghyeok, Dokja has experienced secondhand the trials of repeatedly regressing in time, in search of an end to life-threatening “scenarios” that force people to act out narratives for the amusement of god-like “Constellations.”
After reading 3,149 chapters—long after all other readers lost interest—Dokja finally resigns himself to the story ending. However, he receives an enigmatic message from the author, stating that the story will soon be monetized, before his surroundings suddenly go dark.
He swiftly realizes that fiction has become reality and he is now living through TWSA. Although he is the singular omniscient reader of the events yet to come, his success in the scenarios is not guaranteed—but perhaps his advantage will empower him to step into the protagonist role that never suited him before.`,
			ImageUrlPath: "https://static.wikia.nocookie.net/omniscient-readers-viewpoint/images/4/43/Kim_Dokja_Cover.jpg",
			UrlPath:      "/omniscient-readers-viewpoint",
		},
		{
			Name:          "Solo Leveling",
			AuthorName:    "Chugong",
			TotalChapters: 270,
			Info: `In this world where Hunters with various magical powers battle monsters from invading the defenceless humanity, Seong Jin-Woo was the weakest of all the Hunters, barely able to make a living.
However, a mysterious System grants him the power of the ‘Player’, setting him on a course for an incredible and often times perilous Journey.
Follow Sung Jin-Woo as he embarks on an adventure to become an unparalleled existence through his “Level-Up” system - the only one in the entire world!`,
			ImageUrlPath: "https://images.immediate.co.uk/production/volatile/sites/3/2025/03/solo-levelling-season-2-finale-993c204.jpg",
			UrlPath:      "/solo-leveling",
		},
	}
	for i := range webNovels {
		if err := db.InsertWebNovel(webNovels[i]); err != nil {
			log.Fatalf("Error inserting Webnovel: %s", err)
		}
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

	app.Get("/add", func(c *fiber.Ctx) error {
		return c.Render("add", fiber.Map{
			"Title": "Dokja - Add",
		})
	})

	app.Post("/chapter/new", func(c *fiber.Ctx) error {
		lastRouteFull := strings.Split(c.Get("Referer"), "/")
		lastRoutePart := lastRouteFull[len(lastRouteFull)-2:]
		lastRoutePath := strings.Join(lastRoutePart, "/")
		lastRoutePath = "/" + lastRoutePath

		chapterInfo := db.Chapter{
			WebNovelUrlPath: lastRoutePath,
			Number:          util.ParseInt(c.FormValue("chapter_no")),
			Title:           c.FormValue("chapter_name"),
			UrlPath:         fmt.Sprintf("%s/chapters/%s", lastRoutePath, c.FormValue("chapter_no")),
			Text:            c.FormValue("chapter_text"),
		}
		if err := db.InsertChapter(chapterInfo); err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to insert %s", chapterInfo.Title))
		}

		return c.Redirect(lastRoutePath)
	})

	app.Post("/novels/new", func(c *fiber.Ctx) error {
		webNovel := db.WebNovel{
			Name:          c.FormValue("name"),
			AuthorName:    c.FormValue("author_name"),
			TotalChapters: util.ParseInt(c.FormValue("total_chapters")),
			Info:          c.FormValue("info"),
			ImageUrlPath:  c.FormValue("image_url_path"),
			UrlPath:       util.Sluggify(c.FormValue("name")),
		}
		if err := db.InsertWebNovel(webNovel); err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to insert %s", webNovel.Name))
		}
		return c.Redirect("/")
	})

	app.Get("/novels/:name?", func(c *fiber.Ctx) error {
		webNovel, chapters, err := db.FindWebNovelBasedOnUrlParam(c.Params("name"))
		if err != nil {
			log.Fatalf("Error Getting Novel: %s", err)
		}
		if len(webNovel.Name) == 0 {
			return c.Redirect("/")
		}
		return c.Render("individual_novel_page", fiber.Map{
			"Title":    fmt.Sprintf("Dokja - %s", webNovel.Name),
			"WebNovel": webNovel,
			"Chapters": chapters,
		})
	})

	log.Fatal(app.Listen(":6969"))
}
