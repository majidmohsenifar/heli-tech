package jwt

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type Service struct {
	privateKey []byte
	publicKey  []byte
	passphrase []byte
	ttl        int
}

func (s *Service) GenerateToken(username string) (string, error) {
	rsaPrivateKey, err := ssh.ParseRawPrivateKeyWithPassphrase(s.privateKey, []byte(s.passphrase))
	if err != nil {
		return "", err
	}
	t := jwt.New(jwt.SigningMethodRS256)
	return t.SignedString(rsaPrivateKey)
}

func (s *Service) GetUsernameFromToken(signedToken string) (string, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(s.publicKey)
		return rsaPublicKey, err
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//TODO: check if this is true
		return claims["sub"].(string), nil
	}
	return "", errors.New("invalid token")
}

func NewService(viper *viper.Viper) (*Service, error) {
	privateKeyPath := viper.GetString("jwt.privatekey")
	fmt.Println("privateKeyPath", privateKeyPath)
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	publicKeyPath := viper.GetString("jwt.publickey")
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	passphrase := viper.GetString("jwt.passphrase")
	ttl := viper.GetInt("jwt.ttl")
	return &Service{
		privateKey: privateKey,
		publicKey:  publicKey,
		passphrase: []byte(passphrase),
		ttl:        ttl,
	}, nil
}
