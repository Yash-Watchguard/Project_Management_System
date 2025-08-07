package interfaces

import (
	

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

type ManagerRepository interface {
	ViewAllEmployee()([]user.User,error)
}

