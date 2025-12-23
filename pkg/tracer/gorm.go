package tracer

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	sqltrace "github.com/lsgndln/dd-trace-go/contrib/database/sql"
	gormtrace "github.com/lsgndln/dd-trace-go/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormPostgresTracer(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName(os.Getenv("DD_SERVICE")+"-db"))
	sqlDb, err := sqltrace.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDb}), cfg)
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
