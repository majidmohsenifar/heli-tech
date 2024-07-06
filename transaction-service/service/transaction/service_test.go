package transaction_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/helper"
	"github.com/majidmohsenifar/heli-tech/transaction-service/logger"
	"github.com/majidmohsenifar/heli-tech/transaction-service/mocks"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alicebob/miniredis/v2"
)

func TestService_Withdraw_CannotObtainLock(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	//first we try to obtain lock so the logic would encounter error
	lock, err := redisLocker.Obtain(ctx, fmt.Sprintf("transaction:%d", 1), time.Duration(30*time.Second), nil)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, transaction.ErrOngoingRequest)
	lock.Release(ctx)
	repo.AssertExpectations(t)
}

func TestService_Withdraw_NoBalanceInDB(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserBalanceByUserID(
		mock.Anything,
		mock.Anything,
		int64(1),
	).Once().Return(repository.UserBalance{}, pgx.ErrNoRows)
	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, transaction.ErrInsufficientBalance)
	repo.AssertExpectations(t)
}

func TestService_Withdraw_BalanceExistInDB_ButNotSufficient(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserBalanceByUserID(
		mock.Anything,
		mock.Anything,
		int64(1),
	).Once().Return(repository.UserBalance{Amount: pgtype.Numeric{Int: big.NewInt(50), Valid: true}}, nil)
	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, transaction.ErrInsufficientBalance)
	repo.AssertExpectations(t)
}

func TestService_Withdraw_CannotCreateTransaction(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserBalanceByUserID(
		mock.Anything,
		mock.Anything,
		int64(1),
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(100), Valid: true}}, nil)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			if p.Kind != repository.KindWITHDRAW {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{}, errors.New("db error"))

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, errors.New("cannot create transaction"))
	repo.AssertExpectations(t)
}

func TestService_Withdraw_CannotInsertOrIncreaseUserBalance(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserBalanceByUserID(
		mock.Anything,
		mock.Anything,
		int64(1),
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(100), Valid: true}}, nil)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			if p.Kind != repository.KindWITHDRAW {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{ID: 1}, nil)

	repo.EXPECT().CreateUserBalanceOrDecreaseAmount(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserBalanceOrDecreaseAmountParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			return true
		}),
	).Once().Return(repository.UserBalance{}, errors.New("db error"))

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, errors.New("cannot update user balance"))
	repo.AssertExpectations(t)
}

func TestService_Withdraw_Successful(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserBalanceByUserID(
		mock.Anything,
		mock.Anything,
		int64(1),
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(100), Valid: true}}, nil)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			if p.Kind != repository.KindWITHDRAW {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{
		ID:        1,
		Amount:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
		UserID:    1,
		Kind:      repository.KindWITHDRAW,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}, nil)

	repo.EXPECT().CreateUserBalanceOrDecreaseAmount(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserBalanceOrDecreaseAmountParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			return true
		}),
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(120), Valid: true}}, nil)

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionEventManager.EXPECT().PublishTransactionCreatedEvent(
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(transaction.TransactionCreatedEventParams)
			if p.UserID != 1 {
				return false
			}
			if p.TransactionID != 1 {
				return false
			}
			if p.Kind != "WITHDRAW" {
				return false
			}
			if p.Amount != 100 {
				return false
			}
			if p.Balance != 120 {
				return false
			}
			if p.CreatedAt <= 0 {
				return false
			}
			return true
		}),
	).Once().Return()
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	res, err := transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Nil(err)
	assert.Equal(res.ID, int64(1))
	assert.Equal(res.Amount, 100.0)
	assert.Equal(res.NewBalance, 120.0)
	assert.Greater(res.CreatedAt, int64(0))
	repo.AssertExpectations(t)
}

func TestService_Deposit_CannotObtainLock(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	//first we try to obtain lock so the logic would encounter error
	lock, err := redisLocker.Obtain(ctx, fmt.Sprintf("transaction:%d", 1), time.Duration(30*time.Second), nil)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, transaction.ErrOngoingRequest)
	lock.Release(ctx)
	repo.AssertExpectations(t)
}

func TestService_Deposit_CannotCreateTransaction(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			if p.Kind != repository.KindDEPOSIT {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{}, errors.New("db error"))

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, errors.New("cannot create transaction"))
	repo.AssertExpectations(t)
}

