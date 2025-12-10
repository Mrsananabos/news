package models

// News представляет строку таблицы News
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
	Title      *string   `json:"title" validate:"omitempty" example:"Medicine"`
	Content    *string   `json:"content" validate:"omitempty"`
	Categories *[]uint64 `json:"categories" validate:"omitempty" `
}

type NewsListsResponse struct {
	Success bool
	News    []NewsWithCategories
}
