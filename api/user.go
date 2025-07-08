package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email" `
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var input createUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       input.Username,
		HashedPassword: hashedPassword,
		FullName:       input.FullName,
		Email:          input.Email,
	}
	account, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			//fmt.Println("check error code: ", pgErr.Code)
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	rsp := createUserResponse{
		Username:          account.Username,
		FullName:          account.FullName,
		Email:             account.Email,
		CreatedAt:         account.CreatedAt,
		PasswordChangedAt: account.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
