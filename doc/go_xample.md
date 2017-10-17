# Go Xample

## Motivation

GoXample is the main service logic in this project. GoXample contains all business logics. It is located in [go_xample.go](https://github.com/bukalapak/go-xample/blob/master/go_xample.go).

Note that there is only two package in root folder. The first is the project package. In this case, the package is named `go_xample`. The second package is its test, that is `go_xample_test`.

Why?
Because now we are developing using package driven paradigm. Everything is treated as package. It makes sense that in one folder only exists one package. When developers import this project, then they will import `github.com/bukalapak/go-xample`. What developers hope is that they import **only** what they want. In this case, they just want to import nothing but all GoXample's functionalities. So, anything unrelated to the package's functionality doesn't have the right to be placed in the same package. So, in simpler word, what GoXample can do will be placed in go_xample package and what GoXample depends on will be placed in different package.

Because dependency can be so vast, GoXample needs to define the dependency functionality. It will be done by implementing interface for any dependency that GoXample uses. By doing so, the contract will be clear and developers will easily know what a dependency does just by looking at the interface.

## Content

In [go_xample.go](https://github.com/bukalapak/go-xample/blob/master/go_xample.go), there are some predefined contents.

### Main Struct

There is one main struct, in this example GoXample, to organize the service logic. This struct couples all dependencies to itself. As an example, GoXample needs database (MySQL), messenger (RabbitMQ), and connection to third party service (EmailChecker). Therefore, GoXample must have all dependencies in itself.

```golang
type GoXample struct {
  database   Database
  messenger  Messenger
  connection Connection
}
```

The dependencies' implementation is assigned to each dependent in different package.

### Interface

Interface is a contract for dependency. GoXample knows what it needs from other package. So, GoXample defines the requirements in interface.

Why?
By implementing interface, the contract is clear. Client will just follow the interface to implement the requirements. By using interface, GoXample will not depend on a specific client rather the client must behave just as GoXample wants by satisfying the interface.

A good example is database package. It is located in folder [database](https://github.com/bukalapak/go-xample/tree/master/database). GoXample just wants to be able to store its data to database but GoXample doesn't care what RDBMS is used. It can be MySQL, MongoDB, Cassandra, Redis, or anything. Therefore, GoXample defines a contract for database. The contract contains all functionalities that  database must do.

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

Any client (MySQL, MongoDB, etc) just has to satisfiy the interface so that GoXample can use it. If MySQL is needed, then we just build a MySQL client that implements all the requirements. The same thing goes with MongoDB or other RDBMS. In this project, we use MySQL and we created MySQL client in [mysql.go](https://github.com/bukalapak/go-xample/blob/master/database/mysql.go). It satisfies `Database` so GoXample can use it as database client.

An advantage of using interface is it is easy to mock. We just create a mock that satisfies the interface, then, voila, we have a useful mock. The example can be found in [go_xample_mock_test.go](https://github.com/bukalapak/go-xample/blob/master/go_xample_mock_test.go)

### Service Function

Service Function is an exported GoXample function. An exported GoXample function must implement exactly one business logic. It implements `Single Responsibility Principle`. We can compare it with `Service Object` in Rails or any other framework. Basically, they do the same task.

## Dependency

Dependencies for GoXample are located in different packages.

Why?
Because dependencies are not in the project's scope. They are basically the other packages that need GoXample or are needed by GoXample. Then, it makes sense to separate dependencies package from `go_xample`.

Dependencies must satisfy interface that is defined by GoXample.

Why?
As stated above, a dependency must satisfy the interface so it can be used by GoXample. This way, abstraction is clear, implementation is guided well, and contract is obvious.

Since `go_xample` is the main package and the interfaces are defined in `go_xample`, dependency should import `go_xample`.

## FAQ

### Is there any guide on how to make an interface?

Read [interface](https://github.com/bukalapak/go-xample/blob/master/doc/interface.md)

### What is context?

Read [context](https://github.com/bukalapak/go-xample/blob/master/doc/context.md)
