package lists

import (
	"api/app/communication"
	"api/app/database"
	"api/app/tasks"

	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type ListRequest struct {
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	SortOrder	uint	`json:"sortOrder"`
}

func GetAllListsAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var lists Lists

	// Open DB
	if db, err = database.OpenDB(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	defer db.Close()

	// Get all lists from DB
	var query = "SELECT * FROM lists ORDER BY sort_order ASC"
	if err = db.Select(&lists, query); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Load relations of tasks
	if err = lists.LoadTasks(db); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	communication.ResponseOK(w, lists)
}

func GetListAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var list List
	reqUuid := mux.Vars(r)["uuid"]

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

	// Get list from DB
	var query = "SELECT * FROM lists WHERE uuid = ?"
	if err = db.Get(&list, query, reqUuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			communication.ResponseNotFound(w)
			return
		}

		communication.ResponseInternalServerError(w, err)
		return
	}

	// Load relation of tasks
	if err = list.LoadTasks(db); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	communication.ResponseOK(w, list)
}

func CreateListAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	var list List
	var req ListRequest

	// Read list from request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		communication.ResponseBadRequest(w, err)
		return
	}

	list = List{
		Uuid:		uuid.New(),
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
	var query = "INSERT INTO lists (uuid, title, description, sort_order) VALUES (?, ?, ?, ?)"
	args := []interface{}{list.Uuid, list.Title, list.Description, list.SortOrder}
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

	// Get created list from DB
	query = "SELECT * FROM lists WHERE uuid = ?"
	if err = db.Get(&list, query, list.Uuid); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	list.Tasks = []tasks.Task{} // Relation will be empty to begin with

	communication.ResponseCreated(w, list)
}

func UpdateListsAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	var lists Lists
	var reqs map[string]ListRequest

	// Read lists from request
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
	var rowsAffected int64
	var query string
	var args []interface{}
	if tx, err = db.Beginx(); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	listUuids := make([]string, len(reqs))
	for reqUuid, req := range reqs {
		if _, err = uuid.Parse(reqUuid); err != nil {
			tx.Rollback()
			communication.ResponseBadRequest(w, err)
			return
		}

		listUuids = append(listUuids, reqUuid)
		query = "UPDATE lists SET title = ?, description = ?, sort_order = ? WHERE uuid = ?"
		args = []interface{}{req.Title, req.Description, req.SortOrder, reqUuid}
		if result, err = tx.Exec(query, args...); err != nil {
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
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Get updated lists from DB
	query = "SELECT * FROM lists WHERE uuid IN (?) ORDER BY sort_order ASC"
	args = []interface{}{}
	if query, args, err = sqlx.In(query, listUuids); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}
	if err = db.Select(&lists, query, args...); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Load relation of tasks
	if err = lists.LoadTasks(db); err != nil {
		communication.ResponseInternalServerError(w, err)
		return
	}

	// Map lists by UUID
	listsByUuid := make(ListsMap)
	for _, list := range lists {
		listsByUuid[list.Uuid.String()] = list
	}

	communication.ResponseOK(w, listsByUuid)
}

func DeleteListAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sqlx.DB
	var tx *sqlx.Tx
	reqUuid := mux.Vars(r)["uuid"]

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
	var query = "DELETE FROM lists WHERE uuid = ?"
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

	// Remove uuid from existing tasks related to the list
	query = "DELETE FROM tasks WHERE list_uuid = ?"
	if _, err = tx.Exec(query, reqUuid); err != nil {
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

	communication.ResponseNoContent(w)
}