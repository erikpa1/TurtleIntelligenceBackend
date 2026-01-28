package pipeline

import (
	"turtle/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

type GinPipeline struct {
	Ctx   *gin.Context
	Error error
}

func NewGinPipeline(ctx *gin.Context) *GinPipeline {
	return &GinPipeline{
		Ctx:   ctx,
		Error: nil,
	}
}

func (self *GinPipeline) SetError(err error) bool {
	self.Error = err

	if err != nil {
		lgr.ErrorStack(err.Error())
	}

	return err == nil
}

func (self *GinPipeline) WasError() bool {
	return self.Error != nil
}

func (self *GinPipeline) Return() {
	if self.WasError() {
		tools.Auto500(self.Ctx, self.Error)
	} else {
		tools.AutoReturn(self.Ctx, nil)
	}
}

func (self *GinPipeline) ShouldBindJSON(obj any) bool {
	return self.SetError(self.Ctx.ShouldBindJSON(obj))
}
