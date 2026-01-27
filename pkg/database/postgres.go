package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-starter-kit/internal/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *config.EnvConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}
	return db, nil
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type txManager struct {
	db *sql.DB
}

func NewDBTransaction(db *sql.DB) TxManager {
	return &txManager{db: db}
}

func (t *txManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
