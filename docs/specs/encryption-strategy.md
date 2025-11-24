# Encryption Strategy Specification

## Overview

This document outlines the encryption strategy for the AI-Powered Secrets Manager application. It covers how secrets are encrypted at rest, how encryption keys are managed, and how the system enables searching encrypted data without decrypting the entire database.

---

## 1. Encryption Architecture

### 1.1 Encryption Layers

The application implements a multi-layer encryption approach:

1. **At Rest**: Secrets stored in MongoDB are encrypted using AES-256-GCM
2. **In Transit**: All API communications use HTTPS/TLS
3. **In Use**: Secrets are decrypted only when explicitly requested by authenticated users

### 1.2 Encryption Algorithm: AES-256-GCM

**Why AES-256-GCM?**
- **AES-256**: 256-bit key provides military-grade security
- **GCM Mode**: Galois/Counter Mode provides authenticated encryption
  - Ensures data integrity (detects tampering)
  - Provides authentication tag
  - No manual padding required
  - Suitable for streaming data

**Advantages over AES-CBC:**
- Built-in authentication (prevents tampering)
- Simpler implementation (no manual padding)
- Better performance for modern CPUs
- Recommended by NIST for new applications

---

## 2. Key Management

### 2.1 Master Encryption Key (MEK)

**Purpose**: Encrypts all secrets in the database

**Key Specifications**:
- Length: 256 bits (32 bytes)
- Generation: Cryptographically secure random (crypto/rand in Go)
- Storage: Environment variable `ENCRYPTION_KEY` (loaded at startup)
- Rotation: Manual process (documented separately)

**Storage Strategy**:
```
Development: .env file (git-ignored)
Production: Secure vault (AWS Secrets Manager, HashiCorp Vault, or similar)
```

### 2.2 Key Derivation (Optional Enhancement)

For future multi-user scenarios, consider key derivation:
```
User-Specific Key = HKDF(Master Key, User ID, Salt)
```

This allows per-user encryption without storing multiple keys.

### 2.3 Key Rotation Strategy

**Current Approach**: Manual rotation
1. Generate new MEK
2. Decrypt all secrets with old MEK
3. Re-encrypt all secrets with new MEK
4. Update environment variable
5. Restart application

**Future Enhancement**: Automated key rotation with versioning

---

## 3. Encryption Implementation

### 3.1 Go Implementation (Backend)

#### Encryption Function

```go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
)

// EncryptSecret encrypts a secret value using AES-256-GCM
func EncryptSecret(plaintext string, key []byte) (string, error) {
    // Validate key length (must be 32 bytes for AES-256)
    if len(key) != 32 {
        return "", fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
    }

    // Create AES cipher block
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    // Create GCM cipher mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    // Generate random nonce (12 bytes for GCM)
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", fmt.Errorf("failed to generate nonce: %w", err)
    }

    // Encrypt plaintext
    // gcm.Seal() returns: nonce + ciphertext + authentication tag
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    // Encode to base64 for storage in MongoDB
    encoded := base64.StdEncoding.EncodeToString(ciphertext)
    return encoded, nil
}

// DecryptSecret decrypts an encrypted secret value
func DecryptSecret(encrypted string, key []byte) (string, error) {
    // Validate key length
    if len(key) != 32 {
        return "", fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
    }

    // Decode from base64
    ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64: %w", err)
    }

    // Create AES cipher block
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    // Create GCM cipher mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    // Extract nonce from ciphertext
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    // Decrypt and verify authentication tag
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("decryption failed (data may be tampered): %w", err)
    }

    return string(plaintext), nil
}
```

#### Integration with MongoDB Models

```go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Secret struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    UserID        primitive.ObjectID `bson:"user_id"`
    Name          string             `bson:"name"`
    Type          string             `bson:"type"` // password, token, url, api_key
    EncryptedValue string             `bson:"encrypted_value"` // Base64 encoded
    Category      string             `bson:"category"`
    Tags          []string           `bson:"tags"`
    Metadata      map[string]string  `bson:"metadata"`
    CreatedAt     time.Time          `bson:"created_at"`
    UpdatedAt     time.Time          `bson:"updated_at"`
}

// StoreSecret encrypts and stores a secret
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
```

---

## 4. Searching Encrypted Data

### 4.1 Challenge: Searching Without Decryption

**Problem**: How to search encrypted secrets without decrypting the entire database?

**Solutions**:

#### Option 1: Decrypt on Client (Recommended for MVP)
- Retrieve all secrets for the user
- Decrypt them in memory
- Filter/search on decrypted data
- **Pros**: Simple, secure, no searchable encryption complexity
- **Cons**: Performance impact with large datasets

