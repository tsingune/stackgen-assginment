package repository

import (
	"github.com/tusharsingune/meeting-scheduler/internal/models"
)

// Repository defines the interface for data operations
type Repository interface {
	// Event operations
	CreateEvent(*models.Event) error
	GetEvent(uint) (*models.Event, error)
	UpdateEvent(*models.Event) error
	DeleteEvent(uint) error

	// TimeSlot operations
	CreateTimeSlot(*models.TimeSlot) error
	GetTimeSlots(uint) ([]models.TimeSlot, error)

	// Availability operations
	CreateAvailability(*models.Availability) error
	GetTimeSlotRecommendations(uint) ([]models.TimeSlotRecommendation, error)

	// Participant operations
	CreateParticipant(*models.Participant) error
	GetParticipant(uint) (*models.Participant, error)
}
