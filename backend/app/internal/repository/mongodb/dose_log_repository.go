package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dose-certa-backend/internal/domain/entity"
)

// DoseLogRepository implementa repository.DoseLogRepository sobre MongoDB.
type DoseLogRepository struct {
	col *mongo.Collection
}

func NewDoseLogRepository(db *mongo.Database) *DoseLogRepository {
	return &DoseLogRepository{col: db.Collection("dose_logs")}
}

// MarkTaken faz upsert idempotente por (medicationId, date).
func (r *DoseLogRepository) MarkTaken(ctx context.Context, medicationID, date string) (*entity.DoseLog, error) {
	filter := bson.M{"medicationId": medicationID, "date": date}
	update := bson.M{"$setOnInsert": bson.M{
		"medicationId": medicationID,
		"date":         date,
		"createdAt":    time.Now().UTC(),
	}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	res := r.col.FindOneAndUpdate(ctx, filter, update, opts)
	log, err := decodeDoseLog(res)
	if err != nil {
		return nil, fmt.Errorf("mark dose: %w", err)
	}
	return log, nil
}

// Unmark remove o registro (idempotente — não erra se não existir).
func (r *DoseLogRepository) Unmark(ctx context.Context, medicationID, date string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"medicationId": medicationID, "date": date})
	if err != nil {
		return fmt.Errorf("unmark dose: %w", err)
	}
	return nil
}

func (r *DoseLogRepository) FindByDate(ctx context.Context, date string) ([]*entity.DoseLog, error) {
	cur, err := r.col.Find(ctx, bson.M{"date": date})
	if err != nil {
		return nil, fmt.Errorf("find dose logs: %w", err)
	}
	defer cur.Close(ctx)

	items := make([]*entity.DoseLog, 0)
	for cur.Next(ctx) {
		log, err := decodeDoseLog(cur)
		if err != nil {
			return nil, err
		}
		items = append(items, log)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor dose logs: %w", err)
	}
	return items, nil
}

func decodeDoseLog(d bsonDecoder) (*entity.DoseLog, error) {
	var doc struct {
		ID           primitive.ObjectID `bson:"_id"`
		MedicationID string             `bson:"medicationId"`
		Date         string             `bson:"date"`
		CreatedAt    time.Time          `bson:"createdAt"`
	}
	if err := d.Decode(&doc); err != nil {
		return nil, err
	}
	return &entity.DoseLog{
		ID:           doc.ID.Hex(),
		MedicationID: doc.MedicationID,
		Date:         doc.Date,
		CreatedAt:    doc.CreatedAt,
	}, nil
}
