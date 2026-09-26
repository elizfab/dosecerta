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

// MedicationRepository implementa repository.MedicationRepository sobre MongoDB.
type MedicationRepository struct {
	col *mongo.Collection
}

func NewMedicationRepository(db *mongo.Database) *MedicationRepository {
	return &MedicationRepository{col: db.Collection("medications")}
}

func (r *MedicationRepository) Create(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
	// Deixa o _id em branco para o MongoDB gerar um ObjectID.
	doc := bson.M{
		"name":      m.Name,
		"dosage":    m.Dosage,
		"schedule":  m.Schedule,
		"notes":     m.Notes,
		"active":    m.Active,
		"period":    m.Period,
		"startDate": m.StartDate,
		"endDate":   m.EndDate,
		"stock":     m.Stock,
		"createdAt": m.CreatedAt,
		"updatedAt": m.UpdatedAt,
	}
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("insert medication: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		m.ID = oid.Hex()
	}
	return m, nil
}

func (r *MedicationRepository) FindAll(ctx context.Context) ([]*entity.Medication, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find medications: %w", err)
	}
	defer cur.Close(ctx)

	items := make([]*entity.Medication, 0)
	for cur.Next(ctx) {
		m, err := decodeMedication(cur)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor medications: %w", err)
	}
	return items, nil
}

func (r *MedicationRepository) FindByID(ctx context.Context, id string) (*entity.Medication, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	res := r.col.FindOne(ctx, bson.M{"_id": oid})
	m, err := decodeMedication(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return m, nil
}

func (r *MedicationRepository) Update(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	update := bson.M{"$set": bson.M{
		"name":      m.Name,
		"dosage":    m.Dosage,
		"schedule":  m.Schedule,
		"notes":     m.Notes,
		"active":    m.Active,
		"period":    m.Period,
		"startDate": m.StartDate,
		"endDate":   m.EndDate,
		"stock":     m.Stock,
		"updatedAt": m.UpdatedAt,
	}}
	res := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	updated, err := decodeMedication(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *MedicationRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.ErrInvalidID
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("delete medication: %w", err)
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// singleResult abstrai *mongo.SingleResult e *mongo.Cursor para o decode.
type bsonDecoder interface {
	Decode(v interface{}) error
}

func decodeMedication(d bsonDecoder) (*entity.Medication, error) {
	var doc struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		Dosage    string             `bson:"dosage"`
		Schedule  string             `bson:"schedule"`
		Notes     string             `bson:"notes"`
		Active    bool               `bson:"active"`
		Period    string             `bson:"period"`
		StartDate string             `bson:"startDate"`
		EndDate   string             `bson:"endDate"`
		Stock     int                `bson:"stock"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}
	if err := d.Decode(&doc); err != nil {
		return nil, err
	}
	return &entity.Medication{
		ID:        doc.ID.Hex(),
		Name:      doc.Name,
		Dosage:    doc.Dosage,
		Schedule:  doc.Schedule,
		Notes:     doc.Notes,
		Active:    doc.Active,
		Period:    doc.Period,
		StartDate: doc.StartDate,
		EndDate:   doc.EndDate,
		Stock:     doc.Stock,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}
