// Service url takes URLs, generates random short IDs and stores them.
package todo

import (
	"context"

	"encore.dev/storage/sqldb"
)

type ToDo struct {
	Id     string `json:"id"`      // id of todo
	Title  string `json:"title"`   // desc of todo
	Deadline  string `json:"deadline"`
	IsDone bool   `json:"is_done"` // status of todo
}

type CreateToDoParams struct {
	Id    string
	Title string
	Deadline  string
}

type ToDoListResponse struct {
	ToDo []ToDo
}

type UpdateToDoParams struct {
	Title string
}

type UpdateToDoStatusParam struct {
	IsDone bool
}

var db = sqldb.NewDatabase("todo", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

// encore:api public method=GET path=/todo
func GetToDoList(ctx context.Context) (*ToDoListResponse, error) {
	rows, err := db.Query(ctx, "SELECT * FROM todo")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// // An album slice to hold data from returned rows.
	var todos []ToDo

	// // // Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var item ToDo
		if err := rows.Scan(&item.Id, &item.Title, &item.IsDone); err != nil {
			return nil, err
		}
		todos = append(todos, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &ToDoListResponse{ToDo: todos}, nil
}

//encore:api public method=POST path=/todo
func CreateToDo(ctx context.Context, param *CreateToDoParams) (*ToDo, error) {
	if err := insert(ctx, param.Id, param.Title, param.Deadline); err != nil {
		return nil, err
	}
	return &ToDo{Id: param.Id, Title: param.Title, Deadline: param.Deadline}, nil
}

// insert inserts a URL into the database.
func insert(ctx context.Context, id, title string, deadline string) error {
	_, err := db.Exec(ctx, `
        INSERT INTO todo (id, title)
        VALUES ($1, $2)
    `, id, title)
	return err
}

//encore:api public method=DELETE path=/todo/:id
func DeleteToDo(ctx context.Context, id string) error {
	_, err := db.Exec(ctx, `
        delete from todo where id=$1
    `, id)
	return err
}

//encore:api public method=PUT path=/todo/:id
func UpdateToDo(ctx context.Context, id string, param *UpdateToDoParams) error {
	_, err := db.Exec(ctx, `
        UPDATE todo
        SET title = $2
		WHERE id=$1
    `, id, param.Title)
	return err
}

//encore:api public method=POST path=/todo/status/:id
func UpdateToDoStatus(ctx context.Context, id string, param *UpdateToDoStatusParam) error {
	_, err := db.Exec(ctx, `
        UPDATE todo
        SET is_done = $2
		WHERE id=$1
    `, id, param.IsDone)
	return err
}
