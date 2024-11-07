package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(parcel Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcels (client, status, adress, created_at) VALUES (?, ?, ?, ?)", parcel.Client, parcel.Status, parcel.Address, parcel.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(id int) (Parcel, error) {
	var parcel Parcel
	err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcels WHERE number = ?", id).Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		return parcel, err
	}
	return parcel, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
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

func (s ParcelStore) SetStatus(id int, status string) error {
	_, err := s.db.Exec("UPDATE parcels SET status = ? WHERE number = ?", status, id)
	return err
}

func (s ParcelStore) SetAddress(id int, address string) error {
	_, err := s.db.Exec("UPDATE parcels SET address = ? WHERE number = ?", address, id)
	return err
}

func (s ParcelStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM parcels WHERE number = ?", id)
	return err
}
