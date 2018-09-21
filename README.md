# Savior
Savior is a back-end server that allows handling credit card transactions to be monitored by clients. The Savior server distributes real-time information web clients and mobile devices. As well as writing to reports and other consumers of Savior's transaction data. Savior also is the blocking mechanism to the credit card transaction and will hold transactions waiting to be approved by the client.

This server is made up of separate modules that handle the different areas of the Savior service.

## Savior service

#### Access

Handles login/logout, encryption layer and determining if a user is an administrator.

#### Cache

A general caching module that allows other modules to cache information in the database or elsewhere.

#### Database

Allows Savior to read, write and manage over the database.

#### Logger

Provides logging structure for Savior. Handles output of logs, error output levels and determines how verbose the server is.

#### Messaging

Handles alerts to clients via email, text, push notification, web, et cetera.

#### Metadata

Tracks data from transactions and records them in a way that can be used to determine unusual activity on the account.

#### Rules

Determines if transactions take place automatically or are halted and the user is alerted based on metadata and user configured rules.

#### System

Handles server connecting clients, requests, system information, reading/writing to disk and various other functions.

#### User

Handles user information and requests.

#### Settings

Settings for each modules are located inside the module folder. Option type is declared in the JSON file followed by the option value.

##### Example
```
{"Type":"Value","Type":"Value"}
```

## Coding Standards

