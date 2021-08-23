package database

import (
	"context"
	"fmt"
	"regexp"
	"servus/_core/errorCollector"
)


type UserObject struct {
	userService
}

type userService interface {
	Create(user User) (id string, err error)
	FindBy(user User) error
	Read(user User)
	Update(user User)
	Delete(user User)
}

type userValidateI interface {
	validate() (bool, string) // true means "I have error"
}


type User struct {
	userValidateI
	ID       uint
	Role     string `validate:"alpha,gte=4,lte=16"`
	Username string `validate:"required,alphanum,gte=4,lte=24"`
	Password string `validate:"required,alphanumunicode,gte=8,lte=64"`
	RegIP    string `validate:"ip,omitempty"`
	RegAgent string `validate:"alphanumunicode,gte=4,lte=128,omitempty"`
}


func (obj *UserObject) Create(user User) (id string, err error) {
	user.validate()
	var columns = []string{"role", "username", "password", "reg_ip", "reg_agent"}
	var data = map[string]string{"role": user.Role, "username": user.Username, "password": user.Password, "reg_ip": user.RegIP, "reg_agent": user.RegAgent}
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

func (user *User) validate() (bool, string){
	var ec = errorCollector.New()
	// username and password start //
	username := user.Username
	password := user.Password
	var usernameIssuer = []string{"username"}
	var passwordIssuer = []string{"password"}
	var failed = 0
	if len(username) < 4 || len(username) > 24{
		ec.AddEValidationMinMax(usernameIssuer, 4, 24)
		failed++
	}
	if len(password) < 8 || len(password) > 64{
		ec.AddEValidationMinMax(passwordIssuer, 8, 64)
		failed++
	}
	if failed == 2 { // if username and password failed
		return true, ec.GetErrors()
	}
	isAlphaNumeric := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !isAlphaNumeric.MatchString(username){
		ec.AddEValidationAllowed(usernameIssuer, "allowed symbols like a-zA-Z0-9", []string{"alphanumeric"})
	}
	isPassword := regexp.MustCompile("^[A-Za-z0-9_@!./#&+*%()-]*$")
	if !isPassword.MatchString(password){
		ec.AddEValidationAllowed(passwordIssuer, "allowed symbols like a-zA-Z0-9 and some special chars", []string{"alphanumeric-unicode"})
	}
	// username and password end //
	var role = user.Role
	var roleIssuer = []string{"role"}
	if len(role) > 0{
		if role != "user" && role != "admin"{
			ec.AddEValidationAllowed(roleIssuer, "allowed only user types", []string{"user", "admin"})
		}
	}
	var regIP = user.RegIP
	if len(regIP) > 0{
		err := validate.Var(regIP, "ip")
		if err != nil {
			ec.AddEValidationInvalid([]string{"RegIP"}, "invalid IP address")
		}
	}

	//var regAgent = user.RegAgent
	if ec.HasErrors(){
		println(ec.GetErrors())
		return true, ec.GetErrors()
	}
	return false, ""
}

//func (obj *UserObject) FindBy(user User) error{
//
//}
