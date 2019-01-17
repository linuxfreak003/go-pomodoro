// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/pomodoro.proto

/*
Package pomodoro is a generated protocol buffer package.

It is generated from these files:
	pb/pomodoro.proto

It has these top-level messages:
	Profile
	Timer
*/
package pomodoro

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type State int32

const (
	State_FOCUS State = 0
	State_BREAK State = 1
)

var State_name = map[int32]string{
	0: "FOCUS",
	1: "BREAK",
}
var State_value = map[string]int32{
	"FOCUS": 0,
	"BREAK": 1,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Profile struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *Profile) Reset()                    { *m = Profile{} }
func (m *Profile) String() string            { return proto.CompactTextString(m) }
func (*Profile) ProtoMessage()               {}
func (*Profile) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Profile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Timer struct {
	Nanoseconds int64 `protobuf:"varint,1,opt,name=Nanoseconds" json:"Nanoseconds,omitempty"`
	State       State `protobuf:"varint,2,opt,name=state,enum=State" json:"state,omitempty"`
}

func (m *Timer) Reset()                    { *m = Timer{} }
func (m *Timer) String() string            { return proto.CompactTextString(m) }
func (*Timer) ProtoMessage()               {}
func (*Timer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Timer) GetNanoseconds() int64 {
	if m != nil {
		return m.Nanoseconds
	}
	return 0
}

func (m *Timer) GetState() State {
	if m != nil {
		return m.State
	}
	return State_FOCUS
}

func init() {
	proto.RegisterType((*Profile)(nil), "Profile")
	proto.RegisterType((*Timer)(nil), "Timer")
	proto.RegisterEnum("State", State_name, State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Pomodoro service

type PomodoroClient interface {
	Sync(ctx context.Context, in *Profile, opts ...grpc.CallOption) (*Timer, error)
}

type pomodoroClient struct {
	cc *grpc.ClientConn
}

func NewPomodoroClient(cc *grpc.ClientConn) PomodoroClient {
	return &pomodoroClient{cc}
}

func (c *pomodoroClient) Sync(ctx context.Context, in *Profile, opts ...grpc.CallOption) (*Timer, error) {
	out := new(Timer)
	err := grpc.Invoke(ctx, "/Pomodoro/Sync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Pomodoro service

type PomodoroServer interface {
	Sync(context.Context, *Profile) (*Timer, error)
}

func RegisterPomodoroServer(s *grpc.Server, srv PomodoroServer) {
	s.RegisterService(&_Pomodoro_serviceDesc, srv)
}

func _Pomodoro_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Profile)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PomodoroServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pomodoro/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PomodoroServer).Sync(ctx, req.(*Profile))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pomodoro_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Pomodoro",
	HandlerType: (*PomodoroServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _Pomodoro_Sync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/pomodoro.proto",
}

func init() { proto.RegisterFile("pb/pomodoro.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x48, 0xd2, 0x2f,
	0xc8, 0xcf, 0xcd, 0x4f, 0xc9, 0x2f, 0xca, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe5,
	0x62, 0x0f, 0x28, 0xca, 0x4f, 0xcb, 0xcc, 0x49, 0x15, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d,
	0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0xdc, 0xb9, 0x58, 0x43, 0x32, 0x73,
	0x53, 0x8b, 0x84, 0x14, 0xb8, 0xb8, 0xfd, 0x12, 0xf3, 0xf2, 0x8b, 0x53, 0x93, 0xf3, 0xf3, 0x52,
	0x8a, 0xc1, 0x6a, 0x98, 0x83, 0x90, 0x85, 0x84, 0x64, 0xb8, 0x58, 0x8b, 0x4b, 0x12, 0x4b, 0x52,
	0x25, 0x98, 0x14, 0x18, 0x35, 0xf8, 0x8c, 0xd8, 0xf4, 0x82, 0x41, 0xbc, 0x20, 0x88, 0xa0, 0x96,
	0x2c, 0x17, 0x2b, 0x98, 0x2f, 0xc4, 0xc9, 0xc5, 0xea, 0xe6, 0xef, 0x1c, 0x1a, 0x2c, 0xc0, 0x00,
	0x62, 0x3a, 0x05, 0xb9, 0x3a, 0x7a, 0x0b, 0x30, 0x1a, 0xa9, 0x70, 0x71, 0x04, 0x40, 0x1d, 0x26,
	0x24, 0xc1, 0xc5, 0x12, 0x5c, 0x99, 0x97, 0x2c, 0xc4, 0xa1, 0x07, 0x75, 0x99, 0x14, 0x9b, 0x1e,
	0xd8, 0x11, 0x49, 0x6c, 0x60, 0x37, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x76, 0x0a, 0x8b,
	0x2f, 0xc8, 0x00, 0x00, 0x00,
}
