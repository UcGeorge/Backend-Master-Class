package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	AuthorizationPayloadKey = "authorization_payload"
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
)

func AddTestAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	accessToken, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorization := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(AuthorizationHeaderKey, authorization)
}

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorsResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorsResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorsResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorsResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func errorsResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
