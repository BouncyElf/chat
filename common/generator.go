package common

func NewSnowFlake() string {
	return SnowFlakeNode.Generate().String()
}
