package status

type TaskStatus int

const(
	Pending TaskStatus= iota
	inProgress 
	Done
)