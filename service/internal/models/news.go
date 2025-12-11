package models

//go:generate reform
//reform:news
type News struct {
	ID      int64  `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}

// NewsWithCategories используется для ответа
type NewsWithCategories struct {
	News
	Categories []int64
}

type NewsEditForm struct {
	Title      *string  `json:"title" validate:"omitempty"`
	Content    *string  `json:"content" validate:"omitempty"`
	Categories *[]int64 `json:"categories" validate:"omitempty" `
}

type NewsCreateForm struct {
	Title      string   `json:"title" validate:"omitempty"`
	Content    string   `json:"content" validate:"omitempty"`
	Categories *[]int64 `json:"categories" validate:"omitempty" `
}
