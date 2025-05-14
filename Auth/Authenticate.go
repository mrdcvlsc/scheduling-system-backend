package Auth

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IsAuthSuccess(ctx *gin.Context) bool {
	if os.Getenv("AUTH") == "on" {
		session := sessions.Default(ctx)
		user := session.Get("department_user")

		log.Print(session.Get("department_user"))
		log.Print(reflect.TypeOf(session.Get("department_user")))

		if user == nil {
			ctx.String(http.StatusUnauthorized, "you're not allowed to access that, please login first")
			return false
		}
	}

	return true
}
