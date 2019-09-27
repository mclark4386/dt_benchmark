package middleware

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mclark4386/dt_benchmark/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

const tokenIssuerPrefix = "org.dt_benchmark.jwt"
const envJWTKeyPath = "JWT_KEY_PATH"
const shortExpTime = 10 // minutes
const longExpTime = 30  // days
const CurrentUser = "current_user"

var signingMethod = jwt.SigningMethodHS256

type BackendClaims struct {
	RefreshToken string `json:"rt"`
	Expiration   int64  `json:"accExp"`
	jwt.StandardClaims
}

// TokenMiddleware validates our JWT, handles regeneration as needed
func TokenMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		headToken := GetToken(c)

		if len(headToken) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
		}

		// Check the signing method, parse into token
		token, err := jwt.Parse(headToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey(), nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not parse token, %v", err))
		}

		u := &models.User{}
		tx := c.Value("tx").(*pop.Connection)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}

		uID := claims["jti"].(string)
		shortExp := int64(claims["accExp"].(float64))
		rTok := claims["rt"].(string)
		rExp := int64(claims["exp"].(float64))

		err = tx.Find(u, uID)
		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not identify the user"))
		}

		if int64(shortExp) < time.Now().Unix() {
			if u.RefreshToken.UUID.String() != rTok {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid token, please log in again"))
			}

			newToken, err := BuildJWT(u, rExp)
			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Error building token, %v", err))
			}
			SetToken(c, newToken)
		}

		c.Set(CurrentUser, u)

		spew.Printf("logged in as: %+v\n", u)

		return next(c)
	}
}

// BuildJWT builds our JWT
// standard claims expiration is the short expiration
func BuildJWT(u *models.User, exp int64) (string, error) {
	if exp <= 0 {
		exp = time.Now().AddDate(0, 0, longExpTime).Unix()
	}

	standardClaims := jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    fmt.Sprintf("%s.%s", tokenIssuerPrefix, envy.Get("GO_ENV", "development")),
		Id:        u.ID.String(),
	}

	claims := BackendClaims{
		u.RefreshToken.UUID.String(),
		time.Now().Add(tokenExpiration()).Unix(),
		standardClaims,
	}

	token := jwt.NewWithClaims(signingMethod, claims)

	tokenString, err := token.SignedString(signingKey())

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetToken(c buffalo.Context) string {
	token := c.Request().Header.Get("Authorization")
	if len(token) == 0 {
		sessionToken := c.Session().Get("access_token")
		if sessionToken == nil {
			return ""
		}
		token = sessionToken.(string)
	}
	return token
}

func SetToken(c buffalo.Context, token string) {
	if RequestIsJSON(c) {
		c.Response().Header().Set("Access-Token", token)
	} else {
		c.Session().Set("access_token", token)
	}
}

func RequestIsJSON(c buffalo.Context) bool {
	ct, ok := c.Value("contentType").(string)
	if !ok {
		return false
	}

	ct = strings.ToLower(ct)
	return strings.Contains(ct, "json")
}

func tokenExpiration() time.Duration {
	return shortExpTime * time.Minute
}

func signingKey() []byte {
	path := envy.Get(envJWTKeyPath, "")
	devKey := []byte("developmentSigningKey")

	if path == "" {
		return devKey
	}

	signingKey, err := ioutil.ReadFile(path)
	if err != nil {
		return devKey
	}

	return signingKey
}

func VerifyCode() string {
	i := rand.Intn(1000000)
	code := strconv.Itoa(i)

	for len(code) < 6 {
		code = "0" + code
	}

	return code
}