#### Option 2: Searchable Encryption (Future Enhancement)
- Store encrypted search indexes alongside encrypted data
- Use deterministic encryption for searchable fields
- **Pros**: Fast searches without decryption
- **Cons**: Complex implementation, potential security trade-offs

#### Option 3: Hybrid Approach (Recommended for Production)
- Encrypt sensitive values (passwords, tokens)
- Store metadata (name, category, tags) in plaintext or with deterministic encryption
- Search on metadata, decrypt only matching results

### 4.2 Recommended Approach: Hybrid Encryption

**Implementation Strategy**:

```go
type Secret struct {
    ID                primitive.ObjectID `bson:"_id,omitempty"`
    UserID            primitive.ObjectID `bson:"user_id"`
    Name              string             `bson:"name"` // Plaintext (searchable)
    Type              string             `bson:"type"` // Plaintext (searchable)
    EncryptedValue    string             `bson:"encrypted_value"` // Encrypted
    Category          string             `bson:"category"` // Plaintext (searchable)
    Tags              []string           `bson:"tags"` // Plaintext (searchable)
    Metadata          map[string]string  `bson:"metadata"` // Plaintext
    CreatedAt         time.Time          `bson:"created_at"`
    UpdatedAt         time.Time          `bson:"updated_at"`
}
```

**Search Flow**:
1. User searches for "GitHub" via chatbot or UI
2. Backend queries MongoDB for secrets with `name` or `tags` containing "GitHub"
3. Retrieve matching documents (encrypted values still encrypted)
4. Decrypt only the matching secrets
5. Return decrypted results to user

**MongoDB Query Example**:
```go
// Search by name or tags
filter := bson.M{
    "user_id": userID,
    "$or": []bson.M{
        {"name": bson.M{"$regex": searchTerm, "$options": "i"}},
        {"tags": bson.M{"$in": []string{searchTerm}}},
        {"category": searchTerm},
    },
}

// Execute query
cursor, err := collection.Find(ctx, filter)
```

---

## 5. Chatbot Integration with Encryption

### 5.1 Chatbot Search Flow

```
User Query
    ↓
Ollama/LLM (Parse Intent)
    ↓
Extract Search Terms
    ↓
Query MongoDB (Plaintext Fields)
    ↓
Decrypt Matching Secrets
    ↓
Format Response
    ↓
Return to User
```

### 5.2 Prompt Engineering for Secure Search

**System Prompt**:
```
You are a secure secrets manager assistant. Your role is to:
1. Understand user queries about their stored secrets
2. Extract search terms (names, categories, types)
3. Never suggest revealing secrets in logs or responses
4. Always confirm before returning sensitive data
5. Suggest using categories and tags for better organization

When a user asks for a secret:
- Extract the secret name or category
- Ask for confirmation before revealing
- Return only the requested information
```

### 5.3 Security Considerations

- Never log decrypted secrets
- Confirm sensitive operations before execution
- Implement rate limiting on search queries
- Audit all secret access attempts
- Implement session timeouts

---

## 6. Data Flow Diagrams

### 6.1 Secret Creation Flow

```
User Input (Plaintext)
    ↓
Validate Input
    ↓
Encrypt with AES-256-GCM
    ↓
Base64 Encode
    ↓
Store in MongoDB (encrypted_value field)
    ↓
Audit Log Entry
```

### 6.2 Secret Retrieval Flow

```
User Request (Authenticated)
    ↓
Query MongoDB (by ID or search terms)
    ↓
Retrieve Encrypted Document
    ↓
Base64 Decode
    ↓
Decrypt with AES-256-GCM
    ↓
Verify Authentication Tag
    ↓
Return Plaintext to User
    ↓
Audit Log Entry
```

### 6.3 Search Flow

```
User Search Query
    ↓
Parse with Ollama/LLM
    ↓
Extract Search Terms
    ↓
Query MongoDB (plaintext fields: name, category, tags)
    ↓
For Each Match:
    - Decrypt encrypted_value
    - Verify authentication tag
    ↓
Return Decrypted Results
    ↓
Audit Log Entry
```

---

## 7. Security Best Practices

### 7.1 Key Management

- ✅ Use cryptographically secure random generation (crypto/rand)
- ✅ Store keys in environment variables or secure vaults
- ✅ Never hardcode keys in source code
- ✅ Rotate keys periodically
- ✅ Use different keys for different environments
- ❌ Don't use weak key derivation functions
- ❌ Don't reuse nonces

### 7.2 Encryption

- ✅ Use AES-256-GCM for authenticated encryption
- ✅ Generate unique nonce for each encryption
- ✅ Verify authentication tag on decryption
- ✅ Handle decryption errors gracefully
- ❌ Don't use ECB mode (deterministic, insecure)
- ❌ Don't reuse IVs/nonces

### 7.3 Logging and Monitoring

