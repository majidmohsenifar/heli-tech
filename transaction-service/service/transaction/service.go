package transaction

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/helper"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"

	"github.com/bsm/redislock"
)

var (
	ErrOngoingRequest      = errors.New("there is already an ongoing transaction request")
	ErrInsufficientBalance = errors.New("the balance is not enough for withdraw")
)

type Service struct {
	db           core.PgxInterface
	repo         repository.Querier
	redisLocker  *redislock.Client
	logger       *slog.Logger
	eventManager TransactionEventManager
}

type WithdrawParams struct {
	UserID int64
	Amount float64
}

type TransactionDetail struct {
	ID         int64
	CreatedAt  int64
	Amount     float64
	NewBalance float64
}

type DepositParams struct {
	UserID int64
	Amount float64
}

func (s *Service) Withdraw(ctx context.Context, params WithdrawParams) (TransactionDetail, error) {
	//lock concurrent process
	//we lock for max 30s for each user, because it may have called the deposit and withdraw, and
	//by locking we try to only let each user to do one thing at any given time
	lock, err := s.redisLocker.Obtain(ctx, fmt.Sprintf("transaction:%d", params.UserID), time.Duration(30*time.Second), nil)
	if err != nil && !errors.Is(err, redislock.ErrNotObtained) {
		s.logger.Error("cannot obtain lock", err)
		return TransactionDetail{}, fmt.Errorf("cannot lock order for update")
	}
	if errors.Is(err, redislock.ErrNotObtained) {
		return TransactionDetail{}, ErrOngoingRequest
	}
	defer lock.Release(ctx)

	//first we check if user has balance for the withdraw amount
	ub, err := s.repo.GetUserBalanceByUserID(ctx, s.db, params.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return TransactionDetail{}, fmt.Errorf("cannot check the user balance")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return TransactionDetail{}, ErrInsufficientBalance
	}
	ubAmount, err := helper.PGNumericToFloat64(ub.Amount)
	if err != nil {
		return TransactionDetail{}, fmt.Errorf("cannot check user balance for now")
	}
	if ubAmount < params.Amount {
		return TransactionDetail{}, ErrInsufficientBalance
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("cannot start db transaction", err)
		return TransactionDetail{}, fmt.Errorf("something went wrong")
	}
	amountNumeric, err := helper.Float64ToPGNumeric(params.Amount)
	if err != nil {
		s.logger.Error("cannot convert amount to pgtype numeric", err)
		return TransactionDetail{}, fmt.Errorf("something went wrong with amount")
	}
	tx, err := s.repo.CreateTransaction(ctx, dbTx, repository.CreateTransactionParams{
		UserID: params.UserID,
		Kind:   repository.KindWITHDRAW,
		Amount: amountNumeric,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot create transaction", err)
		return TransactionDetail{}, fmt.Errorf("cannot create transaction")
	}

	balance, err := s.repo.CreateUserBalanceOrDecreaseAmount(ctx, dbTx, repository.CreateUserBalanceOrDecreaseAmountParams{
		UserID: params.UserID,
		Amount: amountNumeric,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot create userBalanceOrUpdateAmount", err)
		return TransactionDetail{}, fmt.Errorf("cannot update user balance")
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot commit transaction", err)
		return TransactionDetail{}, errors.New("cannot store in db")
	}
	newBalance, err := helper.PGNumericToFloat64(balance.Amount)
	if err != nil {
		s.logger.Error("cannot convert PGNumericToFloat64", err)
		//we do not return here as it does not affect our logic
	}

	s.eventManager.PublishTransactionCreatedEvent(ctx, TransactionCreatedEventParams{
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Kind:          string(tx.Kind),
		Amount:        params.Amount,
		Balance:       newBalance,
		CreatedAt:     tx.CreatedAt.Time.Unix(),
	})
	return TransactionDetail{
		ID:         tx.ID,
		CreatedAt:  time.Now().Unix(),
		Amount:     params.Amount,
		NewBalance: newBalance,
	}, nil

}

func (s *Service) Deposit(ctx context.Context, params DepositParams) (TransactionDetail, error) {
	//lock concurrent process
	//we lock for max 30s for each user, because it may have called the deposit and withdraw, and
	//by locking we try to only let each user to do one thing at any given time
	lock, err := s.redisLocker.Obtain(ctx, fmt.Sprintf("transaction:%d", params.UserID), time.Duration(30*time.Second), nil)
	if err != nil && !errors.Is(err, redislock.ErrNotObtained) {
		s.logger.Error("cannot obtain lock", err)
		return TransactionDetail{}, fmt.Errorf("cannot lock order for update")
	}
	if errors.Is(err, redislock.ErrNotObtained) {
		return TransactionDetail{}, ErrOngoingRequest
	}
	defer lock.Release(ctx)

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("cannot start db transaction", err)
		return TransactionDetail{}, fmt.Errorf("something went wrong")
	}
	amountNumeric, err := helper.Float64ToPGNumeric(params.Amount)
	if err != nil {
		s.logger.Error("cannot convert amount to pgtype numeric", err)
		return TransactionDetail{}, fmt.Errorf("something went wrong with amount")
	}
	tx, err := s.repo.CreateTransaction(ctx, dbTx, repository.CreateTransactionParams{
		UserID: params.UserID,
		Kind:   repository.KindDEPOSIT,
		Amount: amountNumeric,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot create transaction", err)
		return TransactionDetail{}, fmt.Errorf("cannot create transaction")
	}
	balance, err := s.repo.CreateUserBalanceOrIncreaseAmount(ctx, dbTx, repository.CreateUserBalanceOrIncreaseAmountParams{
		UserID: params.UserID,
		Amount: amountNumeric,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot create userBalanceOrUpdateAmount", err)
		return TransactionDetail{}, fmt.Errorf("cannot update user balance")
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot commit transaction", err)
		return TransactionDetail{}, errors.New("cannot store in db")
	}
	newBalance, err := helper.PGNumericToFloat64(balance.Amount)
	if err != nil {
		s.logger.Error("cannot convert PGNumericToFloat64", err)
		//we do not return here as it does not affect our logic
	}
	s.eventManager.PublishTransactionCreatedEvent(ctx, TransactionCreatedEventParams{
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Kind:          string(tx.Kind),
		Amount:        params.Amount,
		Balance:       newBalance,
		CreatedAt:     tx.CreatedAt.Time.Unix(),
	})
	return TransactionDetail{
		ID:         tx.ID,
		CreatedAt:  time.Now().Unix(),
		Amount:     params.Amount,
		NewBalance: newBalance,
	}, nil
}

func NewService(
	db core.PgxInterface,
	repo repository.Querier,
	redisLocker *redislock.Client,
	logger *slog.Logger,
	eventManager TransactionEventManager,

) *Service {
	return &Service{
		db:           db,
		repo:         repo,
		redisLocker:  redisLocker,
		logger:       logger,
		eventManager: eventManager,
	}
}
