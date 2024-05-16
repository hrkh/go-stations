package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t := &model.TODO{ID: id}
	if err = s.db.QueryRowContext(ctx, confirm, id).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}

	return t, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read        = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID  = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
		defaultSize = 5
	)

	if size == 0 {
		size = defaultSize
	}

	var rows *sql.Rows
	if prevID == 0 {
		stmt, err := s.db.PrepareContext(ctx, read)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.QueryContext(ctx, size)
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err := s.db.PrepareContext(ctx, readWithID)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	tt := []*model.TODO{}
	for rows.Next() {
		t := &model.TODO{}
		if err := rows.Scan(&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tt = append(tt, t)
	}

	return tt, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}
	nRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if nRows == 0 {
		return nil, &model.ErrNotFound{}
	}

	t := &model.TODO{ID: id}
	if err = s.db.QueryRowContext(ctx, confirm, id).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}

	argIds := []interface{}{}
	for _, id := range ids {
		argIds = append(argIds, id)
	}

	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(argIds)-1)))
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, argIds...)
	if err != nil {
		return err
	}
	nRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if nRows == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
