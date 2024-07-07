package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api/middleware"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/transaction"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"
)

type TransactionHandler struct {
	transactionService *transaction.Service
	validate           *validator.Validate
}

// This endpoint allows user to withdraw
//
//	@Summary		withdraw
//	@Description	allows user to withdraw
//	@Tags			Transaction
//	@ID				Withdraw
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			params	body		transaction.WithdrawParams	false	"Withdraw-Params"
//
// @Success		 200	{object}	ResponseSuccess{data=transaction.TransactionDetail}
//
//	@Failure		400		{object}	ResponseFailure
//	@Failure		403		{object}	ResponseFailure
//	@Failure		422		{object}	ResponseFailure
//	@Failure		500		{object}	ResponseFailure
//	@Router			/api/v1/transactions/withdraw [post]
func (h *TransactionHandler) Withdraw(c *gin.Context) {
	params := transaction.WithdrawParams{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}

	userDataAny, exist := c.Get(middleware.UserDataKey)
	if !exist {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token 1")
		return
	}
	userData, ok := userDataAny.(user.UserData)
	if !ok {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token 2")
		return
	}

	params.UserID = userData.ID
	res, err := h.transactionService.Withdraw(c, params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "successfully withdrawed")
}

// This endpoint allows user to deposit
//
//	@Summary		deposit
//	@Description	allows user to deposit
//	@Tags			Transaction
//	@ID				Deposit
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			params	body		transaction.DepositParams	false	"Deposit-Params"
//
// @Success		 200	{object}	ResponseSuccess{data=transaction.TransactionDetail}
//
//	@Failure		400		{object}	ResponseFailure
//	@Failure		403		{object}	ResponseFailure
//	@Failure		422		{object}	ResponseFailure
//	@Failure		500		{object}	ResponseFailure
//	@Router			/api/v1/transactions/deposit [post]
func (h *TransactionHandler) Deposit(c *gin.Context) {
	params := transaction.DepositParams{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}
	userDataAny, exist := c.Get(middleware.UserDataKey)
	if !exist {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token")
		return
	}
	userData, ok := userDataAny.(user.UserData)
	if !ok {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token")
		return
	}

	params.UserID = userData.ID
	res, err := h.transactionService.Deposit(c, params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "successfully deposited")
}

// This endpoint allows user to see his/her transactions
//
//	@Summary		transactions list
//	@Description	allows user to see his/her transactions
//	@Tags			Transaction
//	@ID				UserTransactions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//
// @Param		    page query int false "Page"
// @Param		    pageSize query int false "PageSize"
//
// @Success		 200	{object}	ResponseSuccess{data=[]transaction.Transaction}
//
//	@Failure		403		{object}	ResponseFailure
//	@Failure		500		{object}	ResponseFailure
//	@Router			/api/v1/transactions [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	params := transaction.GetUserTransactionsParams{}
	err := c.ShouldBindQuery(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}

	userDataAny, exist := c.Get(middleware.UserDataKey)
	if !exist {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token")
		return
	}
	userData, ok := userDataAny.(user.UserData)
	if !ok {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusUnauthorized,
			"cannot get user data from token")
		return
	}

	params.UserID = userData.ID
	res, err := h.transactionService.GetUserTransactions(c, params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "successfully fetched")
}

func NewTransactionHandler(
	transactionService *transaction.Service,
	validate *validator.Validate,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validate:           validate,
	}
}
