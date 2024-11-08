package main

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE parcels (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
		)
	`)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	// prepare
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// get
	gotParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel.Client, gotParcel.Client)
	assert.Equal(t, parcel.Status, gotParcel.Status)
	assert.Equal(t, parcel.Address, gotParcel.Address)
	assert.Equal(t, parcel.CreatedAt, gotParcel.CreatedAt)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// cant get it from db
	_, err = store.Get(id)
	require.Error(t, err)

	//like this?
	require.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows error")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, gotParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, gotParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db := setupTestDB(t)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		assert.Contains(t, parcelMap, parcel.Number)
		assert.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		assert.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		assert.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		assert.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
	}
}
