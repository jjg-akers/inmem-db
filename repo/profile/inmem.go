package profileservice

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/jjg-akers/inmem-db/domain"
)

// DB ...
type DB interface {
	Insert(ctx context.Context, table string, cols []string, vals []string, data []byte) error
	Update(ctx context.Context, table string, col string, val string, data []byte) error
	Get(ctx context.Context, table string, whereCol string, id string) ([][]byte, error)
}

// Service ... the profiles service provides methods for storing and retreiving users
type Service struct {
	db DB
}

// NewService ...
func NewService(db DB) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return &Service{
		db: db,
	}, nil
}

// StoreNewProfile ...
func (u *Service) StoreNewProfile(ctx context.Context, domainProfile domain.Profile) error {

	p := profile{
		ProfileID:   domainProfile.ProfileID,
		ProfileName: domainProfile.ProfileName,
		FirstName:   domainProfile.FirstName,
		LastName:    domainProfile.LastName,
	}

	b, err := marshalProfile(p)
	if err != nil {
		return err
	}

	return u.db.Insert(ctx, "profiles", []string{"id"}, []string{p.ProfileID}, b)
}

// UpdateProfile updates profile's first and last names
func (u *Service) UpdateProfile(ctx context.Context, domainProfile domain.Profile) error {

	b, err := u.db.Get(ctx, "profiles", "id", domainProfile.ProfileID)
	if err != nil {
		return err
	}

	if len(b) < 1 {
		return fmt.Errorf("%w: profile not found", domain.ErrInvalidInput)
	}

	p, err := unmarshalProfile(b[0])
	if err != nil {
		return err
	}

	p.FirstName = domainProfile.FirstName
	p.LastName = domainProfile.LastName

	ub, err := marshalProfile(p)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return u.db.Update(ctx, "profiles", "id", domainProfile.ProfileID, ub)
}

// GetProfile ... retrieves a profile from the db. This satisfies the importstatus.UserService interface
func (u *Service) GetProfile(ctx context.Context, profileID string) (*domain.Profile, error) {

	b, err := u.db.Get(ctx, "profiles", "id", profileID)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	us, err := unmarshalProfile(b[0])
	if err != nil {
		return nil, err
	}

	return &domain.Profile{
		ProfileID:   us.ProfileID,
		ProfileName: us.ProfileName,
		FirstName:   us.FirstName,
		LastName:    us.LastName,
	}, nil
}

func unmarshalProfile(b []byte) (profile, error) {
	us := profile{}
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&us); err != nil {
		return profile{}, fmt.Errorf("invalid data in profile table: %s", string(b))
	}

	return us, nil
}

func marshalProfile(u profile) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(u); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return b.Bytes(), nil
}
