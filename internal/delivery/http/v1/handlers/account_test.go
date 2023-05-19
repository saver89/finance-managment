package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/saver89/finance-management/config"
	"github.com/saver89/finance-management/internal/domain"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	mockdb "github.com/saver89/finance-management/internal/repository/postgres/sqlc/mock"
	"github.com/saver89/finance-management/internal/service/response"
	"github.com/saver89/finance-management/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name       string
		accountID  int64
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// TODO: Add more test cases.
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			log := logger.NewApiLogger(&config.Config{Logger: config.Logger{Level: "panic"}})
			log.InitLogger()
			srv := NewServer(store, log)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/v1/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			srv.router.ServeHTTP(recorder, request)

			tc.checkResp(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:         rand.Int63(),
		OfficeID:   rand.Int63(),
		Name:       gofakeit.Name(),
		CurrencyID: rand.Int63(),
		CreatedBy:  rand.Int63(),
		State:      db.AccountStateActive,
		Balance:    fmt.Sprintf("%f", rand.Float64()),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	balanceFloat, err := strconv.ParseFloat(account.Balance, 64)
	require.NoError(t, err)

	var gotAccountResp response.GetAccountResponse
	err = json.Unmarshal(data, &gotAccountResp)
	require.NoError(t, err)
	require.Equal(t, domain.Account{
		ID:         account.ID,
		OfficeID:   account.OfficeID,
		Name:       account.Name,
		CurrencyID: account.CurrencyID,
		CreatedBy:  account.CreatedBy,
		State:      string(account.State),
		Balance:    balanceFloat,
		CreatedAt:  "0001-01-01 00:00:00",
	}, gotAccountResp.Account)
}
