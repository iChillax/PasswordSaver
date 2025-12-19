package models

import (
	"backend/crypto"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Secret struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name           string             `bson:"name" json:"name"`
	Type           string             `bson:"type" json:"type"` // password, token, url, api_key, account
	EncryptedValue string             `bson:"encrypted_value" json:"-"`
	Category       string             `bson:"category" json:"category"`
	Tags           []string           `bson:"tags" json:"tags"`
	Notes          string             `bson:"notes" json:"notes"`
	Metadata       map[string]string  `bson:"metadata" json:"metadata"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// StoreSecret encrypts and stores a secret value
func (s *Secret) StoreSecret(plainValue string, encryptionKey []byte) error {
	encrypted, err := crypto.EncryptSecret(plainValue, encryptionKey)
	if err != nil {
		return err
	}
	s.EncryptedValue = encrypted
	s.UpdatedAt = time.Now()
	return nil
}

// RetrieveSecret decrypts the stored secret
func (s *Secret) RetrieveSecret(encryptionKey []byte) (string, error) {
	return crypto.DecryptSecret(s.EncryptedValue, encryptionKey)
}

// ValidateRequired checks if required fields are present
func (s *Secret) ValidateRequired() bool {
	return s.Name != "" && s.Type != "" && s.EncryptedValue != ""
}

// ValidateType checks if type is valid
func (s *Secret) ValidateType() bool {
	validTypes := map[string]bool{
		"password": true,
		"token":    true,
		"url":      true,
		"api_key":  true,
		"account":  true,
		"other":    true,
	}
	return validTypes[s.Type]
}
