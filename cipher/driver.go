package cipher

type Driver interface {
	Encode(b []byte, n int)
	Decode(b []byte, n int)
}
