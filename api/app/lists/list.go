package lists

import (
	"api/app/tasks"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type List struct {
	Uuid		uuid.UUID	`json:"uuid" db:"uuid"`
	Title		string		`json:"title" db:"title"`
	Description	string		`json:"description" db:"description"`
	SortOrder	uint		`json:"sortOrder" db:"sort_order"`
	CreatedAt	time.Time	`json:"createdAt" db:"created_at"`
	UpdatedAt	time.Time	`json:"updatedAt" db:"updated_at"`

	// Relations
	Tasks		[]tasks.Task	`json:"tasks"`
}

// Alias for slice of lists
type Lists []List

// Load relation of tasks to a list
func (list *List) LoadTasks(db *sqlx.DB) (err error) {
	list.Tasks = []tasks.Task{}

	// Get related tasks from DB
	var query = "SELECT * FROM tasks WHERE list_uuid = ? ORDER BY sort_order ASC"
	var args []interface{}
	if query, args, err = sqlx.In(query, list.Uuid); err != nil {
		return err
	}
	var tasks []tasks.Task
	if err = db.Select(&tasks, query, args...); err != nil {
		return err
	}

	for _, task := range tasks {
		list.Tasks = append(list.Tasks, task)
	}

	return nil
}

// Load relations of tasks to a slice of lists
func (lists *Lists) LoadTasks(db *sqlx.DB) (err error) {
	if len(*lists) == 0 {
		return nil
	}

	listMap := make(map[string]*List)
	listUids := make([]string, len(*lists))
	for i := range *lists {
		list := &(*lists)[i]
		list.Tasks = []tasks.Task{}

		listMap[list.Uuid.String()] = list
		listUids[i] = list.Uuid.String()
    	}

	// Get related tasks from DB
	var query = "SELECT * FROM tasks WHERE list_uuid IN (?) ORDER BY sort_order ASC"
	var args []interface{}
	if query, args, err = sqlx.In(query, listUids); err != nil {
		return err
	}
	var tasks []tasks.Task
	if err = db.Select(&tasks, query, args...); err != nil {
		return err
	}

	// Map tasks to lists
	for _, task := range tasks {
		if list, ok := listMap[task.ListUuid.String()]; ok {
			list.Tasks = append(list.Tasks, task)
		}
	}

	return nil
}