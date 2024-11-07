package main

import (
	"database/sql"
)

type ParcelStore interface {
	Add(parcel Parcel) (int, error)
	Get(id int) (Parcel, error)
	GetByClient(client int) ([]Parcel, error)
	SetStatus(id int, status string) error
	SetAddress(id int, address string) error
	Delete(id int) error
}

type SQLiteParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return &SQLiteParcelStore{db: db}
}

func (s *SQLiteParcelStore) Add(parcel Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcels (client, status, address, created_at) VALUES (?, ?, ?, ?)", parcel.Client, parcel.Status, parcel.Address, parcel.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SQLiteParcelStore) Get(id int) (Parcel, error) {
	var parcel Parcel
	err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcels WHERE number = ?", id).Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		return parcel, err
	}
	return parcel, nil
}

func (s *SQLiteParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcels WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, parcel)
	}
	return parcels, nil
}

func (s *SQLiteParcelStore) SetStatus(id int, status string) error {
	_, err := s.db.Exec("UPDATE parcels SET status = ? WHERE number = ?", status, id)
	return err
}

func (s *SQLiteParcelStore) SetAddress(id int, address string) error {
	_, err := s.db.Exec("UPDATE parcels SET address = ? WHERE number = ?", address, id)
	return err
}

func (s *SQLiteParcelStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM parcels WHERE number = ?", id)
	return err
}
