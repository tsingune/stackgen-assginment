package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tusharsingune/meeting-scheduler/internal/models"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockRepository) GetEvent(id uint) (*models.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockRepository) UpdateEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockRepository) DeleteEvent(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateTimeSlot(slot *models.TimeSlot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockRepository) GetTimeSlots(eventID uint) ([]models.TimeSlot, error) {
	args := m.Called(eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TimeSlot), args.Error(1)
}

func (m *MockRepository) CreateAvailability(availability *models.Availability) error {
	args := m.Called(availability)
	return args.Error(0)
}

func (m *MockRepository) GetTimeSlotRecommendations(eventID uint) ([]models.TimeSlotRecommendation, error) {
	args := m.Called(eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TimeSlotRecommendation), args.Error(1)
}

func (m *MockRepository) CreateParticipant(participant *models.Participant) error {
	args := m.Called(participant)
	return args.Error(0)
}

func (m *MockRepository) GetParticipant(id uint) (*models.Participant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Participant), args.Error(1)
}

// setupTestHandler creates a handler with a mock repository for testing
func setupTestHandler(mockRepo *MockRepository) *Handler {
	logger, _ := zap.NewDevelopment()
	return &Handler{
		Repo: mockRepo,
		Log:  logger,
	}
}

func TestCreateEvent(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	// Test event
	event := &models.Event{
		Title:       "Test Event",
		Description: "Test Description",
		OrganizerId: 1,
		Duration:    60,
	}

	// Setup expectations
	mockRepo.On("CreateEvent", mock.AnythingOfType("*models.Event")).Return(nil)

	// Create request
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Execute request
	handler.CreateEvent(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify mock
	mockRepo.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	// Test event
	event := &models.Event{
		ID:          1,
		Title:       "Test Event",
		Description: "Test Description",
		OrganizerId: 1,
		Duration:    60,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Setup expectations
	mockRepo.On("GetEvent", uint(1)).Return(event, nil)

	// Create request
	req := httptest.NewRequest("GET", "/events/1", nil)
	w := httptest.NewRecorder()

	// Add URL parameters
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Execute request
	handler.GetEvent(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Event
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, event.ID, response.ID)
	assert.Equal(t, event.Title, response.Title)

	// Verify mock
	mockRepo.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	// Test event
	event := &models.Event{
		ID:          1,
		Title:       "Updated Event",
		Description: "Updated Description",
		OrganizerId: 1,
		Duration:    90,
	}

	// Setup expectations
	mockRepo.On("UpdateEvent", mock.AnythingOfType("*models.Event")).Return(nil)

	// Create request
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/events/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Add URL parameters
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Execute request
	handler.UpdateEvent(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify mock
	mockRepo.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	// Setup expectations
	mockRepo.On("DeleteEvent", uint(1)).Return(nil)

	// Create request
	req := httptest.NewRequest("DELETE", "/events/1", nil)
	w := httptest.NewRecorder()

	// Add URL parameters
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Execute request
	handler.DeleteEvent(w, req)

	// Assert response
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify mock
	mockRepo.AssertExpectations(t)
}

// Add more test cases for time slots and availability
func TestAddTimeSlot(t *testing.T) {
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	timeSlot := &models.TimeSlot{
		EventID:   1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	mockRepo.On("CreateTimeSlot", mock.AnythingOfType("*models.TimeSlot")).Return(nil)

	body, _ := json.Marshal(timeSlot)
	req := httptest.NewRequest("POST", "/events/1/timeslots", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.AddTimeSlot(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetTimeSlots(t *testing.T) {
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	timeSlots := []models.TimeSlot{
		{
			ID:        1,
			EventID:   1,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		},
	}

	mockRepo.On("GetTimeSlots", uint(1)).Return(timeSlots, nil)

	req := httptest.NewRequest("GET", "/events/1/timeslots", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.GetTimeSlots(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetRecommendations(t *testing.T) {
	mockRepo := new(MockRepository)
	handler := setupTestHandler(mockRepo)

	recommendations := []models.TimeSlotRecommendation{
		{
			TimeSlot: models.TimeSlot{
				ID:        1,
				EventID:   1,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			AvailableCount:   3,
			UnavailableCount: 1,
		},
	}

	mockRepo.On("GetTimeSlotRecommendations", uint(1)).Return(recommendations, nil)

	req := httptest.NewRequest("GET", "/events/1/recommendations", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.GetRecommendations(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
