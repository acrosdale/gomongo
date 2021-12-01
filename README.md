# Description
 A golang project that implements a simple product catalog system. Users can do crud operation via the ReST APIs.

# Future Work
    1. Add user ID to the product table to enforce ownership of uploaded catalog.

# Internal Components
    1. controllers: contains the endpoint that process all the incoming requests
        1.1. api

    2. db: contains all db logic, encapsulations and abstraction
        2.1: mgdb: database layer for mongodb

    3. mocks: contains all generated mocks of internal domain interface(eg. MongoQueries) 

    4. services: contains the business logic of the aplication
        4.1 auth: auth domain services
        4.2 api: api domain services
        
    5. utils: contains resuable utils func 

# deployments
development deployment
    
    1. make test-server

production deployment

    1. make prod-server

# Testing
Currently the app support two type of testing unit and integration. to run test suites open the container (go-mongo) CLI and run one of the test suite cmd below. **Integration test only run on {test, prod}-server**

unit test

    1. make run-unit-tests

integration

    1. make run-integration-tests

run all test type in the app (unit, integrations, etc)

    1. make run-all-tests


# Example datatype

A  Product

``` json
{
    "product_name": "test1",
    "price": 250,
    "currency": "JMD",
    "discount": 0,
    "vendor": "BOC LLC",
    "accessories": [
        "charger",
        "gift coupon",
        "subscription"
    ]
}
```

# Endpoints

|       endpoints       | method |                                               description                                              |
|:---------------------:|:------:|:------------------------------------------------------------------------------------------------------:|
| /api/v1/users/        | POST   | create a user                                                                                          |
| /api/v1/auth/         | POST   | authenicate a user. sets x-auth-token in the response header. use the auth-token in subsequent request |
| /api/v1/products      | POST   | create a product                                                                    |
| /api/v1/products/:id  | GET    | get a product                                                                                          |
| /api/v1/products/ :id | DELETE | delete a product                                                
|/api/v1/products/ :id | PUT | update a product                                                                                       |

