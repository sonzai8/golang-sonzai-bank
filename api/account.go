package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var input createAccountRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    input.Owner,
		Balance:  0,
		Currency: input.Currency,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			//fmt.Println("check error code: ", pgErr.Code)
			switch pgErr.Code {
			case "23503", "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	ctx.JSON(http.StatusOK, gin.H{"account": account})
}

func (server *Server) UpdateAccount(ctx *gin.Context) {}

type DeleteAccountRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) DeleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteAccount(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted": true})
}

type ListAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var paging ListAccountRequest
	if err := ctx.ShouldBindQuery(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  paging.PageSize,
		Offset: (paging.PageId - 1) * paging.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type getAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var params getAccountRequest
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, params.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//account = db.Account{}
	ctx.JSON(http.StatusOK, account)
}
