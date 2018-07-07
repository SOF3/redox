// Code generated by protoc-gen-go. DO NOT EDIT.
// source: player.proto

package protocol

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

type Player struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Player) Reset()         { *m = Player{} }
func (m *Player) String() string { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()    {}
func (*Player) Descriptor() ([]byte, []int) {
	return fileDescriptor_player_813aff5d3cc11071, []int{0}
}
func (m *Player) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Player.Unmarshal(m, b)
}
func (m *Player) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Player.Marshal(b, m, deterministic)
}
func (dst *Player) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Player.Merge(dst, src)
}
func (m *Player) XXX_Size() int {
	return xxx_messageInfo_Player.Size(m)
}
func (m *Player) XXX_DiscardUnknown() {
	xxx_messageInfo_Player.DiscardUnknown(m)
}

var xxx_messageInfo_Player proto.InternalMessageInfo

func (m *Player) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Player) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type PlayerLoginEvent struct {
	Player               *Player  `protobuf:"bytes,4,opt,name=player,proto3" json:"player,omitempty"`
	Before               uint64   `protobuf:"varint,8,opt,name=before,proto3" json:"before,omitempty"`
	After                uint64   `protobuf:"varint,9,opt,name=after,proto3" json:"after,omitempty"`
	NodeId               string   `protobuf:"bytes,15,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerLoginEvent) Reset()         { *m = PlayerLoginEvent{} }
func (m *PlayerLoginEvent) String() string { return proto.CompactTextString(m) }
func (*PlayerLoginEvent) ProtoMessage()    {}
func (*PlayerLoginEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_player_813aff5d3cc11071, []int{1}
}
func (m *PlayerLoginEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerLoginEvent.Unmarshal(m, b)
}
func (m *PlayerLoginEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerLoginEvent.Marshal(b, m, deterministic)
}
func (dst *PlayerLoginEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerLoginEvent.Merge(dst, src)
}
func (m *PlayerLoginEvent) XXX_Size() int {
	return xxx_messageInfo_PlayerLoginEvent.Size(m)
}
func (m *PlayerLoginEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerLoginEvent.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerLoginEvent proto.InternalMessageInfo

func (m *PlayerLoginEvent) GetPlayer() *Player {
	if m != nil {
		return m.Player
	}
	return nil
}

func (m *PlayerLoginEvent) GetBefore() uint64 {
	if m != nil {
		return m.Before
	}
	return 0
}

func (m *PlayerLoginEvent) GetAfter() uint64 {
	if m != nil {
		return m.After
	}
	return 0
}

func (m *PlayerLoginEvent) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func init() {
	proto.RegisterType((*Player)(nil), "redox.player.Player")
	proto.RegisterType((*PlayerLoginEvent)(nil), "redox.player.PlayerLoginEvent")
}

func init() { proto.RegisterFile("player.proto", fileDescriptor_player_813aff5d3cc11071) }

var fileDescriptor_player_813aff5d3cc11071 = []byte{
	// 186 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xc8, 0x49, 0xac,
	0x4c, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x29, 0x4a, 0x4d, 0xc9, 0xaf, 0xd0,
	0x83, 0x88, 0x29, 0xe9, 0x70, 0xb1, 0x05, 0x80, 0x59, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9,
	0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x10, 0x1f, 0x17, 0x53, 0x66, 0x81,
	0x04, 0x13, 0x58, 0x84, 0x29, 0xb3, 0x40, 0xa9, 0x8d, 0x91, 0x4b, 0x00, 0xa2, 0xdc, 0x27, 0x3f,
	0x3d, 0x33, 0xcf, 0xb5, 0x2c, 0x35, 0xaf, 0x44, 0x48, 0x87, 0x8b, 0x0d, 0x62, 0x98, 0x04, 0x8b,
	0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0xb2, 0x0d, 0x7a, 0x10, 0xf5, 0x41, 0x50, 0x35, 0x42,
	0x62, 0x5c, 0x6c, 0x49, 0xa9, 0x69, 0xf9, 0x45, 0xa9, 0x12, 0x1c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41,
	0x50, 0x9e, 0x90, 0x08, 0x17, 0x6b, 0x62, 0x5a, 0x49, 0x6a, 0x91, 0x04, 0x27, 0x58, 0x18, 0xc2,
	0x01, 0xa9, 0xce, 0xcb, 0x4f, 0x49, 0xf5, 0x4c, 0x91, 0xe0, 0x07, 0x3b, 0x02, 0xca, 0x73, 0xe2,
	0x8a, 0xe2, 0x00, 0xfb, 0x26, 0x39, 0x3f, 0x27, 0x89, 0x0d, 0xcc, 0x32, 0x06, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x3f, 0xbc, 0x68, 0xf9, 0xe7, 0x00, 0x00, 0x00,
}
