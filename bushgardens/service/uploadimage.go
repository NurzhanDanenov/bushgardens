package service

import (
	"bush/repository"
	"mime/multipart"
)

type UploadImageServise struct {
	repo repository.UploadImage
}

func NewUploadImageService(repo repository.UploadImage) *UploadImageServise {
	return &UploadImageServise{repo: repo}
}

func (s *UploadImageServise) Upload(userId string, file multipart.File) (string, error) {
	return s.repo.UploadImage(userId, file)
}
