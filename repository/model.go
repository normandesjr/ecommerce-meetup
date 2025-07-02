package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Customer struct {
	PK        string             `dynamodbav:"PK"`
	SK        string             `dynamodbav:"SK"`
	Username  string             `dynamodbav:"username"`
	Email     string             `dynamodbav:"email"`
	Name      string             `dynamodbav:"name"`
	Addresses map[string]Address `dynamodbav:"addresses"`
}

type CustomerEmail struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	Username string `dynamodbav:"username"`
}

type Address struct {
	Id            string `dynamodbav:"id"`
	StreetAddress string `dynamodbav:"streetAddress"`
	ZipCode       string `dynamodbav:"zipCode"`
}

type Order struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	Id        string `dynamodbav:"orderId"`
	CreatedAt string `dynamodbav:"createdAt"`
	Status    string `dynamodbav:"status"`
	Total     int    `dynamodbav:"total"`
	ShippedTo string `dynamodbav:"shippedTo"`
	Username  string `dynamodbav:"username"`
	GSI1PK    string `dynamodbav:"GSI1PK"`
	GSI1SK    string `dynamodbav:"GSI1SK"`
}

type OrderItem struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	Id          string `dynamodbav:"itemId"`
	OrderId     string `dynamodbav:"orderId"`
	Description string `dynamodbav:"description"`
	Price       int    `dynamodbav:"price"`
	GSI1PK      string `dynamodbav:"GSI1PK"`
	GSI1SK      string `dynamodbav:"GSI1SK"`
}

type OrderItems []OrderItem

func (customer Customer) GetKey() map[string]types.AttributeValue {
	pk, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", customer.Username))
	if err != nil {
		panic(err)
	}

	sk, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", customer.Username))
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{"PK": pk, "SK": sk}
}

func (items OrderItems) Total() int {
	total := 0
	for _, i := range items {
		total += i.Price
	}

	return total
}
