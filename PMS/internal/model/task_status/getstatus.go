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