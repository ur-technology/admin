package database

type Operator int

const (
	ID = "_id"
)

//go:generate stringer -type=Operator
const (
	EQ Operator = iota
)

type Expression struct {
	Operator Operator    `hcl:"operator"`
	Value    interface{} `hcl:"value"`
}

type Query map[string]Expression

type RetrievalParams struct {
	Page     int
	PageSize int64
}

type Session interface {
	Create(store string, object interface{}) (id string, err error)
	Retrieve(store string, query Query, object interface{}, params ...RetrievalParams) error
	Update(store string, query Query, values map[string]interface{}) error
	Delete(store string, query Query) error
	Close() error
}
