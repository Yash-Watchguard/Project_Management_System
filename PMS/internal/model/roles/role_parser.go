package roles

func RoleParser(role Role)(string){
    switch role{
	case 0:
		return "Admin"
	case 1:
		return "Manager"
	case 2:
		return "Employee"
	}
	return ""
}