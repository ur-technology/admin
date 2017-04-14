package user

// This is sometimes a string and sometimes a bool!
type EventNotificationProcessed interface{}

type Event struct {
	CreatedAt             Timestamp                  `json:"createdAt"`
	MessageText           string                     `json:"messageText"`
	NotificationProcessed EventNotificationProcessed `json:"notificationProcessed"`
	ProfilePhotoURL       string                     `json:"profilePhotoUrl"`
	SourceID              string                     `json:"sourceId"`
	SourceType            string                     `json:"sourceType"`
	Title                 string                     `json:"title"`
	UpdatedAt             Timestamp                  `json:"updatedAt"`
}

type Events map[string]Event
