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

func IsDepartmentAllowed(ctx *gin.Context, target_department_id uint16) bool {
	if os.Getenv("AUTH") == "on" {
		session := sessions.Default(ctx)
		user := session.Get("department_user")

		log.Print(session.Get("department_user"))
		log.Print(reflect.TypeOf(session.Get("department_user")))

		if !(target_department_id == 0 || *(user.(*uint16)) == target_department_id) {
			ctx.String(http.StatusForbidden, "your department is not allowed to edit, update, or add data to other departments for this specific operation")
			return false
		}
	}

	return true
}
