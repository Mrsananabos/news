package service

import (
	"service/internal/models"
	"service/internal/repository"

	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=INewsService --output=mocks --outpkg=mocks --case=snake --with-expecter
type INewsService interface {
	EditNews(newsId uint64, editForm models.NewsEditForm) error
	ListNews(limit, offset uint64) ([]models.NewsWithCategories, error)
}
type NewsService struct {
	repo repository.INewsRepository
	log  *logrus.Logger
}

func NewNewsService(repo repository.INewsRepository, log *logrus.Logger) INewsService {
	return &NewsService{
		repo: repo,
		log:  log,
	}
}

func (s *NewsService) EditNews(newsId uint64, editForm models.NewsEditForm) error {
	updateFields := make(map[string]interface{})
	if editForm.Title != nil {
		updateFields["title"] = editForm.Title
	}
	if editForm.Content != nil {
		updateFields["content"] = editForm.Content
	}

	// Обновляем поля новости
	if len(updateFields) > 0 || editForm.Categories != nil {
		if err := s.repo.UpdateNews(newsId, updateFields, editForm.Categories); err != nil {
			logrus.Error(err)
			return err
		}
	}

	return nil
}

func (s *NewsService) ListNews(limit, offset uint64) ([]models.NewsWithCategories, error) {
	//добавить валидацию лимита и оффсета
	var newsList []models.NewsWithCategories

	newsList, err := s.repo.GetNews(limit, offset)
	if err != nil {
		return newsList, err
	}

	return newsList, nil
}
