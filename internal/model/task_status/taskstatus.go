package Status
import(
	"fmt"
	"database/sql/driver"
)
type TaskStatus int


const(
	Pending TaskStatus= iota
	InProgress 
	Done
)

func (Ts *TaskStatus) Scan(value interface{}) error {
    intVal, ok := value.(int64) // MySQL TINYINT comes as int64
    if !ok {
        return fmt.Errorf("cannot scan Role from %v", value)
    }
    *Ts = TaskStatus(intVal)
    return nil
}

func (Ts TaskStatus) Value() (driver.Value, error) {
    return int64(Ts), nil
}