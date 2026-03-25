package user

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tzincker/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(ctx context.Context, filters Filters) (int64, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := repo.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("user created with id: ", user.ID)
	return user, nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	tx := repo.db.WithContext(ctx).Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&users)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("users got")
	return users, nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{ID: id}

	result := repo.db.WithContext(ctx).First(&user)
	if result.Error != nil {
		repo.log.Println(result.Error)

		if result.Error == gorm.ErrRecordNotFound {
			return nil, &ErrNotFound{UserId: id}
		}

		return nil, result.Error
	}
	repo.log.Println("user found with id: ", user.ID)
	return &user, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	user := domain.User{ID: id}

	result := repo.db.WithContext(ctx).Delete(&user)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("user %s doesn't exists", id)
		return &ErrNotFound{UserId: id}
	}

	repo.log.Println("user deleted with id: ", user.ID)
	return nil
}

func (repo *repo) Update(
	ctx context.Context,
	id string,
	firstName *string,
	lastName *string,
	email *string,
	phone *string,
) error {

	values := make(map[string]any)

	if firstName != nil {
		values["first_name"] = firstName
	}

	if lastName != nil {
		values["last_name"] = lastName
	}

	if email != nil {
		values["email"] = email
	}

	if phone != nil {
		values["phone"] = phone
	}

	result := repo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("user %s doesn't exists", id)
		return &ErrNotFound{UserId: id}
	}

	repo.log.Println("user updated with id: ", id)
	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int64, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	repo.log.Println("users count")
	return count, nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}
