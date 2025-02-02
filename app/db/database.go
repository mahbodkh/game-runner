package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

var Conn *pgx.Conn

func InitDB() {
	connStr := "postgresql://neondb_owner:npg_u7TlqOfk2JtM@ep-patient-dawn-a2goszlp-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require"

	var err error
	Conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	err = Conn.Ping(context.Background())
	if err != nil {
		logrus.Fatalf("Unable to ping database: %v", err)
	}

	dropTables()
	// truncateTables()
	createTables()

	logrus.Println("Database initialized successfully!")
}

func createTables() {

	_, err := Conn.Exec(context.Background(), `
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'session_status') THEN
                CREATE TYPE session_status AS ENUM (
                    'started',
                    'in_progress',
                    'full',
                    'invalid',
                    'error',
                    'waiting',
                    'created',
                    'moved',
                    'completed'
                );
            END IF;
        END
        $$;
    `)
	if err != nil {
		logrus.Fatalf("Unable to create session_status type: %v", err)
	}

	_, err = Conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS players (
            id SERIAL PRIMARY KEY,
            telegram_id INT,
            name VARCHAR(255),
            language_code VARCHAR(10),
            mark CHAR(1),
            created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS sessions (
            id SERIAL PRIMARY KEY,
            status session_status,
            current_player INT,
            created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS session_players (
            session_id INT,
            player_id INT,
            PRIMARY KEY (session_id, player_id)
        );
    `)
	if err != nil {
		logrus.Fatalf("Unable to create tables: %v", err)
	}
}

func truncateTables() {
	_, err := Conn.Exec(context.Background(), `
        TRUNCATE TABLE players, sessions, session_players RESTART IDENTITY CASCADE;
    `)
	if err != nil {
		logrus.Fatalf("Unable to truncate tables: %v", err)
	}
}

func dropTables() {
	_, err := Conn.Exec(context.Background(), `
        DROP TABLE IF EXISTS players, sessions, session_players CASCADE;
        DROP TYPE IF EXISTS session_status;
    `)
	if err != nil {
		logrus.Fatalf("Unable to drop tables: %v", err)
	}
}
