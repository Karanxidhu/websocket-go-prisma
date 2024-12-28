package service

import (
	"context"
	"fmt"

	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/data/response"
	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/repository"
)

type UserRespositoryImpl struct {
	UserRespository repository.UserRepository
}

func NewUserServiceImpl(userRespository repository.UserRepository) *UserRespositoryImpl {
	return &UserRespositoryImpl{UserRespository: userRespository}
}

func (p *UserRespositoryImpl) Save(ctx context.Context, user request.UserCreateReq)(string){
	userData := model.User{
		Username: user.Username,
	}
	fmt.Println(userData)
	token , err := p.UserRespository.Save(ctx, userData)
	if err != nil {	
		panic("unimplemented")
	}
	return token
}

func (p * UserRespositoryImpl) Update(ctx context.Context, user request.UserUpdateReq){
	userData := model.User{
		Id: user.UserId,
		Username: user.Username,
	}
	p.UserRespository.Update(ctx, userData)
}

func (p * UserRespositoryImpl) Delete(ctx context.Context, userId string){
	result, err := p.UserRespository.FindById(ctx, userId)
	if err != nil {
		panic("unimplemented")
	}
	p.UserRespository.Delete(ctx, result.Id)
}

func (p * UserRespositoryImpl) FindById(ctx context.Context, userId string) response.UserResponse{
	result, err := p.UserRespository.FindById(ctx, userId)
	if err != nil {
		fmt.Println(err)
		panic("unimplemented")
	}
	return response.UserResponse{
		Id: result.Id,
		UserName: result.Username,
	}
}

func (p * UserRespositoryImpl) FindAll(ctx context.Context) []response.UserResponse{
	user := p.UserRespository.FindAll(ctx)
	var users []response.UserResponse
	for _, u := range user {
		users = append(users, response.UserResponse{
			Id: u.Id,
			UserName: u.Username,
		})
	}
	return users
}