package data

import (
	"context"
	"time"

	v1 "userservice/api/userservice/v1"
	"userservice/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// DB model
type User struct {
	ID           uint64    `gorm:"primaryKey"`
	Name         string
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash string
	CreatedAt    time.Time
}

// userRepo implements biz.UserRepo
type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) biz.UserRepo {
	return &userRepo{db: db}
}

// ---------- CRUD ----------

func (r *userRepo) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password, // already hashed in usecase
		CreatedAt:    time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return toProto(user), nil
}

func (r *userRepo) Get(ctx context.Context, id uint64) (*v1.User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return toProto(&user), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*v1.User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return toProto(&user), nil
}

func (r *userRepo) Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, req.Id).Error; err != nil {
		return nil, err
	}

	// Handle optional fields (*string)
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.PasswordHash = *req.Password // already hashed in usecase
	}

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}
	return toProto(&user), nil
}

func (r *userRepo) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := r.db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *userRepo) List(ctx context.Context) ([]*v1.User, error) {
	var users []User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]*v1.User, 0, len(users))
	for _, u := range users {
		result = append(result, toProto(&u))
	}
	return result, nil
}

// ---------- Helpers ----------

func toProto(u *User) *v1.User {
	return &v1.User{
		Id:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    timestamppb.New(u.CreatedAt),
	}
}
