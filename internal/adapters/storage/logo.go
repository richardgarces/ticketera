package storage

import (
	"context"
	"ticketera/internal/domain"
)

func (m *MongoRepository) SaveLogoPath(path string) error {
	coll := m.Client.Database("ticketera").Collection("logo")
	_, err := coll.InsertOne(context.TODO(), domain.Logo{FileURL: path})
	return err
}

var DefaultRepo *MongoRepository

func SaveLogoPath(path string) error {
	if DefaultRepo == nil {
		DefaultRepo = NewMongoRepository("mongodb://localhost:27017")
	}
	return DefaultRepo.SaveLogoPath(path)
}
