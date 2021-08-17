package services

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"

	"github.com/ehabterra/flash_api/mocks"

	"github.com/ehabterra/flash_api/internal/models"
)

func TestUsers_GetUserByUsernameOrEmail(t *testing.T) {
	bank := &mocks.Bank{}
	db := &mocks.Database{}

	id := uuid.New().String()
	email := "ehab@email.com"
	db.On("GetUserByUsernameOrEmail", mock.AnythingOfType("string")).Return(func(s string) *models.User {
		return &models.User{
			ID:       id,
			Username: s,
			Email:    email,
			Password: "pass",
			Balance:  10,
		}
	}, nil)
	type fields struct {
		bank Bank
		db   Database
	}
	type args struct {
		usernameOrEmail string
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name:   "Test1",
			fields: fields{bank: bank, db: db},
			args: args{
				"ehab",
			},
			want: &models.User{
				ID:       id,
				Username: "ehab",
				Email:    email,
				Password: "pass",
				Balance:  10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				bank: tt.fields.bank,
				db:   tt.fields.db,
			}
			got, err := u.GetUserByUsernameOrEmail(tt.args.usernameOrEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByUsernameOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByUsernameOrEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
