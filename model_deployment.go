package main

import (
	"errors"
	"fmt"

	r "github.com/dancannon/gorethink"
)

type Deployment struct {
	Id       string   `gorethink:"id" json:"id"`
	GroupId  string   `gorethink:"group_id" json:"group_id"`
	Type     string   `gorethink:"type" json:"type"`
	Name     string   `gorethink:"name" json:"name"`
	Settings Settings `gorethink:"settings" json:"settings"`
	Checks   []Check  `gorethink:"checks" json:"checks"`
}

func (d *Deployment) Validate() error {
	if d.Id == "" {
		return errors.New("missing requiered field, id")
	} else if d.GroupId == "" {
		return errors.New("missing required field, group_id")
	} else if d.Type == "" {
		return errors.New("missing required field, type")
	} else if d.Name == "" {
		return errors.New("missing required field, name")
	}
	return nil
}

func (d *Deployment) Save() error {
	for i, check := range d.Checks {
		d.Checks[i].Id = fmt.Sprintf("%s-%s", d.Type, check.Name)
	}
	resp, err := r.Table("deployments").
		Insert(d, r.InsertOpts{Conflict: "replace"}).
		RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Inserted == 0 && resp.Replaced == 0 {
		return errors.New("Unable to insert/replace Deployment")
	}
	return nil
}

func (d *Deployment) Delete() error {
	resp, err := r.Table("deployments").Get(d.Id).Delete().RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Deleted != 1 {
		return errors.New("Unable to delete Deployment")
	}
	return nil
}

func LookupDeploymentById(deploymentId string) (deployment Deployment, err error) {
	cur, err := r.Table("deployments").Get(deploymentId).Run(session)
	if err != nil {
		return deployment, err
	}
	defer cur.Close()
	if err = cur.One(&deployment); err != nil {
		return deployment, err
	}
	return deployment, nil
}

// TODO: this is all kinds of gross
func (d *Deployment) CurrentChecks() ([]Check, error) {
	checks := make([]Check, 0)
	checksCursor, err := r.Table("checks").GetAllByIndex("type", d.Type).Run(session)
	if err != nil {
		return nil, err
	}
	defer checksCursor.Close()
	var defaultChecks []Check
	checksCursor.All(&defaultChecks)
	for _, dc := range defaultChecks {
		checkFound := false
		for _, c := range d.Checks {
			if c.Name == dc.Name {
				checks = append(checks, c)
				checkFound = true
				break
			}
		}
		if !checkFound {
			checks = append(checks, dc)
		}
	}
	if checksCursor.Err() != nil {
		return nil, checksCursor.Err()
	}
	return checks, nil
}

func (d *Deployment) CheckByName(name string) (check Check, err error) {
	cur, err := r.Branch(
		r.Table("deployments").Get(d.Id).Field("checks").Filter(map[string]interface{}{"name": name}).IsEmpty(),
		r.Table("checks").GetAllByIndex("type_name", []string{d.Type, name}),
		r.Table("deployments").Get(d.Id).Field("checks").Filter(map[string]interface{}{"name": name}),
	).Run(session)
	defer cur.Close()
	err = cur.One(&check)
	if err != nil {
		return check, err
	}
	return check, nil
}

func (d *Deployment) MergedAlertSettings() (settings Settings, err error) {
	cur, err := r.Branch(
		r.Table("groups").Get(d.GroupId),
		r.Table("groups").EqJoin("id", r.Table("deployments"), r.EqJoinOpts{Index: "group_id"}).Zip().Filter(map[string]interface{}{"id": d.Id}).Field("settings"),
		r.Table("deployments").Get(d.Id).Field("settings"),
	).Run(session)
	defer cur.Close()
	err = cur.One(&settings)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
