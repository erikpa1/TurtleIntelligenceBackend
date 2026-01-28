package types

//==============Output==============

type StringOutput struct {
	Data string
}

func (self *StringOutput) SetData(data string) {
	self.Data = data
}

//==============Input==============

type StringInput struct {
	SisterOutput *StringOutput
}

func (self *StringInput) GetData() string {
	return self.SisterOutput.Data
}
