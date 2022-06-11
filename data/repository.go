package data

type TransactionalRepository interface {
	Begin() error
	Commit() error
	Rollback() error
}

type Repository interface {
	TransactionalRepository
	Close()
	Create(value interface{}) error
}

type Predicate struct {
	Query string
	Args  []interface{}
}

func NewPredicate(query string, args ...interface{}) Predicate {
	return Predicate{query, args}
}

type Column struct {
	Name  string
	Value interface{}
}

func NewColumn(name string, value interface{}) Column {
	return Column{name, value}
}
