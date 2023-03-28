package ziface

type IDecoder interface {
	IInterceptor
	GetLengthField() *LengthField
}
