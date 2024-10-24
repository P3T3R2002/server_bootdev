package main

import(
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/P3T3R2002/server_bootdev/internal/auth"
)

func (cfg *apiConfig)handleWebhooks(w http.ResponseWriter, r *http.Request) {
	ApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondError(w, "Something went wrong!", 500)
		return
	}
	if cfg.POLKA_KEY != ApiKey {
		respondError(w, "", 401)
		return
	}
	type data struct {
		User_id string `json:"user_id"`
	}
    type parameters struct {
		Event string `json:"event"`
		Data data `json:"data"`
    }
	
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
	if err != nil {
		respondError(w, "", 404)
		return
	} else if params.Event != "user.upgraded" {
		respondJson(w, nil, 204)
		return
	} else {
		ID, err := uuid.Parse(params.Data.User_id)
		err = cfg.db.UpgradeUser(r.Context(), ID)
		if err != nil {
			respondError(w, "", 404)
			return
		}
		respondJson(w, nil, 204)
		return
	}
}