package dto

type VerificationRequest struct {
	Token string `json:"token" binding:"required"`
}
