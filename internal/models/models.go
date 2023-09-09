package models

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	// ErrNoRecord no record found in database error
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials invalid username/password error
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail duplicate email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrInactiveAccount inactive account error
	ErrInactiveAccount = errors.New("models: Inactive Account")
)

// User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	UserActive  int
	AccessLevel int
	Email       string
	Password    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Preferences map[string]string
}

// Preference model
type Preference struct {
	ID         int
	Name       string
	Preference []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Host is the model for hosts
type Host struct {
	ID            int
	Hostname      string
	CanonicalName string
	URL           string
	IP            string
	IPv6          string
	Location      string
	OS            string
	Active        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	HostServices  []HostService
}

// Service in the model for services
type Service struct {
	ID          int
	ServiceName string
	Active      int
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// HostService is the model for host services
type HostService struct {
	ID             int
	HostID         int
	ServiceID      int
	Active         int
	ScheduleNumber int
	ScheduleUnit   string
	Status         string
	LastCheck      time.Time
	LastMessage    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Service        Service
	Hostname       string
}

// Schedule model for schedules
type Schedule struct {
	ID            int
	EntryID       cron.EntryID
	Entry         cron.Entry
	Host          string
	Service       string
	LastRunFromHS time.Time
	HostServiceID int
	ScheduleText  string
}

// Events model for events
type Events struct {
	ID            int
	EventType     string
	HostServiceID int
	HostID        int
	ServiceName   string
	Hostname      string
	Message       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
