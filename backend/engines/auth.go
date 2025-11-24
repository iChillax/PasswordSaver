package engines

import (
	"backend/models"
	"backend/settings"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Register creates a new user account
func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user model
	user := &models.User{
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate email
	if !user.ValidateEmail() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	// Validate password strength
	if !user.ValidatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password must be at least 8 characters with uppercase, number, and special character",
		})
		return
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		log.Error("Failed to hash password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	// Check if user already exists
	usersCollection := settings.MongoDatabase.Collection("users")
	existingUser := usersCollection.FindOne(ctx, bson.M{"email": user.Email})
	if existingUser.Err() == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	if existingUser.Err() != mongo.ErrNoDocuments {
		log.Error("Database error:", existingUser.Err())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Insert user
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Error("Failed to insert user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Error("Failed to generate token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID.Hex(),
			Email: user.Email,
		},
	})
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	usersCollection := settings.MongoDatabase.Collection("users")
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		log.Error("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Update last login
	user.UpdateLastLogin()
	_, err = usersCollection.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"last_login": user.LastLogin,
			"updated_at": user.UpdatedAt,
		},
	})
	if err != nil {
		log.Error("Failed to update last login:", err)
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Error("Failed to generate token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID.Hex(),
			Email: user.Email,
		},
	})
}
