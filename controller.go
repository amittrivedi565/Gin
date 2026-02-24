package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context, db *DB){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	count, _ := db.Collection.CountDocuments(ctx, user.Email)
	if count > 0 {
		c.JSON(500, gin.H{"error": "Email Already Registered"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user.ID = bson.ObjectID{}
	user.Password = string(hashedPassword)
	
	_, err = db.Collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context, db *DB){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	var input User
	var user User

	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err := db.Collection.FindOne(ctx, bson.M{"email" : input.Email}).Decode(&user)
	if err != nil{
		c.JSON(500, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil{
		c.JSON(500, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}


func GetTodos(c *gin.Context, db *DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todos []Todo

	cursor, err := db.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &todos); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todos)
}

func GetTodo(c *gin.Context, db *DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todo Todo
	id := c.Param("id")

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	err = db.Collection.
		FindOne(ctx, bson.M{"_id": objID}).
		Decode(&todo)

	if err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(200, todo)
}

func CreateTodo(c *gin.Context, db *DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todo Todo

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo.ID = bson.NewObjectID()

	result, err := db.Collection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, result)
}

func DeleteTodo(c *gin.Context, db *DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(200, result)
}

