package storage

import (
	"context"
)

func (m *MongoRepository) DeleteTexto(clave string) error {
	coll := m.Client.Database("ticketera").Collection("textos")
	_, err := coll.DeleteOne(context.TODO(), map[string]interface{}{"clave": clave})
	return err
}
