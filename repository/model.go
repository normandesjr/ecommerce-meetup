package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Customer struct {
	Username  string             `dynamodbav:"username"`
	Email     string             `dynamodbav:"email"`
	Name      string             `dynamodbav:"name"`
	Addresses map[string]Address `dynamodbav:"addresses"`
}

type Address struct {
	Id            string `dynamodbav:"id"`
	StreetAddress string `dynamodbav:"streetAddress"`
	ZipCode       string `dynamodbav:"zipCode"`
}

type Order struct {
	Id        string `dynamodbav:"orderId"`
	CreatedAt string `dynamodbav:"createdAt"`
	Status    string `dynamodbav:"status"`
	Amount    string `dynamodbav:"amount"`
}

type OrderItem struct {
	Id          string `dynamodbav:"itemId"`
	Description string `dynamodbav:"description"`
	Price       string `dynamodbav:"price"`
}

func (customer Customer) GetKey() map[string]types.AttributeValue {
	//TODO: Melhorar essa constante CUSTOMER#
	pk, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", customer.Username))
	if err != nil {
		panic(err)
	}

	//TODO: Melhorar essa constante CUSTOMER#
	sk, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", customer.Username))
	if err != nil {
		panic(err)
	}

	//TODO: Melhorar essas constantes PK e SK
	return map[string]types.AttributeValue{"PK": pk, "SK": sk}
}

func (c Customer) String() string {
	return fmt.Sprintf("Username: %s - Email: %s - Name: %s", c.Username, c.Email, c.Name)
}

func (o Order) String() string {
	return fmt.Sprintf("Id: %s - CreatedAt: %s - Status: %s - Amount: %s", o.Id, o.CreatedAt, o.Status, o.Amount)
}

func (o OrderItem) String() string {
	return fmt.Sprintf("Id: %s - Description: %s - Price: %s", o.Id, o.Description, o.Price)
}
