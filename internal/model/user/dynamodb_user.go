package user

type DynamoDbUserResponse struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	Id          string `dynamodbav:"Id"`
	Name        string `dynamodbav:"Name"`
	Role        string `dynamodbav:"Role"`
	Password    string `dynamodbav:"Password"`
	Status      string `dynamodbav:"Status"`
	Phonenumber string `dynamodbav:"Phonenumber"`
}
