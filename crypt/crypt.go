package crypt

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Crypt interface {
	Encrypt(plaintext string) ([]byte, error)
	Decrypt(ciphertext []byte) (string, error)
	HashSecret(plaintext string) ([]byte, error)
	HashPassword(password string) ([]byte, error)
	GenerateToken(userID string) (string, time.Time, error)
	VerifyToken(tokenString string) (*Claims, error)
	ComparePassword(hashedPassword []byte, password string) error
}

type CryptImplementation struct {
	Token          []byte
	SessionTimeout int
}
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ConfigureCrypt(ctx context.Context, token string, timeout int) (Crypt, error) {
	return &CryptImplementation{
		Token:          []byte(token),
		SessionTimeout: timeout,
	}, nil
}
func (c *CryptImplementation) Encrypt(plaintext string) ([]byte, error) {
	if len(c.Token) != 16 && len(c.Token) != 24 && len(c.Token) != 32 {
		return nil, errors.New("invalid key length; must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(c.Token)
	if err != nil {
		return nil, err
	}

	// Pad the plaintext to ensure it's a multiple of the block size
	paddedText := pad([]byte(plaintext), aes.BlockSize)

	// Create a ciphertext buffer with space for the IV and the padded plaintext
	ciphertext := make([]byte, aes.BlockSize+len(paddedText))

	// Generate a random IV and place it at the beginning of the ciphertext
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create a new CBC encrypter and encrypt the padded plaintext
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedText)

	return ciphertext, nil
}

// pad adds PKCS7 padding to the plaintext to ensure it is a multiple of the block size.
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (c *CryptImplementation) Decrypt(ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(c.Token)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func (c *CryptImplementation) HashSecret(plaintext string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(plaintext))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// hash password and return to store it in db
func (c *CryptImplementation) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}
	return hashedPassword, nil
}

// compare password with hashed password
func (c *CryptImplementation) ComparePassword(hashedPassword []byte, password string) error {

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("error comparing password: %v", err)
	}
	return nil
}

// generate token
func (c *CryptImplementation) GenerateToken(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(c.SessionTimeout) * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.Token))

	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

// verify token
func (c *CryptImplementation) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return c.Token, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("error parsing claims")
	}

	return claims, nil
}
