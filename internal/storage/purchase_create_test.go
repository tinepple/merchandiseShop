//go:build integration

package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_CreatePurchase(t *testing.T) {
	type args struct {
		userID     int
		itemID     int
		newBalance int
	}
	testCases := []struct {
		name          string
		initDbQueries []string
		args          args
		expectedErr   error
	}{
		{
			name: "Успешно создана покупка",
			initDbQueries: []string{
				`
					insert into users 
						(id,username,password) 
					values 
						(1, 'test', 'test')
				`,
				`
					insert into balances
						(user_id,balance)
					values 
						(1, 1000)
    			`,
			},
			args: args{
				userID:     1,
				itemID:     3,
				newBalance: 900,
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

			err = storageRepo.CreatePurchase(context.Background(), tc.args.userID, tc.args.itemID, tc.args.newBalance)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
