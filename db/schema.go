package db

type WebNovel struct {
	Name          string `bson:"name"`
	AuthorName    string `bson:"author_name"`
	TotalChapters int    `bson:"total_chapters"`
	Info          string `bson:"info"`
	ImageUrlPath  string `bson:"image_url_path"`
	UrlPath       string `bson:"url_path"`
}

type Chapters struct {
	WebNovelUrlPath string    `bson:"webnovel_url_path"`
	ChaptersUrlPath []Chapter `bson:"chapters"`
}

type Chapter struct {
	Name    string `bson:"name"`
	UrlPath string `bson:"url_path"`
}
