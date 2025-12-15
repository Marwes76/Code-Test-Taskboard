# Code Test Taskboard

This project is built with a Golang backend (**"api"** service), Angular frontend (**"ui"** service), MySQL database (**"db"** service), and deployed using Docker.

To build and run the project, run the command `docker-compose -f docker-compose.dev.yml up --build`.

To deploy changes to the **api** service, run `docker-compose -f docker-compose.dev.yml up -d --build api` from another command line.

The **ui** service is automatically reloaded in response to changes.

- User interface is located at http://localhost:4200
- A simple web interface for DB is located at http://localhost:8081 (login with user `app`, password `app`, database `code_test_taskboard`)

In the UI, press the plus-sign to create your first list, and then create new tasks by pressing the plus-sign within the list.

---

### TODO

#### api service:
- [ ] Move "CREATE TABLE"-scripts from db and make proper migrations
- [ ] Make generic db-functions where you can just plug in db-table and data model

#### ui service:
- [ ] Fix bug where tasks can't be moved between lists
- [ ] Drag-and-drop
- [ ] Localization
- [ ] Make reusable component for edit/save/delete buttons
- [ ] Alerts for when HTTP-requests go through/fail
- [ ] Lock tasks instead if hiding them altogether when editing parent list
- [ ] Confirmation modals when deleting data

#### General:
- [ ] sortOrder needs to be normalized/adjusted when lists/tasks are moved/removed
