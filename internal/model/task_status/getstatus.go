package Status

import "errors"

func GetStatusString(status TaskStatus) string {
	switch status {
	case Pending:
		return "Pending"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}
func GetStatusFromString(statusStr string) (TaskStatus, error) {
	switch statusStr {
	case "Pending":
		return Pending, nil
	case "In Progress":
		return InProgress, nil
	case "Done":
		return Done, nil
	default:
		return Pending, errors.New("invalid status")
	}
}

func GetStatusFromString1(statusStr string) (TaskStatus, error) {
	switch statusStr {
	case "Pending":
		return Pending, nil
	case "In Progress":
		return InProgress, nil
	case "Done":
		return Done, nil
	default:
		return Pending, errors.New("invalid status")
	}
}
