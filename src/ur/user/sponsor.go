package user

type Sponsor struct {
	AnnouncementTransactionConfirmed bool   `json:"announcementTransactionConfirmed"`
	Name                             string `json:"name"`
	ProfilePhotoURL                  string `json:"profilePhotoUrl"`
	UserID                           string `json:"userId"`
}
