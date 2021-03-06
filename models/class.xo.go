// Package models contains the types for schema 'public'.
package models

// GENERATED BY XO. DO NOT EDIT.

import (
	"errors"

	"github.com/google/uuid"
)

// Class represents a row from 'public.classes'.
type Class struct {
	ID          uuid.UUID `json:"id"`           // id
	Name        string    `json:"name"`         // name
	CurrentUnit uuid.UUID `json:"current_unit"` // current_unit
	Active      bool      `json:"-"`       // active

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Class exists in the database.
func (c *Class) Exists() bool {
	return c._exists
}

// Deleted provides information if the Class has been deleted from the database.
func (c *Class) Deleted() bool {
	return c._deleted
}

// Insert inserts the Class to the database.
func (c *Class) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.classes (` +
		`id, name, current_unit, active` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`)`

	// run query
	XOLog(sqlstr, c.ID, c.Name, c.CurrentUnit, c.Active)
	err = db.QueryRow(sqlstr, c.ID, c.Name, c.CurrentUnit, c.Active).Scan(&c.ID)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Update updates the Class in the database.
func (c *Class) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if c._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.classes SET (` +
		`name, current_unit, active` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	// run query
	XOLog(sqlstr, c.Name, c.CurrentUnit, c.Active, c.ID)
	_, err = db.Exec(sqlstr, c.Name, c.CurrentUnit, c.Active, c.ID)
	return err
}

// Save saves the Class to the database.
func (c *Class) Save(db XODB) error {
	if c.Exists() {
		return c.Update(db)
	}

	return c.Insert(db)
}

// Upsert performs an upsert for Class.
//
// NOTE: PostgreSQL 9.5+ only
func (c *Class) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.classes (` +
		`id, name, current_unit, active` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, current_unit, active` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.current_unit, EXCLUDED.active` +
		`)`

	// run query
	XOLog(sqlstr, c.ID, c.Name, c.CurrentUnit, c.Active)
	_, err = db.Exec(sqlstr, c.ID, c.Name, c.CurrentUnit, c.Active)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Delete deletes the Class from the database.
func (c *Class) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return nil
	}

	// if deleted, bail
	if c._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.classes WHERE id = $1`

	// run query
	XOLog(sqlstr, c.ID)
	_, err = db.Exec(sqlstr, c.ID)
	if err != nil {
		return err
	}

	// set deleted
	c._deleted = true

	return nil
}

// ClassByID retrieves a row from 'public.classes' as a Class.
//
// Generated from index 'classes_pkey'.
func ClassByID(db XODB, id uuid.UUID) (*Class, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, current_unit, active ` +
		`FROM public.classes ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	c := Class{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&c.ID, &c.Name, &c.CurrentUnit, &c.Active)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
