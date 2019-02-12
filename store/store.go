package store

import (
	"context"
	"time"

	"github.com/bluehawk27/medcost-go/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const driver string = "mysql"
const ISODateFormat = time.RFC3339 //"2006-01-28T15:04:05Z"

//  StoreType : Interface exposes methods that can later be used for mocking calls
type StoreType interface {
	GetProviderByID(ctx context.Context, id int) (*Provider, error)
	InsertProvider(ctx context.Context, p *Provider) (*Provider, error)
	GetProviderByProviderID(ctx context.Context, id *int) (*Provider, error)
}

// Store : Represents the DB object could add a config here to read in driver and address
type Store struct {
	db *sqlx.DB
}

// NewStore : New DB Connection
func NewStore() (StoreType, error) {
	connection, err := config.DBConnectionString()
	if err != nil {
		return nil, err
	}
	// "admin:Hawk27396!!@tcp(medcost-dev.crjriys0s5jc.us-west-2.rds.amazonaws.com:3306)/medcost-dev?parseTime=true"
	db, err := sqlx.Connect(driver, connection)
	if err != nil {
		return nil, err
	}
	s := &Store{
		db: db,
	}

	return s, nil
}

// GetProviderByID : Get All Provider by pk
func (s *Store) GetProviderByID(ctx context.Context, id int) (*Provider, error) {
	p := Provider{}

	if err := s.db.Select(&p, GetProviderByID, id); err != nil {
		return nil, err
	}

	return &p, nil
}

// GetProviderByProviderID : Getprovider by provider_id
func (s *Store) GetProviderByProviderID(ctx context.Context, id *int) (*Provider, error) {
	p := Provider{}

	if err := s.db.Select(&p, GetProviderByProviderID, id); err != nil {
		return nil, err
	}

	return &p, nil
}

// InsertProvider : Insert Provider
func (s *Store) InsertProvider(ctx context.Context, p *Provider) (*Provider, error) {
	now := time.Now().UTC().Format(ISODateFormat)

	result := s.db.MustExec(InsertProviderNolatlong, p.ProviderID, p.Name, p.Street, p.City, p.Zipcode, p.State, p.HrrDescription, now, p.Latitude, p.Longitude)
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	prov, err := s.GetProviderByID(ctx, int(id))
	if err != nil {
		return nil, err
	}

	return prov, nil
}
