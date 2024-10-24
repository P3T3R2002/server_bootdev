package main

import(
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/P3T3R2002/server_bootdev/internal/auth"
	"github.com/P3T3R2002/server_bootdev/internal/database"
)

func (cfg *apiConfig)registerUsers(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    } else {
		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondError(w, "Something went wrong", 400) 
			return
		}
		user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			ID: uuid.New(),
			Email: 			params.Email,
			HashedPassword:	hashedPassword,

		})
		if err != nil {
			respondError(w, "Something went wrong", 500)
		}
		respondUser(w, user, 201)
	}
}

//------------------------------------------------------------

func (cfg *apiConfig)updateUsers(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondError(w, "Something went wrong", 400)
		return
    } else {
		tokenString, err := auth.GetBearerToken(r.Header)
		ID, err := auth.ValidateJWT(tokenString, cfg.JWT_secret)
		if err != nil {
			respondError(w, "", 401) 
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondError(w, "Something went wrong", 400) 
			return
		}
		
		user, err := cfg.db.UpdatePassword(r.Context(), database.UpdatePasswordParams{
			ID: ID,
			Email: params.Email,
			HashedPassword:	hashedPassword,

		})
		if err != nil {
			respondError(w, "Something went wrong", 500)
			return
		}

		respondUser(w, user, 200)
	}
}

//------------------------------------------------------------

func (cfg *apiConfig)loginUser(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    } else {
		user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondError(w, "Something went wrong", 500)
			return
		}

		err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
		if err != nil {
			respondError(w, "Wrong password!", 401) 
			return
		}
		
		accessToken, err := auth.MakeJWT(user.ID, cfg.JWT_secret, time.Hour)
		if err != nil {
			respondError(w, "Problem creating JTW!", 401) 
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondError(w, "Problem creating reftesh token!", 401) 
			return
		}

		_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token: 			refreshToken,
			UserID: 		user.ID,
			ExpiresAt:		time.Now().Add(time.Hour*24*60),
		})
		if err != nil {
			respondError(w, "Something went wrong", 500)
			return
		}

		respondUserWithTokens(w, user, accessToken, refreshToken, 200)
	}
}

//****************//

func respondUser(w http.ResponseWriter, u database.User, code int) {
    type returnVals struct {
		ID        	uuid.UUID 	`json:"id"`
		CreatedAt 	time.Time 	`json:"created_at"`
		UpdatedAt 	time.Time 	`json:"updated_at"`
		Email     	string 		`json:"email"`
		IsChirpyRed	bool		`json:"is_chirpy_red"`
    }
	
    respBody := returnVals{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
    }
	respondJson(w, respBody, code)
}

//**********

func respondUserWithTokens(w http.ResponseWriter, u database.User, accessToken string, refreshToken string, code int) {
    type returnVals struct {
		ID        		uuid.UUID 	`json:"id"`
		CreatedAt 		time.Time 	`json:"created_at"`
		UpdatedAt 		time.Time 	`json:"updated_at"`
		Email     		string 		`json:"email"`
		IsChirpyRed		bool		`json:"is_chirpy_red"`
		Token	  		string		`json:"token"`
		RefreshToken	string  	`json:"refresh_token"`
    }
	
    respBody := returnVals{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
		Token: accessToken,
		RefreshToken: refreshToken,
    }
	respondJson(w, respBody, code)
}