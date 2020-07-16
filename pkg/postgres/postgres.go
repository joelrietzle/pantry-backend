package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pantry struct {
	QR   string
	UUID string
}

type Merchant struct {
	FirstName string
	LastName  string
	Email     string
}

type Customer struct {
	FirstName string
	LastName  string
	Email     string
}

type PantryStore struct {
	pool *pgxpool.Pool
}

type MerchantStore struct {
	pool *pgxpool.Pool
}

type CustomerStore struct {
	pool *pgxpool.Pool
}

type QrList struct {
	pool *pgxpool.Pool
	QR   string
}

func NewPantryStore(pool *pgxpool.Pool) *PantryStore {
	return &PantryStore{
		pool: pool,
	}
}
func NewMerchantStore(pool *pgxpool.Pool) *MerchantStore {
	return &MerchantStore{
		pool: pool,
	}
}

func NewCustomerStore(pool *pgxpool.Pool) *CustomerStore {
	return &CustomerStore{
		pool: pool,
	}
}

func (p *PantryStore) ListPantries(ctx context.Context) ([]Pantry, error) {
	rows, err := p.pool.Query(ctx, "SELECT qr, uuid FROM pantry")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pantries := []Pantry{}

	for rows.Next() {
		pantry := Pantry{}
		if err := rows.Scan(&pantry.QR, &pantry.UUID); err != nil {
			return nil, err
		}

		pantries = append(pantries, pantry)
	}

	return pantries, nil
}

func (q *PantryStore) ListDeviceByQR(ctx context.Context, qrString string) ([]Pantry, error) {
	fmt.Println("QR SHUDOASHID:", qrString)
	rows, err := q.pool.Query(ctx, "SELECT uuid FROM pantry WHERE qr = $1", qrString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uuids := []Pantry{}

	for rows.Next() {
		uuid := Pantry{}
		if err := rows.Scan(&uuid.UUID); err != nil {
			return nil, err
		}

		uuids = append(uuids, uuid)

	}

	return uuids, err
}

func (m *MerchantStore) ListMerchants(ctx context.Context) ([]Merchant, error) {
	rows, err := m.pool.Query(ctx, "SELECT FirstName, LastName, Email FROM merchant")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	merchants := []Merchant{}

	for rows.Next() {
		merchant := Merchant{}
		if err := rows.Scan(&merchant.FirstName, &merchant.LastName, &merchant.Email); err != nil {
			return nil, err
		}

		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (c *CustomerStore) ListCustomers(ctx context.Context, firstName string, lastName string, email string) ([]Customer, error) {
	rows, err := c.pool.Query(ctx, "SELECT FirstName, LastName, Email FROM customer WHERE FirstName = $1 OR LastName = $2 OR Email = $3", firstName, lastName, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		customer := Customer{}
		if err := rows.Scan(&customer.FirstName, &customer.LastName, &customer.Email); err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
