package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"

	"github.com/gorilla/mux"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	
	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	room.CreatorID = claims.UserID
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	db := database.GetDB()
	err := db.QueryRow(`
		INSERT INTO rooms (name, description, creator_id, private, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		room.Name, room.Description, room.CreatorID,
		room.Private, room.CreatedAt, room.UpdatedAt).Scan(&room.ID)

	if err != nil {
		http.Error(w, "Error creating room", http.StatusInternalServerError)
		return
	}

	// Add creator as room admin
	_, err = db.Exec(`
		INSERT INTO room_members (room_id, user_id, role, joined_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		room.ID, claims.UserID, "admin", time.Now(), time.Now())

	if err != nil {
		http.Error(w, "Error adding room member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT r.id, r.name, r.description, r.creator_id, r.private, r.created_at, r.updated_at
		FROM rooms r
		LEFT JOIN room_members rm ON r.id = rm.room_id
		WHERE r.private = false OR rm.user_id = $1
		GROUP BY r.id`, claims.UserID)

	if err != nil {
		http.Error(w, "Error fetching rooms", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description,
			&room.CreatorID, &room.Private,
			&room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning rooms", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	json.NewEncoder(w).Encode(rooms)
}

func JoinRoomHTTP(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	roomID := vars["id"]

	db := database.GetDB()
	var room models.Room
	err := db.QueryRow("SELECT id, private FROM rooms WHERE id = $1", roomID).Scan(&room.ID, &room.Private)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.Private {
		http.Error(w, "Cannot join private room", http.StatusForbidden)
		return
	}

	_, err = db.Exec(`
		INSERT INTO room_members (room_id, user_id, role, joined_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id, user_id) DO NOTHING`,
		room.ID, claims.UserID, "member", time.Now(), time.Now())

	if err != nil {
		http.Error(w, "Error joining room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
