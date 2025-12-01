package scamDetection

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

/*
TODO
Toto by malo zobrat mena ludi z organizacie pravidla a podobne
a pomocou AIcka vyhodnotit scame a navigovat uzivatela

AI by malo byt naviazane na knowlege base a tahat data z neho
*/

func InitScamDetection(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "security", "scam_detection")
}
