type Action struct {
	ActionID    uint      `json:"action_id"`
	UserID      uint     `json:"user_id"`
	ActionType  string    `json:"action_type"`
	TargetID    uint     `json:"target_id"`
	TargetType  string    `json:"target_type"`
	ActionDate  time.Time `json:"action_date" gorm:"column:action_date;type:timestamp;default:CURRENT_TIMESTAMP"`
}