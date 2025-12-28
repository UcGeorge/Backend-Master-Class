package api

import (
	"net/http"

	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/middleware"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505", "23503":
				ctx.JSON(http.StatusForbidden, errorsResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorsResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorsResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsForUserParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		Username: authPayload.Username,
	}

	accounts, err := server.store.ListAccountsForUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) setupAccountRoutes(router gin.IRoutes) {
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
}
