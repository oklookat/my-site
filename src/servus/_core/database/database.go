package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"servus/_core/logger"
)

var pLogger *logger.Logger
var pConnection *pgx.Conn

func New(connectionStr string, _logger *logger.Logger) *DB {
	pLogger = _logger
	db, err := pgx.Connect(context.Background(), connectionStr)
	if err != nil {
		pLogger.Error(fmt.Sprintf("%v\n", err))
		os.Exit(1)
	}
	pConnection = db
	return &DB{
		Connection: pConnection,
		User:       UserObject{},
	}
}

type DB struct {
	Connection *pgx.Conn
	User       UserObject
}

// UserObject logic start
type UserObject struct {
	userService
}

type userService interface {
	Create(user User) error
	FindBy(user User) error
	Read(user User)
	Update(user User)
	Delete(user User)
}

type User struct {
	ID       uint
	Role     string
	Username string
	Password string
	RegIP    string
	RegAgent string
}

func (obj *UserObject) Create(user User) error {
	if (User{}) == user {
		return structEmptyErr()
	}
	var qu = fmt.Sprintf("insert into users (role, username, password, reg_ip, reg_agent) values ($1,  $2,  $3,  $4, $5) RETURNING id")
	result, err := pConnection.Exec(context.Background(), qu, user.Role, user.Username, user.Password, user.RegIP, user.RegAgent)
	println(result)
	if !checkErr(err){
		return err
	}
	return nil
}
func (obj *UserObject) FindBy(user User) error{
	if (User{}) == user {
		return structEmptyErr()
	}
	for _, d:= range user {
		//do something with the d
	}
}

// User logic end

// Service start
func structEmptyErr() error {
	return errors.New("STRUCT_EMPTY")
}

func checkErr(err error) bool {
	if err != nil {
		pLogger.Error(fmt.Sprintf("%v\n", err))
		return false
	} else {
		return true
	}
}
// Service end
