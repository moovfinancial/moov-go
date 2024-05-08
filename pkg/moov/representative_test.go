package moov_test

import (
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Representatives(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	create := moov.CreateRepresentative{
		Name: moov.Name{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
		},
		Phone: &moov.Phone{
			Number:      "5555555555",
			CountryCode: "1",
		},
		Email: faker.Email(),
		BirthDate: &moov.Date{
			Year:  1980,
			Month: 1,
			Day:   1,
		},
		Address: &moov.Address{
			AddressLine1:    "123 Main St",
			City:            "Anytown",
			StateOrProvince: "CA",
			Country:         "US",
			PostalCode:      "90210",
		},
		Responsibilities: &moov.Responsibilities{
			IsController:        true,
			IsOwner:             true,
			OwnershipPercentage: 50,
			JobTitle:            "CEO",
		},
	}

	resp, err := mc.CreateRepresentative(BgCtx(), account.AccountID, create)

	t.Run("create", func(t *testing.T) {
		NoResponseError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, create.Name, resp.Name)
		require.Equal(t, create.Phone, resp.Phone)
		require.Equal(t, create.Email, resp.Email)
		require.Equal(t, true, resp.BirthDateProvided)
		require.Equal(t, create.Address, resp.Address)
		require.Equal(t, create.Responsibilities, resp.Responsibilities)
	})

	t.Run("get", func(t *testing.T) {
		rep, err := mc.GetRepresentative(BgCtx(), account.AccountID, resp.RepresentativeID)
		NoResponseError(t, err)
		require.NotNil(t, rep)
	})

	t.Run("list", func(t *testing.T) {
		representatives, err := mc.ListRepresentatives(BgCtx(), account.AccountID)

		NoResponseError(t, err)
		require.NotEmpty(t, representatives)
	})

	update := moov.UpdateRepresentative{
		Name: moov.Name{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
		},
		Phone: &moov.Phone{
			Number:      "1111111111",
			CountryCode: "1",
		},
		Email: faker.Email(),
		BirthDate: &moov.Date{
			Year:  1999,
			Month: 2,
			Day:   2,
		},
		Address: &moov.Address{
			AddressLine1:    "321 Main St",
			City:            "Anytown",
			StateOrProvince: "UT",
			Country:         "US",
			PostalCode:      "84096",
		},
		Responsibilities: &moov.Responsibilities{
			IsController:        false,
			IsOwner:             true,
			OwnershipPercentage: 25,
			JobTitle:            "CTO",
		},
	}

	t.Run("update", func(t *testing.T) {
		actual, err := mc.UpdateRepresentative(BgCtx(), account.AccountID, resp.RepresentativeID, update)

		NoResponseError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, update.Name, actual.Name)
		require.Equal(t, update.Phone, actual.Phone)
		require.Equal(t, update.Email, actual.Email)
		require.Equal(t, update.Address, actual.Address)
		require.Equal(t, update.Responsibilities, actual.Responsibilities)
	})

	t.Run("delete", func(t *testing.T) {
		err = mc.DeleteRepresentative(BgCtx(), account.AccountID, resp.RepresentativeID)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
	})
}
