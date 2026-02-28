package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	conf "san/internal/config"
	db "san/internal/db/sqlc"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Pool    *pgxpool.Pool
	Config  conf.Config
	Queries *db.Queries
}

func NewDatabase(cfg conf.Config) *Database {
	cfg.DBUrl = buildDBURL(cfg)
	pool := Connect(cfg)

	return &Database{
		Pool:    pool,
		Config:  cfg,
		Queries: db.New(pool),
	}
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

func (d *Database) Migrate() {
	AutoMigrate(d.Config)
}

func (d *Database) HealthCheck(ctx context.Context) error {
	if d.Pool == nil {
		return fmt.Errorf("database connection pool is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return d.Pool.Ping(ctx)
}

func buildDBURL(cfg conf.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHostname,
		cfg.DBPort,
		cfg.DBName,
	)
}

func Connect(config conf.Config) *pgxpool.Pool {
	config.DBUrl = buildDBURL(config)
	poolConfig, err := pgxpool.ParseConfig(config.DBUrl)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v\n", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return pool
}

func Close(pool *pgxpool.Pool) {
	pool.Close()
}

func AutoMigrate(config conf.Config) {
	path := config.MigrationPath
	if !strings.HasPrefix(path, "file://") {
		path = fmt.Sprintf("file://%s", path)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHostname,
		config.DBPort,
		config.DBName,
	)

	m, err := migrate.New(path, dsn)
	if err != nil {
		log.Fatalf("unable to create migration: %v\n", err)
	}

	if config.DBRecreate {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatalf("unable to drop database: %v\n", err)
			}
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("unable to migrate database: %v\n", err)
	}
}

func Drop(config conf.Config) {
	path := fmt.Sprintf("file://%s", config.MigrationPath)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHostname,
		config.DBPort,
		config.DBName,
	)

	m, err := migrate.New(path, dsn)
	if err != nil {
		log.Fatalf("unable to create migration: %v\n", err)
	}
	if err := m.Down(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("unable to drop database: %v\n", err)
		}
	}
}
