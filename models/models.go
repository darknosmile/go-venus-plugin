package models

import ()

type PacketHead struct {
	PacketLength    int32
	PacketVersion   int32
	CommondType     []byte
	SerializeType   byte
	Flags           byte
	ClientId        int32
	ClientRequestId int32
	DateLength      int32
}

type PacketHandShake struct {
	Capabilities        int32
	SupportAuthenMethod int32
	Challenge           string
	Version             string
}

type FindNameData struct {
	UserName string `json:"userName"`
}
