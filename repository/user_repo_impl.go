package repository

import (
	"context"
	"fmt"

	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/prisma/db"
)

type UserRepositoryImpl struct {
	Db *db.PrismaClient
}

func NewUserRespository(db *db.PrismaClient) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Db: db,
	}
}

func (p *UserRepositoryImpl) Delete(ctx context.Context, userId string) {
	result, err := p.Db.User.FindUnique(db.User.ID.Equals(userId)).Delete().Exec(ctx)
	if err != nil {
		panic("unimplemented")
	}
	fmt.Println("Rows affected: ", result)
}

func (p *UserRepositoryImpl) Save(ctx context.Context, user model.User) {
	result, err := p.Db.User.CreateOne(
		db.User.Username.Set(user.Username),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		panic("unimplemented")
	}
	fmt.Println("Rows affected: ", result)
}

func (p *UserRepositoryImpl) Update(ctx context.Context, user model.User) {
	result, err := p.Db.User.FindUnique(db.User.ID.Equals(user.Id)).Update(
		db.User.Username.Set(user.Username),
	).Exec(ctx)
	if err != nil {
		panic("unimplemented")
	}
	fmt.Println("Rows affected: ", result)
}

func (p *UserRepositoryImpl) FindById(ctx context.Context, userId string) (model.User, error) {
	// Query the database
	result, err := p.Db.User.FindFirst(db.User.ID.Equals(userId)).Exec(ctx)
	if err != nil {
		fmt.Println("error is", err)
		return model.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Check if result is nil
	if result == nil {
		return model.User{}, fmt.Errorf("user not found")
	}

	// Map database result to the User model
	userData := model.User{
		Id:       result.ID,
		Username: result.Username,
	}

	return userData, nil
}


func (p *UserRepositoryImpl) FindAll(ctx context.Context) []model.User {	
	result, err := p.Db.User.FindMany().Exec(ctx)
	if err != nil {
		panic("unimplemented")
	}
	var users []model.User
	for _, user := range result {
		users = append(users, model.User{
			Id: user.ID,
			Username: user.Username,
		})
	}
	return users
}

