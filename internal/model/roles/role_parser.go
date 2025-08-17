package roles
import(
	"fmt"
	"database/sql/driver"
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