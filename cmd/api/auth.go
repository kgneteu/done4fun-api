package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os/user"
	"strconv"
	"time"
)

const (
	accessTokenCookieName = "X-Auth"
	jwtSecretKey          = "some-secret-key"
)



//Auth Middleware
func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenCookie, _ := c.Cookie(accessTokenCookieName)
		//todo support for auth headers
		//tokenHeader.

		claims := &jwt.StandardClaims{}
		if tokenCookie != nil {
			//token = decodeJWTToken
			_, err := jwt.ParseWithClaims(tokenCookie.Value, claims,
				func(token *jwt.Token) (interface{}, error) {
					return getJWTSigningKey(), nil
				})
			if err != nil {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "invalid token"})
			}

			c.Set("User", map[string]interface{}{
				"Id":    claims.Id,
				"Token": tokenCookie.Value,
			})
		}
		return next(c)
	}
}

func getJWTSigningKey() []byte {
	return []byte(jwtSecretKey)
}

func createToken(id uint) (signedString string, err error) {
	//todo: secret for jwt & expiration
	claims := &jwt.StandardClaims{
		Id:        strconv.FormatUint(uint64(id), 10),
		IssuedAt: time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err = token.SignedString(getJWTSigningKey())
	return
}

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}






























func checkAdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("User")
		if user != nil {
			//role:=user.
			return next(c)
		}
		return c.JSON(http.StatusForbidden, echo.Map{"message": "invalid token"})
	}
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.

type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func GetJWTSecret() string {
	return jwtSecretKey
}

// GenerateTokensAndSetCookies generates jwt token and saves it to the http-only cookie.
func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)

	return nil
}

func generateAccessToken(user *user.User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *user.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// Purpose of this cookie is to store the user's name.
func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(err error, c echo.Context) error {
	// Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
