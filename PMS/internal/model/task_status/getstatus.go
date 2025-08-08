package status

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
func GetStatusFromString(statusStr string) TaskStatus {
	switch statusStr {
	case "pending":
		return Pending
	case "in progress":
		return InProgress
	case "done":
		return Done
	default:
		return Pending
	}
}
