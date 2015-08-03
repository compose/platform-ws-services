package main

import (
	"errors"
	"fmt"

	r "github.com/dancannon/gorethink"
)

type Check struct {
	Id          string `gorethink:"id" json:"id"`
	Type        string `gorethink:"type,omitempty" json:"type,omitempty"`
	Name        string `gorethink:"name" json:"name"`
	Level       int    `gorethink:"level" json:"level"`
	Title       string `gorethink:"title" json:"title"`
	Description string `gorethink:"description" json:"description"`
}

func (c *Check) Validate() error {
	if c.Type == "" {
		return errors.New("missing required field, type")
	} else if c.Name == "" {
		return errors.New("missing required field, name")
	} else if c.Level == 0 {
		return errors.New("invalid field, level, must be greater than 0")
	} else if c.Level > 3 {
		return errors.New("invalid field, level, must be less than 3")
	}
	return nil
}

func (c *Check) Save() error {
	c.Id = fmt.Sprintf("%s-%s", c.Type, c.Name)
	resp, err := r.Table("checks").
		Insert(c, r.InsertOpts{Conflict: "replace"}).
		RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Inserted == 0 && resp.Replaced == 0 {
		return errors.New("Unable to insert/replace Check")
	}
	return nil
}

func (c *Check) Delete() error {
	resp, err := r.Table("checks").
		Get(c.Id).
		Delete().
		RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return errors.New("Unable to delete Check")
	}
	return nil
}

func GetChecksByType(check_type string) ([]Check, error) {
	cur, err := r.Table("checks").GetAllByIndex("type", check_type).Run(session)
	if err != nil {
		return nil, err
	}
	defer cur.Close()
	var checks []Check
	err = cur.All(&checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func GetCheckByTypeName(check_type, name string) (check Check, err error) {
	cur, err := r.Table("checks").GetAllByIndex("type_name", []string{check_type, name}).Run(session)
	if err != nil {
		return check, err
	}
	defer cur.Close()
	err = cur.One(&check)
	if err != nil {
		return check, err
	}
	return check, nil
}
