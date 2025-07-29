package storage

import (
	"context"
	"ticketera/internal/domain"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) SaveTicket(t *domain.Ticket) error {
	coll := m.Client.Database("ticketera").Collection("tickets")
	// Obtener el último número correlativo
	var last domain.Ticket
	opts := options.FindOne().SetSort(map[string]interface{}{"numero": -1})
	_ = coll.FindOne(context.TODO(), map[string]interface{}{}, opts).Decode(&last)
	t.Numero = last.Numero + 1
	_, err := coll.InsertOne(context.TODO(), t)
	return err
}
