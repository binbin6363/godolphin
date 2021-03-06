// Code generated by protoc-gen-go.
// source: base.proto
// DO NOT EDIT!

/*
Package com_kook_proto is a generated protocol buffer package.

It is generated from these files:
	base.proto
	message.proto

It has these top-level messages:
	BaseRequest
	BaseResult
	MsgInfo
	MsgSessionInfo
	SingleMsgRequest
	SingleMsgResult
	SingleMsgNotice
	SyncSingleMsgRequest
	SyncSingleMsgResult
	GetMsgSessionInfoRequest
	GetMsgSessionInfoResult
	DelMsgSessionInfoRequest
	DelMsgSessionInfoResult
	UpdateMsgReadStatusRequest
	UpdateMsgReadStatusResult
	GetMsgReadStatusRequest
	GetMsgReadStatusResult
*/
package com_kook_proto

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

// 基础的请求数据
type BaseRequest struct {
	TransId string `protobuf:"bytes,1,opt,name=trans_id,json=transId" json:"trans_id,omitempty"`
	Cid     uint32 `protobuf:"varint,2,opt,name=cid" json:"cid,omitempty"`
	Uid     uint64 `protobuf:"varint,3,opt,name=uid" json:"uid,omitempty"`
}

func (m *BaseRequest) Reset()                    { *m = BaseRequest{} }
func (m *BaseRequest) String() string            { return proto.CompactTextString(m) }
func (*BaseRequest) ProtoMessage()               {}
func (*BaseRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// 基础的返回结果数据
type BaseResult struct {
	TransId string `protobuf:"bytes,1,opt,name=trans_id,json=transId" json:"trans_id,omitempty"`
	Cid     uint32 `protobuf:"varint,2,opt,name=cid" json:"cid,omitempty"`
	Uid     uint64 `protobuf:"varint,3,opt,name=uid" json:"uid,omitempty"`
	RetCode uint32 `protobuf:"varint,4,opt,name=ret_code,json=retCode" json:"ret_code,omitempty"`
	RetMsg  string `protobuf:"bytes,5,opt,name=ret_msg,json=retMsg" json:"ret_msg,omitempty"`
}

func (m *BaseResult) Reset()                    { *m = BaseResult{} }
func (m *BaseResult) String() string            { return proto.CompactTextString(m) }
func (*BaseResult) ProtoMessage()               {}
func (*BaseResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// 消息数据结构定义
type MsgInfo struct {
	SrvMsgId uint64   `protobuf:"varint,1,opt,name=srv_msg_id,json=srvMsgId" json:"srv_msg_id,omitempty"`
	CliMsgId uint64   `protobuf:"varint,2,opt,name=cli_msg_id,json=cliMsgId" json:"cli_msg_id,omitempty"`
	FromUid  uint64   `protobuf:"varint,3,opt,name=from_uid,json=fromUid" json:"from_uid,omitempty"`
	ToUid    uint64   `protobuf:"varint,4,opt,name=to_uid,json=toUid" json:"to_uid,omitempty"`
	GroupId  uint32   `protobuf:"varint,5,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	MsgType  uint32   `protobuf:"varint,6,opt,name=msg_type,json=msgType" json:"msg_type,omitempty"`
	MsgTime  uint32   `protobuf:"varint,7,opt,name=msg_time,json=msgTime" json:"msg_time,omitempty"`
	Msg      []byte   `protobuf:"bytes,8,opt,name=msg,proto3" json:"msg,omitempty"`
	ReadUids []uint64 `protobuf:"varint,9,rep,packed,name=read_uids,json=readUids" json:"read_uids,omitempty"`
}

func (m *MsgInfo) Reset()                    { *m = MsgInfo{} }
func (m *MsgInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgInfo) ProtoMessage()               {}
func (*MsgInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// 消息会话信息中只保存一条最新的消息以及未读过的消息数
type MsgSessionInfo struct {
	Msg          *MsgInfo `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	UnreadMsgNum uint32   `protobuf:"varint,2,opt,name=unread_msg_num,json=unreadMsgNum" json:"unread_msg_num,omitempty"`
}

func (m *MsgSessionInfo) Reset()                    { *m = MsgSessionInfo{} }
func (m *MsgSessionInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgSessionInfo) ProtoMessage()               {}
func (*MsgSessionInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MsgSessionInfo) GetMsg() *MsgInfo {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*BaseRequest)(nil), "com.kook.proto.BaseRequest")
	proto.RegisterType((*BaseResult)(nil), "com.kook.proto.BaseResult")
	proto.RegisterType((*MsgInfo)(nil), "com.kook.proto.MsgInfo")
	proto.RegisterType((*MsgSessionInfo)(nil), "com.kook.proto.MsgSessionInfo")
}

