package filestore

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// ErrNotFound is returned when the ID can't be found.
var ErrNotFound = errors.New("not found")

type db struct {
	Contents map[string]string
}

func open() (*db, error) {
	contents, err := ioutil.ReadFile("db")
	if os.IsNotExist(err) {
		return &db{}, nil
	}

	var db db
	if err := json.Unmarshal(contents, &db.Contents); err != nil {
		return nil, err
	}

	return &db, nil
}

func (d *db) save() error {
	out, err := json.Marshal(d.Contents)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("db", out, 0666)
}

func (d *db) writeDocument(id, contents string) {
	if d.Contents == nil {
		d.Contents = make(map[string]string)
	}
	d.Contents[id] = contents
}

// Read out the specified type with the given ID to the "db".
// Returns ErrNotFound if there's no document with the specified ID.
func Read(id string, out interface{}) error {
	d, err := open()
	if err != nil {
		return err
	}

	doc := d.Contents[id]
	if doc == "" {
		return ErrNotFound
	}

	return json.Unmarshal([]byte(doc), &out)
}

// Write out the specified type with the given ID to the "db".
func Write(id string, doc interface{}) error {
	d, err := open()
	if err != nil {
		return err
	}

	contents, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	d.writeDocument(id, string(contents))
	return d.save()
}

// Delete deletes the document with the specified id.
func Delete(id string) error {
	d, err := open()
	if err != nil {
		return err
	}

	delete(d.Contents, id)
	return d.save()
}
