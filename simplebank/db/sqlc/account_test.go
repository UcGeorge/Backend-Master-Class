package db

import (
	"context"
	"testing"

	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createAccountFromArg(t *testing.T, arg CreateAccountParams) Account {
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomAccount(t *testing.T, owner User) Account {
	arg := CreateAccountParams{
		Owner:    owner.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account := createAccountFromArg(t, arg)
	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createRandomAccount(t, user)
}

func TestGetAccount(t *testing.T) {
	user := createRandomUser(t)
	account1 := createRandomAccount(t, user)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 0)
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)
	account1 := createRandomAccount(t, user)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 0)
}

func TestDeleteAccount(t *testing.T) {
	user := createRandomUser(t)
	account1 := createRandomAccount(t, user)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for range 10 {
		user := createRandomUser(t)
		createRandomAccount(t, user)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestListAccountsForUser(t *testing.T) {
	user := createRandomUser(t)

	accountUSD := createAccountFromArg(t, CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomAmount(),
		Currency: util.USD,
	})
	accountEUR := createAccountFromArg(t, CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomAmount(),
		Currency: util.EUR,
	})
	accountCAD := createAccountFromArg(t, CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomAmount(),
		Currency: util.CAD,
	})

	arg := ListAccountsForUserParams{
		Limit:    5,
		Offset:   0,
		Username: user.Username,
	}

	accounts, err := testQueries.ListAccountsForUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 3)

	for _, account := range accounts {
		require.NotEmpty(t, account)

		switch account.Currency {
		case util.USD:
			require.Equal(t, account.ID, accountUSD.ID)
		case util.EUR:
			require.Equal(t, account.ID, accountEUR.ID)
		case util.CAD:
			require.Equal(t, account.ID, accountCAD.ID)
		}
	}
}
