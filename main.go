package main

import (
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

var db = initDbContext()

func main() {
	router := gin.Default()
	router.GET("/token", genInviteToken)
	router.GET("/tokens", seeAllTokens)
	router.Run("localhost:8080")
}

// GET: Generate and return a token with expiration date
func genInviteToken(c *gin.Context) {
	// Create a write transaction
	txn := db.Txn(true)

	// Generate a token
	randomToken := Token{Token: uniuri.NewLen(6), ExpirationDate: (time.Now().AddDate(0, 0, 7)).Format(time.UnixDate)}

	// Add token to DB
	if err := txn.Insert("token", randomToken); err != nil {
		panic(err)
	}

	// Commit the change
	txn.Commit()

	// Return that token
	c.IndentedJSON(http.StatusOK, randomToken)
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

// Backbone stuff goes here. Create database with schema and return the context
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
