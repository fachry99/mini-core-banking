package repository

import "github.com/jmoiron/sqlx"

type IdempotencyRepository struct {
	DB *sqlx.DB
}

func NewIdempotencyRepository(db *sqlx.DB) *IdempotencyRepository {
	return &IdempotencyRepository{DB: db}
}

func (r *IdempotencyRepository) Get(key string) ([]byte, error) {
	var resp []byte
	err := r.DB.Get(&resp,
		`SELECT response FROM idempotency_keys WHERE key = $1`,
		key,
	)
	return resp, err
}

func (r *IdempotencyRepository) Save(
	key string,
	requestHash string,
	response []byte,
) error {
	_, err := r.DB.Exec(`
        INSERT INTO idempotency_keys (key, request_hash, response)
        VALUES ($1, $2, $3)
    `, key, requestHash, response)
	return err
}
