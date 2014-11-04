package main

func Noise(x, y int) float64 {
	n := x + y*57
	n = (n << 13) ^ n
	return (1.0 - float64((n*(n*n*15731.0+789221.0)+1376312589.0)&0x7fffffff)/1073741824.0)
}

func SmoothedNoise(x, y int) float64 {
	corners := float64(( Noise(x-1, y-1)+Noise(x+1, y-1)+Noise(x-1, y+1)+Noise(x+1, y+1) )) / 16
	sides := float64(( Noise(x-1, y)  +Noise(x+1, y)  +Noise(x, y-1)  +Noise(x, y+1) )) /  8
	center := float64(Noise(x, y)) / 4
	return corners + sides + center
}

/*
func LinearInterpolate(a, b int, x float) []float {

}

func CosineInterpolate(a, b int, x float) []float {

}*/