func TestService_Deposit_CannotInsertOrIncreaseUserBalance(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			if p.Kind != repository.KindDEPOSIT {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{ID: 1}, nil)

	repo.EXPECT().CreateUserBalanceOrIncreaseAmount(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserBalanceOrIncreaseAmountParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			return true
		}),
	).Once().Return(repository.UserBalance{}, errors.New("db error"))

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, errors.New("cannot update user balance"))
	repo.AssertExpectations(t)
}

func TestService_Deposit_Successful(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().CreateTransaction(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateTransactionParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			if p.Kind != repository.KindDEPOSIT {
				return false
			}
			return true
		}),
	).Once().Return(repository.Transaction{
		ID:        1,
		Amount:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
		UserID:    1,
		Kind:      repository.KindDEPOSIT,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}, nil)

	repo.EXPECT().CreateUserBalanceOrIncreaseAmount(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserBalanceOrIncreaseAmountParams)
			if p.UserID != 1 {
				return false
			}
			amount, err := helper.PGNumericToFloat64(p.Amount)
			if err != nil {
				return false
			}
			if amount != 100 {
				return false
			}
			return true
		}),
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(120), Valid: true}}, nil)

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionEventManager.EXPECT().PublishTransactionCreatedEvent(
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(transaction.TransactionCreatedEventParams)
			if p.UserID != 1 {
				return false
			}
			if p.TransactionID != 1 {
				return false
			}
			if p.Kind != "DEPOSIT" {
				return false
			}
			if p.Amount != 100 {
				return false
			}
			if p.Balance != 120 {
				return false
			}
			if p.CreatedAt <= 0 {
				return false
			}
			return true
		}),
	).Once().Return()
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	res, err := transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Nil(err)
	assert.Equal(res.ID, int64(1))
	assert.Equal(res.Amount, 100.0)
	assert.Equal(res.NewBalance, 120.0)
	assert.Greater(res.CreatedAt, int64(0))
	repo.AssertExpectations(t)
}

func TestService_GetUserTransactions_DBError(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserTransactionsByPagination(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.GetUserTransactionsByPaginationParams)
			if p.UserID != 1 {
				return false
			}
			if p.Limit != 10 {
				return false
			}
			if p.Offset != 10 {
				return false
			}
			return true
		}),
	).Once().Return(nil, errors.New("dbError"))

	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		nil,
		logger,
		transactionEventManager,
	)
	_, err = transactionService.GetUserTransactions(ctx, transaction.GetUserTransactionsParams{
		UserID:   1,
		Page:     1,
		PageSize: 10,
	})
	assert.Equal(err, errors.New("cannot get user transactions"))

	repo.AssertExpectations(t)
}

func TestService_GetUserTransactions_Successful(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	assert.Nil(err)
	defer dbMock.Close()
	ctx := context.Background()

	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserTransactionsByPagination(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.GetUserTransactionsByPaginationParams)
			if p.UserID != 1 {
				return false
			}
			if p.Limit != 10 {
				return false
			}
			if p.Offset != 10 {
				return false
			}
			return true
		}),
	).Once().Return([]repository.Transaction{
		{
			ID:        2,
			UserID:    1,
			Kind:      repository.KindDEPOSIT,
			Amount:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        1,
			UserID:    1,
			Kind:      repository.KindWITHDRAW,
			Amount:    pgtype.Numeric{Int: big.NewInt(50), Valid: true},
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}, nil)

	logger := logger.NewLogger()
	transactionEventManager := new(mocks.MockTransactionEventManager)
	transactionService := transaction.NewService(
		dbMock,
		repo,
		nil,
		logger,
		transactionEventManager,
	)
	txs, err := transactionService.GetUserTransactions(ctx, transaction.GetUserTransactionsParams{
		UserID:   1,
		Page:     1,
		PageSize: 10,
	})
	assert.Nil(err)

	tx1 := txs[0]
	assert.Equal(tx1.ID, int64(2))
	assert.Equal(tx1.Amount, 100.0)
	assert.Equal(tx1.Kind, "DEPOSIT")
	assert.Greater(tx1.CreatedAt, int64(0))
	tx2 := txs[1]
	assert.Equal(tx2.ID, int64(1))
	assert.Equal(tx2.Amount, 50.0)
	assert.Equal(tx2.Kind, "WITHDRAW")
	assert.Greater(tx2.CreatedAt, int64(0))

	repo.AssertExpectations(t)
}
