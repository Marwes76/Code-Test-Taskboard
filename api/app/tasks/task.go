package tasks

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Uuid		uuid.UUID	`json:"uuid" db:"uuid"`
	ListUuid	uuid.UUID	`json:"listUuid" db:"list_uuid"`
	Title		string		`json:"title" db:"title"`
	Description	string		`json:"description" db:"description"`
	SortOrder	uint		`json:"sortOrder" db:"sort_order"`
	CreatedAt	time.Time	`json:"createdAt" db:"created_at"`
	UpdatedAt	time.Time	`json:"updatedAt" db:"updated_at"`
}

// Alias for slice of tasks
type Tasks []Task

// Alias for map of tasks indexed by UUID
type TasksMap map[string]Task