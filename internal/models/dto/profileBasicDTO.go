package dto

type ProfileBasicDTO struct {
	ProfileID      uint    `json:"profileId"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
}
