package Priority

import(
    "fmt"
    "database/sql/driver"
)

type Priority int

const (
    Low Priority = iota
    Medium
    High
)
func (p *Priority) Scan(value interface{}) error {
    intVal, ok := value.(int64) // MySQL TINYINT comes as int64
    if !ok {
        return fmt.Errorf("cannot scan Role from %v", value)
    }
    *p = Priority(intVal)
    return nil
}

func (p Priority) Value() (driver.Value, error) {
    return int64(p), nil
}

// func GetPriority(priority Priority)string{
//       switch priority{
//       case 0:
//         return "Pending"
//       case 1:
//         return "In progress"
//       case 2:
//         return "Done"
//       }
//       return ""
// }
