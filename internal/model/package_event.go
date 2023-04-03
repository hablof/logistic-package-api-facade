package model

import (
	"encoding/json"
	"time"
)

type EventType uint8

func (e EventType) MarshalJSON() ([]byte, error) {
	var s string
	switch e {
	case Created:
		s = "Created"
	case Updated:
		s = "Updated"
	case Removed:
		s = "Removed"
	}

	return json.Marshal(s)
}

type EventStatus uint8

type Payload []byte

func (p Payload) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return json.Marshal([]byte(p))
	}
	return []byte(p), nil
}

const (
	_ EventType = iota
	Created
	Updated
	Removed
)

const (
	_ EventStatus = iota
	Locked
	Unlocked
)

type PackageEvent struct {
	ID        uint64      // `db:"package_event_id"`
	PackageID uint64      // `db:"package_id"`
	Type      EventType   // `db:"event_type"`
	Status    EventStatus // `db:"event_status"` // Это даже нигде не используется, кроме как внутри базы
	Created   time.Time   // `db:"created_at"`
	Payload   Payload     // `db:"payload"`
}
