package main

import(
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/P3T3R2002/server_bootdev/internal/auth"
	"github.com/P3T3R2002/server_bootdev/internal/database"
)

func (cfg *apiConfig)postChirps(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body 	string 		`json:"body"`
    }
	tokenString, err := auth.GetBearerToken(r.Header)
	ID, err := auth.ValidateJWT(tokenString, cfg.JWT_secret)
	if err != nil {
		respondError(w, "", 401) 
		return
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    }
	if len(params.Body) > 140 {
		respondError(w, "Chirp is too long", 400)
		return
	} else {
		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			ID: 		uuid.New(),
			Body: 		params.Body,
			UserID:		ID,
		})
		if err != nil {
			respondError(w, "Something went wrong", 500)
			return
		}
		chirp_arr := []database.Chirp{chirp}
		respondChirps(w, chirp_arr, 201)
	}
}

//------------------------------------------------------------

func (cfg *apiConfig)deleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	ID, err := auth.ValidateJWT(tokenString, cfg.JWT_secret)
	if err != nil {
		respondError(w, "", 401) 
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}

	if chirp.UserID == ID {
		err = cfg.db.DeleteChirp(r.Context(), chirpID)
		if err != nil {
			respondError(w, "No chirp", 404)
			return
		}
	} else {
		respondError(w, "", 403)
		return
	}
	
	respondJson(w, nil, 204)
}

//------------------------------------------------------------

func (cfg *apiConfig)getChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondError(w, "", 404)
		return
	}

	chirp_arr := []database.Chirp{chirp}
	respondChirps(w, chirp_arr, 200)
}

//------------------------------------------------------------

func (cfg *apiConfig)getChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")
	ID, err := uuid.Parse(author)

	if order == "" || order == "asc" {
		if author == "" {
			cfg.getChirpsASC(w, r)
		} else {
			cfg.getChirpsByAuthorID_ASC(w, r, ID)
		}	
	} else if order == "desc" {
		if author == "" {
			cfg.getChirpsDESC(w, r)
		} else {
			cfg.getChirpsByAuthorID_DESC(w, r, ID)
		}
	} else {
		if err != nil {
			respondError(w, "Wrong sort parameter!", 400)
			return
		}
	}
}

//****************//

func respondChirps(w http.ResponseWriter, chirp []database.Chirp, code int) {
    type returnVals struct {
		ID        	uuid.UUID 	`json:"id"`
		CreatedAt 	time.Time 	`json:"created_at"`
		UpdatedAt 	time.Time 	`json:"updated_at"`
        Body 		string 		`json:"body"`
		User_id 	uuid.UUID 	`json:"user_id"`
    }
	respBody := []returnVals{}
	for _, c := range chirp {
		item := returnVals{
			ID: 		c.ID,
			CreatedAt: 	c.CreatedAt,
			UpdatedAt: 	c.UpdatedAt,
			Body: 		c.Body,
			User_id: 	c.UserID,
		}
		respBody = append(respBody, item)
	}
	if len(respBody) == 1 {
		respondJson(w, respBody[0], code)
		return
	} else {
		respondJson(w, respBody, code)
	}
}

//------------------------------------------------------------

func (cfg *apiConfig)getChirpsASC(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps_ASC(r.Context())
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}
	respondChirps(w, chirps, 200)
}

func (cfg *apiConfig)getChirpsDESC(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps_DESC(r.Context())
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}
	respondChirps(w, chirps, 200)
}

func (cfg *apiConfig)getChirpsByAuthorID_ASC(w http.ResponseWriter, r *http.Request, ID uuid.UUID) {
	chirps, err := cfg.db.GetChirpsByAuthorID_ASC(r.Context(), ID)
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}
	respondChirps(w, chirps, 200)
}

func (cfg *apiConfig)getChirpsByAuthorID_DESC(w http.ResponseWriter, r *http.Request, ID uuid.UUID) {
	chirps, err := cfg.db.GetChirpsByAuthorID_DESC(r.Context(), ID)
	if err != nil {
		respondError(w, "Something went wrong", 500)
		return
	}
	respondChirps(w, chirps, 200)
}