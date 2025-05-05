package dto

type ProfileBasicDTO struct {
	UserID         uint    `json:"userId"`
	ProfileID      uint    `json:"profileId"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
}
