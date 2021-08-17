package jwt_service

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/energy-uktc/eventpool-api/config"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/utils"
)

const (
	SCOPE_BASIC = "Basic"
)

var (
	jwtPrivateKey *rsa.PrivateKey
	jwtPublicKey  *rsa.PublicKey
)

type customClaims struct {
	*jwt.StandardClaims
	CustomerInfo customerInfo `json:"user"`
	Scopes       []string     `json:"scopes"`
}
type customerInfo struct {
	Id string `json:"id"`
}

func (c customClaims) Valid() error {
	return c.StandardClaims.Valid()
}

func init() {
	privateKeyPath := readJwtConfig("privateKey").(string)
	publicKeyPath := readJwtConfig("publicKey").(string)

	signBytes, err := ioutil.ReadFile(privateKeyPath)
	fatal(err)

	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	fatal(err)

	jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func CreateToken(userId string, scopes []string) (models.GeneratedToken, error) {
	if len(scopes) == 0 {
		scopes = append(scopes, SCOPE_BASIC)
	}
	token := jwt.New(jwt.SigningMethodRS256)
	tokenId := utils.GenerateString(32, utils.ALPHA_NUMERIC)
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	token.Claims = &customClaims{
		&jwt.StandardClaims{
			Id:        tokenId,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt,
		},
		customerInfo{Id: userId},
		scopes,
	}
	tokenString, err := token.SignedString(jwtPrivateKey)
	if err != nil {
		return models.GeneratedToken{}, err
	}
	return models.GeneratedToken{
		Token:          tokenString,
		TokenType:      "bearer",
		TokenId:        tokenId,
		ExpirationTime: expiresAt,
		Scope:          strings.Join(scopes, ","),
	}, nil
}

func VerifyToken(tokenString string, skipClaimsValidation bool) (*jwt.Token, error) {
	parser := &jwt.Parser{SkipClaimsValidation: skipClaimsValidation, ValidMethods: []string{jwt.SigningMethodRS256.Alg()}}
	token, err := parser.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtPublicKey, nil
	})
	return token, err
}

func GetClaims(token *jwt.Token) *customClaims {
	return token.Claims.(*customClaims)
}
func readJwtConfig(key string) interface{} {
	value := config.GetPath("certificates.jwt." + key)
	if value == nil {
		log.Fatalf("JWT %s is not configured", key)
	}
	return value
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
