package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-memdb"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

type Token struct {
	Token          string `json:"token"`
	ExpirationDate string `json:"expirationDate"`
}

// Pointer to the current Db Context
var db = initDbContext()

func main() {
	// Initial Setup
	router := gin.Default()
	insertPredefinedTokens()

	// Closed Endpoints
	closedRouter := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin123",
	}))
	closedRouter.GET("/token", genInviteToken)
	closedRouter.GET("/tokens", seeAllTokens)
	closedRouter.POST("/revoke", revokeToken)

	// Public Endpoints
	router.POST("/validate", validateToken)

	router.Run(":8080")
}

// GET: Generate and return a token with expiration date
func genInviteToken(c *gin.Context) {
	// Create a write transaction
	txn := db.Txn(true)

	// Generate a token
	var randomToken Token
	baseSeed := time.Now().UTC().Second()
	if baseSeed%2 == 0 {
		randomToken = Token{Token: uniuri.NewLen(6), ExpirationDate: (time.Now().AddDate(0, 0, 7)).UTC().Format(time.UnixDate)}
	} else {
		randomToken = Token{Token: uniuri.NewLen(12), ExpirationDate: (time.Now().AddDate(0, 0, 7)).UTC().Format(time.UnixDate)}
	}

	// Add token to DB
	if err := txn.Insert("token", randomToken); err != nil {
		panic(err)
	}

	// Commit the change
	txn.Commit()

	// Return that token
	c.IndentedJSON(http.StatusOK, randomToken)
}

// POST: Revoke token
func revokeToken(c *gin.Context) {
	txn := db.Txn(true)
	defer txn.Abort()

	// Get token from traidtional POST param query
	token, hasValue := c.GetQuery("token")
	tokenLength := len(token)
	if tokenLength != 6 && tokenLength != 12 {
		c.IndentedJSON(http.StatusBadRequest, "Invalid Token")
		return
	}
	if hasValue {
		lookUp, lookUpErr := txn.First("token", "id", token)
		if lookUpErr != nil {
			panic(lookUpErr)
		}
		if lookUp != nil {
			txn.Delete("token", lookUp)
			txn.Commit()
			c.IndentedJSON(http.StatusOK, "Token revoked")
		} else {
			c.IndentedJSON(http.StatusOK, "No token found")
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, "No token passed in")
	}
}

// GET: Returns all generated tokens
func seeAllTokens(c *gin.Context) {
	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	tokens, err := txn.Get("token", "id")
	if err != nil {
		panic(err)
	}

	tempSlice := make([]Token, 1)
	for obj := tokens.Next(); obj != nil; obj = tokens.Next() {
		tempToken := obj.(Token)
		tempSlice = append(tempSlice, tempToken)
	}

	c.IndentedJSON(http.StatusOK, tempSlice)
}

// POST: Validate the passed in token
func validateToken(c *gin.Context) {
	txn := db.Txn(false)
	defer txn.Abort()
	/*
		// Get token from Request Body
		var token Token
		decoder := json.NewDecoder(c.Request.Body)
		decodeErr := decoder.Decode(&token)
		if decodeErr != nil {
			panic(decodeErr)
		}
		token = token.Token
	*/
	// Get token from Header
	token := c.GetHeader("token")
	tokenLength := len(token)
	if tokenLength != 6 && tokenLength != 12 {
		c.IndentedJSON(http.StatusBadRequest, "Invalid Token")
		return
	}
	lookUp, lookUpErr := txn.First("token", "id", token)
	if lookUpErr != nil {
		panic(lookUpErr)
	}
	if lookUp != nil {
		expDate, err := time.Parse(time.UnixDate, lookUp.(Token).ExpirationDate)
		if err != nil {
			panic(err)
		}
		elapsedTime := time.Since(expDate).Seconds() * (-1)
		fmt.Println(elapsedTime)
		if elapsedTime < (7 * 24 * 60 * 60) {
			c.IndentedJSON(http.StatusOK, "Authenticated")
		} else {
			c.IndentedJSON(http.StatusForbidden, "Token Expired")
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, "Token not found")
	}
}

// Foundational parts goes from here. Create database with schema and return the context
func initDbContext() *memdb.MemDB {
	// Create the DB schema
	var schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"token": {
				Name: "token",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Token"},
					},
					"expDate": {
						Name:    "expDate",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ExpirationDate"},
					},
				},
			},
		},
	}

	// Create a new database using defined schema
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return db
}

func insertPredefinedTokens() {
	// Create a write transaction
	txn := db.Txn(true)
	defer txn.Abort()

	// Define tokens
	presetTokens := []Token{
		{Token: "abc123", ExpirationDate: time.Date(2022, 12, 17, 0, 0, 0, 651387237, time.Local).Format(time.UnixDate)},
		{Token: "abc456", ExpirationDate: time.Date(2022, 3, 30, 0, 0, 0, 651387237, time.Local).Format(time.UnixDate)},
		{Token: "xyz123456abc", ExpirationDate: time.Date(2022, 3, 30, 0, 0, 0, 651387237, time.Local).Format(time.UnixDate)},
	}

	for _, token := range presetTokens {
		if err := txn.Insert("token", token); err != nil {
			panic(err)
		}
	}

	txn.Commit()
}
