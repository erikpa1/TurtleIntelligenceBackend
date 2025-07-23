package fn

type FnVarType int8

const (
	FN_VT_BOOL               FnVarType = 0
	FN_VT_FLOAT              FnVarType = 1
	FN_VT_INT8               FnVarType = 2
	FN_VT_INT16              FnVarType = 3
	FN_VT_INT32              FnVarType = 4
	FN_VT_INT64              FnVarType = 5
	FN_VT_UINT8              FnVarType = 6
	FN_VT_UINT16             FnVarType = 7
	FN_VT_UINT32             FnVarType = 8
	FN_VT_UINT64             FnVarType = 9
	FN_VT_FLOAT32            FnVarType = 10
	FN_VT_FLOAT64            FnVarType = 11
	FN_VT_STRING             FnVarType = 12
	FN_VT_ARRAY              FnVarType = 13
	FN_VT_ARRAY_BYTE_ARRAY   FnVarType = 14
	FN_VT_ARRAY_INT8_ARRAY   FnVarType = 15
	FN_VT_ARRAY_INT16_ARRAY  FnVarType = 16
	FN_VT_ARRAY_INT32_ARRAY  FnVarType = 17
	FN_VT_ARRAY_INT64_ARRAY  FnVarType = 18
	FN_VT_ARRAY_UINT8_ARRAY  FnVarType = 19
	FN_VT_ARRAY_UINT16_ARRAY FnVarType = 20
)
