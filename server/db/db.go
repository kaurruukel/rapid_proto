package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Events struct {
	db *sql.DB
}

type Event struct {
	Id           string `json:"id"`
	EventType    string `json:"type"`
	Acknowledged bool   `json:"acknowledged"`
	Date         string `json:"date"`
}

type UpdateEvent struct {
	Acknowledged bool `json:"acknowledged"`
}

func CreateSchema() (*Events, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS events (id TEXT PRIMARY KEY, event_type TEXT, acknowledged BOOLEAN, date STRING)")
	if err != nil {
		return nil, err
	}
	return &Events{db: db}, nil
}

func (e *Events) Close() error {
	return e.db.Close()
}

func (e *Events) AddEvent(event Event) (string, error) {
	newID := uuid.New().String()
	_, err := e.db.Exec("INSERT INTO events (id, event_type, acknowledged, date) VALUES (?, ?, ?, ?)", newID, event.EventType, event.Acknowledged, time.Now())
	return newID, err
}

func (e *Events) DeleteEvents() error {
	_, err := e.db.Exec("DELETE FROM events");
	return err;
}

func (e *Events) GetEvents() ([]Event, error) {
	rows, err := e.db.Query("SELECT id, event_type, acknowledged, date FROM events ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.EventType, &event.Acknowledged, &event.Date); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Events) GetEventById(id string) (Event, error) {
	var event Event
	err := e.db.QueryRow("SELECT id, event_type, acknowledged, date FROM events WHERE id = ?", id).Scan(&event.Id, &event.EventType, &event.Acknowledged, &event.Date)
	return event, err
}

func (e *Events) UpdateEvent(event UpdateEvent, id string) error {
	_, err := e.db.Exec("UPDATE events SET  acknowledged = ? WHERE id = ?", event.Acknowledged, id)
	return err
}
