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

// ExamRepository implementa repository.ExamRepository sobre MongoDB.
type ExamRepository struct {
	col *mongo.Collection
}

func NewExamRepository(db *mongo.Database) *ExamRepository {
	return &ExamRepository{col: db.Collection("exams")}
}

func examFields(e *entity.Exam) bson.M {
	return bson.M{
		"name":         e.Name,
		"value":        e.Value,
		"unit":         e.Unit,
		"referenceMin": e.ReferenceMin,
		"referenceMax": e.ReferenceMax,
		"date":         e.Date,
		"notes":        e.Notes,
		"updatedAt":    e.UpdatedAt,
	}
}

func (r *ExamRepository) Create(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
	doc := examFields(e)
	doc["createdAt"] = e.CreatedAt
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("insert exam: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		e.ID = oid.Hex()
	}
	return e, nil
}

func (r *ExamRepository) FindAll(ctx context.Context) ([]*entity.Exam, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find exams: %w", err)
	}
	defer cur.Close(ctx)

	items := make([]*entity.Exam, 0)
	for cur.Next(ctx) {
		e, err := decodeExam(cur)
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor exams: %w", err)
	}
	return items, nil
}

func (r *ExamRepository) FindByID(ctx context.Context, id string) (*entity.Exam, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	res := r.col.FindOne(ctx, bson.M{"_id": oid})
	e, err := decodeExam(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *ExamRepository) Update(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrInvalidID
	}
	res := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": examFields(e)},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	updated, err := decodeExam(res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *ExamRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.ErrInvalidID
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("delete exam: %w", err)
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func decodeExam(d bsonDecoder) (*entity.Exam, error) {
	var doc struct {
		ID           primitive.ObjectID `bson:"_id"`
		Name         string             `bson:"name"`
		Value        float64            `bson:"value"`
		Unit         string             `bson:"unit"`
		ReferenceMin float64            `bson:"referenceMin"`
		ReferenceMax float64            `bson:"referenceMax"`
		Date         string             `bson:"date"`
		Notes        string             `bson:"notes"`
		CreatedAt    time.Time          `bson:"createdAt"`
		UpdatedAt    time.Time          `bson:"updatedAt"`
	}
	if err := d.Decode(&doc); err != nil {
		return nil, err
	}
	return &entity.Exam{
		ID:           doc.ID.Hex(),
		Name:         doc.Name,
		Value:        doc.Value,
		Unit:         doc.Unit,
		ReferenceMin: doc.ReferenceMin,
		ReferenceMax: doc.ReferenceMax,
		Date:         doc.Date,
		Notes:        doc.Notes,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}, nil
}
