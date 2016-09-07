// github.com/hduplooy/gonorm
// Author: Hannes du Plooy
// Revision Date: 7 Sep 2016
// This is not a object relational mapping in the normal sense.
// The functions take a sql query and a struct value and then populate or do with the data what it must.
// So these are more help functions so that one doesn't have to manually fill an array with individual sql row scans.
package gonorm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

// Norm is a structure holding the necessary data that we are working with
type Norm struct {
	// The link to the opened database connection
	DB *sql.DB
	// The driver string with which the connecion was opened
	Driver string
	// The connection string with which the database was opened
	Connect string
}

// NewNorm just connects to a database with the necessary details and return a *Norm and an error if any
// The driver and connect strings are kept in the returned structure
func NewNorm(driver, connect string) (*Norm, error) {
	db, err := sql.Open(driver, connect)
	if err != nil {
		return nil, err
	}
	return &Norm{DB: db, Driver: driver, Connect: connect}, nil
}

// GetRows will execute the sql query and then based on the type of val generate a slice with entries populated from each row returned
// sql is the sql query (select query normally)
// val is the val (an empty structure) which is used as template to generate the slice
func (ent *Norm) GetRows(sql string, val interface{}) (interface{}, error) {
	strflds := make(map[string]int) // Here a mapping is kept from the field name to the position in the struct
	valr := reflect.TypeOf(val)
	slice := reflect.MakeSlice(reflect.SliceOf(valr), 0, 100) // Our slice where we will put our results
	for i := 0; i < valr.NumField(); i++ {
		fld := valr.Field(i)
		// Get the name of the field
		nm := strings.ToLower(fld.Name)
		// If there is a fldnm tag then rather use that
		if fldnm := fld.Tag.Get("fldnm"); len(strings.Trim(fldnm, " \t")) > 0 {
			nm = strings.Trim(fldnm, " \t")
		}
		// Save name to index mapping
		strflds[nm] = i
	}
	// Get rows
	rows, err := ent.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the column names as returned by the query
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// colmap is a query column to field index mapping
	colmap := make([]int, len(cols))
	for i, val := range cols {
		// Find the column in the name to index mapping
		if ind, ok := strflds[val]; ok {
			colmap[i] = ind // Save the query column to field index
		} else {
			// If not found then no match and it cannot be handled
			return nil, errors.New(fmt.Sprintf("No match for column %d: %s", i, val))
		}
	}

	// Process the individual rows returned
	for rows.Next() {
		// Create a new empty copy of the provided structure
		newval := reflect.Indirect(reflect.New(valr))
		// Save addresses of field of structure
		adds := make([]interface{}, len(colmap))
		for i, ind := range colmap {
			adds[i] = newval.Field(ind).Addr().Interface()
		}
		// Scan the row and save to the empty copy
		rows.Scan(adds...)
		// Append it to our slice
		slice = reflect.Append(slice, newval)
	}
	// Return the slice
	return slice.Interface(), nil
}

// GetRow similar to GetRows except it is for only the first row returned
// sql is the sql query (select query normally)
// val is the val (an empty structure) which is used as template to generate the slice
func (ent *Norm) GetRow(sql string, val interface{}) (interface{}, error) {
	strflds := make(map[string]int) // Here a mapping is kept from the field name to the position in the struct
	valr := reflect.TypeOf(val)
	for i := 0; i < valr.NumField(); i++ {
		fld := valr.Field(i)
		// Get the name of the field
		nm := strings.ToLower(fld.Name)
		// If there is a fldnm tag then rather use that
		if fldnm := fld.Tag.Get("fldnm"); len(strings.Trim(fldnm, " \t")) > 0 {
			nm = strings.Trim(fldnm, " \t")
		}
		// Save name to index mapping
		strflds[nm] = i
	}
	// Get rows
	rows, err := ent.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the column names as returned by the query
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// colmap is a query column to field index mapping
	colmap := make([]int, len(cols))
	for i, val := range cols {
		// Find the column in the name to index mapping
		if ind, ok := strflds[val]; ok {
			colmap[i] = ind
		} else {
			// If not found then no match and it cannot be handled
			return nil, errors.New(fmt.Sprintf("No match for column %d: %s", i, val))
		}
	}

	// Get first row
	if rows.Next() {
		// Create a new empty copy of the provided structure
		newval := reflect.Indirect(reflect.New(valr))
		// Save addresses of field of structure
		adds := make([]interface{}, len(colmap))
		for i, ind := range colmap {
			adds[i] = newval.Field(ind).Addr().Interface()
		}
		// Scan the row and save to the empty copy
		rows.Scan(adds...)
		return newval, nil
	}
	return nil, nil
}

// GetRowsJson get the rows from the sql query map them to the struct type of val and convert the slice to json
// sql is the sql query (select query normally)
// val is the val (an empty structure) which is used as template to generate the slice
func (ent *Norm) GetRowsJson(sql string, val interface{}) (string, error) {
	// Get the rows
	vals, err := ent.GetRows(sql, val)
	// If an error was experienced, return and empty json array with the error
	if err != nil {
		return "[]", err
	}
	// Generate the json representation and return
	buf, err := json.Marshal(vals)
	return string(buf), nil
}

// GetRowJson get the first row from the sql query map it to the struct type of val and convert to json
// sql is the sql query (select query normally)
// val is the val (an empty structure) which is used as template
func (ent *Norm) GetRowJson(sql string, val interface{}) (string, error) {
	// Get the row
	vals, err := ent.GetRow(sql, val)
	// If an error was experienced, return and empty string with the error
	if err != nil {
		return "", err
	}
	// Generate the json representation and return
	buf, err := json.Marshal(vals)
	return string(buf), nil
}
