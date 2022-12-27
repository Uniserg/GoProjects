package services

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/services/jwtutils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type User map[string]interface{}

type AuthService interface {
	Save(email, plainPassword, name string) (User, error)

	FindOneByEmailAndPassword(email string, password string) (User, error)

	ExtractUserId(bearer string) (string, error)
}

type neo4jAuthService struct {
	loader     *fixtures.FixtureLoader
	driver     neo4j.Driver
	jwtSecret  string
	saltRounds int
}

func NewAuthService(loader *fixtures.FixtureLoader, driver neo4j.Driver, jwtSecret string, saltRounds int) AuthService {
	return &neo4jAuthService{
		loader:     loader,
		driver:     driver,
		jwtSecret:  jwtSecret,
		saltRounds: saltRounds,
	}
}

// Save should create a new User node in the database with the email and name
// provided, along with an encrypted version of the password and a `userId` property
// generated by the server.
//
// The properties also be used to generate a JWT `token` which should be included
// with the returned user.
// tag::register[]
func (as *neo4jAuthService) Save(email, plainPassword, name string) (_ User, err error) {

	session := as.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	encryptedPassword, err := encryptPassword(plainPassword, as.saltRounds)
	if err != nil {
		return nil, err
	}

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
		CREATE (u:User {
			userId: randomUuid(),
			email: $email,
			password: $encrypted,
			name: $name
		})
		RETURN u { .userId, .name, .email } as u`,
			map[string]interface{}{
				"email":     email,
				"encrypted": encryptedPassword,
				"name":      name,
			})

		record, err := result.Single()

		if err != nil {
			return nil, err
		}

		user, _ := record.Get("u")

		return user, nil
	})

	if neo4jError, ok := err.(*neo4j.Neo4jError); ok && neo4jError.Title() == "ConstraintValidationFailed" {
		return nil, NewDomainError(
			422,
			fmt.Sprintf("An account already exists with the email address %s", email),
			map[string]interface{}{
				"email": "Email address taken",
			},
		)
	}

	user := result.(map[string]interface{})
	subject := user["userId"].(string)
	token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)
	if err != nil {
		return nil, err
	}

	return userWithToken(user, token), nil
}

// end::register[]

// tag::authenticate[]
func (as *neo4jAuthService) FindOneByEmailAndPassword(email string, password string) (_ User, err error) {
	// TODO: Authenticate the user from the database
	session := as.driver.NewSession(neo4j.SessionConfig{})

	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User {email:$email}) RETURN u",
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Single()

		if err != nil {
			return nil, fmt.Errorf("account not found or incorrect password")
		}

		user, _ := record.Get("u")

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	userNode := result.(neo4j.Node)
	user := userNode.Props

	if !verifyPassword(password, user["password"].(string)) {
		return nil, fmt.Errorf("account not found or incorrect password")
	}

	subject := user["userId"].(string)
	token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)

	if err != nil {
		return nil, err
	}

	return userWithToken(user, token), nil
}

// end::authenticate[]

func (as *neo4jAuthService) ExtractUserId(bearer string) (string, error) {
	if bearer == "" {
		return "", nil
	}
	userId, err := jwtutils.ExtractToken(bearer, as.jwtSecret, func(token *jwt.Token) interface{} {
		claims := token.Claims.(jwt.MapClaims)
		return claims["sub"]
	})
	if err != nil {
		return "", err
	}
	return userId.(string), nil
}

func encryptPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}

func userToClaims(user User) map[string]interface{} {
	return map[string]interface{}{
		"sub":    user["userId"],
		"userId": user["userId"],
		"name":   user["name"],
	}
}

func userWithToken(user User, token string) User {
	return map[string]interface{}{
		"token":  token,
		"userId": user["userId"],
		"email":  user["email"],
		"name":   user["name"],
	}
}
