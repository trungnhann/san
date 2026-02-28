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
	"github.com/jackc/pgx/v4"
)

type Database struct {
	Conn    *pgx.Conn
	Config  conf.Config
	Queries *db.Queries
}

func NewDatabase(cfg conf.Config) *Database {
	cfg.DBUrl = buildDBURL(cfg)
	conn := Connect(cfg)

	return &Database{
		Conn:    conn,
		Config:  cfg,
		Queries: db.New(conn),
	}
}

func (d *Database) Close() {
	if d.Conn != nil {
		Close(d.Conn)
	}
}

func (d *Database) Migrate() {
	AutoMigrate(d.Config)
}

func (d *Database) HealthCheck(ctx context.Context) error {
	if d.Conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return d.Conn.Ping(ctx)
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

func Connect(config conf.Config) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config.DBUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}

func Close(conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		log.Fatalf("Unable to close connection: %v\n", err)
	}
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
