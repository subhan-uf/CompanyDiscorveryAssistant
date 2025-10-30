package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dsn string) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }
    pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
    if err != nil {
        return nil, err
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, err
    }
    return pool, nil
}
