package main

import (
	"time"
	"net/http"
	"github.com/P3T3R2002/server_bootdev/internal/auth"
)

func (cfg *apiConfig)handlerRefresh(w http.ResponseWriter, r *http.Request) {
    type response struct {
		Token string `json:"token"`
    }

	refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    }

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
    if err != nil {
		respondError(w, "", 401) 
		return
    }

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWT_secret, time.Hour)
    if err != nil {
		respondError(w, "", 401) 
		return
    }


	respondJson(w, response{Token: accessToken}, 200)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    }

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
    if err != nil {
		respondError(w, "Something went wrong", 400) 
		return
    }

	w.WriteHeader(204)
}