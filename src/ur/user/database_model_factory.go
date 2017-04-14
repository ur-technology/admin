package user

import (
	"ur/database"
	"ur/errors"
)

const (
	databaseStore = "users"
)

type databaseModelFactory struct {
	sessionFactory database.SessionFactory
}

func NewDatabaseModelFactory(sf database.SessionFactory) ModelFactory {
	return &databaseModelFactory{
		sessionFactory: sf,
	}
}

func (f *databaseModelFactory) FetchOne(id string) (*Model, error) {

	session, err := f.sessionFactory.NewSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	models := make(map[string]*Model)

	if err := session.Retrieve(databaseStore, database.Query{database.ID: database.Expression{database.EQ, id}}, &models); err != nil {
		return nil, err
	}

	m, found := models[id]

	if !found {
		return nil, errors.NewUser("User not found")
	}

	m.UID = id

	return m, nil
}
