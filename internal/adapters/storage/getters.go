package storage

import (
	"context"
	"ticketera/internal/domain"
)

func (m *MongoRepository) GetTextos() (map[string]string, error) {
	coll := m.Client.Database("ticketera").Collection("textos")
	cur, err := coll.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	result := make(map[string]string)
	for cur.Next(context.TODO()) {
		var t domain.Texto
		if err := cur.Decode(&t); err == nil {
			result[t.Clave] = t.Valor
		}
	}
	return result, nil
}

func (m *MongoRepository) GetLogoPath() (string, error) {
	coll := m.Client.Database("ticketera").Collection("logo")
	var l domain.Logo
	err := coll.FindOne(context.TODO(), map[string]interface{}{}).Decode(&l)
	if err != nil {
		return "", err
	}
	return l.FileURL, nil
}
