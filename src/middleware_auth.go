package main

import (
	"fmt"
	"net/http"

	"github.com/CLinnSheng/RSSagg/auth"
	"github.com/CLinnSheng/RSSagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// middleware function to get user through API Key
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldnt get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
