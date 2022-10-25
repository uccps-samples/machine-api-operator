// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/uccps-samples/api/cloudnetwork/v1/generated.proto

package v1

import (
	fmt "fmt"

	io "io"

	proto "github.com/gogo/protobuf/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func (m *CloudPrivateIPConfig) Reset()      { *m = CloudPrivateIPConfig{} }
func (*CloudPrivateIPConfig) ProtoMessage() {}
func (*CloudPrivateIPConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_454253a7ab01c6d0, []int{0}
}
func (m *CloudPrivateIPConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CloudPrivateIPConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *CloudPrivateIPConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloudPrivateIPConfig.Merge(m, src)
}
func (m *CloudPrivateIPConfig) XXX_Size() int {
	return m.Size()
}
func (m *CloudPrivateIPConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CloudPrivateIPConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CloudPrivateIPConfig proto.InternalMessageInfo

func (m *CloudPrivateIPConfigList) Reset()      { *m = CloudPrivateIPConfigList{} }
func (*CloudPrivateIPConfigList) ProtoMessage() {}
func (*CloudPrivateIPConfigList) Descriptor() ([]byte, []int) {
	return fileDescriptor_454253a7ab01c6d0, []int{1}
}
func (m *CloudPrivateIPConfigList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CloudPrivateIPConfigList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *CloudPrivateIPConfigList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloudPrivateIPConfigList.Merge(m, src)
}
func (m *CloudPrivateIPConfigList) XXX_Size() int {
	return m.Size()
}
func (m *CloudPrivateIPConfigList) XXX_DiscardUnknown() {
	xxx_messageInfo_CloudPrivateIPConfigList.DiscardUnknown(m)
}

var xxx_messageInfo_CloudPrivateIPConfigList proto.InternalMessageInfo

func (m *CloudPrivateIPConfigSpec) Reset()      { *m = CloudPrivateIPConfigSpec{} }
func (*CloudPrivateIPConfigSpec) ProtoMessage() {}
func (*CloudPrivateIPConfigSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_454253a7ab01c6d0, []int{2}
}
func (m *CloudPrivateIPConfigSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CloudPrivateIPConfigSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *CloudPrivateIPConfigSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloudPrivateIPConfigSpec.Merge(m, src)
}
func (m *CloudPrivateIPConfigSpec) XXX_Size() int {
	return m.Size()
}
func (m *CloudPrivateIPConfigSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_CloudPrivateIPConfigSpec.DiscardUnknown(m)
}

var xxx_messageInfo_CloudPrivateIPConfigSpec proto.InternalMessageInfo

func (m *CloudPrivateIPConfigStatus) Reset()      { *m = CloudPrivateIPConfigStatus{} }
func (*CloudPrivateIPConfigStatus) ProtoMessage() {}
func (*CloudPrivateIPConfigStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_454253a7ab01c6d0, []int{3}
}
func (m *CloudPrivateIPConfigStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CloudPrivateIPConfigStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *CloudPrivateIPConfigStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloudPrivateIPConfigStatus.Merge(m, src)
}
func (m *CloudPrivateIPConfigStatus) XXX_Size() int {
	return m.Size()
}
func (m *CloudPrivateIPConfigStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_CloudPrivateIPConfigStatus.DiscardUnknown(m)
}

var xxx_messageInfo_CloudPrivateIPConfigStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CloudPrivateIPConfig)(nil), "github.com.uccps-samples.api.cloudnetwork.v1.CloudPrivateIPConfig")
	proto.RegisterType((*CloudPrivateIPConfigList)(nil), "github.com.uccps-samples.api.cloudnetwork.v1.CloudPrivateIPConfigList")
	proto.RegisterType((*CloudPrivateIPConfigSpec)(nil), "github.com.uccps-samples.api.cloudnetwork.v1.CloudPrivateIPConfigSpec")
	proto.RegisterType((*CloudPrivateIPConfigStatus)(nil), "github.com.uccps-samples.api.cloudnetwork.v1.CloudPrivateIPConfigStatus")
}

func init() {
	proto.RegisterFile("github.com/uccps-samples/api/cloudnetwork/v1/generated.proto", fileDescriptor_454253a7ab01c6d0)
}

