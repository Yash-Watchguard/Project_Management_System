package repository

import (
	"encoding/json"

	"errors"
	"fmt"

	"os"
	"sync"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

var (
	userFile = "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/user.json"
	mu       sync.Mutex
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (repo *UserRepo) SaveUser(newUser *user.User) error {
	mu.Lock()
	defer mu.Unlock()

	var users []user.User

	// Read existing users
	data, err := os.ReadFile(userFile)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &users)
		if err != nil {
			fmt.Println(" Error parsing JSON:", err)
		}
	}

	// Append new user
	users = append(users, *newUser)

	// Marshal and write
	out, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) IsUserPresent(name, email, password string) (*user.User, error) {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Name == name && u.Email == email {
			// compare stored hashed password with plain password
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err == nil {
				return &u, nil
			}
		}
	}

	return nil, errors.New("invalid details")
}
func (repo *UserRepo) ViewProfile(userId string) error {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return errors.New("error in readfile")
	}
	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return errors.New("error in unmarshal")
	}
	for _, user := range users {
		if user.Id == userId {
			color.Cyan("----------- %s Profile -----------", user.Role)
			color.Yellow("ID     : %d", user.Id)
			color.Yellow("Name   : %s", user.Name)
			color.Yellow("Email  : %s", user.Email)
			color.Yellow("Role   : %s", user.Role)
			color.Cyan("-------------------------------------")
			fmt.Println("Press ENTER to return to dashboard...")
			fmt.Scanln()
		}
	}
	return nil
}

func (ur *UserRepo) GetAllUsers() []user.User {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return []user.User{}
	}

	var users []user.User
	_ = json.Unmarshal(data, &users)
	return users
}

func (ur *UserRepo) DeleteUserById(userId string) error {
	data, err := os.ReadFile(userFile)

	if err != nil {
		return err
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	flag := false
	newUsers := []user.User{}
	for _, user := range users {
		if user.Id != userId {
			newUsers = append(newUsers, user)
		} else {
			flag = true
		}
	}

	if !flag {
		return errors.New("user not found")
	}

	updatedData, err := json.MarshalIndent(newUsers, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
func (ur *UserRepo) GetAllManager() error {
	Data, err := os.ReadFile(userFile)

	if err != nil {
		return err
	}

	var users []user.User
	_ = json.Unmarshal(Data, &users)
	counter := 1
	flag := false
	for _, user := range users {
		if user.Role == roles.Manager { // TODO: Use enum
			flag = true
			fmt.Printf("%d. Name :%s , UserId : %s \n", counter, user.Name, user.Id)
			counter++
		}
	}
	if !flag {
		return errors.New("no manager found")
	}
	return nil
}
