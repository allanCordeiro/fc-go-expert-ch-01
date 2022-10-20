package service

import (
	"context"
	"database/sql"
	"log"
	"server/entity"
)

type Data struct {
	Bid entity.Bid
}

type Repository interface {
	InsertBID(ctx context.Context, db *sql.DB, data entity.Bid) error
}

func (d *Data) InsertBID(ctx context.Context, db *sql.DB) error {
	log.Println("Inserting row")
	stmt, err := db.Prepare("INSERT INTO bid (value) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, d.Bid.BidValue)
	if err != nil {
		return err
	}
	return nil
}
