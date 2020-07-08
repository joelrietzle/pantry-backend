package db

import (
	//"database/sql"
	"context"
	"fmt"
	"os"
	"time"


	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	qr_row, err := pool.Query(ctx, "select qr from pantry")
	uuid_row, err := pool.Query(ctx, "select uuid from pantry")
	//err = pool.QueryRow(ctx, "select qr from pantry").Scan(&qr)
	for qr_row.Next() && uuid_row.Next() {
		qr, err := qr_row.Values()
		uuid, err := uuid_row.Values()
		fmt.Print(qr)
		fmt.Println(uuid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to obtain qr row value: %v\n", err)
		}
	}
}