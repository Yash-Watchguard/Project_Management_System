package roles

import (
	"database/sql/driver"
	"fmt"
)

func RoleParser(role Role)(string){
    switch role{
	case 0:
		return "Admin"
	case 1:
		return "Manager"
	case 2:
		return "Employee"
	}
	return ""
}
func RoleParserStringToRole(role string)(Role){
    switch role{
	case "Admin":
		return Role(0)
	case "Manager":
		return Role(1)
	case "Employee":
		return Role(2)
	}
	return 8
}
func (r *Role) Scan(value interface{}) error {
    intVal, ok := value.(int64) // MySQL TINYINT comes as int64
    if !ok {
        return fmt.Errorf("cannot scan Role from %v", value)
    }
    *r = Role(intVal)
    return nil
}

func (r Role) Value() (driver.Value, error) {
    return int64(r), nil
}