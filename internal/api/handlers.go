package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tusharsingune/meeting-scheduler/internal/logger"
	"github.com/tusharsingune/meeting-scheduler/internal/models"
	"github.com/tusharsingune/meeting-scheduler/internal/repository"
	"go.uber.org/zap"
)

// Handler handles HTTP requests
type Handler struct {
	Repo repository.Repository
	Log  *zap.Logger
}

// NewHandler creates a new Handler instance
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
		Log:  logger.GetLogger(),
	}
}

// RegisterHandlers registers all API routes
func RegisterHandlers(r *mux.Router, repo repository.Repository) {
	h := NewHandler(repo)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Debug endpoint to check database connectivity
	r.HandleFunc("/debug/db", func(w http.ResponseWriter, r *http.Request) {
		// Try to get a participant to test DB connection
		_, err := repo.GetParticipant(1)
		if err != nil {
			// Create a test participant if not found
			testParticipant := &models.Participant{
				Name:  "Test User",
				Email: "test@example.com",
			}
			err = repo.CreateParticipant(testParticipant)
			if err != nil {
				respondWithJSON(w, http.StatusInternalServerError, map[string]string{
					"status":  "error",
					"message": "Database connection failed: " + err.Error(),
				})
				return
			}
		}

		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Database connection successful",
		})
	})

	// API v1 subrouter
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Event endpoints
	events := v1.PathPrefix("/events").Subrouter()
	events.HandleFunc("", h.CreateEvent).Methods(http.MethodPost)
	events.HandleFunc("/{id}", h.GetEvent).Methods(http.MethodGet)
	events.HandleFunc("/{id}", h.UpdateEvent).Methods(http.MethodPut)
	events.HandleFunc("/{id}", h.DeleteEvent).Methods(http.MethodDelete)

	// Time slots
	events.HandleFunc("/{id}/timeslots", h.AddTimeSlot).Methods(http.MethodPost)
	events.HandleFunc("/{id}/timeslots", h.GetTimeSlots).Methods(http.MethodGet)

	// Availability
	events.HandleFunc("/{id}/availability", h.SubmitAvailability).Methods(http.MethodPost)
	events.HandleFunc("/{id}/recommendations", h.GetRecommendations).Methods(http.MethodGet)

	// Participant endpoints
	participants := v1.PathPrefix("/participants").Subrouter()
	participants.HandleFunc("", h.CreateParticipant).Methods(http.MethodPost)
	participants.HandleFunc("/{id}", h.GetParticipant).Methods(http.MethodGet)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// CreateEvent handles the creation of a new event
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.Repo.CreateEvent(&event); err != nil {
		h.Log.Error("Failed to create event", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Event created successfully", zap.Uint("event_id", event.ID))
	respondWithJSON(w, http.StatusCreated, event)
}

// GetEvent handles the retrieval of an event by ID
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	event, err := h.Repo.GetEvent(uint(id))
	if err != nil {
		h.Log.Error("Failed to get event", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Event retrieved successfully", zap.Uint("event_id", event.ID))
	respondWithJSON(w, http.StatusOK, event)
}

// UpdateEvent handles updating an existing event
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var event models.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	event.ID = uint(id)
	if err := h.Repo.UpdateEvent(&event); err != nil {
		h.Log.Error("Failed to update event", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Event updated successfully", zap.Uint("event_id", event.ID))
	respondWithJSON(w, http.StatusOK, event)
}

// DeleteEvent handles deleting an event
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := h.Repo.DeleteEvent(uint(id)); err != nil {
		h.Log.Error("Failed to delete event", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Event deleted successfully", zap.Uint64("event_id", id))
	w.WriteHeader(http.StatusNoContent)
}

// AddTimeSlot handles adding a time slot to an event
func (h *Handler) AddTimeSlot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var timeSlot models.TimeSlot
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&timeSlot); err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	timeSlot.EventID = uint(eventID)
	if err := h.Repo.CreateTimeSlot(&timeSlot); err != nil {
		h.Log.Error("Failed to create time slot", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Time slot created successfully",
		zap.Uint("event_id", timeSlot.EventID),
		zap.Uint("timeslot_id", timeSlot.ID))
	respondWithJSON(w, http.StatusCreated, timeSlot)
}

// GetTimeSlots handles retrieving all time slots for an event
func (h *Handler) GetTimeSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	timeSlots, err := h.Repo.GetTimeSlots(uint(eventID))
	if err != nil {
		h.Log.Error("Failed to get time slots", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Time slots retrieved successfully", zap.Uint64("event_id", eventID))
	respondWithJSON(w, http.StatusOK, timeSlots)
}

// SubmitAvailability handles submitting availability for a participant
func (h *Handler) SubmitAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var availability models.Availability
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&availability); err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	availability.ID = uint(eventID)
	if err := h.Repo.CreateAvailability(&availability); err != nil {
		h.Log.Error("Failed to create availability", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Availability submitted successfully",
		zap.Uint("event_id", availability.ID),
		zap.Uint("participant_id", availability.ParticipantID))
	respondWithJSON(w, http.StatusCreated, availability)
}

// GetRecommendations handles retrieving time slot recommendations for an event
func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid event ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	recommendations, err := h.Repo.GetTimeSlotRecommendations(uint(eventID))
	if err != nil {
		h.Log.Error("Failed to get recommendations", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Recommendations retrieved successfully", zap.Uint64("event_id", eventID))
	respondWithJSON(w, http.StatusOK, recommendations)
}

// CreateParticipant handles creating a new participant
func (h *Handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	var participant models.Participant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&participant); err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.Repo.CreateParticipant(&participant); err != nil {
		h.Log.Error("Failed to create participant", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Participant created successfully", zap.Uint("participant_id", participant.ID))
	respondWithJSON(w, http.StatusCreated, participant)
}

// GetParticipant handles retrieving a participant by ID
func (h *Handler) GetParticipant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.Log.Error("Invalid participant ID", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid participant ID")
		return
	}

	participant, err := h.Repo.GetParticipant(uint(id))
	if err != nil {
		h.Log.Error("Failed to get participant", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Info("Participant retrieved successfully", zap.Uint("participant_id", participant.ID))
	respondWithJSON(w, http.StatusOK, participant)
}
