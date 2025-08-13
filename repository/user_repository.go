package repository

import(
	"context" //for timeouts, authentication, cancellation, deadlines
	"hope/db"
	"gorm.io/gorm"
)

//Repository Interface
type UserRepository interface{
	//Methods for CRUD
	Create(ctx context.Context, user *db.User) error
	FindByEmail(ctx context.Context, email string) (*db.User, error)
	FindByID(ctx context.Context, id string)(*db.User, error)
	Update(ctx context.Context, user *db.User)error
	Delete(ctx context.Context, id string) error
}


// implementing repositories through userRepository struct
type userRepository struct{
	db *gorm.DB  //gorm DB object
}


func (r *userRepository) Create(ctx context.Context, user *db.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	var user db.User
	err := r.db.WithContext(ctx).Where("email=?",email).First(&user).Error
	return &user, err
}

func(r *userRepository) FindByID(ctx context.Context, id string) (*db.User, error){
	var user db.User
	err := r.db.WithContext(ctx).Where("id=?", id).First(&user).Error
	return &user, err
}



func (r *userRepository)Update(ctx context.Context, user *db.User) error{
	return r.db.WithContext(ctx).Save(user).Error
}


func(r *userRepository)Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&db.User{},"id = ?", id).Error
}


