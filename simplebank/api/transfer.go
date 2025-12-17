package api

import (
	"fmt"
	"net/http"

	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id"  binding:"required"`
	Amount        int64  `json:"amount"  binding:"required,min=1"`
	Currency      string `json:"currency"  binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorsResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return false
	}

	// Check account currency
	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %v vs %v", accountID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return false
	}

	return true
}

func (server *Server) setupTransferRoutes() {
	server.router.POST("/transfer", server.CreateTransfer)
}
