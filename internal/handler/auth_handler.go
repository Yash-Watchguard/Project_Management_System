package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)
var GenerateUUID = util.GenerateUniqueUUID
var ValidEmail = util.ValidateEmail
var ValidPhoneNumber = util.ValidateMobileNumber

var TimeParser = util.ParseDate
// allow overriding JWT generator in tests
var GenerateJwt = util.GenerateJwt

var validate *validator.Validate
type authHandler struct{
	authService service1.AuthServiceInterface
	userService service1.UserServiceInterface
}
func NewAuthHandler(authService service1.AuthServiceInterface,userService service1.UserServiceInterface)*authHandler{
	return &authHandler{authService: authService,userService: userService}
}
func init(){
	validate=validator.New()
}
func(au *authHandler)Signup(w http.ResponseWriter,r *http.Request){
    // dto at client side

	if r.Method != http.MethodPost {
        logger.Error("Invalid HTTP method")
        response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
        return
    }
	var newUser struct {
		Name           string `json:"name" validate:"required"`
		Email          string `json:"email" validate:"required,email"`
		Password       string `json:"password" validate:"required,min=8"`
		PhoneNumber    string  `json:"phonenumber" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		logger.Error("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}
	

    
	flag :=au.userService.CheckUserExist(newUser.Email)
	if flag{
		logger.Error("User already exists")
		response.ErrorResponse(w, http.StatusConflict, "User already exists", 1009)
		return
	}

	// varify the password
	
	if err := util.ValidatePassword(newUser.Password); err != nil {
		logger.Error("Invalid password")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid password", 1005)
		return
	}

	hhashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error hashing password", 1006)
		//http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser.Password=string(hhashedPassword)

	// save the user

	NewUser:=&user.User{
		Id: GenerateUUID(),
        Name: newUser.Name,
		Email: newUser.Email,
		Password: newUser.Password,
        Role: roles.Employee,
		PhoneNumber: newUser.PhoneNumber,
	}

	err =au.authService.Signup(NewUser)

	if err != nil {
		logger.Error("Error creating user")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error creating user", 1006)
		return
	}
    userDto:=model.UserDto{
		Id: NewUser.Id,
		Name: NewUser.Name,
		Email: NewUser.Email,
		PhoneNumber: NewUser.PhoneNumber,
		Role: roles.RoleParser(NewUser.Role),
	}

	logger.Info("User created sucessfully")
	response.SuccessResponse(w, userDto, "User created successfully", http.StatusCreated)

}

func(au * authHandler)Login(w http.ResponseWriter,r * http.Request)  {
	if r.Method != http.MethodGet {
        logger.Error("Invalid HTTP method")
        response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
        return
    }

	var userInput struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

	

    var err error

	if err=json.NewDecoder(r.Body).Decode(&userInput);err!=nil{
		logger.Error("Invalid Input")
		response.ErrorResponse(w,http.StatusBadRequest,"Invalid input",1001)
		return 
	}

	// validate the userinput

	if err=validate.Struct(userInput); err!=nil{
		logger.Error("Validation error")
		response.ErrorResponse(w,http.StatusBadRequest,"Invalid Request Body",1001)
		return
	}

	var person *user.User
       
	person,err =au.userService.IsUserPresent(userInput.Name,userInput.Email,userInput.Password)

	person.Status=user.Active

	if err!=nil{
		logger.Error("Invalid email or password")
		response.ErrorResponse(w,http.StatusUnauthorized,"Invalid email or password",1005)

		return
	}

	// generate JWT token
    
	var jwtTokenString string
	jwtTokenString, err = GenerateJwt(person.Id, person.Role)

	if err!=nil{
		logger.Error("Error generating the token")
		response.ErrorResponse(w,http.StatusInternalServerError,"Error generating the token",1006)

		return
	}

	// return the jwt token as json
	logger.Info(fmt.Sprintf("token generated for userId:%v",person.Id))
	response.SuccessResponse(w,map[string]interface{}{"token":jwtTokenString,"UserId":person.Id},"Token generated Successfully",http.StatusCreated)

}

