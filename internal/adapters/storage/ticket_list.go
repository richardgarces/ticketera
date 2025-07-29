package storage

import (
	"context"
	"ticketera/internal/domain"
)

func (m *MongoRepository) GetTickets() ([]domain.Ticket, error) {
	coll := m.Client.Database("ticketera").Collection("tickets")
	cur, err := coll.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	var tickets []domain.Ticket
	for cur.Next(context.TODO()) {
		var t domain.Ticket
		if err := cur.Decode(&t); err == nil {
			tickets = append(tickets, t)
		}
	}
	return tickets, nil
}
