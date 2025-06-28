## Delete Table

- go run main.go delete-table

## Create Table

- go run main.go create-table-gsi

## Create Customers

- go run main.go create-customer --email normandes@email.com --username normandesjr --name "Normandes Junior"
- go run main.go create-customer --email sarah@email.com --username sarahmamede --name "Sarah Mamede"

## Add Addresses

- go run main.go update-customer --customer normandesjr --addressId home --address "Al Qwerty 256"
- go run main.go update-customer --customer normandesjr --addressId work --address "Av Rondon 1700"

## Create Orders

- go run main.go create-order --customer normandesjr --amount 100
- go run main.go create-order --customer normandesjr --amount 200
- go run main.go create-order --customer normandesjr --amount 300

## Search Orders

- go run main.go search-customer-orders --customer normandesjr --limit 2

## Create Order with Items

- go run main.go create-order --customer normandesjr --amount 250 --add-items true

## Search Order with items

- go run main.go search-order-items --order-id <orderId>
