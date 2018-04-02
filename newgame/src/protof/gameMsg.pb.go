// Code generated by protoc-gen-go.
// source: gameMsg.proto
// DO NOT EDIT!

/*
Package protof is a generated protocol buffer package.

It is generated from these files:
	gameMsg.proto

It has these top-level messages:
	Message1
*/
package protof

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message1 struct {
	Ping             *Message1_PING            `protobuf:"bytes,100,opt,name=ping" json:"ping,omitempty"`
	CsLoginMessage   *Message1_CS_LoginMessage `protobuf:"bytes,101,opt,name=cs_login_message" json:"cs_login_message,omitempty"`
	ScLoginMessage   *Message1_SC_LoginMessage `protobuf:"bytes,102,opt,name=sc_login_message" json:"sc_login_message,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *Message1) Reset()                    { *m = Message1{} }
func (m *Message1) String() string            { return proto.CompactTextString(m) }
func (*Message1) ProtoMessage()               {}
func (*Message1) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message1) GetPing() *Message1_PING {
	if m != nil {
		return m.Ping
	}
	return nil
}

func (m *Message1) GetCsLoginMessage() *Message1_CS_LoginMessage {
	if m != nil {
		return m.CsLoginMessage
	}
	return nil
}

func (m *Message1) GetScLoginMessage() *Message1_SC_LoginMessage {
	if m != nil {
		return m.ScLoginMessage
	}
	return nil
}

type Message1_PING struct {
	Time             *float32 `protobuf:"fixed32,1,req,name=time" json:"time,omitempty"`
	SerTime          *int32   `protobuf:"varint,2,opt,name=serTime" json:"serTime,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Message1_PING) Reset()                    { *m = Message1_PING{} }
func (m *Message1_PING) String() string            { return proto.CompactTextString(m) }
func (*Message1_PING) ProtoMessage()               {}
func (*Message1_PING) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Message1_PING) GetTime() float32 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Message1_PING) GetSerTime() int32 {
	if m != nil && m.SerTime != nil {
		return *m.SerTime
	}
	return 0
}

type Message1_CS_LoginMessage struct {
	Id               *int32  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Name             *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	Opt              *int32  `protobuf:"varint,3,req,name=opt" json:"opt,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Message1_CS_LoginMessage) Reset()                    { *m = Message1_CS_LoginMessage{} }
func (m *Message1_CS_LoginMessage) String() string            { return proto.CompactTextString(m) }
func (*Message1_CS_LoginMessage) ProtoMessage()               {}
func (*Message1_CS_LoginMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Message1_CS_LoginMessage) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Message1_CS_LoginMessage) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Message1_CS_LoginMessage) GetOpt() int32 {
	if m != nil && m.Opt != nil {
		return *m.Opt
	}
	return 0
}

type Message1_SC_LoginMessage struct {
	Code             *int32  `protobuf:"varint,1,req,name=code" json:"code,omitempty"`
	Name             *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	LoginTime        *int32  `protobuf:"varint,3,req,name=loginTime" json:"loginTime,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Message1_SC_LoginMessage) Reset()                    { *m = Message1_SC_LoginMessage{} }
func (m *Message1_SC_LoginMessage) String() string            { return proto.CompactTextString(m) }
func (*Message1_SC_LoginMessage) ProtoMessage()               {}
func (*Message1_SC_LoginMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 2} }

func (m *Message1_SC_LoginMessage) GetCode() int32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *Message1_SC_LoginMessage) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Message1_SC_LoginMessage) GetLoginTime() int32 {
	if m != nil && m.LoginTime != nil {
		return *m.LoginTime
	}
	return 0
}

func init() {
	proto.RegisterType((*Message1)(nil), "protof.Message1")
	proto.RegisterType((*Message1_PING)(nil), "protof.Message1.PING")
	proto.RegisterType((*Message1_CS_LoginMessage)(nil), "protof.Message1.CS_LoginMessage")
	proto.RegisterType((*Message1_SC_LoginMessage)(nil), "protof.Message1.SC_LoginMessage")
}

func init() { proto.RegisterFile("gameMsg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 220 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x4f, 0xcc, 0x4d,
	0xf5, 0x2d, 0x4e, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0x69, 0x4a, 0xa7,
	0x98, 0xb8, 0x38, 0x7c, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x0d, 0x85, 0x94, 0xb9, 0x58, 0x0a,
	0x32, 0xf3, 0xd2, 0x25, 0x52, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0x44, 0x21, 0x4a, 0xd3, 0xf4, 0x60,
	0xf2, 0x7a, 0x01, 0x9e, 0x7e, 0xee, 0x42, 0x56, 0x5c, 0x02, 0xc9, 0xc5, 0xf1, 0x39, 0xf9, 0xe9,
	0x99, 0x79, 0xf1, 0xb9, 0x10, 0x19, 0x89, 0x54, 0xb0, 0x06, 0x05, 0x0c, 0x0d, 0xce, 0xc1, 0xf1,
	0x3e, 0x20, 0x85, 0x50, 0x01, 0x90, 0xde, 0xe2, 0x64, 0x34, 0xbd, 0x69, 0x38, 0xf4, 0x06, 0x3b,
	0xa3, 0xe8, 0x95, 0x52, 0xe5, 0x62, 0x01, 0xdb, 0xcf, 0xc3, 0xc5, 0x52, 0x92, 0x99, 0x9b, 0x2a,
	0xc1, 0xa8, 0xc0, 0xa4, 0xc1, 0x24, 0xc4, 0xcf, 0xc5, 0x5e, 0x9c, 0x5a, 0x14, 0x02, 0x12, 0x60,
	0x52, 0x60, 0xd4, 0x60, 0x95, 0xb2, 0xe0, 0xe2, 0x47, 0xb7, 0x95, 0x8b, 0x8b, 0x29, 0x33, 0x05,
	0xac, 0x9e, 0x15, 0xa4, 0x3b, 0x2f, 0x11, 0xac, 0x98, 0x49, 0x83, 0x53, 0x88, 0x9b, 0x8b, 0x39,
	0xbf, 0xa0, 0x44, 0x82, 0x19, 0x24, 0x25, 0xe5, 0xc0, 0xc5, 0x8f, 0x66, 0x27, 0x48, 0x75, 0x72,
	0x7e, 0x4a, 0x2a, 0x56, 0xbd, 0x82, 0x5c, 0x9c, 0x60, 0x8f, 0x80, 0xed, 0x06, 0x9b, 0x00, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xc9, 0x12, 0x8b, 0xe4, 0x64, 0x01, 0x00, 0x00,
}
