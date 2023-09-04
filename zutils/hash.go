package zutils

const (
	Prime   = 16777619
	HashVal = 2166136261
)

type IHash interface {
	Sum(string) uint32
}

type Fnv32Hash struct{}

func DefaultHash() IHash {
	return &Fnv32Hash{}
}

// fnv32 algorithm
func (f *Fnv32Hash) Sum(key string) uint32 {
	hashVal := uint32(HashVal)
	prime := uint32(Prime)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hashVal *= prime
		hashVal ^= uint32(key[i])
	}
	return hashVal
}
