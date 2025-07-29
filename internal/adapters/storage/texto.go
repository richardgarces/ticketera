package storage

import (
	"context"
	"ticketera/internal/domain"
)

func (m *MongoRepository) SaveTexto(clave, valor string) error {
	coll := m.Client.Database("ticketera").Collection("textos")
	_, err := coll.InsertOne(context.TODO(), domain.Texto{Clave: clave, Valor: valor})
	return err
}
