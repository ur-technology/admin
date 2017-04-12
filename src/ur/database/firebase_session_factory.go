package database

import (
	"io/ioutil"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type firebaseSessionFactory struct {
	ref *firego.Firebase
}

func NewFirebaseSessionFactory(endpoint, jsonFile string) (SessionFactory, error) {

	d, err := ioutil.ReadFile(jsonFile)

	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(
		d,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		return nil, err
	}

	fb := firego.New(endpoint, conf.Client(oauth2.NoContext))

	return &firebaseSessionFactory{
		ref: fb,
	}, nil
}

func (f *firebaseSessionFactory) NewSession() (Session, error) {
	return newFireBaseSession(f.ref), nil
}

func (f *firebaseSessionFactory) Close() error {
	return nil
}
