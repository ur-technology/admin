package user

type Registration struct {
	AnnouncementFinalizedAt Timestamp `json:"announcementFinalizedAt"`
	Status                  string    `json:"status"`
}
