package lists

import (
	"api/app/communication"
	"api/app/database"

	"github.com/jmoiron/sqlx"
	"net/http"
)

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
	var query = "SELECT * FROM lists"
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
}

func CreateListAPI(w http.ResponseWriter, r *http.Request) {
}

func UpdateListAPI(w http.ResponseWriter, r *http.Request) {
}

func DeleteListAPI(w http.ResponseWriter, r *http.Request) {
}