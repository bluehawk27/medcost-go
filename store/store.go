package store

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/config"

	// driver required by sqlx
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const driver string = "mysql"
const ISODateFormat = time.RFC3339 //"2006-01-28T15:04:05Z"

//  StoreType Interface exposes methods that can later be used for mocking calls
type StoreType interface {
	GetProviderByID(ctx context.Context, id int) (*[]Provider, error)
	GetProviderByProviderID(ctx context.Context, id *int) (*[]Provider, error)
	InsertProvider(ctx context.Context, p *Provider) (*[]Provider, error)
	GetDrgByID(ctx context.Context, id int) (*[]DRG, error)
	GetDrgByDRGID(ctx context.Context, id *int) (*[]DRG, error)
	InsertDrg(ctx context.Context, drg *DRG) (*[]DRG, error)
	GetInpatientServiceByID(ctx context.Context, id int) (*[]Inpatient, error)
	GetInpatientServiceByProviderID(ctx context.Context, id *int) (*[]Inpatient, error)
	GetInpatientServiceByProvIDDrgIDYear(ctx context.Context, provID *int, drgID *int, year *int) (*[]Inpatient, error)
	InsertInpatientService(ctx context.Context, in *Inpatient) (*[]Inpatient, error)
}

// Store : Represents the DB object could add a config here to read in driver and address
type Store struct {
	db *sqlx.DB
}

// NewStore : New DB Connection
func NewStore() (StoreType, error) {
	connection, err := config.DBConnectionString()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	db, err := sqlx.Connect(driver, connection)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s := &Store{
		db: db,
	}

	return s, nil
}

// GetProviderByID : Get All Provider by pk
func (s *Store) GetProviderByID(ctx context.Context, id int) (*[]Provider, error) {
	p := []Provider{}

	if err := s.db.Select(&p, getProviderByID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &p, nil
}

// GetProviderByProviderID : Getprovider by provider_id
func (s *Store) GetProviderByProviderID(ctx context.Context, id *int) (*[]Provider, error) {
	p := []Provider{}

	if err := s.db.Select(&p, getProviderByProviderID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &p, nil
}

// InsertProvider : Insert Provider
func (s *Store) InsertProvider(ctx context.Context, p *Provider) (*[]Provider, error) {
	now := time.Now().UTC().Format(ISODateFormat)
	p.CreatedAt = &now

	result, err := s.db.NamedExec(insertProvider, *p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	prov, err := s.GetProviderByID(ctx, int(id))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return prov, nil
}

// GetDrgByID : Get DRG by pk
func (s *Store) GetDrgByID(ctx context.Context, id int) (*[]DRG, error) {
	drg := []DRG{}

	if err := s.db.Select(&drg, getDrgByID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &drg, nil
}

// GetDrgByDRGID : Get DRG by code
func (s *Store) GetDrgByDRGID(ctx context.Context, id *int) (*[]DRG, error) {
	drg := []DRG{}

	if err := s.db.Select(&drg, getDrgByDRGID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &drg, nil
}

// InsertDrg : Insert Diagnostic Related Group
func (s *Store) InsertDrg(ctx context.Context, drg *DRG) (*[]DRG, error) {
	result, err := s.db.NamedExec(insertDrg, *drg)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	drgInserted, err := s.GetDrgByID(ctx, int(id))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return drgInserted, nil
}

// GetInpatientServiceByID : Get All Provider by pk
func (s *Store) GetInpatientServiceByID(ctx context.Context, id int) (*[]Inpatient, error) {
	in := []Inpatient{}

	if err := s.db.Select(&in, getInpatientServiceByID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &in, nil
}

// GetInpatientServiceByProviderID : Getprovider by provider_id
func (s *Store) GetInpatientServiceByProviderID(ctx context.Context, id *int) (*[]Inpatient, error) {
	in := []Inpatient{}

	if err := s.db.Select(&in, getInpatientServiceByProviderID, id); err != nil {
		log.Error(err)
		return nil, err
	}

	return &in, nil
}

// GetInpatientServiceByProviderID : Getprovider by provider_id
func (s *Store) GetInpatientServiceByProvIDDrgIDYear(ctx context.Context, provID *int, drgID *int, year *int) (*[]Inpatient, error) {
	in := []Inpatient{}

	if err := s.db.Select(&in, getInpatientServiceByProvIDDrgIDYear, provID, drgID, year); err != nil {
		log.Error(err)
		return nil, err
	}

	return &in, nil
}

// InsertInpatientService : Insert Provider
func (s *Store) InsertInpatientService(ctx context.Context, in *Inpatient) (*[]Inpatient, error) {
	now := time.Now().UTC().Format(ISODateFormat)
	in.CreatedAt = &now

	result, err := s.db.NamedExec(insertInpatientService, *in)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	inp, err := s.GetInpatientServiceByID(ctx, int(id))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return inp, nil
}
