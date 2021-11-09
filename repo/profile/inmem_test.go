package profileservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jjg-akers/inmem-db/db/inmem"
	ps "github.com/jjg-akers/inmem-db/repo/profile"

	"github.com/jjg-akers/inmem-db/domain"

	"github.com/stretchr/testify/assert"
)

func TestProfileService_NewService(t *testing.T) {
	type args struct {
		userStores ps.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "should fail",
			wantErr: fmt.Errorf("db is nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := ps.NewService(tt.args.userStores)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestProfileService_StoreNewUser(t *testing.T) {
	type args struct {
		u domain.Profile
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "should store a profile with all fields",
			args: args{
				u: domain.Profile{
					ProfileID:   "00000000-0000-0000-0000-000000000001",
					ProfileName: "username",
					FirstName:   "test-firstName",
					LastName:    "test-lastName",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := ps.NewService(inmem.NewDB([]inmem.Table{
				{
					Name: "profiles",
				},
			}))

			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.wantErr, s.StoreNewProfile(context.Background(), tt.args.u))
		})
	}
}

func TestProfileService_UpdateUser(t *testing.T) {
	type args struct {
		u domain.Profile
	}
	tests := []struct {
		name    string
		args    args
		setup   func(s *ps.Service) error
		wantErr error
	}{
		{
			name: "should successfully update a profile",
			args: args{
				u: domain.Profile{
					ProfileID: "00000000-0000-0000-0000-000000000001",
					FirstName: "test-firstName",
					LastName:  "test-lastName",
				},
			},
			setup: func(s *ps.Service) error {
				u := domain.Profile{
					ProfileID:   "00000000-0000-0000-0000-000000000001",
					ProfileName: "test-username",
				}

				return s.StoreNewProfile(context.Background(), u)
			},
		},
		{
			name: "should fail due to profile missing",
			args: args{
				u: domain.Profile{
					ProfileID: "00000000-0000-0000-0000-000000000001",
					FirstName: "test-firstName",
					LastName:  "test-lastName",
				},
			},
			setup: func(s *ps.Service) error {
				u := domain.Profile{
					ProfileID:   "00000000-0000-0000-0000-000000000009",
					ProfileName: "test-username",
				}

				return s.StoreNewProfile(context.Background(), u)
			},
			wantErr: fmt.Errorf("%w: profile not found", domain.ErrInvalidInput),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := inmem.NewDB([]inmem.Table{

				{
					Name: "profiles",
				},
			})

			s, err := ps.NewService(db)
			if err != nil {
				t.Fatal(err)
			}

			if tt.setup != nil {
				if err := tt.setup(s); err != nil {
					t.Fatal(err)
				}
			}

			assert.Equal(t, tt.wantErr, s.UpdateProfile(context.Background(), tt.args.u))
		})
	}
}

func TestProfileService_GetUser(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *domain.Profile
		setup        func(s *ps.Service) error
		wantErr      error
	}{
		{
			name: "should successfully get a profile",
			args: args{
				userID: "00000000-0000-0000-0000-000000000001",
			},
			expectedUser: &domain.Profile{
				ProfileID:   "00000000-0000-0000-0000-000000000001",
				ProfileName: "test-username",
			},
			setup: func(s *ps.Service) error {
				u := domain.Profile{
					ProfileID:   "00000000-0000-0000-0000-000000000001",
					ProfileName: "test-username",
				}

				return s.StoreNewProfile(context.Background(), u)
			},
		},
		{
			name: "should return nil if profile not found",
			args: args{
				userID: "00000000-0000-0000-0000-000000000002",
			},
			setup: func(s *ps.Service) error {
				u := domain.Profile{
					ProfileID:   "00000000-0000-0000-0000-000000000001",
					ProfileName: "test-username",
				}

				return s.StoreNewProfile(context.Background(), u)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := inmem.NewDB([]inmem.Table{

				{
					Name: "profiles",
				},
			})

			s, err := ps.NewService(db)
			if err != nil {
				t.Fatal(err)
			}

			if tt.setup != nil {
				if err := tt.setup(s); err != nil {
					t.Fatal(err)
				}
			}

			gotUser, gotErr := s.GetProfile(context.Background(), tt.args.userID)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, gotErr)
				return
			}
			if !assert.Nil(t, gotErr) {
				return
			}

			assert.Equal(t, tt.expectedUser, gotUser)
		})
	}
}
