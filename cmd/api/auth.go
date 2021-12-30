package main

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

//todo secret key .env
const (
	jwtSecretKey = "some-secret-key"
)

func getJWTSigningKey() []byte {
	return []byte(jwtSecretKey)
}

func createToken(id uint) (signedString string, err error) {
	//todo: secret for jwt & expiration .env
	claims := &jwt.StandardClaims{
		Id:        strconv.FormatUint(uint64(id), 10),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute*1 + 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err = token.SignedString(getJWTSigningKey())
	return
}

func createRefreshToken(id uint, email string, password string) (signedString string, err error) {

	type RefreshClaims struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	claims := RefreshClaims{
		email,
		password,
		jwt.StandardClaims{
			Id:        strconv.FormatUint(uint64(id), 10),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			//ExpiresAt: time.Now().Add(time.Hour * 180).Unix(), never expires
		},
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

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
//
//type Claims struct {
//	Name string `json:"name"`
//	jwt.StandardClaims
//}
//
//func GetJWTSecret() string {
//	return jwtSecretKey
//}
//
//// GenerateTokensAndSetCookies generates jwt token and saves it to the http-only cookie.
//func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {
//	accessToken, exp, err := generateAccessToken(user)
//	if err != nil {
//		return err
//	}
//
//	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
//	setUserCookie(user, exp, c)
//
//	return nil
//}
//
//func generateAccessToken(user *user.User) (string, time.Time, error) {
//	// Declare the expiration time of the token (1h).
//	expirationTime := time.Now().Add(1 * time.Hour)
//
//	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
//}
//
//// Pay attention to this function. It holds the main JWT token generation logic.
//func generateToken(user *user.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
//	// Create the JWT claims, which includes the username and expiry time.
//	claims := &Claims{
//		Name: user.Name,
//		StandardClaims: jwt.StandardClaims{
//			// In JWT, the expiry time is expressed as unix milliseconds.
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}
//
//	// Declare the token with the HS256 algorithm used for signing, and the claims.
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	// Create the JWT string.
//	tokenString, err := token.SignedString(secret)
//	if err != nil {
//		return "", time.Now(), err
//	}
//
//	return tokenString, expirationTime, nil
//}
//
//// Here we are creating a new cookie, which will store the valid JWT token.
//func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
//	cookie := new(http.Cookie)
//	cookie.Name = name
//	cookie.Value = token
//	cookie.Expires = expiration
//	cookie.Path = "/"
//	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
//	cookie.HttpOnly = true
//
//	c.SetCookie(cookie)
//}
//
//// Purpose of this cookie is to store the user's name.
//func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
//	cookie := new(http.Cookie)
//	cookie.Name = "user"
//	cookie.Value = user.Name
//	cookie.Expires = expiration
//	cookie.Path = "/"
//	c.SetCookie(cookie)
//}
//
//// JWTErrorChecker will be executed when user try to access a protected path.
//func JWTErrorChecker(err error, c echo.Context) error {
//	// Redirects to the signIn form.
//	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
//}
