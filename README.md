# Saviour

## Description:

Savior is a back-end server that allows handling credit card transactions to be
monitored by clients. The Saviour server distributes real-time information web
clients and mobile devices. As well as writing to reports and other consumers
of Saviour's transaction data. Saviour also is the blocking mechanism to the
credit card transaction and will hold transactions waiting to be approved by
the client.

This server is made up of separate modules that handle the different areas of
the Saviour service. These include:

* Access

Handles login/logout, encryption layer, & determining if a user is an Administrator

* Cache

A general caching module that allows other modules to cache information in the
database or elsewhere.

* Database

Allows Saviour to read, write, & manage over the database

* Logger

Provides logging structure for Saviour. Handles output of logs, Error output
levels and determine how verbose the server is.

* Messaging

Handles alerts to clients via Email, Text, Notification, Web, etc.

* Metadata

Tracks data from transactions and records them in a way that can be used
to determine unusual activity on the account.

* Rules

Determines if transactions take place automatically or are halted and the
user is alerted based on metadata and user configured rules.

* System

Handles server connecting clients, requests, system information,
reading/writing to disk, and various other functions

* User

Handles user information and requests

### Settings

Settings for each modules are located inside the module folder. Option type is
declared in the json file followed by the option value.

Example:

```
{"Type":"Value","Type":"Value"}
```


### Coding Standards

This program is written in go (www.golang.org) following the durpal coding
standards (https://www.drupal.org/docs/develop/standards/coding-standards)

Exceptions:

*Double quotes should always be used in go
*Required_Once() does not apply
*Semicolons are not used in go
*Naming convention is camel case instead of underscore

Comments/GoDoc:
GoDoc parses lines that "//Name" proceed any type, variable, constant, function,or
package and generates a html file with documentation

Example:

```
// Function Description
// More Text Here
func Function() {
...
}
```

### API Documentation

The following contains examples of json transactions between the client and server. Both successful
and failure.

## General

Standard Response to a packet that is not recognized or contains incorrect elements for the current
transaction.

{
  "login": {},
    "saviour": {
      "status": 400,
      "message": "InvalidPacket"
    }
}

## User registration

user registration requests are sent as a post to http://domain.name:8080/register

JSON Registration Client Packet Example:

{
"login":
  {
    "user":"User1","pass":"Password","email":"user@somewhere.com"
  }
}

Array Login Contains Three Elements:

user contains the user name for the new user
pass contains the plain text password for the new user
email contains the email address for the new user

The correct server response from a successful user creation:

{
  "login": {},
    "saviour": {
      "username": "User1",
      "status": 200,
      "message": "UserCreationSucsessful::User1"
    }
}

The correct server response from a duplicate user creation:

{
  "login": {},
    "saviour": {
      "username": "User1",
      "status": 400,
      "message": "DuplicateUser::UserCreationFailed"
    }
}

The correct server response from a empty name field:

{
  "login": {},
    "saviour": {
      "status": 400,
      "message": "NameEntryIsEmpty::UserCreationFailed"
    }
}

The correct server response from a empty password field:

{
  "login": {},
    "saviour": {
      "username": "User1",
      "status": 400,
      "message": "PasswordEntryIsEmpty::UserCreationFailed"
    }
}

The correct server response from a empty email field:

{
  "login": {},
    "saviour": {
      "username": "User12",
      "status": 400,
      "message": "EmailEntryIsEmpty::UserCreationFailed"
    }
}

Array Saviour Contains Three Elements:
username contains the user name thats been requested to be created
status contains current http status
message contains current server message

## Login transaction

login requests are sent as a post request to http://domain.name:8080/login

JSON Login Client Packet:

{
  "login":
  {
  "user":"Admin","pass":"Password"
  }
}

Array Login Contains Two Elements:

user contains the user name for this login session
pass  contains the password for this login session

JSON Login Server Response:

The correct server response from a successful login:

{
  "login": {},
    "saviour": {
      "username": "Admin",
      "status": 200,
      "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU=",
      "message": "LoginSuccessful"
    }
}

The correct server response to a failed login attempt:

{
  "login": {},
    "saviour": {
      "username": "Admin",
      "status": 400,
      "message": "UserNotFound"
    }
}

Array Saviour Contains Four Elements:

username contains the current users name
status contains current http status code
token contains the current authentication token
message contains the current server message

### Logoff Transaction

Logoff request to the are send as a POST request to http://domain.name/request/logoff

JSON Client Logoff Packet:

{
  "login": {},
    "saviour": {
      "username": "Admin",
      "status": 200,
      "token": "M8YQ6Ez_-c9wyzBJ362l2Kqi8B2SJ0GxwBW_JZiVbaU="
    }
}

Array saviour contains the following elements:

username contains the name of the current username
status contains the current html status
token contains the authentication token for the current sesssion

JSON Server Logoff Response:

The correct server response to a successful logoff:

{
  "login": {},
    "saviour": {
      "username": "Admin",
      "status": 200,
      "message": "LogOff::Sucsessful"
    }
}
