package user

import "fmt"

type Model struct {
	UID                  string               `json:"-"`
	Admin                bool                 `json:"admin"`
	CountryCode          string               `json:"countryCode"`
	CreatedAt            Timestamp            `json:"createdAt"`
	DownlineLevel        int                  `json:"downlineLevel"`
	DownlineSize         int                  `json:"downlineSize"`
	Email                string               `json:"email"`
	Events               Events               `json:"events"`
	IDHash               string               `json:"idHash"`
	IDRecognitionStatus  string               `json:"idRecognitionStatus"`
	IDUploaded           bool                 `json:"idUploaded"`
	IsEmailVerified      bool                 `json:"isEmailVerified"`
	LastName             string               `json:"lastName"`
	MiddleName           string               `json:"middleName"`
	Name                 string               `json:"name"`
	Phone                string               `json:"phone"`
	ProfilePhotoURL      string               `json:"profilePhotoUrl"`
	ReferralCode         string               `json:"referralCode"`
	Registration         Registration         `json:"registration"`
	SelfieConfidence     int                  `json:"selfieConfidence"`
	SelfieMatchStatus    string               `json:"selfieMatchStatus"`
	SelfieMatched        bool                 `json:"selfieMatched"`
	ServerHashedPassword string               `json:"serverHashedPassword"`
	SignUpBonusApproved  bool                 `json:"signUpBonusApproved"`
	SignedUpAt           Timestamp            `json:"signedUpAt"`
	Sponsor              Sponsor              `json:"sponsor"`
	Transactions         Transactions         `json:"transactions"`
	VerificationAttempts VerificationAttempts `json:"verificationAttempts"`
	VerificationCode     string               `json:"verificationCode"`
	Wallet               Wallet               `json:"wallet"`
}

func (m Model) String() string {
	return fmt.Sprintf("(model.User #%s: %s, %s)", m.UID, m.Phone, m.Name)
}
