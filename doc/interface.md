# Interface

## Motivation

Interface is just a set of methods. It defines behaviours for its implementer. Interface makes the code more flexible, scalable, and obvious.

## Usage

In this project, interface is used to defines dependency's behaviours. Take a look at `Database` below.

```golang
type Database interface {
  InsertUser(context.Context, User) error
  FindUserByID(context.Context, int) (User, error)
  FindUserByCredential(context.Context, User) (User, error)
  FindInactiveUsers(context.Context) ([]User, error)
  InsertLoginHistory(context.Context, LoginHistory) error
  DeactivateUsers(context.Context, []User) error
}
```

Any client who wants to act as GoXample database must implement all those behaviours. An example is provided in [mysql.go](https://github.com/bukalapak/go-xample/blob/master/database/mysql.go).

## Defining An Interface

These are the steps to define an interface:

- For a project, define what the project's needs from its dependencies.

  For example, GoXample needs to store its data to database, sends its data via AMQP, and does a remote call to another service.
  For simplicity, the rest of the explanation will only use database as example.

- Then, list all its requirements.

  We take database as an example. The database should be able to do the following:

  - Insert user data.
  - Get active user data by ID.
  - Get active user data by username and password.
  - Get inactive user data.
  - Insert login history data.
  - Update users when their last login is at least 30 days ago. The update function should change field `active` from `true` to `false`.


- From the requirements, make some single-responsibility methods.

  - To insert user data, database needs a well-defined user data as its parameter. Because insert operation can be unsuccessful, the method should return error. Knowing this, we decide to make a method named `InsertUser` that takes context and user data as its parameter and returns error.
  - To update users when user's last login is at least 30 days ago, database needs a well-defined array of user data as its parameter. Because update operation can be unsuccessful, the method should return error. Then, we make a method named `DeactivateUsers` that takes context and array of user data as its parameter and returns error.
  - To retrieve user data by ID, database needs to know the ID. So, the ID should be one of parameters. Because select operation can be unsuccessful, the method should return error alongside user data. Then, we make a method named `FindUserByID` that takes context and ID as its parameters and returns two values, those are user data and error.

- Congratz!

## Client

Since interface is only a set of methods, it literally does nothing. We need an implementer so it can be used. The implementer, or we should call it client, must implement all methods to satisfy the interface.

In this project, we choose MySQL as RDBMS. So, all we do is create a MySQL client that satisfies GoXample's `Database`.

If in the future we want to change RDBMS from MySQL to MongoDB, we don't need to change GoXample. That is one advantage! We only need to create MongoDB client that satisfies `Database`.

Let's look at GoXample struct.

```golang
type GoXample struct {
  database   Database
  messenger  Messenger
  connection Connection
}
```

Field `database` is an interface. It never say MySQL or MongoDB. So, whatever RDBMS that we use, as long as it satisfies `Database`, it can be used in GoXample. Assuming we have MySQL client and MongoDB client that satisfy `Database`, the snapshot code below is valid.

```golang
package main

import (
  gx "github.com/bukalapak/go-xample"
  "github.com/bukalapak/go-xample/database"
)

func main() {
  // use MySQL
  mysql := database.NewMySQL(database.Option{})
  goXampleMySQL := gx.GoXample{database: mysql}

  // use MongoDB
  mongodb := database.NewMongoDB(database.Option{})
  goXampleMongoDB := gx.GoXample{database: mongodb}
}
```

We have seen that by using interface, our code is flexible. The example proves that the code doesn't need major changes if we want to change the client.

## Rule of Thumb

- A method that more likely will call any dependency should only take 2 parameters **at most**. The first parameter is always context. The second parameter is well-defined main data. Term "should only 2 ... at most" means that it is allowed to have less than 2.
- A method should only return 2 values **at most**. The first return value is a well-defined data. The second parameter is error. Term "should only 2 ... at most" means that it is allowed to have less than 2.
- A method should return error as the last return value unless there is a condition that error will never occur.
- The `Database` above already implements these rule of thumbs.

## Tips

- When your struct or package depends on another package or something else, it is more likely that it needs interface.
- Develop your codes by abstracting them first. Then, slowly but surely make them obvious but flexible by detailing the requirements and behaviours. Finally, write your codes!
- Architecting codes is hard but it only in the beginning. Make the processes above as your behaviour and you'll see that it is an easy task!
