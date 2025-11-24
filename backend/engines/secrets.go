package engines

import (
	"backend/models"
	"backend/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateSecretRequest struct {
	Name     string            `json:"name" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Value    string            `json:"value" binding:"required"`
	Category string            `json:"category"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

type UpdateSecretRequest struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Value    string            `json:"value"`
	Category string            `json:"category"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

type SecretResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Category  string            `json:"category"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type SecretDetailResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Value     string            `json:"value"`
	Category  string            `json:"category"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// getEncryptionKey retrieves the encryption key from environment
func getEncryptionKey() ([]byte, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		return nil, nil
	}
	// For now, assume it's a 32-byte hex string or base64
	// In production, use proper key derivation
	if len(keyStr) != 64 {
		log.Warn("ENCRYPTION_KEY should be 64 hex characters (32 bytes)")
	}
	// Convert hex string to bytes
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		var b byte
		fmt.Sscanf(keyStr[i*2:i*2+2], "%02x", &b)
		key[i] = b
	}
	return key, nil
}

// CreateSecret creates a new secret
func CreateSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req CreateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get encryption key
	key, err := getEncryptionKey()
	if err != nil || key == nil {
		log.Error("Failed to get encryption key:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption configuration error"})
		return
	}

	// Create secret model
	secret := &models.Secret{
		UserID:   userID.(primitive.ObjectID),
		Name:     req.Name,
		Type:     req.Type,
		Category: req.Category,
		Tags:     req.Tags,
		Metadata: req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate type
	if !secret.ValidateType() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret type"})
		return
	}

	// Encrypt and store secret value
	if err := secret.StoreSecret(req.Value, key); err != nil {
		log.Error("Failed to encrypt secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt secret"})
		return
	}

	// Insert into database
	secretsCollection := settings.MongoDatabase.Collection("secrets")
	result, err := secretsCollection.InsertOne(ctx, secret)
	if err != nil {
		log.Error("Failed to insert secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create secret"})
		return
	}

	secret.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, SecretResponse{
		ID:        secret.ID.Hex(),
		Name:      secret.Name,
		Type:      secret.Type,
		Category:  secret.Category,
		Tags:      secret.Tags,
		Metadata:  secret.Metadata,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	})
}

// ListSecrets returns all secrets for the authenticated user
func ListSecrets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get pagination parameters
	limit := int64(10)
	offset := int64(0)
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = parsed
		}
	}

	secretsCollection := settings.MongoDatabase.Collection("secrets")
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := secretsCollection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		log.Error("Failed to query secrets:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve secrets"})
		return
	}
	defer cursor.Close(ctx)

	var secrets []SecretResponse
	if err := cursor.All(ctx, &secrets); err != nil {
		log.Error("Failed to decode secrets:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode secrets"})
		return
	}

	if secrets == nil {
		secrets = []SecretResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(secrets),
		"data":  secrets,
	})
}

// GetSecret returns a single secret with decrypted value
func GetSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	secretID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(secretID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret ID"})
		return
	}

	// Get encryption key
	key, err := getEncryptionKey()
	if err != nil || key == nil {
		log.Error("Failed to get encryption key:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption configuration error"})
		return
	}

	secretsCollection := settings.MongoDatabase.Collection("secrets")
	var secret models.Secret
	err = secretsCollection.FindOne(ctx, bson.M{
		"_id":     objID,
		"user_id": userID,
	}).Decode(&secret)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
			return
		}
		log.Error("Failed to query secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve secret"})
		return
	}

	// Decrypt secret value
	decryptedValue, err := secret.RetrieveSecret(key)
	if err != nil {
		log.Error("Failed to decrypt secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt secret"})
		return
	}

	c.JSON(http.StatusOK, SecretDetailResponse{
		ID:        secret.ID.Hex(),
		Name:      secret.Name,
		Type:      secret.Type,
		Value:     decryptedValue,
		Category:  secret.Category,
		Tags:      secret.Tags,
		Metadata:  secret.Metadata,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	})
}

// UpdateSecret updates an existing secret
func UpdateSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	secretID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(secretID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret ID"})
		return
	}

	var req UpdateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get encryption key
	key, err := getEncryptionKey()
	if err != nil || key == nil {
		log.Error("Failed to get encryption key:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption configuration error"})
		return
	}

	secretsCollection := settings.MongoDatabase.Collection("secrets")

	// Find existing secret
	var secret models.Secret
	err = secretsCollection.FindOne(ctx, bson.M{
		"_id":     objID,
		"user_id": userID,
	}).Decode(&secret)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
			return
		}
		log.Error("Failed to query secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve secret"})
		return
	}

	// Update fields
	if req.Name != "" {
		secret.Name = req.Name
	}
	if req.Type != "" {
		if !secret.ValidateType() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret type"})
			return
		}
		secret.Type = req.Type
	}
	if req.Category != "" {
		secret.Category = req.Category
	}
	if req.Tags != nil {
		secret.Tags = req.Tags
	}
	if req.Metadata != nil {
		secret.Metadata = req.Metadata
	}

	// Encrypt new value if provided
	if req.Value != "" {
		if err := secret.StoreSecret(req.Value, key); err != nil {
			log.Error("Failed to encrypt secret:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt secret"})
			return
		}
	}

	secret.UpdatedAt = time.Now()

	// Update in database
	_, err = secretsCollection.UpdateByID(ctx, objID, bson.M{
		"$set": secret,
	})
	if err != nil {
		log.Error("Failed to update secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update secret"})
		return
	}

	c.JSON(http.StatusOK, SecretResponse{
		ID:        secret.ID.Hex(),
		Name:      secret.Name,
		Type:      secret.Type,
		Category:  secret.Category,
		Tags:      secret.Tags,
		Metadata:  secret.Metadata,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	})
}

// DeleteSecret deletes a secret
func DeleteSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	secretID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(secretID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret ID"})
		return
	}

	secretsCollection := settings.MongoDatabase.Collection("secrets")
	result, err := secretsCollection.DeleteOne(ctx, bson.M{
		"_id":     objID,
		"user_id": userID,
	})

	if err != nil {
		log.Error("Failed to delete secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete secret"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "secret deleted successfully"})
}

// SearchSecrets searches for secrets by name, category, or tags
func SearchSecrets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query required"})
		return
	}

	// Get pagination parameters
	limit := int64(10)
	offset := int64(0)
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = parsed
		}
	}

	secretsCollection := settings.MongoDatabase.Collection("secrets")
	filter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"category": bson.M{"$regex": query, "$options": "i"}},
			{"tags": bson.M{"$in": []string{query}}},
		},
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := secretsCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Error("Failed to search secrets:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search secrets"})
		return
	}
	defer cursor.Close(ctx)

	var secrets []SecretResponse
	if err := cursor.All(ctx, &secrets); err != nil {
		log.Error("Failed to decode secrets:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode secrets"})
		return
	}

	if secrets == nil {
		secrets = []SecretResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(secrets),
		"data":  secrets,
	})
}
