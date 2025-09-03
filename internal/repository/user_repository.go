package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepoInterface interface{
	 Exec(query string, args ...any) (sql.Result, error)
	 Query(query string, args ...any) (*sql.Rows, error)
	 QueryRow(query string, args ...any) *sql.Row
}

type UserRepo struct{
	db UserRepoInterface
}

func NewUserRepo(db UserRepoInterface) *UserRepo {
	return &UserRepo{db: db}
}


func (repo *UserRepo) SaveUser(newUser *user.User) error {
    columns := []string{"id", "role", "name", "password", "phone_number", "email"}
    query := config.InsertQuery("users", columns)

    _, err := repo.db.Exec(query,
        newUser.Id,
        newUser.Role,
        newUser.Name,
        newUser.Password,
        newUser.PhoneNumber,
        newUser.Email,
    )

    if err != nil {
        if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
            if strings.Contains(mysqlErr.Message, "phone_number") {
                return errors.New("phone number already exists")
            }
            if strings.Contains(mysqlErr.Message, "email") {
                return errors.New("email already exists")
            }
            return errors.New("duplicate entry")
        }
        return err
    }

    return nil
}


func (repo *UserRepo) IsUserPresent(name, email, password string) (*user.User, error) {
	
	row:=repo.db.QueryRow(config.SelectQuery("users",[]string{"id","name","email","password","role","phone_number"},"name","email"),name,email)

	var user user.User

	err:=row.Scan(&user.Id,&user.Name,&user.Email,&user.Password,&user.Role,&user.PhoneNumber)

	if err!=nil{
		if err==sql.ErrNoRows{
			return nil,errors.New("no user found")
		}
		return nil,err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user,nil

}
func (repo *UserRepo) ViewProfile(userId string) ([]user.User, error) {
	
	query := config.SelectQuery("users", []string{"id", "name", "email", "role", "phone_number"}, "id")

	row := repo.db.QueryRow(query, userId)

	var u user.User
	var user []user.User
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Role, &u.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
    user = append(user, u)
	return user, nil
}


func (repo *UserRepo) GetAllUsers() ([]user.User, error) {
	query := `SELECT id, name, email, role, phone_number FROM users`
	// `SELECT id, name, email, role, phone_number FROM users`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Role, &u.PhoneNumber)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		return nil, errors.New("no users present")
	}

	return users, nil
}


func (ur *UserRepo) DeleteUserById(userId string) error {
	// Build delete query
	query := config.DeleteQuery("users", []string{"id"}) // DELETE FROM users WHERE id = ?

	// Execute query
	result, err := ur.db.Exec(query, userId)
	if err != nil {
		return errors.New("please enter valid user id")
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("please enter valid user id")
	}

	return nil
}

func (ur *UserRepo) GetAllManager()([]user.User, error) {
	var managers []user.User

    query := config.SelectQuery("users", []string{"id", "name"}, "role")

	rows, err := ur.db.Query(query,roles.Manager)
	if err != nil {
		return  nil,err
	}
	defer rows.Close()
    counter:=1
	for rows.Next(){
        // var id,name string
        var u user.User
		err:=rows.Scan(&u.Id,&u.Name)
		// color.Yellow("%d. Name: %s, UserId: %s\n", counter, name, id)
        counter++
		if err!=nil{
			return nil,err
		}
		managers = append(managers, u)
	}

	return managers,nil
	
}

func (ur *UserRepo) UpdateProfile(userId string, updates map[string]interface{}) error {

	query:=config.SelectQuery("users",[]string{"id"},"id")
	_,err:=ur.db.Query(query,userId)
	if err!=nil{
		return errors.New("invalid id")
	}
    allowedFields := map[string]bool{
        "name":         true,
        "email":        true,
        "password":     true,
        "phone_number": true,
		"status"      : true,
    }
    setClauses := []string{}
    args := []interface{}{}

    for field, value := range updates {
        if !allowedFields[field] {
            return fmt.Errorf("invalid field update: %s", field)
        }
        setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
        args = append(args, value)
    }

    if len(setClauses) == 0 {
        return errors.New("no valid fields to update")
    }

    query = fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))
    args = append(args, userId)

    result, err := ur.db.Exec(query, args...)
    if err != nil {
        // Handle unique constraint errors
        if strings.Contains(err.Error(), "Duplicate entry") {
            if strings.Contains(err.Error(), "email") {
                return errors.New("email already exists")
            }
            if strings.Contains(err.Error(), "phone_number") {
                return errors.New("phone number already exists")
            }
        }
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return errors.New("user not found")
    }

    return nil
}

func (ur *UserRepo) PromoteEmployee(employeeId string) error {

	query := "UPDATE users SET role = ? WHERE id = ?"
    result, err := ur.db.Exec(query, roles.Manager, employeeId)
	if err != nil {
		return err
	}

	// Check if any row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
func (ur *UserRepo) ViewAllEmployee() ([]user.User, error) {
	// Build SELECT query to get all employees
	query := config.SelectQuery("users",[]string{"id", "name", "email", "role", "phone_number", "password"},"role")

	
	rows, err := ur.db.Query(query, roles.Employee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Role, &u.PhoneNumber, &u.Password)
		if err != nil {
			return nil, err
		}
		employees = append(employees, u)
	}

	if len(employees) == 0 {
		return nil, errors.New("no employees found")
	}

	return employees, nil
}

func(ur *UserRepo)GetUserByEmail(email string)(*user.User,error){
      query := config.SelectQuery("users", []string{"id", "name", "email", "role", "phone_number"}, "email")

	row := ur.db.QueryRow(query, email)

	var u user.User
	var user []user.User
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Role, &u.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
    user = append(user, u)
	return &user[0], nil
}

