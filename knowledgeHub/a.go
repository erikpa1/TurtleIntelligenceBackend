package knowledgeHub

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/artifact"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/domains"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/knowledge"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/node"
	"github.com/gin-gonic/gin"
)

func InitKnowledgeHubApi(r *gin.Engine) {
	artifact.InitArtifactApi(r)
	domains.InitDomainApi(r)
	knowledge.InitKnowledgeApi(r)
	node.InitKnowledgeNodeApi(r)
}
