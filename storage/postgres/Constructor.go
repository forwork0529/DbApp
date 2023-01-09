package postgres


import (

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct{
	db *pgxpool.Pool
}

func NewDb(ConnString string)(*Storage, error){
	db, err := pgxpool.Connect(context.Background(), ConnString)
	if err != nil{
		return nil, err
	}
	return &Storage{db : db}, nil
}

