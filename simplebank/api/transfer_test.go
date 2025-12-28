package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/mock"
	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount := randomAccount()
	toAccount := randomAccount()

	fromAccount.Currency = util.USD
	toAccount.Currency = util.USD

	amount := util.RandomInt(1, 1000)

	fromEntry := db.Entry{
		ID:        1,
		AccountID: fromAccount.ID,
		Amount:    -amount,
		CreatedAt: time.Now(),
	}

	toEntry := db.Entry{
		ID:        2,
		AccountID: toAccount.ID,
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	transfer := db.Transfer{
		ID:            0,
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}

	transferTxResult := db.TransferTxResult{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Transfer:    transfer,
		FromEntry:   fromEntry,
		ToEntry:     toEntry,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// Happy Path (200 OK)
		{
			name: "OK",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          amount,
				"currency":        toAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect check for FromAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)

				// Expect check for ToAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)

				// Expect the actual Transfer transaction
				arg := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(transferTxResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransferTxResponse(t, recorder, transferTxResult)
			},
		},

		// Gin Validation Errors (400 Bad Request)
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          amount,
				"currency":        "XYZ",
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          -amount,
				"currency":        toAccount.Currency,
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ZeroAmount",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          0,
				"currency":        toAccount.Currency,
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingToID",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"amount":          amount,
				"currency":        toAccount.Currency,
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingAmount",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"currency":        toAccount.Currency,
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingCurrency",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          amount,
			},
			buildStubs: expectNoAction,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// Business Logic Errors (404 / 400)
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          amount,
				"currency":        toAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)

				// Expect no Transfer transaction
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          amount,
				"currency":        toAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)

				// Expect no Transfer transaction
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Create a specific FromAccount for this test that has the wrong currency
				wrongCurrencyAccount := db.Account(fromAccount)
				wrongCurrencyAccount.Currency = util.EUR

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(wrongCurrencyAccount, nil)

				// Expect no check for ToAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(0)

				// Expect no Transfer transaction
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Return(fromAccount, nil)

				// Create a specific ToAccount for this test that has the wrong currency
				wrongCurrencyAccount := toAccount
				wrongCurrencyAccount.Currency = util.EUR

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(wrongCurrencyAccount, nil)

				// Expect no Transfer transaction
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: Transaction / Internal Errors (500)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfer"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func expectNoAction(store *mockdb.MockStore) {
	// Expect no check for FromAccount
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Any()).
		Times(0)

	// Expect no check for ToAccount
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Any()).
		Times(0)

	// Expect no Transfer transaction
	store.EXPECT().
		TransferTx(gomock.Any(), gomock.Any()).
		Times(0)
}

func requireBodyMatchTransferTxResponse(t *testing.T, body *httptest.ResponseRecorder, expected db.TransferTxResult) {
	data, err := io.ReadAll(body.Body)
	require.NoError(t, err)

	var got db.TransferTxResult
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	// 1. Check timestamps specifically to ensure they are within a valid range (e.g. 1 second)
	// This ensures the dates are correct without tripping over the Monotonic clock.
	require.WithinDuration(t, expected.Transfer.CreatedAt, got.Transfer.CreatedAt, time.Second)
	require.WithinDuration(t, expected.FromEntry.CreatedAt, got.FromEntry.CreatedAt, time.Second)
	require.WithinDuration(t, expected.ToEntry.CreatedAt, got.ToEntry.CreatedAt, time.Second)
	require.WithinDuration(t, expected.FromAccount.CreatedAt, got.FromAccount.CreatedAt, time.Second)
	require.WithinDuration(t, expected.ToAccount.CreatedAt, got.ToAccount.CreatedAt, time.Second)

	// 2. Normalize the timestamps
	// Overwrite the 'got' timestamps with the 'expected' ones.
	// Since we already verified they are close enough in step 1, we can now make them identical
	// to allow require.Equal to check all other fields (ID, Amount, Balance, etc.) automatically.
	got.Transfer.CreatedAt = expected.Transfer.CreatedAt
	got.FromEntry.CreatedAt = expected.FromEntry.CreatedAt
	got.ToEntry.CreatedAt = expected.ToEntry.CreatedAt
	got.FromAccount.CreatedAt = expected.FromAccount.CreatedAt
	got.ToAccount.CreatedAt = expected.ToAccount.CreatedAt

	// 3. Perform deep comparison on the rest of the struct
	require.Equal(t, expected, got)
}
