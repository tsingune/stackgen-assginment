package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a scheduled event
type Event struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	OrganizerId uint           `json:"organizer_id" gorm:"not null"`
	Duration    int            `json:"duration" gorm:"not null"` // in minutes
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
	TimeSlots   []TimeSlot     `json:"time_slots,omitempty" gorm:"foreignKey:EventID"`
}

// TimeSlot represents a potential time for an event
type TimeSlot struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	EventID   uint           `json:"event_id" gorm:"not null"`
	StartTime time.Time      `json:"start_time" gorm:"not null"`
	EndTime   time.Time      `json:"end_time" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// Participant represents a user who can participate in events
type Participant struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null;unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// Availability represents a participant's availability for a time slot
type Availability struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ParticipantID uint           `json:"participant_id" gorm:"not null"`
	TimeSlotID    uint           `json:"time_slot_id" gorm:"not null"`
	IsAvailable   bool           `json:"is_available" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

// TimeSlotRecommendation represents a recommended time slot with availability information
type TimeSlotRecommendation struct {
	TimeSlot         TimeSlot      `json:"time_slot"`
	AvailableCount   int           `json:"available_count"`
	UnavailableCount int           `json:"unavailable_count"`
	AvailableUsers   []Participant `json:"available_users,omitempty"`
	UnavailableUsers []Participant `json:"unavailable_users,omitempty"`
}
