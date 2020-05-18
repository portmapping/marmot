package lurker

import (
	"encoding/json"
	"errors"
	"net"
)

// HandshakeStatus ...
type HandshakeStatus int

// Version ...
type Version string

// RequestType ...
type RequestType int

// RequestTypePing ...
const RequestTypePing RequestType = 0x01

// RequestTypeConnect ...
const RequestTypeConnect RequestType = 0x02

// RequestTypeAdapter ...
const RequestTypeAdapter RequestType = 0x03

// Handshake ...
type Handshake struct {
	Type RequestType `json:"type"`
}

// HandshakeAble ...
type HandshakeAble interface {
	Ping() error
	Connect() error
	Adapter() error
}

// HandshakeRequest ...
type HandshakeRequest struct {
	ProtocolVersion Version `json:"protocol_version"`
	Data            []byte  `json:"data"`
}

// HandshakeResponse ...
type HandshakeResponse struct {
	Status HandshakeStatus `json:"status"`
	Data   []byte          `json:"data"`
}

// Service ...
type Service struct {
	ID          string `json:"id"`
	Addr        []Addr `json:"addr"`
	ISP         net.IP `json:"isp"`
	Local       net.IP `json:"local"`
	PortUDP     int    `json:"port_udp"`
	PortHole    int    `json:"port_hole"`
	PortTCP     int    `json:"port_tcp"`
	KeepConnect bool   `json:"keep_connect"`
}

// ParseHandshake ...
func ParseHandshake(data []byte) (Handshake, error) {
	var h Handshake
	err := json.Unmarshal(data, &h)
	if err != nil {
		return Handshake{}, err
	}
	return h, nil
}

// DecodeHandshakeRequest ...
func DecodeHandshakeRequest(data []byte) (Service, error) {
	var r HandshakeRequest
	err := json.Unmarshal(data, &r)
	if err != nil {
		return Service{}, err
	}
	return decodeHandshakeRequestV1(&r)
}

func decodeHandshakeRequestV1(request *HandshakeRequest) (Service, error) {
	var s Service
	err := json.Unmarshal(request.Data, &s)
	if err != nil {
		return Service{}, err
	}
	return s, nil
}

// EncodeHandshakeRequest ...
func EncodeHandshakeRequest(service Service) ([]byte, error) {
	return encodeHandshakeRequestV1(&HandshakeRequest{
		ProtocolVersion: "v0.0.1",
		Data:            service.JSON(),
	})
}
func encodeHandshakeRequestV1(request *HandshakeRequest) ([]byte, error) {
	return json.Marshal(request)
}

// EncodeHandshakeResponse ...
func EncodeHandshakeResponse(ver Version) ([]byte, error) {
	var r HandshakeResponse
	return encodeHandshakeResponseV1(&r)
}

func encodeHandshakeResponseV1(r *HandshakeResponse) ([]byte, error) {
	return json.Marshal(r)
}

// JSON ...
func (h Handshake) JSON() []byte {
	marshal, err := json.Marshal(h)
	if err != nil {
		return nil
	}
	return marshal
}

// ProcessHandshake ...
func (h Handshake) ProcessHandshake(able HandshakeAble) (er error) {
	switch h.Type {
	case RequestTypePing:
		return able.Ping()
	case RequestTypeConnect:
		return able.Connect()
	case RequestTypeAdapter:
		return able.Adapter()
	default:
	}
	return errors.New("wrong type")
}
