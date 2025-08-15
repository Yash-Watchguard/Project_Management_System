package status

type TaskStatus int

const(
	Pending TaskStatus= iota
	InProgress 
	Done
)