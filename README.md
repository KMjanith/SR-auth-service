1. install dependencies
```
go mod tidy
```

2. Compile the protobuf file
   ```
   protoc --go_out=. --go_opt=paths=source_relative spec/apiMessages.proto
   ```
3. This service is a part of a microservice basded application. You can find details in [here](https://github.com/KMjanith/SR-service-runner/blob/main/Readme.md) to run this with other services.
4. Basically what this service does is get the username and password from the user and save them in the mongodb databse and pass a jwt token to the api-gateway.
