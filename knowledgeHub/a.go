package knowledgeHub

import (
	"github.com/gin-gonic/gin"
	"turtle/knowledgeHub/artifact"
	"turtle/knowledgeHub/domains"
	"turtle/knowledgeHub/knowledge"
	"turtle/knowledgeHub/node"
)

func InitKnowledgeHubApi(r *gin.Engine) {
	artifact.InitArtifactApi(r)
	domains.InitDomainApi(r)
	knowledge.InitKnowledgeApi(r)
	node.InitKnowledgeNodeApi(r)
}
