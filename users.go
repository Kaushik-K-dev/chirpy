package main
import ("fmt"; "encoding/json"; "net/http"; "time"; "strconv"; "crypto/rand"; "encoding/hex")
import "golang.org/x/crypto/bcrypt"
import "github.com/golang-jwt/jwt/v5"

type User struct{
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"hashed_pass"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiration time.Time `json:"refresh_expiration"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type UserResp struct{
	Id int `json:"id"`
	Email string `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type LoginResp struct{
	Id int `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, password string) (*User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return nil, err}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return nil, fmt.Errorf("Email already exists")
		}
	}

	id := len(dbStruct.Users) + 1
	hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {return nil, err}
	hashedPass := string(hashedPassBytes)

	user := User{
		Id:    id,
		Email: email,
		Password: hashedPass,
		RefreshToken: "",
		RefreshTokenExpiration: time.Time{},
		IsChirpyRed: false,
	}
	dbStruct.Users[id] = &user

	err = db.writeDB(dbStruct)
	if err != nil {return nil, err}
	return &user, nil
}

func (db *DB) GetUser(id int) (*User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return nil, err}

	user, ok := dbStruct.Users[id]
	if !ok {return nil, fmt.Errorf("Error: User doesn't exist")}

	return user, nil
}

func (db *DB) GetUserbyEmail(email string) (*User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return nil, err}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("Error: User doesn't exist")
}

func (db *DB) GetUserbyRefreshToken(token string) (*User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return nil, err}

	for _, user := range dbStruct.Users {
		if user.RefreshToken == token {
			return user, nil
		}
	}
	return nil, fmt.Errorf("Error: User doesn't exist or Refresh Token doesn't exist or has expired")
}

func (db *DB) StoreRefreshToken(id int, refreshToken string, expiration time.Time) error {
	dbStruct, err := db.loadDB()
	if err != nil {return err}

	user, ok := dbStruct.Users[id]
    if !ok {return fmt.Errorf("Error: User not found")}

    user.RefreshToken = refreshToken
    user.RefreshTokenExpiration = expiration

    dbStruct.Users[id] = user
    return db.writeDB(dbStruct)
}

func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Couldn't create User")
		return
	}

	respJson(w, http.StatusCreated, UserResp{
		Id:    user.Id,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

type loginRequest struct {
	Password 		 string `json:"password"`
    Email    		 string `json:"email"`
	ExpiresInSeconds *int64  `json:"expires_in_seconds"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
    var loginReq loginRequest
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&loginReq); err != nil {
        respError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    user, err := cfg.DB.GetUserbyEmail(loginReq.Email)
    if err != nil {
        respError(w, http.StatusUnauthorized, "User not found")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
        respError(w, http.StatusUnauthorized, "Invalid password")
        return
    }

	// var expirationTime time.Duration
	// if loginReq.ExpiresInSeconds == nil {
	// 	expirationTime = 24 * time.Hour
	// } else {
	// 	expirationTime = time.Duration(*loginReq.ExpiresInSeconds) * time.Second
	// 	if expirationTime > 24 * time.Hour {
	// 		expirationTime = 24 * time.Hour
	// 	}
	// }

	expiration := 1 * time.Hour 
	tokenString, err := jwtTokenGen(user.Id, cfg.jwtSecret, expiration)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
    if err != nil {
        respError(w, http.StatusInternalServerError, "Could not generate refresh token")
        return
    }
    refreshToken := hex.EncodeToString(refreshTokenBytes)

	RefreshTokenExpiration := time.Now().Add(60 * 24 * time.Hour)
	err = cfg.DB.StoreRefreshToken(user.Id, refreshToken, RefreshTokenExpiration)
    if err != nil {
        respError(w, http.StatusInternalServerError, "Could not store refresh token")
        return
    }

    respJson(w, http.StatusOK, LoginResp{
		Id:    user.Id,
		Email: user.Email,
		Token: tokenString,
		RefreshToken: refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (db *DB) UpdateUser(user *User) error {
    dbStruct, err := db.loadDB()
    if err != nil {return err}

    dbStruct.Users[user.Id] = user

    return db.writeDB(dbStruct)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
    tokenString := TokenfromHeader(w, req)
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.jwtSecret), nil
    })
    if err != nil || !token.Valid {
        respError(w, http.StatusUnauthorized, "Invalid or expired token")
        return
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok {
        respError(w, http.StatusUnauthorized, "Failed to extract claims from token")
        return
    }

    userId := claims.Subject

	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

    update := parameters{}
    err = json.NewDecoder(req.Body).Decode(&update)
    if err != nil {
        respError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if update.Email == "" || update.Password == "" {
        respError(w, http.StatusBadRequest, "Email and password cannot be empty")
        return
    }

	userIdint, _ := strconv.Atoi(userId)
    user, err := cfg.DB.GetUser(userIdint)
    if err != nil {
        respError(w, http.StatusNotFound, "User not found")
        return
    }

    user.Email = update.Email
    hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(update.Password), bcrypt.DefaultCost)
    if err != nil {
        respError(w, http.StatusInternalServerError, "Failed to hash password")
        return
    }
    user.Password = string(hashedPassBytes)

    err = cfg.DB.UpdateUser(user)
    if err != nil {
        respError(w, http.StatusInternalServerError, "Failed to update user")
        return
    }

    respJson(w, http.StatusOK, UserResp{
		Id:    user.Id,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	refreshToken := TokenfromHeader(w, req)

	user, err := cfg.DB.GetUserbyRefreshToken(refreshToken)
    if err != nil || user.RefreshTokenExpiration.Before(time.Now()) {
        respError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
        return
    }

	expiration := 1 * time.Hour 
	tokenString, err := jwtTokenGen(user.Id, cfg.jwtSecret, expiration)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respJson(w, http.StatusOK, map[string]string{
        "token": tokenString,
    })
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {
	refreshToken := TokenfromHeader(w, req)

	user, err := cfg.DB.GetUserbyRefreshToken(refreshToken)
    if err != nil || user.RefreshTokenExpiration.Before(time.Now()) {
        respError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
        return
    }

	user.RefreshToken = ""
    user.RefreshTokenExpiration = time.Time{}
	err = cfg.DB.UpdateUser(user)
    if err != nil {
        respError(w, http.StatusInternalServerError, "Failed to revoke token")
        return
    }
	w.WriteHeader(http.StatusNoContent)
}