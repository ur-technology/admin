package database

import (
	"fmt"
	"strings"

	"github.com/zabawaba99/firego"
)

func newFireBaseSession(root *firego.Firebase) Session {
	return &firebaseSession{
		root: root,
	}
}

type firebaseSession struct {
	root *firego.Firebase
}

///////////////////////////////////////////////////////////////////////////////
// Interface functions
///////////////////////////////////////////////////////////////////////////////

func (f *firebaseSession) Create(store string, object interface{}) (string, error) {

	ref, err := f.root.Ref(store)

	if err != nil {
		return "", err
	}

	finalRef, err := ref.Push(object)

	if err != nil {
		return "", err
	}

	return f.idFromRef(finalRef), nil
}

func (f *firebaseSession) Retrieve(store string, query Query, object interface{}, params ...RetrievalParams) error {

	ref, err := f.query(store, query, params...)

	if err != nil {
		return err
	}

	return ref.Value(object)
}

func (f *firebaseSession) Update(store string, query Query, values map[string]interface{}) error {

	idExpression, ok := query[ID]

	if !ok {
		return fmt.Errorf("Firebase only supports updates on one record at a time")
	}

	id, ok := idExpression.Value.(string)

	if !ok {
		return fmt.Errorf("Firebase only supports string IDs")
	}

	ref, err := f.root.Ref(store + "/" + id)

	if err != nil {
		return err
	}

	if err := ref.Update(values); err != nil {
		return err
	}

	return nil
}

func (f *firebaseSession) Delete(store string, query Query) error {

	ref, err := f.query(store, query)

	if err != nil {
		return err
	}

	return ref.Remove()
}

func (f *firebaseSession) Close() error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Utility functions
///////////////////////////////////////////////////////////////////////////////

func (f *firebaseSession) query(store string, query Query, params ...RetrievalParams) (*firego.Firebase, error) {

	ref, err := f.root.Ref(store)

	if err != nil {
		return nil, err
	}

	for k, v := range query {

		// A special case for ID (or _id)
		if k == ID {

			sid, ok := v.Value.(string)

			if !ok {
				return nil, fmt.Errorf("ID values must be strings")
			}

			ref = ref.OrderBy("$key").EqualTo(sid)

			continue
		}

		ref = ref.OrderBy(k)

		switch v.Operator {
		case EQ:
			ref = ref.EqualToValue(v.Value)

		default:
			return nil, fmt.Errorf("Operator %d not supported", v.Operator)
		}
	}

	// Firebase doesn't support pagination!
	if len(params) > 0 {
		ref = ref.LimitToFirst(params[0].PageSize)
	}

	return ref, nil
}

func (f *firebaseSession) idFromRef(ref *firego.Firebase) string {

	raw := strings.TrimSuffix(ref.String(), "/.json")
	parts := strings.Split(raw, "/")

	return parts[len(parts)-1]
}
