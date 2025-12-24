package security

import (
	"github.com/gin-gonic/gin"
	"turtle/security/loginPenetration"
	"turtle/security/scamDetection"
)

func InitSecurityApi(r *gin.Engine) {
	loginPenetration.InitLoginPenetration(r)
	scamDetection.InitScamDetection(r)
}
