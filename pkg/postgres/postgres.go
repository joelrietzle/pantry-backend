package postgres

import (
	"context"

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

func (q *PantryStore) ListDevice(ctx context.Context, qrString string, UUID string) ([]Pantry, []Pantry, error) {
	rows, err := q.pool.Query(ctx, "SELECT qr, uuid FROM pantry WHERE lower(qr) = $1 OR lower(uuid) = $2", qrString, UUID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	uuids := []Pantry{}
	qrs := []Pantry{}

	for rows.Next() {
		uuid := Pantry{}
		qr := Pantry{}
		if err := rows.Scan(&qr.QR, &uuid.UUID); err != nil {
			return nil, nil, err
		}

		qrs = append(qrs, qr)
		uuids = append(uuids, uuid)

	}

	return qrs, uuids, err
}

func (m *MerchantStore) ListMerchants(ctx context.Context, firstName string, lastName string, email string) ([]Merchant, []Merchant, []Merchant, error) {
	rows, err := m.pool.Query(ctx, "SELECT FirstName, LastName, Email FROM merchant WHERE lower(FirstName) = $1 OR lower(LastName) = $2 OR lower(Email) = $3", firstName, lastName, email)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	firstnames := []Merchant{}
	lastnames := []Merchant{}
	emails := []Merchant{}

	for rows.Next() {
		firstname := Merchant{}
		lastname := Merchant{}
		email := Merchant{}
		if err := rows.Scan(&firstname.FirstName, &lastname.LastName, &email.Email); err != nil {
			return nil, nil, nil, err
		}

		firstnames = append(firstnames, firstname)
		lastnames = append(lastnames, lastname)
		emails = append(emails, email)
	}

	return firstnames, lastnames, emails, nil
}

func (c *CustomerStore) ListCustomers(ctx context.Context, firstName string, lastName string, email string) ([]Customer, []Customer, []Customer, error) {
	rows, err := c.pool.Query(ctx, "SELECT FirstName, LastName, Email FROM customer WHERE lower(FirstName) = $1 OR lower(LastName) = $2 OR lower(Email) = $3", firstName, lastName, email)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	firstnames := []Customer{}
	lastnames := []Customer{}
	emails := []Customer{}

	for rows.Next() {
		firstname := Customer{}
		lastname := Customer{}
		email := Customer{}
		if err := rows.Scan(&firstname.FirstName, &lastname.LastName, &email.Email); err != nil {
			return nil, nil, nil, err
		}

		firstnames = append(firstnames, firstname)
		lastnames = append(lastnames, lastname)
		emails = append(emails, email)
	}

	return firstnames, lastnames, emails, nil
}

func (c *PantryStore) LockUnlockQR(ctx context.Context, qrString string) ([]Pantry, error) {
	rows, err := c.pool.Query(ctx, "SELECT uuid FROM pantry WHERE qr = $1", qrString)
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

	return uuids, nil
}

func (c *PantryStore) LockUnlockUUID(ctx context.Context, uuidString string) ([]Pantry, error) {
	rows, err := c.pool.Query(ctx, "SELECT qr FROM pantry WHERE uuid = $1", uuidString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	qrs := []Pantry{}

	for rows.Next() {
		qr := Pantry{}

		if err := rows.Scan(&qr.QR); err != nil {
			return nil, err
		}

		qrs = append(qrs, qr)
	}

	return qrs, nil
}
