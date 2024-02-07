package handlers

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	JwtExpiresIn int
}

// func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
// 	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
// 	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

// 	r.Context().Value("token")

// 	var user dto.GetJWTInput
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	u, err := h.UserDB.FindByEmail(user.Email)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		err := Error{Message: err.Error()}
// 		json.NewEncoder(w).Encode(err)
// 		return
// 	}

// 	if !u.ValidatePassword(user.Password) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	_, tokenString, _ := jwt.Encode(map[string]interface{}{
// 		"sub": u.ID.String(),
// 		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
// 	})
// 	acessToken := dto.GetJWTOutPut{AccessToken: tokenString}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(acessToken)

// 	// token := jwtauth.New("HS256", []byte("secret"), nil)
// 	// _, tokenString, _ := token.Encode(map[string]interface{}{"id": u.ID})
// 	//w.Write([]byte(tokenString))
// }