var fileDescriptor_454253a7ab01c6d0 = []byte{
	// 482 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x3f, 0x6f, 0xd3, 0x40,
	0x18, 0xc6, 0x7d, 0x69, 0x5a, 0x95, 0x2b, 0x20, 0x64, 0x31, 0x58, 0x19, 0xae, 0x51, 0xa6, 0x2c,
	0xdc, 0x91, 0x0a, 0xa1, 0x0e, 0x88, 0xc1, 0x61, 0xa9, 0xc4, 0x9f, 0xca, 0x6c, 0x88, 0x81, 0xcb,
	0xf9, 0xe2, 0x1c, 0xa9, 0xef, 0x2c, 0xdf, 0xd9, 0x88, 0x8d, 0x8f, 0xc0, 0x77, 0xe0, 0xcb, 0x64,
	0x60, 0xe8, 0xd8, 0x85, 0x8a, 0x98, 0x2f, 0x82, 0xee, 0xe2, 0x24, 0x56, 0x9b, 0xaa, 0x91, 0xb2,
	0xf9, 0x7d, 0xed, 0xe7, 0xf9, 0xbd, 0xef, 0xf3, 0x1a, 0x9e, 0x26, 0xc2, 0x4c, 0x8a, 0x11, 0x66,
	0x2a, 0x25, 0x2a, 0xe3, 0x52, 0x4f, 0xc4, 0xd8, 0x10, 0x9a, 0x09, 0xc2, 0x2e, 0x54, 0x11, 0x4b,
	0x6e, 0xbe, 0xa9, 0x7c, 0x4a, 0xca, 0x01, 0x49, 0xb8, 0xe4, 0x39, 0x35, 0x3c, 0xc6, 0x59, 0xae,
	0x8c, 0xf2, 0xfb, 0x6b, 0x25, 0x5e, 0x29, 0x31, 0xcd, 0x04, 0x6e, 0x2a, 0x71, 0x39, 0xe8, 0x3c,
	0x6b, 0x30, 0x12, 0x95, 0x28, 0xe2, 0x0c, 0x46, 0xc5, 0xd8, 0x55, 0xae, 0x70, 0x4f, 0x0b, 0xe3,
	0xce, 0x8b, 0xe9, 0xa9, 0xc6, 0x42, 0xd9, 0x21, 0x52, 0xca, 0x26, 0x42, 0xf2, 0xfc, 0x3b, 0xc9,
	0xa6, 0x89, 0x6d, 0x68, 0x92, 0x72, 0x43, 0x37, 0x8c, 0xd3, 0x21, 0x77, 0xa9, 0xf2, 0x42, 0x1a,
	0x91, 0xf2, 0x5b, 0x82, 0x97, 0xf7, 0x09, 0x34, 0x9b, 0xf0, 0x94, 0xde, 0xd4, 0xf5, 0x7e, 0xb7,
	0xe0, 0xd3, 0xa1, 0xdd, 0xf0, 0x3c, 0x17, 0x25, 0x35, 0xfc, 0xec, 0x7c, 0xa8, 0xe4, 0x58, 0x24,
	0xfe, 0x17, 0x78, 0x68, 0x87, 0x8b, 0xa9, 0xa1, 0x01, 0xe8, 0x82, 0xfe, 0xd1, 0xc9, 0x73, 0xbc,
	0x60, 0xe0, 0x26, 0x03, 0x67, 0xd3, 0xc4, 0x36, 0x34, 0xb6, 0x5f, 0xe3, 0x72, 0x80, 0x3f, 0x8c,
	0xbe, 0x72, 0x66, 0xde, 0x71, 0x43, 0x43, 0x7f, 0x76, 0x7d, 0xec, 0x55, 0xd7, 0xc7, 0x70, 0xdd,
	0x8b, 0x56, 0xae, 0x7e, 0x0c, 0xdb, 0x3a, 0xe3, 0x2c, 0x68, 0x39, 0xf7, 0x10, 0x6f, 0x7b, 0x01,
	0xbc, 0x69, 0xde, 0x8f, 0x19, 0x67, 0xe1, 0xc3, 0x9a, 0xd7, 0xb6, 0x55, 0xe4, 0xdc, 0xfd, 0x0b,
	0x78, 0xa0, 0x0d, 0x35, 0x85, 0x0e, 0xf6, 0x1c, 0xe7, 0xcd, 0x8e, 0x1c, 0xe7, 0x15, 0x3e, 0xae,
	0x49, 0x07, 0x8b, 0x3a, 0xaa, 0x19, 0xbd, 0x3f, 0x00, 0x06, 0x9b, 0x64, 0x6f, 0x85, 0x36, 0xfe,
	0xe7, 0x5b, 0x91, 0xe2, 0xed, 0x22, 0xb5, 0x6a, 0x17, 0xe8, 0x93, 0x1a, 0x7b, 0xb8, 0xec, 0x34,
	0xe2, 0x64, 0x70, 0x5f, 0x18, 0x9e, 0xea, 0xa0, 0xd5, 0xdd, 0xeb, 0x1f, 0x9d, 0xbc, 0xde, 0x6d,
	0xcf, 0xf0, 0x51, 0x8d, 0xda, 0x3f, 0xb3, 0xa6, 0xd1, 0xc2, 0xbb, 0xf7, 0x6a, 0xf3, 0x7a, 0x36,
	0x6f, 0xbf, 0x0b, 0xdb, 0x52, 0xc5, 0xdc, 0xad, 0xf6, 0x60, 0x7d, 0x8b, 0xf7, 0x2a, 0xe6, 0x91,
	0x7b, 0xd3, 0xfb, 0x05, 0x60, 0xe7, 0xee, 0x50, 0xef, 0x37, 0xf0, 0x19, 0x84, 0x4c, 0xc9, 0x58,
	0x18, 0xa1, 0xe4, 0x72, 0x51, 0xb2, 0x5d, 0x86, 0xc3, 0xa5, 0x6e, 0xfd, 0x57, 0xae, 0x5a, 0x3a,
	0x6a, 0xd8, 0x86, 0xfd, 0xd9, 0x1c, 0x79, 0x97, 0x73, 0xe4, 0x5d, 0xcd, 0x91, 0xf7, 0xa3, 0x42,
	0x60, 0x56, 0x21, 0x70, 0x59, 0x21, 0x70, 0x55, 0x21, 0xf0, 0xb7, 0x42, 0xe0, 0xe7, 0x3f, 0xe4,
	0x7d, 0x6a, 0x95, 0x83, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa5, 0xd1, 0x85, 0x90, 0x6f, 0x04,
	0x00, 0x00,
}

