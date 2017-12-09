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