The Savior service is written in the [Go programming language](https://www.golang.org) following the [Drupal coding standards](https://www.drupal.org/docs/develop/standards/coding-standards).

#### Exceptions

- Double quotes should always be used in Go
- Required_Once() does not apply
- Semicolons are not used in Go
- Naming convention is camelCase instead of underscore

#### Comments/GoDoc

GoDoc parses lines that "//Name" proceed any type, variable, constant, function or package and generates a HTML file with documentation.

##### Example
```
// Function name
// and description
func Foo() {
    // Code here
}
```

## Savior API

The following documentation contains examples of JSON transactions between the client and server. Examples include successes and failures for reference.

All objects below are to be assumed to be formatted with standard JSON.

### General

Standard response to a packet that is not recognized or contains incorrect elements for the current transaction.

```
{
  "status": "fail",
  "code": 400,
  "message": "invalid"
  "result": {}
}
```

### User registration

User registration requests are sent as a POST request to http://domain.name:8080/savior/user/register.

#### User registration object

The **userRegister** object contains two main values and four sub-values:
- *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
- *{user}* contains the user information
    - *{id}* contains the integer ID number for the new user
    - *{username}* contains the plain text username for the new user
    - *{password}* contains the plain text password for the new user
    - *{email}* contains the email address for the new user

##### User registration POST example

```
{
  "type": "userRegister",
  "user": {
    "id": 123,
    "username": "User1",
    "password": "Password",
    "email": "user@somewhere.com"
  }
}
```

#### Response
The response object contains four main values and three sub-values:
- *{status}* contains current HTTP status in human-readable format
- *{code}* contains current HTTP status code
- *{message}* contains current server message
- *{result}* contains the user information
    - *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
    - *{id}* contains the integer ID number for the new user
    - *{username}* contains the plain text username for the new user

##### Expected response for successful user creation
```
{
  "status": "ok",
  "code": 200,
  "message": "success",
  "result": {
    "type": "userRegister",
    "id": "123",
    "username": "User1"
  }
}
```

##### Expected response for duplicate user
```
{
  "status": "fail",
  "code": 400,
  "message": "duplicate",
  "result": {
    "type": "userRegister",
    "id": "123",
    "username": "User1"
  }
}
```

##### Expected response when a sent value is *NULL*
```
{
  "status": "fail",
  "code": 400,
  "message": "invalid",
  "result": {
    "type": "userRegister",
    "id": "123",
    "username": "User1"
  }
}
```

### User login

Login requests are sent as a POST request to http://domain.name:8080/savior/user/login.

#### User login object

The **userLogin** object contains two main values and four sub-values:
- *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
- *{user}* contains the user information
    - *{id}* contains the integer ID number for the existing user
    - *{username}* contains the plain text username for the existing user
    - *{password}* contains the plain text password for the existing user
    - *{email}* contains the email address for the existing user

##### User login POST example
```
{
  "type": "userLogin",
  "user": {
    "id": 123,
    "username": "User1",
    "password": "Password",
    "email": "user@somewhere.com"
  }
}
```

#### Response

The response object contains five main values and three sub-values:
- *{status}* contains current HTTP status in human-readable format
- *{code}* contains current HTTP status code
- *{token}* contains unique token for transaction
- *{message}* contains current server message
- *{result}* contains the user information
    - *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
    - *{id}* contains the integer ID number for the new user
    - *{username}* contains the plain text username for the new user

##### Expected response for successful user login
```
{
  "status": "ok",
  "code": 200,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "success",
  "result": {
    "type": "userLogin",
    "id": "123",
    "username": "User1"
  }
}
```

##### Expected response for failed login attempt
```
{
  "status": "fail",
  "code": 400,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "invalid",
  "result": {
    "type": "userLogin",
    "id": "123",
    "username": "User1"
  }
}
```

### Change password

Change password requests to the server are sent as a POST request to http://domain.name:8080/savior/user/change-password.

#### User change password object

The **userChangePassword** object contains two main values and four sub-values:
- *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
- *{user}* contains the user information
    - *{id}* contains the integer ID number for the existing user
    - *{username}* contains the plain text username for the existing user
    - *{password}* contains the plain text password for the existing user
    - *{email}* contains the email address for the existing user

##### User change password POST example
```
{
  "type": "userChangePassword",
  "user": {
    "id": 123,
    "username": "User1",
    "password": "newPassword",
    "email": "user@somewhere.com"
  }
}
```
#### Response

The response object contains five main values and three sub-values:
- *{status}* contains current HTTP status in human-readable format
- *{code}* contains current HTTP status code
- *{token}* contains unique token for transaction
- *{message}* contains current server message
- *{result}* contains the user information
    - *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
    - *{id}* contains the integer ID number for the new user
    - *{username}* contains the plain text username for the new user

##### Expected response for successful user password change
```
{
  "status": "ok",
  "code": 200,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "success",
  "result": {
    "type": "userChangePassword",
    "id": "123",
    "username": "User1"
  }
}
```

##### Expected response for failed user password change
```
{
  "status": "fail",
  "code": 400,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "invalid",
  "result": {
    "type": "userChangePassword",
    "id": "123",
    "username": "User1"
  }
}
```

### Remove user

User removal requests to the server are sent as a POST request to http://domain.name:8080/savior/user/remove.

#### Remove user password object

The **userRemove** object contains two main values and four sub-values:
- *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
- *{user}* contains the user information
    - *{id}* contains the integer ID number for the existing user
    - *{username}* contains the plain text username for the existing user
    - *{password}* contains the plain text password for the existing user
    - *{email}* contains the email address for the existing user

##### Remove user POST example
```
{
  "type": "userRemove",
  "user": {
    "id": 123,
    "username": "User1",
    "password": "Password",
    "email": "user@somewhere.com"
  }
}
```
#### Response

The response object contains five main values and two sub-values:
- *{status}* contains current HTTP status in human-readable format
- *{code}* contains current HTTP status code
- *{token}* contains unique token for transaction
- *{message}* contains current server message
- *{result}* contains the user information
    - *{type}* contains the transaction type for processing (e.g. userRegister, userLogin, etc)
    - *{id}* contains the integer ID number for the new user
    - *{username}* contains the plain text username for the new user

##### Expected response for successful user removal
```
{
  "status": "ok",
  "code": 200,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "success",
  "result": {
    "type": "userRemove",
    "id": "123",
    "username": "User1"
  }
}
```

##### Expected response for failed user removal
```
{
  "status": "fail",
  "code": 400,
  "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
  "message": "invalid",
  "result": {
    "type": "userRemove",
    "id": "123",
    "username": "User1"
  }
}
```

### Credits

**Compiled by:** Reg Proctor and Ian Bartelds
**Refactored by:** Michael Gilardi
