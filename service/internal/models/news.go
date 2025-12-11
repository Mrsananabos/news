package models

import "gopkg.in/reform.v1"

type News struct {
	ID      int64  `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}

func (News) Schema() string {
	return ""
}

// Table возвращает имя таблицы в БД
func (News) Table() string {
	return "news"
}

// PKColumnName возвращает имя колонки primary key
func (News) PKColumnName() string {
	return "id"
}

// Values возвращает значения всех полей в порядке колонок
// Нужен для INSERT и UPDATE
func (n News) Values() []interface{} {
	return []interface{}{
		n.ID,      // id
		n.Title,   // title
		n.Content, // content
	}
}

// HasPK - говорит Reform, есть ли уже PK у записи
func (n *News) HasPK() bool {
	return n.ID != 0 // если ID == 0, значит новая запись
}

// PKValue и PKPointer - для работы с primary key
func (n *News) PKValue() interface{}   { return n.ID }
func (n *News) PKPointer() interface{} { return &n.ID }

// Pointers возвращает указатели на все поля в порядке колонок
// Нужен для SELECT (сканирование результатов)
func (n *News) Pointers() []interface{} {
	return []interface{}{
		&n.ID,      // id
		&n.Title,   // title
		&n.Content, // content
	}
}

// View возвращает View, если это view (у нас таблица, поэтому nil)
func (News) View() reform.View {
	return nil
}

// String возвращает строковое представление (для логирования)
func (n News) String() string {
	return "News: " + n.Title
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

type NewsCreateForm struct {
	Title      string    `json:"title" validate:"omitempty" example:"Medicine"`
	Content    string    `json:"content" validate:"omitempty"`
	Categories *[]uint64 `json:"categories" validate:"omitempty" `
}
