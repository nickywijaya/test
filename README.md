# Go-Xample

## Description

Go-Xample is a complete example of building microservice by following all microservice standards using Golang.

In this example, we create some made-up requirements:

- Build a user service
- There are http endpoins for login, logout, and registration
- When users register, check whether their email is valid by calling another service
- When users login, enqueue a job to log their login time
- The job mentioned above is part of this project
- There is a cron job to mark inactive users

## SLO and SLI

- There is no SLO / SLI

## Architecture Diagram

![go-xample-architecture](https://user-images.githubusercontent.com/4661221/31042306-000ece7e-a5cf-11e7-8745-6c0874b75565.png)

## Owner

SRE

## Contact and On-Call Information

See [Contact and On-Call Information](https://bukalapak.atlassian.net/wiki/display/INF/Contact+and+On-Call+Information)

## Links

- [Project Structure - Golang](https://bukalapak.atlassian.net/wiki/spaces/INF/pages/106987662/Project+Structure+-+Golang)
- [Microservice Standardizations](https://bukalapak.atlassian.net/wiki/spaces/INF/pages/100328043/Standardization)

## Onboarding and Development Guide

### Prerequisite

- Read all documentations in doc folder
- Git
- Go 1.9 or later

### Setup

- Install Git

  See [Git Installation](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

- Install Go (Golang)

  See [Golang Installation](https://golang.org/doc/install)

- Clone this repo in `$GOPATH/src/github.com/bukalapak`

  If you have not set your GOPATH, set it using [this](https://golang.org/doc/code.html#GOPATH) guide.
  If you don't have directory `src`, `github.com`, or `bukalapak` in your GOPATH, please make them.

  ```sh
  git clone git@github.com:bukalapak/go-xample.git
  ```

- Go to Go-Xample directory, then install `govendor` and sync the vendor file

  ```sh
  cd $GOPATH/src/github.com/bukalapak/go-xample
  go get -u github.com/kardianos/govendor
  govendor fetch -v +outside
  ```
- Install dependencies

  ```sh
  make dep
  ```

### Development

- Make new branch with descriptive name about the change(s) and checkout to the new branch

  ```sh
  git checkout -b branch-name
  ```

- Make your change(s) and make the test(s)

- Beautify and review the codes by running this command (note: **this is a must!**)

  ```sh
  make pretty
  ```

  If there are any errors after executing the command, please fix them first

- Save dependencies

  ```sh
  govendor add +external
  govendor fetch -v +outside
  ```

- Commit and push your change to upstream repository

  ```sh
  git commit -m "a meaningful commit message"
  git push origin branch-name
  ```

- Open Pull Request in Repository

- Pull request should be merged only if review phase is passed


### Docker Image

Docker image can be found in `registry.bukalapak.io/bukalapak/go-xample`

## Request Flow and Endpoint

### Request Flow

Imagine your own version of this project's request flow :D

or just read architecture section above

### Endpoint

There are some HTTP endpoints:

- #### Healthz

  ```
  GET /healthz
  ```

  Example

  ```sh
  curl -X GET "http://localhost:1234/healthz"
  ```

  Output (application/json)

  ```json
  ok
  ```

- #### Register

  ```
  POST /users
  ```

  Example

  ```sh
  curl -X POST "http://localhost:1234/users" -d <user-data>
  ```

  Output (application/json)

  ```json
  <api v4 standard response>
  ```

- #### Get User

  ```
  GET /users/:id
  ```

  Example

  ```sh
  curl -X GET "http://localhost:1234/users/:id"
  ```

  Output (application/json)

  ```json
  <api v4 standard response>
  ```

- #### Login

  ```
  POST /login
  ```

  Example

  ```sh
  curl -X POST "http://localhost:1234/login" -d <credential-data>
  ```

  Output (application/json)

  ```json
  <api v4 standard response>
  ```

- #### Logout

  ```
  GET /logout
  ```

  Example

  ```sh
  curl -X GET "http://localhost:1234/logout"
  ```

  Output (application/json)

  ```json
  <api v4 standard response>
  ```

## FAQ

#### When do we have to implement this structure?

NOW!
