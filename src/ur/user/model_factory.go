package user

type ModelFactory interface {
	FetchOne(id string) (*Model, error)
}
