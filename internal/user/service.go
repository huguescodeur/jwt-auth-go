package user

import "context"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}
