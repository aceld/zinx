package ziface

type IDecoder interface {
	Interceptor
	GetLengthField() *LengthField
}
