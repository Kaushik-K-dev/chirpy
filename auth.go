package main
import ("strings"; "strconv"; "time"; "net/http")
import "github.com/golang-jwt/jwt/v5"

func jwtTokenGen(Id int, jwtSecret string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),// expirationTime of 1 hour
		Subject:   strconv.Itoa(Id),
	})

	return token.SignedString([]byte(jwtSecret))
}

func TokenfromHeader(w http.ResponseWriter, req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        respError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
        return ""
    }

	return strings.TrimPrefix(authHeader, "Bearer ")
}