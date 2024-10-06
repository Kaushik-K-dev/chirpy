package main

// import ("net/http"; "encoding/json"; "strings")

// func (cfg *apiConfig) polkaHandler(w http.ResponseWriter, req *http.Request) {
// 	authHeader := req.Header.Get("Authorization")
//     if authHeader == "" || !strings.HasPrefix(authHeader, "ApiKey ") {
//         respError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
//         return
//     }
// 	providedKey := strings.TrimPrefix(authHeader, "ApiKey ")
// 	if providedKey != cfg.polkaKey {
//         respError(w, http.StatusUnauthorized, "Invalid API key")
//         return
//     }

// 	type WebhookRequest struct {
//         Event string `json:"event"`
//         Data  struct {
//             Id int `json:"user_id"`
//         } `json:"data"`
//     }

// 	var webhookRequest WebhookRequest
//     err := json.NewDecoder(req.Body).Decode(&webhookRequest)
//     if err != nil {
//         respError(w, http.StatusBadRequest, "Invalid request body")
//         return
//     }

// 	if webhookRequest.Event != "user.upgraded" {
//         w.WriteHeader(http.StatusNoContent)
//         return
//     }

// 	userId := webhookRequest.Data.Id
// 	user, err := cfg.DB.GetUser(userId)
// 	if err != nil {
//         respError(w, http.StatusNotFound, "User not found")
//         return
//     }

// 	user.IsChirpyRed = true
//     err = cfg.DB.UpdateUser(user)
//     if err != nil {
//         respError(w, http.StatusInternalServerError, "Failed to update user")
//         return
//     }

// 	w.WriteHeader(http.StatusNoContent)
// }