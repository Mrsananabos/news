package models

//go:generate reform
//reform:news_categories
type NewsCategory struct {
	NewsId     int64 `reform:"news_id,pk"`
	CategoryId int64 `reform:"category_id"`
}
