package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// SettingRepository implementa repository.SettingRepository sobre MongoDB.
// Usa a própria chave como _id (single-user), então cada chave é única.
type SettingRepository struct {
	col *mongo.Collection
}

func NewSettingRepository(db *mongo.Database) *SettingRepository {
	return &SettingRepository{col: db.Collection("settings")}
}

func (r *SettingRepository) Get(ctx context.Context, key string) (*entity.Setting, error) {
	res := r.col.FindOne(ctx, bson.M{"_id": key})
	var doc struct {
		Key   string `bson:"_id"`
		Value string `bson:"value"`
	}
	if err := res.Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get setting: %w", err)
	}
	return &entity.Setting{Key: doc.Key, Value: doc.Value}, nil
}

func (r *SettingRepository) Set(ctx context.Context, s *entity.Setting) (*entity.Setting, error) {
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": s.Key},
		bson.M{"$set": bson.M{"value": s.Value}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return nil, fmt.Errorf("set setting: %w", err)
	}
	return s, nil
}
