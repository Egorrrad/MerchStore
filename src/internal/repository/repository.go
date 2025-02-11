package repository

type Database interface {
}

type Repository struct {
	database Database
}

func NewRepository(database Database) Repository {
	return Repository{
		database: database,
	}
}
