package main

import (
	"errors"

	r "github.com/dancannon/gorethink"
)

type Group struct {
	Id       string   `gorethink:"id" json:"id"`
	Settings Settings `gorethink:"settings" json:"settings"`
}

type Settings struct {
	PagerdutyKey string `gorethink:"pagerduty_key,omitempty" json:"pagerduty_key"`
	Slack        Slack  `gorethink:"slack,omitempty" json:"slack"`
}

func (g *Group) Validate() error {
	if g.Id == "" {
		return errors.New("missing id field")
	}
	return nil
}

func (g *Group) Save() error {
	resp, err := r.Table("groups").
		Insert(g, r.InsertOpts{Conflict: "replace"}).
		RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Inserted == 0 && resp.Replaced == 0 {
		return errors.New("Unable to insert/replace Group")
	}
	return nil
}
