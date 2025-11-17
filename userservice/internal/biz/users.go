package biz

import (
	"context"
	"errors"
	"time"

	v1 "userservice/api/userservice/v1"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key") // ⚠️ In production, load from config/env

// ---------------- Repo Interface ----------------
// Must match what data/user.go implements
type UserRepo interface {
	Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error)
	Get(ctx context.Context, id uint64) (*v1.User, error)
	GetByEmail(ctx context.Context, email string) (*v1.User, error)
	Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	List(ctx context.Context) ([]*v1.User, error)
}

// ---------------- Usecase ----------------
type UserUsecase struct {
	repo UserRepo
	
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// ---------- Registration ----------
func (uc *UserUsecase) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(hashedPassword) // store hashed password
	return uc.repo.Create(ctx, req)
}

// ---------- Login ----------
func (uc *UserUsecase) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.AuthReply, error) {
	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// JWT claims
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expires in 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &v1.AuthReply{
		Token: tokenString,
		User:  user,
	}, nil
}

// ---------- Other CRUD ----------
func (uc *UserUsecase) Get(ctx context.Context, id uint64) (*v1.User, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *UserUsecase) Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	// re-hash password if provided
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		*req.Password = string(hashedPassword)
	}
	return uc.repo.Update(ctx, req)
}

func (uc *UserUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	return uc.repo.Delete(ctx, id)
}

func (uc *UserUsecase) List(ctx context.Context) ([]*v1.User, error) {
	return uc.repo.List(ctx)
}
