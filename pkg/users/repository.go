package users

import (
	"context"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	Admin  = "admin"
	Member = "member"
)

type ErrEmail string

func (e ErrEmail) Error() string {
	return string(e) + " has already been registered"
}

// Validate func for validating user input
func (u *UserDTO) Validate() error {
	return ozzo.ValidateStruct(u,
		ozzo.Field(&u.FirstName, ozzo.Required),
		ozzo.Field(&u.LastName, ozzo.Required),
		ozzo.Field(&u.EmailAdress, ozzo.Required, is.Email),
		ozzo.Field(&u.Password, ozzo.Required, is.Alphanumeric, ozzo.Length(8, 20)),
		ozzo.Field(&u.IsAdmin, ozzo.Required, ozzo.In(Admin, Member)),
	)
}

// Repository for mongo repository
type Repository struct {
	coll *mongo.Collection
}

// NewRepository helps to create new mongo repository
func NewRepository(ctx context.Context, db *mongo.Database) (*Repository, error) {
	coll := db.Collection("users")

	indexes := []mongo.IndexModel{
		{Keys: bson.M{"email_address": 1}, Options: options.Index().SetUnique(true)},
	}

	_, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}

	return &Repository{coll}, nil
}

func (r *Repository) Create(ctx context.Context, dto UserDTO) (*User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:          uuid.New().String(),
		CreadAt:     time.Now(),
		UpdatedAt:   time.Now(),
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		EmailAdress: dto.EmailAdress,
		IsAdmin:     dto.IsAdmin,
		Password:    password,
	}

	_, err = r.coll.InsertOne(ctx, user)
	if err != nil {
		writeErr, ok := err.(mongo.WriteException)
		if ok && writeErr.WriteErrors[0].Code == 11000 {
			return nil, ErrEmail(dto.EmailAdress)
		} else {
			return nil, err
		}
	}

	return user, err
}
