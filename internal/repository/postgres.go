package repository

import (
	"fmt"
	"time"

	"github.com/tusharsingune/meeting-scheduler/internal/config"
	"github.com/tusharsingune/meeting-scheduler/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresRepository implements the Repository interface
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresDB creates a new PostgreSQL repository
func NewPostgresDB(cfg config.DatabaseConfig) (Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.Event{},
		&models.TimeSlot{},
		&models.Participant{},
		&models.Availability{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	repo := &PostgresRepository{db: db}

	return repo, nil
}

// CreateEvent creates a new event
func (r *PostgresRepository) CreateEvent(event *models.Event) error {
	return r.db.Create(event).Error
}

// GetEvent retrieves an event by ID
func (r *PostgresRepository) GetEvent(id uint) (*models.Event, error) {
	var event models.Event
	if err := r.db.Preload("TimeSlots").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// UpdateEvent updates an existing event
func (r *PostgresRepository) UpdateEvent(event *models.Event) error {
	return r.db.Save(event).Error
}

// DeleteEvent deletes an event by ID
func (r *PostgresRepository) DeleteEvent(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

// CreateTimeSlot creates a new time slot
func (r *PostgresRepository) CreateTimeSlot(slot *models.TimeSlot) error {
	return r.db.Create(slot).Error
}

// GetTimeSlots retrieves all time slots for an event
func (r *PostgresRepository) GetTimeSlots(eventID uint) ([]models.TimeSlot, error) {
	var slots []models.TimeSlot
	if err := r.db.Where("event_id = ?", eventID).Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

// CreateAvailability creates a new availability record
func (r *PostgresRepository) CreateAvailability(availability *models.Availability) error {
	return r.db.Create(availability).Error
}

// GetTimeSlotRecommendations returns recommended time slots for an event
func (r *PostgresRepository) GetTimeSlotRecommendations(eventID uint) ([]models.TimeSlotRecommendation, error) {
	var recommendations []models.TimeSlotRecommendation

	// Get all time slots for the event
	var timeSlots []models.TimeSlot
	if err := r.db.Where("event_id = ?", eventID).Find(&timeSlots).Error; err != nil {
		return nil, err
	}

	for _, slot := range timeSlots {
		var recommendation models.TimeSlotRecommendation
		recommendation.TimeSlot = slot

		// Get all availabilities for this time slot
		var availabilities []models.Availability
		if err := r.db.Where("time_slot_id = ?", slot.ID).Find(&availabilities).Error; err != nil {
			return nil, err
		}

		// Count available and unavailable participants
		for _, availability := range availabilities {
			if availability.IsAvailable {
				recommendation.AvailableCount++
				var participant models.Participant
				if err := r.db.First(&participant, availability.ParticipantID).Error; err != nil {
					return nil, err
				}
				recommendation.AvailableUsers = append(recommendation.AvailableUsers, participant)
			} else {
				recommendation.UnavailableCount++
				var participant models.Participant
				if err := r.db.First(&participant, availability.ParticipantID).Error; err != nil {
					return nil, err
				}
				recommendation.UnavailableUsers = append(recommendation.UnavailableUsers, participant)
			}
		}

		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

// CreateParticipant creates a new participant
func (r *PostgresRepository) CreateParticipant(participant *models.Participant) error {
	return r.db.Create(participant).Error
}

// GetParticipant retrieves a participant by ID
func (r *PostgresRepository) GetParticipant(id uint) (*models.Participant, error) {
	var participant models.Participant
	if err := r.db.First(&participant, id).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}
