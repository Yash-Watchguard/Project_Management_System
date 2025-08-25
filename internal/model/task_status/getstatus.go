package status

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
	case "pending":
		return Pending, nil
	case "in progress":
		return InProgress, nil
	case "done":
		return Done, nil
	default:
		return Pending, errors.New("invalid status")
	}
}
