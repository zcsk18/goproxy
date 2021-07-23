package cipher


type Caesar struct {
	dis byte
}

func (this *Caesar) SetDis(dis byte) {
	this.dis = dis
}

func (this *Caesar) Encode(b []byte, n int) {
	const MaxValue = 255

	for i:=0; i<n ;i++  {
		if (b[i] + this.dis > MaxValue) {
			b[i] = MaxValue - b[i] + (this.dis - 1)
		} else {
			b[i] += this.dis
		}
	}
}

func (this *Caesar) Decode(b []byte, n int) {
	const MaxValue = 255

	for i:=0; i<n ;i++  {
		if (b[i] - this.dis < 0) {
			b[i] = MaxValue - (this.dis-1) - b[i]
		} else {
			b[i] -= this.dis
		}
	}
}