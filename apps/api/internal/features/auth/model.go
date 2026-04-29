package auth

type LoginRequest struct { Email string `json:"email"`; Password string `json:"password"` }
type SignupRequest struct { Email string `json:"email"`; Username string `json:"username"`; Password string `json:"password"`; Keywords []string `json:"keywords"` }
