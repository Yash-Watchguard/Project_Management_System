package handler

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)
var GenerateUUID = util.GenerateUniqueUUID
var ValidEmail = util.ValidateEmail
var ValidPhoneNumber = util.ValidateMobileNumber
// var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)
var TimeParser = util.ParseDate

// func SetInputReader(r io.Reader) {
// 	inputReader = bufio.NewReader(r)
// }

type authHandler struct{
	userService *service1.AuthService
}
func NewAuthHandler(userService *service1.AuthService)*authHandler{
	return &authHandler{userService: userService}
}
func(au *authHandler)Signup() error {
	
	 
		name ,err:=GetInput("Enter name : ")
	    if err != nil {
	    return err
	   }

	var email string
	for {
		email, err = GetValidEmail()
		if err != nil {
			color.Red("Please Enter Valid Email Address : ")
			continue
		} else {
			break
		}

	}

	var password string
	for {
		password, err = GetValidPassword()
		if err != nil {
			color.Red("%v (Abccgg@#12345)",err)
			continue
		} else {
			break
		}
	}

	var phoneNumber string
	for {
		phoneNumber, err = GetValidPhoneNumber()
		if err != nil {
			color.Red("Please Enter Valid Phone number : ")
			continue
		} else {
			break
		}

	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := user.User{
		Id:          GenerateUUID(),
		Name:        name,
		Email:       email,
		Password:    string(hashedPassword),
		PhoneNumber: phoneNumber,
		Role:        2,
	}

	if err := au.userService.Signup(&user); err != nil {
		color.Red("Signup failed: %v", err)
		return err
	}

	color.Green("✅ User signed up successfully!")
	return nil
}

func GetInput(prompt string) (string, error) {
	fmt.Print(color.RedString(prompt))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}


func GetValidEmail() (string, error) {
	email, err := GetInput("Enter Email : ")
	if err != nil {
		return "", err
	}
	if err := ValidEmail(email); err != nil {

		return "", err
	}
	return email, nil
}

func GetValidPhoneNumber() (string, error) {
	number, err := GetInput("Enter Phone Number : ")
	if err != nil {
		return "", err
	}
	number = strings.TrimSpace(number)
	if err := ValidPhoneNumber(number); err != nil {
		return "", err
	}
	return number, nil
}

func GetValidPassword() (string, error) {
	fmt.Print(color.RedString("Enter Password : "))

	
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() 
	if err != nil {
		return "", err
	}

	password := string(bytePassword)

	if err := util.ValidatePassword(password); err != nil {
		return "", err
	}

	return password, nil
}

func(au * authHandler)Login() (*user.User,error) {

	
	name, err := GetInput("Enter name : ")
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}
	mailId, err := GetValidEmail()
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}
	password, err := GetValidPassword()
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}

	

	user, err := au.userService.Login(name, mailId, password)
	if err != nil {
		color.Red("----------Invalid details,Login Faild----------------")
		return nil,err
	}
	// 
	return user,err
}
