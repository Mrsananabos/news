package repository

import (
	"context"
	_ "embed"
	"fmt"
	"service/internal/models"
	"strings"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

var (
	//go:embed sql/select_news_by_limit_and_offset.sql
	SqlSelectNewsByLimitAndOffset string
)

//go:generate mockery --name=INewsRepository --output=mocks --outpkg=mocks --case=snake --with-expecter
type INewsRepository interface {
	GetNews(limit, offset uint64) ([]models.NewsWithCategories, error)
	//CreateNews(limit int) ([]models.Order, error)
	UpdateNews(newsId uint64, updateFields map[string]interface{}, categories *[]uint64) error
}

type NewsRepository struct {
	db  *reform.DB
	log *logrus.Logger
	ctx context.Context
}

func NewNewsRepository(db *reform.DB, log *logrus.Logger, ctx context.Context) INewsRepository {
	return &NewsRepository{
		db:  db,
		log: log,
		ctx: ctx,
	}
}

func (r *NewsRepository) GetNews(limit, offset uint64) ([]models.NewsWithCategories, error) {
	opt := "repository.news.getNewsByLimitAndOffset"
	var newsList []models.NewsWithCategories
	rows, err := r.db.QueryContext(r.ctx, SqlSelectNewsByLimitAndOffset, limit, offset)

	if err != nil {
		return newsList, fmt.Errorf("%s: could not get news by limit and offset: %w", opt, err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.NewsWithCategories
		var categories []int64
		if err = rows.Scan(&n.ID, &n.Title, &n.Content, pq.Array(&categories)); err != nil {
			return nil, fmt.Errorf("%s: could not get news by limit and offset: %w", opt, err)
		}

		n.Categories = categories

		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: could not get news by limit and offset: %w", opt, err)
	}
	return newsList, nil
}

func (r *NewsRepository) UpdateNews(newsId uint64, updateFields map[string]interface{}, categories *[]uint64) error {
	opt := "repository.news.updateNews"
	tx, err := r.db.Begin()
	if err != nil {
		r.log.Errorf("%s: could not begin transaction: %w", opt, err)
		return fmt.Errorf("%s: could not update news: %w", opt, err)
	}
	defer tx.Rollback()

	if categories != nil {
		_, err = tx.Exec("DELETE FROM news_categories WHERE news_id = $1", newsId)
		if err != nil {
			r.log.Errorf("%s: could not delete news categories: %w", opt, err)
			return fmt.Errorf("%s: could not update news: %w", opt, err)
		}

		// Добавляем новые категории
		for _, categoryID := range *categories {
			_, err = tx.Exec("INSERT INTO news_categories (news_id, category_id) VALUES ($1, $2)", newsId, categoryID)
			if err != nil {
				r.log.Errorf("Ошибка добавления категории %d для новости ID=%d: %v", categoryID, newsId, err)
				return err
			}
		}
	}

	if len(updateFields) != 0 {

		setClauses := []string{}
		args := []interface{}{}
		i := 1

		for column, value := range updateFields {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
			args = append(args, value)
			i++
		}

		args = append(args, newsId) // в конце условие по id

		query := fmt.Sprintf("UPDATE news SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)

		_, err = tx.Exec(query, args...)
		if err != nil {
			r.log.Errorf("%s: could not update news categories: %w", opt, err)
			return fmt.Errorf("%s: could not update news: %w", opt, err)
		}
	}

	if err = tx.Commit(); err != nil {
		r.log.Errorf("Ошибка коммита транзакции: %v", err)
		return err
	}

	r.log.Infof("successfully updated news categories with id=%d", newsId)

	return nil

}

//func (r *NewsRepository) GetCategoriesForNews(newsID int64) ([]int64, error) {
//	rows, err := r.db.Query("SELECT category_id FROM news_categories WHERE news_id = $1", newsID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var categories []int64
//	for rows.Next() {
//		var categoryID int64
//		if err := rows.Scan(&categoryID); err != nil {
//			r.log.Errorf("Ошибка сканирования категории: %v", err)
//			continue
//		}
//		categories = append(categories, categoryID)
//	}
//
//	return categories, nil
//}
//
//func (r *NewsRepository) UpdateCategories(newsID int64, categories []int64) error {
//	tx, err := r.db.Begin()
//	if err != nil {
//		r.log.Errorf("Ошибка начала транзакции: %v", err)
//		return err
//	}
//	defer tx.Rollback()
//
//	// Удаляем старые категории
//	_, err = tx.Exec("DELETE FROM news_categories WHERE news_id = $1", newsID)
//	if err != nil {
//		r.log.Errorf("Ошибка удаления категорий для новости ID=%d: %v", newsID, err)
//		return err
//	}
//
//	// Добавляем новые категории
//	for _, categoryID := range categories {
//		_, err = tx.Exec("INSERT INTO news_categories (news_id, category_id) VALUES ($1, $2)", newsID, categoryID)
//		if err != nil {
//			r.log.Errorf("Ошибка добавления категории %d для новости ID=%d: %v", categoryID, newsID, err)
//			return err
//		}
//	}
//
//	if err := tx.Commit(); err != nil {
//		r.log.Errorf("Ошибка коммита транзакции: %v", err)
//		return err
//	}
//
//	return nil
//}