- ✅ Log all secret access attempts
- ✅ Log encryption/decryption errors
- ✅ Monitor for suspicious patterns
- ✅ Alert on failed decryption attempts
- ❌ Never log plaintext secrets
- ❌ Never log encryption keys

### 7.4 Access Control

- ✅ Authenticate all requests (JWT)
- ✅ Authorize access to user's own secrets only
- ✅ Implement rate limiting
- ✅ Use HTTPS for all communications
- ✅ Implement session timeouts
- ❌ Don't trust client-side validation alone

---

## 8. Implementation Checklist

### Phase 1: MVP
- [ ] Implement AES-256-GCM encryption/decryption functions
- [ ] Create Secret model with encrypted_value field
- [ ] Implement secret creation with encryption
- [ ] Implement secret retrieval with decryption
- [ ] Add encryption key management (environment variable)
- [ ] Implement basic search on plaintext fields
- [ ] Add audit logging for secret access

### Phase 2: Chatbot Integration
- [ ] Integrate Ollama for intent parsing
- [ ] Implement chatbot search endpoint
- [ ] Add decryption to chatbot response flow
- [ ] Implement confirmation prompts for sensitive operations
- [ ] Add rate limiting on search queries

### Phase 3: Enhancement
- [ ] Implement key rotation strategy
- [ ] Add searchable encryption for specific fields (optional)
- [ ] Implement advanced audit logging
- [ ] Add monitoring and alerting
- [ ] Performance optimization for large datasets

### Phase 4: Production
- [ ] Integrate with secure key management system
- [ ] Implement comprehensive security testing
- [ ] Add penetration testing
- [ ] Document security procedures
- [ ] Implement disaster recovery procedures

---

## 9. Testing Strategy

### 9.1 Unit Tests

```go
// Test encryption/decryption roundtrip
func TestEncryptDecryptRoundtrip(t *testing.T) {
    key := generateTestKey()
    plaintext := "my-secret-password"
    
    encrypted, err := EncryptSecret(plaintext, key)
    assert.NoError(t, err)
    
    decrypted, err := DecryptSecret(encrypted, key)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}

// Test tampering detection
func TestTamperingDetection(t *testing.T) {
    key := generateTestKey()
    plaintext := "my-secret-password"
    
    encrypted, _ := EncryptSecret(plaintext, key)
    
    // Tamper with ciphertext
    tampered := tamperWithCiphertext(encrypted)
    
    _, err := DecryptSecret(tampered, key)
    assert.Error(t, err) // Should fail
}

// Test wrong key rejection
func TestWrongKeyRejection(t *testing.T) {
    key1 := generateTestKey()
    key2 := generateTestKey()
    plaintext := "my-secret-password"
    
    encrypted, _ := EncryptSecret(plaintext, key1)
    
    _, err := DecryptSecret(encrypted, key2)
    assert.Error(t, err) // Should fail
}
```

### 9.2 Integration Tests

- Test secret creation and retrieval from MongoDB
- Test search functionality with encrypted data
- Test chatbot integration with decryption
- Test audit logging

### 9.3 Security Tests

- Test key management procedures
- Test access control enforcement
- Test rate limiting
- Test session timeout

---

## 10. Future Enhancements

1. **Searchable Encryption**: Implement deterministic encryption for searchable fields
2. **Key Rotation**: Automated key rotation with versioning
3. **Field-Level Encryption**: Different encryption for different field types
4. **Homomorphic Encryption**: Search without decryption (advanced)
5. **Hardware Security Module (HSM)**: Store keys in HSM for production
6. **Multi-Key Support**: Different keys for different users or data classifications
7. **Encryption at Application Level**: Encrypt before sending to database

---

## 11. References

- [Go crypto/aes Documentation](https://pkg.go.dev/crypto/aes)
- [Go crypto/cipher Documentation](https://pkg.go.dev/crypto/cipher)
- [NIST SP 800-38D: GCM Mode](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf)
- [OWASP Cryptographic Failures](https://owasp.org/Top10/A02_2021-Cryptographic_Failures/)
- [MongoDB Queryable Encryption](https://www.mongodb.com/docs/manual/core/queryable-encryption/)

---

## 12. Appendix: Environment Configuration

### Development (.env)
```
ENCRYPTION_KEY=your-32-byte-base64-encoded-key-here
MONGODB_URI=mongodb://localhost:27017
OLLAMA_API_URL=http://localhost:11434
```

### Production (.env.production)
```
ENCRYPTION_KEY=${AWS_SECRETS_MANAGER_KEY}
MONGODB_URI=${MONGODB_ATLAS_URI}
OLLAMA_API_URL=${OLLAMA_ENDPOINT}
```

---

**Document Version**: 1.0
**Last Updated**: November 2025
**Status**: Draft
