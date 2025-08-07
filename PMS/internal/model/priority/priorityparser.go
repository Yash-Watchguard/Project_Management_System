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