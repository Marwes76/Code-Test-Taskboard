package tasks

import (
	"api/app/communication"
	"api/app/database"
	"api/app/search"

	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strings"
)

type TaskRequest struct {
	ListUuid	string		`json:"listUuid"`
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	SortOrder	uint		`json:"sortOrder"`
}

func SearchTasksAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tasks Tasks
	searchString := mux.Vars(r)["searchString"]
	listUuid := r.URL.Query().Get("listUuid")
	orderBy := r.URL.Query().Get("orderBy")

	// Validate
	if listUuid != "" {
		if _, err = uuid.Parse(listUuid); err != nil {
			communication.ResponseBadRequest(w, err)
			return
		}
	}
	if orderBy != "" {
		if !search.IsValidOrderBy(orderBy) {
			communication.ResponseBadRequest(w, errors.New("Invalid orderBy"))
			return
		}
	}

	// Open DB
	if db, err = database.OpenDB(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Build search query
	var query = "SELECT * FROM tasks"
	var conditions []string
	var args []interface{}
	if searchString != "" {
		conditions = append(conditions, "(title LIKE ? OR title LIKE ? OR title LIKE ?)")
		args = append(args, searchString + "%")
		args = append(args, "%" + searchString + "%")
		args = append(args, "%" + searchString)
	}
	if listUuid != "" {
		conditions = append(conditions, "list_uuid = ?")
		args = append(args, listUuid)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if orderBy != "" {
		if search.OrderBy(orderBy) == search.OrderBySortOrder {
			query += " ORDER BY sort_order ASC"
		} else if search.OrderBy(orderBy) == search.OrderByAlphabetical {
			query += " ORDER BY title ASC"
		} else if search.OrderBy(orderBy) == search.OrderByCreatedAt {
			query += " ORDER BY created_at DESC"
		} else if search.OrderBy(orderBy) == search.OrderByUpdatedAt {
			query += " ORDER BY updated_at DESC"
		}
	} else {
		query += " ORDER BY sort_order ASC"
	}

	// Search for tasks
	if err = db.Select(&tasks, query, args...); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	if tasks == nil {
		tasks = Tasks{}
	}

	communication.ResponseOK(w, tasks)
}

func GetTaskAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var task Task
	reqUuid := mux.Vars(r)["uuid"]

	// Validate
	if _, err = uuid.Parse(reqUuid); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	// Open DB
	if db, err = database.OpenDB(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Get task from DB
	var query = "SELECT * FROM tasks WHERE uuid = ?"
	if err = db.Get(&task, query, reqUuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			communication.ResponseNotFound(w)
			return
		}

		communication.ResponseInternalServerError(w, err)
		return
	}

	communication.ResponseOK(w, task)
}

func CreateTaskAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	var task Task
	var listUuid uuid.UUID
	var req TaskRequest

	// Read task from request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	// Validate
	if listUuid, err = uuid.Parse(req.ListUuid); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	task = Task{
		Uuid:		uuid.New(),
		ListUuid:	listUuid,
		Title:		req.Title,
		Description:	req.Description,
		SortOrder:	req.SortOrder,
	}

	// Open DB
	if db, tx, err = database.OpenDBWithTX(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Save to DB
	var query = "INSERT INTO tasks (uuid, list_uuid, title, description, sort_order) VALUES (?, ?, ?, ?, ?)"
	args := []interface{}{task.Uuid, task.ListUuid, task.Title, task.Description, task.SortOrder}
	if _, err = tx.Exec(query, args...); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Get created task from DB
	query = "SELECT * FROM tasks WHERE uuid = ?"
	if err = db.Get(&task, query, task.Uuid); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	communication.ResponseCreated(w, task)
}

func UpdateTasksAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	var tasks Tasks
	var reqs map[string]TaskRequest

	// Read tasks from request
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	// Open DB
	if db, tx, err = database.OpenDBWithTX(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Update in DB
	// Since we're only dealing with very little data, we do one UPDATE-query at a time, but for the future we might want to
	// scale it up to a more flexible query that can handle multiple rows
	var result sql.Result
	var query string
	var args []interface{}
	if tx, err = db.Beginx(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	taskUuids := make([]string, len(reqs))
	for reqUuid, req := range reqs {
		// Validate
		if _, err = uuid.Parse(reqUuid); err != nil {
			tx.Rollback()
			communication.ResponseBadRequest(w, err)
			return
		}
		if _, err = uuid.Parse(req.ListUuid); err != nil {
			tx.Rollback()
			communication.ResponseBadRequest(w, err)
			return
		}

		taskUuids = append(taskUuids, reqUuid)
		query = "UPDATE tasks SET list_uuid = ?, title = ?, description = ?, sort_order = ? WHERE uuid = ?"
		args = []interface{}{req.ListUuid, req.Title, req.Description, req.SortOrder, reqUuid}
		if result, err = tx.Exec(query, args...); err != nil {
			tx.Rollback()
			communication.ResponseInternalServerError(w, err)
			return
		}
		if _, err = result.RowsAffected(); err != nil {
			tx.Rollback()
			communication.ResponseInternalServerError(w, err)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Get updated tasks from DB
	query = "SELECT * FROM tasks WHERE uuid IN (?) ORDER BY sort_order ASC"
	args = []interface{}{}
	if query, args, err = sqlx.In(query, taskUuids); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	if err = db.Select(&tasks, query, args...); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Map tasks by UUID
	tasksByUuid := make(TasksMap)
	for _, task := range tasks {
		tasksByUuid[task.Uuid.String()] = task
	}

	communication.ResponseOK(w, tasksByUuid)
}

func DeleteTaskAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	reqUuid := mux.Vars(r)["uuid"]

	// Validate
	if _, err = uuid.Parse(reqUuid); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	// Open DB
	if db, tx, err = database.OpenDBWithTX(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Remove from DB
	var result sql.Result
	var rowsAffected int64
	var query = "DELETE FROM tasks WHERE uuid = ?"
	if result, err = tx.Exec(query, reqUuid); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}
	if rowsAffected == 0 {
		tx.Rollback()
		communication.ResponseNotFound(w)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}

	communication.ResponseNoContent(w)
}