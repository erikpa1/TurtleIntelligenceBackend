package security

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/security/loginPenetration"
	"github.com/erikpa1/TurtleIntelligenceBackend/security/scamDetection"
	"github.com/gin-gonic/gin"
)

func InitSecurityApi(r *gin.Engine) {
	loginPenetration.InitLoginPenetration(r)
	scamDetection.InitScamDetection(r)
}
