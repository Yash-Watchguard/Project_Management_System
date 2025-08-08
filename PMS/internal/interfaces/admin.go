package interfaces

type AdminRepository interface{
   PromoteEmployee(employeeId string)error
}