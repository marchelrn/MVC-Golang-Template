package dto

import "time"

type UserRequest struct {
	Id             uint   `json:"id,omitempty"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	Role           string `json:"role,omitempty"`
	AvatarUrl      string `json:"avatar_url,omitempty"`
	Status         string `json:"status,omitempty"`
	OldPassword    string `json:"old_password,omitempty"`
	NewPassword    string `json:"new_password,omitempty"`
	DeviceId       string `json:"device_id,omitempty"`
	RequestingRole string `json:"-"`
}

type UserResponse struct {
	StatusCode int        `json:"status_code"`
	Message    string     `json:"message"`
	Data       UserData   `json:"data,omitempty"`
	Tokens     *TokenPair `json:"tokens,omitempty"`
}

type UserAllResponse struct {
	Data   UserData   `json:"data,omitempty"`
	Tokens *TokenPair `json:"tokens,omitempty"`
}

type UserData struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status,omitempty"`
	AvatarUrl string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUpdateData struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status,omitempty"`
}

type UserDeleteData struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type UserUpdateResponse struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Data       UserUpdateData `json:"data"`
}

type UserDeleteResponse struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Data       UserDeleteData `json:"data"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