func init() { proto.RegisterFile("base.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x51, 0x4f, 0x4f, 0xfa, 0x40,
	0x10, 0x0d, 0x50, 0x5a, 0x18, 0xf8, 0x91, 0x5f, 0x9a, 0x18, 0x6a, 0xf4, 0x40, 0x1a, 0x0f, 0x7a,
	0xe9, 0x41, 0xbf, 0x81, 0x9e, 0x48, 0xd4, 0x43, 0x95, 0x73, 0x53, 0xda, 0x85, 0x34, 0xb0, 0x9d,
	0xba, 0x7f, 0x4c, 0x3c, 0xf9, 0xb1, 0xbd, 0x3a, 0xb3, 0x0b, 0x12, 0xcf, 0x9e, 0x76, 0xde, 0x7b,
	0xb3, 0x6f, 0xde, 0xec, 0x02, 0xac, 0x4b, 0x2d, 0xb2, 0x4e, 0xa1, 0xc1, 0x78, 0x56, 0xa1, 0xcc,
	0x76, 0x88, 0x3b, 0x8f, 0xd3, 0x47, 0x98, 0xdc, 0x93, 0x9a, 0x8b, 0x37, 0x2b, 0xb4, 0x89, 0xcf,
	0x61, 0x64, 0x54, 0xd9, 0xea, 0xa2, 0xa9, 0x93, 0xde, 0xa2, 0x77, 0x3d, 0xce, 0x23, 0x87, 0x97,
	0x75, 0xfc, 0x1f, 0x06, 0x15, 0xb1, 0x7d, 0x62, 0xff, 0xe5, 0x5c, 0x32, 0x63, 0x89, 0x19, 0x10,
	0x13, 0xe4, 0x5c, 0xa6, 0x9f, 0x00, 0xde, 0x4d, 0xdb, 0xfd, 0x5f, 0xcd, 0xf8, 0xba, 0x12, 0xa6,
	0xa8, 0xb0, 0x16, 0x49, 0xe0, 0x1a, 0x23, 0xc2, 0x0f, 0x04, 0xe3, 0x39, 0x70, 0x59, 0x48, 0xbd,
	0x4d, 0x86, 0xce, 0x38, 0x24, 0xf8, 0xa4, 0xb7, 0xe9, 0x57, 0x0f, 0x22, 0x3a, 0x97, 0xed, 0x06,
	0xe3, 0x4b, 0x00, 0xad, 0xde, 0xb9, 0xe9, 0x18, 0x20, 0xc8, 0x47, 0xc4, 0xb0, 0x5e, 0xb3, 0x5a,
	0xed, 0x9b, 0xa3, 0xda, 0xf7, 0x2a, 0x31, 0x5e, 0xa5, 0xd9, 0x1b, 0x85, 0xb2, 0x38, 0x45, 0x8a,
	0x18, 0xaf, 0x28, 0xd6, 0x19, 0x84, 0x06, 0x9d, 0x10, 0x38, 0x61, 0x68, 0x70, 0xe5, 0xd3, 0x6e,
	0x15, 0xda, 0x8e, 0xdd, 0x86, 0x3e, 0xad, 0xc3, 0xde, 0x8c, 0xc7, 0x98, 0x8f, 0x4e, 0x24, 0xa1,
	0x97, 0x08, 0xbf, 0x12, 0xfc, 0x91, 0x1a, 0x29, 0x92, 0xe8, 0x24, 0x11, 0xe4, 0x07, 0xe1, 0xfd,
	0x46, 0xc4, 0x4e, 0x73, 0x2e, 0xe3, 0x0b, 0x18, 0x2b, 0x51, 0xd6, 0x3c, 0x5b, 0x27, 0xe3, 0xc5,
	0x80, 0x13, 0x33, 0x41, 0xe3, 0x75, 0x5a, 0xc2, 0x8c, 0xa2, 0xbf, 0x08, 0xad, 0x1b, 0x6c, 0xdd,
	0xfe, 0x37, 0xde, 0x80, 0x17, 0x9f, 0xdc, 0xce, 0xb3, 0xdf, 0x1f, 0x9f, 0x1d, 0x5e, 0xc9, 0x3b,
	0x5f, 0xc1, 0xcc, 0xb6, 0xce, 0x9b, 0xd3, 0xb4, 0x56, 0x1e, 0x7e, 0x66, 0xea, 0x59, 0xea, 0x7d,
	0xb6, 0x72, 0x1d, 0xba, 0x9b, 0x77, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xeb, 0xa6, 0x83, 0x13,
	0x50, 0x02, 0x00, 0x00,
}
