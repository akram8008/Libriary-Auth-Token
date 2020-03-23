package users

import (
	"auth/pkg/crud/models"
	"context"
	"errors"
	"fmt"
	jwt "github.com/akram8008/jwt/pkg/cmd"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UsersSvc struct {
	secret jwt.Secret
	pool *pgxpool.Pool
}


func NewUserSvc(secret jwt.Secret,pool *pgxpool.Pool) *UsersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil")) // <- be accurate
	}
	return &UsersSvc{secret: secret, pool: pool}
}


var ErrInvalidLoginOrPassword = errors.New("login or password is wrong")


func (service *UsersSvc) AddNewUser(ctx context.Context, model models.User) (err error) {
	log.Print("Adding a new user")
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't execute pool: %v",err)
		return errors.New(fmt.Sprintf("can't execute pool: %v", err))
	}
	defer conn.Release()

	log.Print("Checking login for not repeating")
	rows, err := conn.Query(ctx, selectLoginById,&model.Login)
	if err != nil {
		return errors.New(fmt.Sprintf("can't execute a querry: %v", err))
	}
	defer rows.Close()
	log.Print("Login is repeating")
	if rows.Next() {
		return errors.New("login is repeating")
	}


	log.Printf("register a  new user ")
	password, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	_, err = conn.Exec(ctx, insertInUsers,model.Name, model.Login,password, model.Role)

	if err != nil {
		log.Printf("can't register a  new user ")
		return errors.New(fmt.Sprintf("can't save a new user: %v ", err))
	}
	log.Printf("new user successufuly added")
	return nil
}

func (service *UsersSvc)Login(ctx context.Context, model models.User) (response ResponseDTO, err error) {
	var pass,name string
	var id int64
	err = service.pool.QueryRow(ctx, selectUserByLogin, &model.Login).Scan(&id,&name,&pass)
	if err != nil {
		return ResponseDTO{}, ErrInvalidLoginOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(pass))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Printf("Error password %s %s",model.Password,pass)
		return ResponseDTO{}, ErrInvalidLoginOrPassword
	}

	response.Token, err = jwt.Encode(Payload{Id: id,  Exp:   time.Now().Add(time.Hour).Unix(), Roles: []string{"ROLE_USER"}, 	}, service.secret)
    response.Id = id
    response.Name = name
	return response,nil
}

