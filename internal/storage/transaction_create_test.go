//go:build integration

package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_CreateTransaction(t *testing.T) {
	type args struct {
		transaction        Transaction
		NewBalanceUserFrom int
		NewBalanceUserTo   int
	}
	testCases := []struct {
		name          string
		initDbQueries []string
		args          args
		expectedErr   error
	}{
		{
			name: "Успешно переданы монетки ",
			initDbQueries: []string{
				`
					insert into users 
						(id,username,password) 
					values 
						(1, 'test', 'test'),
						(2, 'user', 'user')
				`,
				`
					insert into balances
						(user_id,balance) 
					values 
						(1, 1000),
						(2, 1000)
    			`,
			},
			args: args{
				transaction: Transaction{
					UserIDFrom: 1,
					UserIDTo:   2,
					Amount:     100,
				},
				NewBalanceUserFrom: 900,
				NewBalanceUserTo:   1100,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := clearTableUsers()
			if err != nil {
				t.Fatal(err)
			}

			if err = prepareDB(tc.initDbQueries...); err != nil {
				t.Fatal(err)
			}

			err = storageRepo.CreateTransaction(context.Background(), tc.args.transaction, tc.args.NewBalanceUserFrom, tc.args.NewBalanceUserTo)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
