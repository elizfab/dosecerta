package usecase

import (
	"context"
	"errors"
	"testing"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// mock em memória que satisfaz repository.SettingRepository.
type settingRepoMock struct {
	getFn func(ctx context.Context, key string) (*entity.Setting, error)
	setFn func(ctx context.Context, s *entity.Setting) (*entity.Setting, error)
}

func (r *settingRepoMock) Get(ctx context.Context, key string) (*entity.Setting, error) {
	return r.getFn(ctx, key)
}
func (r *settingRepoMock) Set(ctx context.Context, s *entity.Setting) (*entity.Setting, error) {
	return r.setFn(ctx, s)
}

func TestSettingSet_Success(t *testing.T) {
	var saved *entity.Setting
	repo := &settingRepoMock{
		setFn: func(ctx context.Context, s *entity.Setting) (*entity.Setting, error) {
			saved = s
			return s, nil
		},
	}
	uc := NewSettingUseCase(repo)

	out, err := uc.Set(context.Background(), "weightGoal", "68")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if out.Key != "weightGoal" || out.Value != "68" {
		t.Fatalf("configuração inesperada: %+v", out)
	}
	if saved == nil {
		t.Fatal("repo.Set deveria ter sido chamado")
	}
}

func TestSettingSet_RequiresKey(t *testing.T) {
	repo := &settingRepoMock{
		setFn: func(ctx context.Context, s *entity.Setting) (*entity.Setting, error) {
			t.Fatal("repo não deveria ser chamado com chave vazia")
			return nil, nil
		},
	}
	uc := NewSettingUseCase(repo)
	if _, err := uc.Set(context.Background(), "", "x"); !errors.Is(err, entity.ErrSettingKeyRequired) {
		t.Fatalf("esperava ErrSettingKeyRequired, obtive %v", err)
	}
}

func TestSettingGet_NotFound(t *testing.T) {
	repo := &settingRepoMock{
		getFn: func(ctx context.Context, key string) (*entity.Setting, error) {
			return nil, repository.ErrNotFound
		},
	}
	uc := NewSettingUseCase(repo)
	if _, err := uc.Get(context.Background(), "weightGoal"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound, obtive %v", err)
	}
}
