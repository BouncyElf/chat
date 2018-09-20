package common

import uuid "github.com/satori/go.uuid"

func NewUUID() string {
	res, _ := uuid.NewV4()
	return res.String()
}

func NewSnowFlake() string {
	return SnowFlakeNode.Generate().String()
}

func NewSHAUUID(uid, name string) string {
	v := uuid.FromStringOrNil(uid)
	return uuid.NewV5(v, name).String()
}
