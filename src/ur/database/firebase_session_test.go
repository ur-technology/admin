package database

import (
	"fmt"
	"testing"
	"ur/database/internal/firetest"

	"github.com/stretchr/testify/require"
	"github.com/zabawaba99/firego"
)

func newMockFirebaseServer(t *testing.T) (fb *firego.Firebase, ft *firetest.Firetest, fn func()) {

	ft = firetest.New()
	ft.Start()

	fb = firego.New(ft.URL, nil)

	fn = func() {
		ft.Close()
	}

	return
}

func newMockFirebaseSession(t *testing.T) (*firebaseSession, *firetest.Firetest, func()) {

	fb, ft, closeFn := newMockFirebaseServer(t)

	return newFireBaseSession(fb).(*firebaseSession), ft, closeFn
}

type mockModelType struct {
	A1 float64 `json:"A,omitempty"`
	B2 string  `json:"2,omitempty"`
}

func Test_AFirebaseSessionCanExtractAnIdFromARef(t *testing.T) {
	fb, _, closeFn := newMockFirebaseSession(t)
	defer closeFn()

	for cycle, test := range []struct {
		description string
		ref         string
		id          string
	}{
		{
			description: "Expected result",
			ref:         "test/~MTQ5MTk4ODM5NTIwMzc3ODg5MQ==",
			id:          "~MTQ5MTk4ODM5NTIwMzc3ODg5MQ==",
		},
	} {
		t.Logf("Cycle %d: %s", cycle, test.description)

		ref, err := fb.root.Ref(test.ref)
		require.Nil(t, err)

		require.Equal(t, test.id, fb.idFromRef(ref))
	}
}

func Test_AFirebaseSessionCanCreateRecords(t *testing.T) {

	t.Parallel()

	fb, ft, closeFn := newMockFirebaseSession(t)
	defer closeFn()

	record := mockModelType{
		A1: 1.0,
		B2: "B",
	}

	id, err := fb.Create("test", record)

	require.Nil(t, err)
	require.NotEmpty(t, id)

	require.Equal(t, map[string]interface{}{"2": "B", "A": 1.0}, ft.Get("test/"+id))
}

func Test_AFirebaseSessionCanRetrieveRecords(t *testing.T) {

	t.Parallel()

	fb, _, closeFn := newMockFirebaseSession(t)
	defer closeFn()

	record := mockModelType{
		A1: 1.0,
		B2: "B",
	}

	id1, err := fb.Create("test", record)
	require.Nil(t, err)
	require.NotEmpty(t, id1)

	for cycle, test := range []struct {
		description string

		// Input
		store  string
		query  Query
		params []RetrievalParams

		// Output
		results map[string]mockModelType
		outcome error
	}{
		{
			description: "Success: single result",
			store:       "test",
			query:       Query{ID: Expression{EQ, id1}},
			results: map[string]mockModelType{
				id1: {
					A1: 1.0,
					B2: "B",
				},
			},
		},
		{
			description: "Failure: invalid ID",
			store:       "test",
			query:       Query{ID: Expression{EQ, 1}},
			outcome:     fmt.Errorf("ID values must be strings"),
		},
	} {
		t.Logf("Cycle %d: %s", cycle, test.description)

		results := map[string]mockModelType{}
		outcome := fb.Retrieve(test.store, test.query, &results, test.params...)

		if test.outcome == nil {
			require.Nil(t, outcome)
			require.Equal(t, test.results, results)
		} else {
			require.NotNil(t, outcome)
			require.Equal(t, test.outcome.Error(), outcome.Error())
		}
	}
}

func Test_AFirebaseSessionCanUpdateRecords(t *testing.T) {

	t.Parallel()

	fb, _, closeFn := newMockFirebaseSession(t)
	defer closeFn()

	record0 := mockModelType{
		A1: 1.0,
		B2: "B",
	}

	record1 := mockModelType{
		A1: 2.0,
	}

	record2 := mockModelType{
		A1: 2.0,
		B2: "B",
	}

	id, err := fb.Create("test", record0)
	require.Nil(t, err)
	require.NotEmpty(t, id)

	require.Nil(t, fb.Update("test", Query{ID: Expression{EQ, id}}, map[string]interface{}{"A": record1.A1}))

	var final map[string]mockModelType
	require.Nil(t, fb.Retrieve("test", Query{ID: Expression{EQ, id}}, &final))

	require.Equal(t, record2, final[id])
}

func Test_AFirebaseSessionCanDeleteRecords(t *testing.T) {

	t.Parallel()

	fb, _, closeFn := newMockFirebaseSession(t)
	defer closeFn()

	record0 := mockModelType{
		A1: 1.0,
		B2: "B",
	}

	id, err := fb.Create("test", record0)
	require.Nil(t, err)
	require.NotEmpty(t, id)

	require.Nil(t, fb.Delete("test", Query{ID: Expression{EQ, id}}))

	var final map[string]mockModelType
	require.Nil(t, fb.Retrieve("test", Query{ID: Expression{EQ, id}}, &final))

	require.Len(t, final, 0)
}
