package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// WeightRepository implementa repository.WeightRecordRepository sobre MongoDB.
// Guarda peso + circunferências corporais ("medidas") por data.
type WeightRepository struct {
	col *mongo.Collection
}

func NewWeightRepository(db *mongo.Database) *WeightRepository {
	return &WeightRepository{col: db.Collection("weight_records")}
}

// fields monta o documento/$set com todos os campos de medida.
func measurementFields(w *entity.WeightRecord) bson.M {
	return bson.M{
		"date":         w.Date,
		"weightKg":     w.WeightKg,
		"waist":        w.Waist,
		"hip":          w.Hip,
		"glute":        w.Glute,
		"chest":        w.Chest,
		"armRight":     w.ArmRight,
		"armLeft":      w.ArmLeft,
		"forearmRight": w.ForearmRight,
		"forearmLeft":  w.ForearmLeft,
		"thighRight":   w.ThighRight,
		"thighLeft":    w.ThighLeft,
		"calfRight":    w.CalfRight,
		"calfLeft":     w.CalfLeft,
		"neck":         w.Neck,
		"shoulders":    w.Shoulders,
		"bodyFatPct":   w.BodyFatPct,
	}
}

func (r *WeightRepository) Create(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	doc := measurementFields(w)
	doc["createdAt"] = w.CreatedAt
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("insert measurement: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		w.ID = oid.Hex()
	}
	return w, nil
}

func (r *WeightRepository) FindAll(ctx context.Context) ([]*entity.WeightRecord, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find measurements: %w", err)
	}
	defer cur.Close(ctx)

	items := make([]*entity.WeightRecord, 0)
	for cur.Next(ctx) {
		w, err := decodeWeight(cur)
		if err != nil {
			return nil, err
		}
		items = append(items, w)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor measurements: %w", err)
	}
	return items, nil
}

func (r *WeightRepository) FindByID(ctx context.Context, id string) (*entity.WeightRecord, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	res := r.col.FindOne(ctx, bson.M{"_id": oid})
	w, err := decodeWeight(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return w, nil
}

func (r *WeightRepository) Update(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	res := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": measurementFields(w)},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	updated, err := decodeWeight(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *WeightRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.ErrInvalidID
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("delete measurement: %w", err)
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func decodeWeight(d bsonDecoder) (*entity.WeightRecord, error) {
	var doc struct {
		ID           primitive.ObjectID `bson:"_id"`
		Date         string             `bson:"date"`
		WeightKg     float64            `bson:"weightKg"`
		Waist        float64            `bson:"waist"`
		Hip          float64            `bson:"hip"`
		Glute        float64            `bson:"glute"`
		Chest        float64            `bson:"chest"`
		ArmRight     float64            `bson:"armRight"`
		ArmLeft      float64            `bson:"armLeft"`
		ForearmRight float64            `bson:"forearmRight"`
		ForearmLeft  float64            `bson:"forearmLeft"`
		ThighRight   float64            `bson:"thighRight"`
		ThighLeft    float64            `bson:"thighLeft"`
		CalfRight    float64            `bson:"calfRight"`
		CalfLeft     float64            `bson:"calfLeft"`
		Neck         float64            `bson:"neck"`
		Shoulders    float64            `bson:"shoulders"`
		BodyFatPct   float64            `bson:"bodyFatPct"`
		CreatedAt    time.Time          `bson:"createdAt"`
	}
	if err := d.Decode(&doc); err != nil {
		return nil, err
	}
	return &entity.WeightRecord{
		ID:           doc.ID.Hex(),
		Date:         doc.Date,
		WeightKg:     doc.WeightKg,
		Waist:        doc.Waist,
		Hip:          doc.Hip,
		Glute:        doc.Glute,
		Chest:        doc.Chest,
		ArmRight:     doc.ArmRight,
		ArmLeft:      doc.ArmLeft,
		ForearmRight: doc.ForearmRight,
		ForearmLeft:  doc.ForearmLeft,
		ThighRight:   doc.ThighRight,
		ThighLeft:    doc.ThighLeft,
		CalfRight:    doc.CalfRight,
		CalfLeft:     doc.CalfLeft,
		Neck:         doc.Neck,
		Shoulders:    doc.Shoulders,
		BodyFatPct:   doc.BodyFatPct,
		CreatedAt:    doc.CreatedAt,
	}, nil
}
