package Priority
import ("errors")

func PriorityParser(input string) (Priority, error) {
	switch input {
	case "Low":
		return Low, nil
	case "Medium":
		return Medium, nil
	case "High":
		return High, nil
	default:
		return 0, errors.New("invalid priority")
	}
}
func GetPriority(priority Priority) (string) {
	switch priority {
	case 0:
		return "Low"
	case 1:
		return "Medium"
	case 2:
		return "High"
	}
	return  ""
}