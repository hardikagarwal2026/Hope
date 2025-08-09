package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"time"
)

type UserService interface{
    CreateUser(ctx context.Context, user *db.User) error
    GetUserByID(ctx context.Context, id string) (*db.User, error)
    GetUserByEmail(ctx context.Context, email string) (*db.User, error)
    UpdateUser(ctx context.Context, user *db.User) error
    DeleteUser(ctx context.Context, id string) error
    UpdateLastSeen(ctx context.Context, id string) error
}

type userService struct{
    userRepo repository.UserRepository
}


func NewUserService(userRepo repository.UserRepository) UserService{
    return &userService{userRepo: userRepo}
}

// naya user account banane ke liye
// email and name required
// checks if email already exist in the db
//if not, then it sets the last seen as current time
// and create the user
func(s userService) CreateUser(ctx context.Context, user *db.User) error{
    if user.Email == "" || user.Name == "" {
        return errors.New("name and Email Required")
    }

    existing, _ := s.userRepo.FindByEmail(ctx, user.Email)
    if existing != nil && existing.ID == ""{
        return errors.New("user already exist with this Email")
    }

    user.LastSeen = time.Now()
    return s.userRepo.Create(ctx, user)
}


//to get the user by id
// returns the user and errors(if any)
func(s userService) GetUserByID(ctx context.Context, id string) (*db.User, error){
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil{
        return nil, err
    }

    if user == nil || user.ID == ""{
        return nil, errors.New("user Not Found")
    }
    return user, nil
}


// to get the user info by its emailo
func(s userService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil{
        return nil, err
    }

    if user == nil || user.Email == ""{
        return nil, errors.New("user Not Found")
    } 
    return user, nil
}

// to update the users details like name, email, geohash
func (s userService) UpdateUser(ctx context.Context, user *db.User) error{
    curr, err := s.userRepo.FindByID(ctx, user.ID)
    if curr == nil || err != nil {
        return errors.New("user not foumd")
    }

    if user.Name != "" {
        curr.Name = user.Name
    }

    if user.Email != "" {
        curr.Email = user.Email
    }

    if user.Geohash != "" {
        curr.Geohash = user.Geohash
    }

    return s.userRepo.Update(ctx, curr)
}

// to delete the user from the db
func (s userService) DeleteUser(ctx context.Context, id string) error{
    return s.userRepo.Delete(ctx, id)
}

// to update the user's last seen to time.Now()
func (s userService) UpdateLastSeen(ctx context.Context, id string) error{
    user, err := s.userRepo.FindByID(ctx, id)
    if user == nil || err != nil {
        return errors.New("user Not Found")
    }

    user.LastSeen = time.Now()
    return s.userRepo.Update(ctx, user)
}
