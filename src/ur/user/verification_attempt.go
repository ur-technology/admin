package user

type VerificationAttempt struct {
	CreatedAt        Timestamp `json:"createdAt"`
	VerificationCode string    `json:"verificationCode"`
}

type VerificationAttempts map[string]VerificationAttempt
