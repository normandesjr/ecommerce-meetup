package cmd

import "fmt"

type Customer struct {
	Username string `dynamodbav:"Username"`
	Email    string `dynamodbav:"Email"`
	Name     string `dynamodbav:"Name"`
}

type Order struct {
	Id        string `dynamodbav:"OrderId"`
	CreatedAt string `dynamodbav:"CreatedAt"`
	Status    string `dynamodbav:"Status"`
	Amount    string `dynamodbav:"Amount"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Username: %s - Email: %s - Name: %s", c.Username, c.Email, c.Name)
}

func (o Order) String() string {
	return fmt.Sprintf("Id: %s - CreatedAt: %s - Status: %s - Amount: %s", o.Id, o.CreatedAt, o.Status, o.Amount)
}
