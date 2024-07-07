# Heli-Tech test project

# Services
- [x] Gateway-Service
- [x] User-Service
- [x] Transaction-Service
- [x] Notification-Service

# Gateway-Service
- [x] Gateway-Service is a service that acts as a gateway for all the services.
it exposes a rest API and calls other services using GRPC.
- the rest API is developed using golang gin router.
- this service also serves the swagger.

# User-Service (AKA Auth-Service)
- [x] this is a service that manages user related APIs.
- it serves requests using grpc.
- it also handle the auth using JWT.

# Transaction-Service 
- [x] this is a service that manages deposits and withdraws.
- it serves requests using grpc.
- it also sends events to kafka after each deposit or withdraw.

# Notification-Service 
- [x] this is a service that manages sending notifications to users.
- it is a consumer of kafka messages.
- messages related to deposits and withdraws.

## How to run the project
you need to have docker and docker-compose installed on your machine.
- clone the project
- run docker-compose up -d
- open the swagger in address http://localhost:8081/swagger/index.html


## how to check the flow using swagger
- you need to register first
- then login
- then authorize in swagger
- then you can call the transaction routes like deposit, withdraw, and transactions-list


## locking mechanism
as the deposits and withdraws can be called in parallel, the Transaction-Service uses a locking mechanism.
the locking mechanism is provided by redis.


## Other considerations
- if you have new macbook with arm chip, you may change the kafka image in docker-compose.yml line 95.
- each services directory structure consists 3 main layer, handler, services and repository.
* handler layer is responsible for the requests and response and also consuming from message brokers.
* service layer is responsible for the bussiness logic and can be called by handler layer.
* repositry layer is responsible for data access and read and write to database, the repo methods are called by service layer.
