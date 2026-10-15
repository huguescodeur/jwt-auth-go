package auth

import (
	"context"
	"errors"
	"jwt-auth/internal/user"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       user.Repository
	jwtManager *JWTManager
}

func NewService(r user.Repository, j *JWTManager) *Service {
	return &Service{repo: r, jwtManager: j}
}

func (s *Service) Register(ctx context.Context, u *user.User) (*user.User, error) {
	fullName := strings.TrimSpace(u.FullName)
	email := strings.TrimSpace(u.Email)
	password := strings.TrimSpace(u.PasswordHash)
	if len(fullName) == 0 || len(email) == 0 || len(password) == 0 {
		return nil, user.ErrAllFieldsRequired
	}

	existingUser, err := s.repo.GetByEmail(ctx, u.Email)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, user.ErrAlreadyExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = string(hashedPassword)

	return s.repo.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, u *user.User) (string, *user.User, error) {
	email := strings.TrimSpace(u.Email)
	password := strings.TrimSpace(u.PasswordHash)
	if len(email) == 0 || len(password) == 0 {
		return "", nil, user.ErrAllFieldsRequired
	}

	existingUser, err := s.repo.GetByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {

			return "", nil, user.ErrInvalidCredential
		}

		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(u.PasswordHash)); err != nil {
		return "", nil, user.ErrInvalidCredential
	}

	//? Generate jwt
	token, err := s.jwtManager.Generate(existingUser.ID)
	if err != nil {
		return "", nil, err
	}

	return token, existingUser, nil
}
