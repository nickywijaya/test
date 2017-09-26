package go_xample_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	gx "github.com/bukalapak/go-xample"
)

type GoXampleSuite struct {
	suite.Suite
	goXample *gx.GoXample
}

func TestGoXampleSuite(t *testing.T) {
	suite.Run(t, &GoXampleSuite{})
}

func (g *GoXampleSuite) SetupSuite() {
	db := &MySQLMock{}
	msgr := &RabbitMQMock{}
	conn := &EmailCheckerMock{}

	g.goXample = gx.NewGoXample(db, msgr, conn)
}

func (g *GoXampleSuite) TestNewGoXample() {
	db := &MySQLMock{}
	msgr := &RabbitMQMock{}
	conn := &EmailCheckerMock{}

	goXample := gx.NewGoXample(db, msgr, conn)

	assert.NotNil(g.T(), goXample, "GoXample instance should not be nil!")
}

func (g *GoXampleSuite) TestCreateUser() {
	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Nanosecond)
	user, err := g.goXample.CreateUser(ctx, gx.User{})

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// db error
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Username: "bad-user"}
	user, err = g.goXample.CreateUser(ctx, user)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// email checker error
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Email: "bad@email.com"}
	user, err = g.goXample.CreateUser(ctx, user)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// invalid email
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Email: "invalid@email.com"}
	user, err = g.goXample.CreateUser(ctx, user)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// normal
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Username: "user1"}
	user, err = g.goXample.CreateUser(ctx, user)

	assert.NotNil(g.T(), user, "User instance should not be nil!")
	assert.Nil(g.T(), err, "Error should be nil!")
}

func (g *GoXampleSuite) TestGetUserByID() {
	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Nanosecond)
	user, err := g.goXample.GetUserByID(ctx, 1)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// db error
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err = g.goXample.GetUserByID(ctx, 0)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// normal
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err = g.goXample.GetUserByID(ctx, 1)

	assert.NotNil(g.T(), user, "User instance should not be nil!")
	assert.Nil(g.T(), err, "Error should be nil!")
}

func (g *GoXampleSuite) TestGetUserByCredential() {
	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Nanosecond)
	user, err := g.goXample.GetUserByCredential(ctx, gx.User{})

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// db error
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Username: "bad-user", Password: "user1"}
	user, err = g.goXample.GetUserByCredential(ctx, user)

	assert.Equal(g.T(), gx.User{}, user, "User instance should be empty!")
	assert.NotNil(g.T(), err, "Error should not be nil!")

	// normal
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user = gx.User{Username: "user1", Password: "user1"}
	user, err = g.goXample.GetUserByCredential(ctx, user)

	assert.NotNil(g.T(), user, "User instance should not be nil!")
	assert.Nil(g.T(), err, "Error should be nil!")
}

func (g *GoXampleSuite) TestSaveLoginHistory() {
	// db error
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lh := gx.LoginHistory{Username: "bad-user"}
	err := g.goXample.SaveLoginHistory(ctx, lh)

	assert.NotNil(g.T(), err, "Error should not be nil!")

	// normal
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lh = gx.LoginHistory{Username: "user1"}
	err = g.goXample.SaveLoginHistory(ctx, lh)

	assert.Nil(g.T(), err, "Error should be nil!")
}

func (g *GoXampleSuite) TestDeactivateInactiveUsers() {
	// db error in find inactive users
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx = context.WithValue(ctx, "Error", true)
	defer cancel()

	err := g.goXample.DeactivateInactiveUsers(ctx)

	assert.NotNil(g.T(), err, "Error should not be nil!")

	// db error in deactivate users
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	ctx = context.WithValue(ctx, "Users", 0)
	defer cancel()

	err = g.goXample.DeactivateInactiveUsers(ctx)

	assert.NotNil(g.T(), err, "Error should not be nil!")

	// normal
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = g.goXample.DeactivateInactiveUsers(ctx)

	assert.Nil(g.T(), err, "Error should be nil!")
}
