package database

import (
	"context"
	"fmt"
	"strings"
)

// UserObject logic start
type UserObject struct {
	userService
}

type userService interface {
	Create(data map[string]string) (id string, err error)
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

func (obj *UserObject) Create(data map[string]string) (id string, err error) {
	columns := []string{"role", "username", "password", "reg_ip", "reg_agent"}
	var builded = queryBuilder(data, columns)
	var qu = fmt.Sprintf("insert into users (%v) values (%v) RETURNING id", builded.fields, builded.dollars)
	lastInsertId := "error"
	err = pConnection.QueryRow(context.Background(), qu, builded.values...).Scan(&lastInsertId)
	err = errorHandler(err)
	if err != nil {
		println(err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

//func (obj *UserObject) FindBy(user User) error{
//
//}

// User logic end
type queryBuilded struct {
	fields string
	dollars string
	values []interface{}
}

func queryBuilder(userData map[string]string, columns []string) queryBuilded{
	var fields []string
	var dollars []string
	var values []interface{}
	i := 1
	for _, k := range columns {
		if _, ok := userData[k]; ok {
			values = append(values, userData[k])
			dollars = append(dollars, fmt.Sprintf("$%v", i))
			fields = append(fields, fmt.Sprintf("%s", k))
			i++
		}
	}
	return queryBuilded{fields: strings.Join(fields,", "), dollars: strings.Join(dollars,", "), values: values}
}
