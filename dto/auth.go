package dto

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceId     string `json:"device_id,omitempty"`
}

type TokenResponse struct {
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
	Tokens     TokenPair `json:"tokens"`
}

type LoginResponse struct {
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
	Data       LoginData `json:"data"`
	Tokens     TokenPair `json:"tokens"`
}

type LoginData struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status,omitempty"`
}

type BasicResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
