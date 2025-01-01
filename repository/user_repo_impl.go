package repository

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
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

var secretKey = []byte("hjgdsajfsdakfhasdhfiao@!#!@$!$231231")

func (p *UserRepositoryImpl) Delete(ctx context.Context, userId string) {
	result, err := p.Db.User.FindUnique(db.User.ID.Equals(userId)).Delete().Exec(ctx)
	if err != nil {
		panic("unimplemented")
	}
	fmt.Println("Rows affected: ", result)
}

func (p *UserRepositoryImpl) Save(ctx context.Context, user model.User) (string, error) {
	result, err := p.Db.User.CreateOne(
		db.User.Username.Set(user.Username),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		panic("unimplemented")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": result.ID,
	})
	return token.SignedString(secretKey)
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
			Id:       user.ID,
			Username: user.Username,
		})
	}
	return users
}

func FindById(userId string, ctx context.Context, p *db.PrismaClient) (model.User, error) {
	result, err := p.User.FindFirst(db.User.ID.Equals(userId)).Exec(ctx)
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

// AddRoomMember adds a user to the members of a specified room. If the room
// does not exist, it creates the room first. It uses the roomName to find 
// the room and appends the userName to the room's members. Any errors 
// encountered during the process are logged to the console.
// func AddRoomMember(userName string, roomName string, ctx context.Context, p *db.PrismaClient) {
// 	result, err := p.Room.FindFirst(db.Room.Name.Equals(roomName)).Exec(ctx)
// 	if err != nil {
// 		fmt.Println("error is", err)
// 	}
// 	if result == nil {
// 		result, err = p.Room.CreateOne(
// 			db.Room.Name.Set(roomName),
// 		).Exec(ctx)
// 		fmt.Println("room created")
// 		if err != nil {
// 			fmt.Println("failed to create room", err)
// 		}
// 	}
// 	result, err = p.Room.FindUnique(
// 		db.Room.ID.Equals(result.ID), // Find the room by its unique ID
// 	).Update(
// 		db.Room.Members.Push([]string{userName}), // Append the new userId to the Members array
// 	).Exec(ctx)
// 	if err != nil {
// 		fmt.Println("error is", err)
// 		panic("unimplemented")
// 	}
// 	fmt.Println("Rows affected: ", result)
// }

// func RemoveMember(userName string, roomName string, ctx context.Context, p *db.PrismaClient) {
// 	result, err := p.Room.FindFirst(
// 		db.Room.Name.Equals(roomName), // Find the room by its unique ID
// 	).Exec(ctx)
	
// 	if err != nil {
// 		// Handle the error
// 		fmt.Printf("Error finding room: %v\n", err)
// 		return
// 	}
	
// 	// Remove the target userID from the Members array
// 	updatedMembers := []string{}
// 	for _, member := range result.Members {
// 		if member != userName {
// 			updatedMembers = append(updatedMembers, member)
// 		}
// 	}
	
// 	// Update the Members array in the database
// 	result, err = p.Room.FindUnique(
// 		db.Room.ID.Equals(result.ID),
// 	).Update(
// 		db.Room.Members.Set(updatedMembers), // Replace the Members array with the updated array
// 	).Exec(ctx)
	
// 	if err != nil {
// 		// Handle the error
// 		fmt.Printf("Error updating room members: %v\n", err)
// 		return
// 	}
	
// 	fmt.Println("Updated room members:", result)
	

// }