func (m *CloudPrivateIPConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CloudPrivateIPConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CloudPrivateIPConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Spec.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.ObjectMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *CloudPrivateIPConfigList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CloudPrivateIPConfigList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CloudPrivateIPConfigList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Items) > 0 {
		for iNdEx := len(m.Items) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Items[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.ListMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *CloudPrivateIPConfigSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CloudPrivateIPConfigSpec) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CloudPrivateIPConfigSpec) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i -= len(m.Node)
	copy(dAtA[i:], m.Node)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Node)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *CloudPrivateIPConfigStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CloudPrivateIPConfigStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CloudPrivateIPConfigStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Conditions) > 0 {
		for iNdEx := len(m.Conditions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Conditions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	i -= len(m.Node)
	copy(dAtA[i:], m.Node)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Node)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CloudPrivateIPConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Spec.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Status.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *CloudPrivateIPConfigList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ListMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func (m *CloudPrivateIPConfigSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Node)
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *CloudPrivateIPConfigStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Node)
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Conditions) > 0 {
		for _, e := range m.Conditions {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *CloudPrivateIPConfig) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&CloudPrivateIPConfig{`,
		`ObjectMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ObjectMeta), "ObjectMeta", "v1.ObjectMeta", 1), `&`, ``, 1) + `,`,
		`Spec:` + strings.Replace(strings.Replace(this.Spec.String(), "CloudPrivateIPConfigSpec", "CloudPrivateIPConfigSpec", 1), `&`, ``, 1) + `,`,
		`Status:` + strings.Replace(strings.Replace(this.Status.String(), "CloudPrivateIPConfigStatus", "CloudPrivateIPConfigStatus", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *CloudPrivateIPConfigList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForItems := "[]CloudPrivateIPConfig{"
	for _, f := range this.Items {
		repeatedStringForItems += strings.Replace(strings.Replace(f.String(), "CloudPrivateIPConfig", "CloudPrivateIPConfig", 1), `&`, ``, 1) + ","
	}
	repeatedStringForItems += "}"
	s := strings.Join([]string{`&CloudPrivateIPConfigList{`,
		`ListMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ListMeta), "ListMeta", "v1.ListMeta", 1), `&`, ``, 1) + `,`,
		`Items:` + repeatedStringForItems + `,`,
		`}`,
	}, "")
	return s
}
func (this *CloudPrivateIPConfigSpec) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&CloudPrivateIPConfigSpec{`,
		`Node:` + fmt.Sprintf("%v", this.Node) + `,`,
		`}`,
	}, "")
	return s
}
func (this *CloudPrivateIPConfigStatus) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForConditions := "[]Condition{"
	for _, f := range this.Conditions {
		repeatedStringForConditions += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForConditions += "}"
	s := strings.Join([]string{`&CloudPrivateIPConfigStatus{`,
		`Node:` + fmt.Sprintf("%v", this.Node) + `,`,
		`Conditions:` + repeatedStringForConditions + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *CloudPrivateIPConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CloudPrivateIPConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CloudPrivateIPConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spec", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Spec.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CloudPrivateIPConfigList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CloudPrivateIPConfigList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CloudPrivateIPConfigList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ListMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, CloudPrivateIPConfig{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CloudPrivateIPConfigSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CloudPrivateIPConfigSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CloudPrivateIPConfigSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CloudPrivateIPConfigStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CloudPrivateIPConfigStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CloudPrivateIPConfigStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Conditions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Conditions = append(m.Conditions, v1.Condition{})
			if err := m.Conditions[len(m.Conditions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
