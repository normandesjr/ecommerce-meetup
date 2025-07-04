# ECommerce DynamoDB Modeling

## AWS Profile

You should use an AWS Profile to connect to AWS and use the commands.

You can pass to CLI using `-profile`, `-p` or via environment variable `CDAY_PROFILE`.

All the commands bellow expect the environment variable to be set, be sure to export it.

```sh
export CDAY_PROFILE=profile
```
## Delete Table

```sh
go run main.go delete-table
```

## Create Table

The table will be named "CloudDayTable" unless you define the flag `-table` or via environment variable `CDAY_TABLE`.

```sh
go run main.go create-table
```

## Create Customers

```sh
go run main.go create-customer --email normandes@email.com --username normandesjr --name "Normandes Junior"
go run main.go create-customer --email sarah@email.com --username sarahmamede --name "Sarah Mamede"
```

## Add Addresses

```sh
go run main.go update-address --username normandesjr --id home --street-address "Al Qwerty 256" --zip-code "38400-111"

go run main.go update-address --username normandesjr --id work --street-address "Av Rondon 1700" --zip-code "38400-222"

go run main.go update-address --username sarahmamede --id home --street-address "Al Qwerty 156" --zip-code "38400-111"
```

## Update Address

```sh
go run main.go update-address --username normandesjr --id home --street-address "Al Qwerty 156" --zip-code "38400-111"
```

## Delete Address

```sh
go run main.go remove-address --username normandesjr --id work
```

## Create Orders

```sh
go run main.go create-order --username normandesjr --ship-address home --items 1,3,5

go run main.go create-order --username normandesjr --ship-address home --items 2,4,6

go run main.go create-order --username sarahmamede --ship-address home --items 5,6,7,8,9
``` 

## Search Orders

```sh
go run main.go get-order --username normandesjr

go run main.go get-order --username sarahmamede
```

## Search Order with items

```sh
go run main.go get-order-items --order-id <orderId>
```

## Update order status

```sh
go run main.go update-status-order --order-id <orderId> --status shipped
```

### Remove user address

```sh
go run main.go remove-address --username normandesjr --id home
```