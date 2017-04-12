package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ALiveDatabaseConnectionWorks(t *testing.T) {

	t.Skip()

	sf, err := NewFirebaseSessionFactory("https://ur-money-paul.firebaseio.com/", "fixtures/serviceAccountCredentials.ur-money-paul.json")
	require.Nil(t, err)

	session, err := sf.NewSession()
	require.Nil(t, err)

	var users map[string]interface{}
	require.Nil(t, session.Retrieve("users", Query{ID: Expression{Operator: EQ, Value: "-KhS35nLXdlYwYT62f9b"}}, &users))

	fmt.Println(users)

	require.Nil(t, session.Update("users", Query{ID: Expression{Operator: EQ, Value: "-KhS35nLXdlYwYT62f9b"}}, map[string]interface{}{"test": "Ok!"}))
}
