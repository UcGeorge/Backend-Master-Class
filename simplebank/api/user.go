package api

import (
	"net/http"
	"time"

	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required,alphanum"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorsResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorsResponse(err))
		return
	}

	rsp := CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) setupUserRoutes() {
	server.router.POST("/users", server.createUser)
	// server.router.GET("/accounts/:id", server.getUser)
	// server.router.GET("/accounts", server.listUsers)
}
