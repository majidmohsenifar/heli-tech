package transaction_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
	)
	_, err = transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Equal(err, transaction.ErrOngoingRequest)
	lock.Release(ctx)
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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
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
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(120), Valid: true}}, nil)

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
	)
	res, err := transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Nil(err)
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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
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
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
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
	).Once().Return(repository.UserBalance{ID: 1, Amount: pgtype.Numeric{Int: big.NewInt(120), Valid: true}}, nil)

	s := miniredis.RunT(t)
	redisClient, err := core.NewRedisClient(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(err)
	redisLocker := core.NewRedisLocker(redisClient)
	logger := logger.NewLogger()
	transactionService := transaction.NewService(
		dbMock,
		repo,
		redisLocker,
		logger,
	)
	res, err := transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: 1,
		Amount: 100,
	})
	assert.Nil(err)
	assert.Equal(res.Amount, 100.0)
	assert.Equal(res.NewBalance, 120.0)
	assert.Greater(res.CreatedAt, int64(0))
	repo.AssertExpectations(t)
}
