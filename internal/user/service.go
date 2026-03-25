package user

import (
	"context"
	"log"

	"github.com/tzincker/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(ctx context.Context, filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error) {
	log.Println("Create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	u, err := s.repo.Create(ctx, &user)

	if err != nil {
		s.log.Println(err)
	}

	return u, err
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	log.Println("Get all users service")

	users, err := s.repo.GetAll(ctx, filters, offset, limit)

	if err != nil {
		s.log.Println(err)
	}

	return users, err
}

func (s service) Get(ctx context.Context, id string) (*domain.User, error) {
	log.Println("Get user service")

	u, err := s.repo.Get(ctx, id)

	if err != nil {
		s.log.Println(err)
	}

	return u, err
}

func (s service) Delete(ctx context.Context, id string) error {
	log.Println("Delete user service")

	err := s.repo.Delete(ctx, id)

	if err != nil {
		s.log.Println(err)
		return err
	}

	return nil
}

func (s service) Update(
	ctx context.Context,
	id string,
	firstName *string,
	lastName *string,
	email *string,
	phone *string,
) error {
	log.Println("Update user service")
	err := s.repo.Update(ctx, id, firstName, lastName, email, phone)
	if err != nil {
		s.log.Println(err)
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int64, error) {
	log.Println("Get all users count service")
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.log.Println(err)
	}

	return count, err
}
