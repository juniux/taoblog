// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocols/post.proto

package protocols

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Post struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Date                 int32    `protobuf:"varint,2,opt,name=date,proto3" json:"date,omitempty"`
	Modified             int32    `protobuf:"varint,3,opt,name=modified,proto3" json:"modified,omitempty"`
	Title                string   `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Content              string   `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Slug                 string   `protobuf:"bytes,6,opt,name=slug,proto3" json:"slug,omitempty"`
	Type                 string   `protobuf:"bytes,7,opt,name=type,proto3" json:"type,omitempty"`
	Category             int64    `protobuf:"varint,8,opt,name=category,proto3" json:"category,omitempty"`
	Status               string   `protobuf:"bytes,9,opt,name=status,proto3" json:"status,omitempty"`
	PageView             int64    `protobuf:"varint,10,opt,name=page_view,json=pageView,proto3" json:"page_view,omitempty"`
	CommentStatus        bool     `protobuf:"varint,11,opt,name=comment_status,json=commentStatus,proto3" json:"comment_status,omitempty"`
	Comments             int64    `protobuf:"varint,12,opt,name=comments,proto3" json:"comments,omitempty"`
	Metas                string   `protobuf:"bytes,13,opt,name=metas,proto3" json:"metas,omitempty"`
	Source               string   `protobuf:"bytes,14,opt,name=source,proto3" json:"source,omitempty"`
	SourceType           string   `protobuf:"bytes,15,opt,name=source_type,json=sourceType,proto3" json:"source_type,omitempty"`
	Tags                 []string `protobuf:"bytes,16,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{0}
}

func (m *Post) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Post.Unmarshal(m, b)
}
func (m *Post) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Post.Marshal(b, m, deterministic)
}
func (m *Post) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Post.Merge(m, src)
}
func (m *Post) XXX_Size() int {
	return xxx_messageInfo_Post.Size(m)
}
func (m *Post) XXX_DiscardUnknown() {
	xxx_messageInfo_Post.DiscardUnknown(m)
}

var xxx_messageInfo_Post proto.InternalMessageInfo

func (m *Post) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Post) GetDate() int32 {
	if m != nil {
		return m.Date
	}
	return 0
}

func (m *Post) GetModified() int32 {
	if m != nil {
		return m.Modified
	}
	return 0
}

func (m *Post) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Post) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Post) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *Post) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Post) GetCategory() int64 {
	if m != nil {
		return m.Category
	}
	return 0
}

func (m *Post) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Post) GetPageView() int64 {
	if m != nil {
		return m.PageView
	}
	return 0
}

func (m *Post) GetCommentStatus() bool {
	if m != nil {
		return m.CommentStatus
	}
	return false
}

func (m *Post) GetComments() int64 {
	if m != nil {
		return m.Comments
	}
	return 0
}

func (m *Post) GetMetas() string {
	if m != nil {
		return m.Metas
	}
	return ""
}

func (m *Post) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Post) GetSourceType() string {
	if m != nil {
		return m.SourceType
	}
	return ""
}

func (m *Post) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type GetPostRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	WithSource           bool     `protobuf:"varint,2,opt,name=with_source,json=withSource,proto3" json:"with_source,omitempty"`
	WithContent          bool     `protobuf:"varint,3,opt,name=with_content,json=withContent,proto3" json:"with_content,omitempty"`
	WithTags             bool     `protobuf:"varint,4,opt,name=with_tags,json=withTags,proto3" json:"with_tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostRequest) Reset()         { *m = GetPostRequest{} }
func (m *GetPostRequest) String() string { return proto.CompactTextString(m) }
func (*GetPostRequest) ProtoMessage()    {}
func (*GetPostRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{1}
}

func (m *GetPostRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostRequest.Unmarshal(m, b)
}
func (m *GetPostRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostRequest.Marshal(b, m, deterministic)
}
func (m *GetPostRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostRequest.Merge(m, src)
}
func (m *GetPostRequest) XXX_Size() int {
	return xxx_messageInfo_GetPostRequest.Size(m)
}
func (m *GetPostRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostRequest proto.InternalMessageInfo

func (m *GetPostRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetPostRequest) GetWithSource() bool {
	if m != nil {
		return m.WithSource
	}
	return false
}

func (m *GetPostRequest) GetWithContent() bool {
	if m != nil {
		return m.WithContent
	}
	return false
}

func (m *GetPostRequest) GetWithTags() bool {
	if m != nil {
		return m.WithTags
	}
	return false
}

type UpdatePostRequest struct {
	Post                 *Post                  `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
	UpdateMask           *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *UpdatePostRequest) Reset()         { *m = UpdatePostRequest{} }
func (m *UpdatePostRequest) String() string { return proto.CompactTextString(m) }
func (*UpdatePostRequest) ProtoMessage()    {}
func (*UpdatePostRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{2}
}

func (m *UpdatePostRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatePostRequest.Unmarshal(m, b)
}
func (m *UpdatePostRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatePostRequest.Marshal(b, m, deterministic)
}
func (m *UpdatePostRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatePostRequest.Merge(m, src)
}
func (m *UpdatePostRequest) XXX_Size() int {
	return xxx_messageInfo_UpdatePostRequest.Size(m)
}
func (m *UpdatePostRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatePostRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatePostRequest proto.InternalMessageInfo

func (m *UpdatePostRequest) GetPost() *Post {
	if m != nil {
		return m.Post
	}
	return nil
}

func (m *UpdatePostRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type DeletePostRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePostRequest) Reset()         { *m = DeletePostRequest{} }
func (m *DeletePostRequest) String() string { return proto.CompactTextString(m) }
func (*DeletePostRequest) ProtoMessage()    {}
func (*DeletePostRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{3}
}

func (m *DeletePostRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePostRequest.Unmarshal(m, b)
}
func (m *DeletePostRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePostRequest.Marshal(b, m, deterministic)
}
func (m *DeletePostRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePostRequest.Merge(m, src)
}
func (m *DeletePostRequest) XXX_Size() int {
	return xxx_messageInfo_DeletePostRequest.Size(m)
}
func (m *DeletePostRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePostRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePostRequest proto.InternalMessageInfo

func (m *DeletePostRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetPostCommentsCountRequest struct {
	PostId               int64    `protobuf:"varint,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostCommentsCountRequest) Reset()         { *m = GetPostCommentsCountRequest{} }
func (m *GetPostCommentsCountRequest) String() string { return proto.CompactTextString(m) }
func (*GetPostCommentsCountRequest) ProtoMessage()    {}
func (*GetPostCommentsCountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{4}
}

func (m *GetPostCommentsCountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostCommentsCountRequest.Unmarshal(m, b)
}
func (m *GetPostCommentsCountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostCommentsCountRequest.Marshal(b, m, deterministic)
}
func (m *GetPostCommentsCountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostCommentsCountRequest.Merge(m, src)
}
func (m *GetPostCommentsCountRequest) XXX_Size() int {
	return xxx_messageInfo_GetPostCommentsCountRequest.Size(m)
}
func (m *GetPostCommentsCountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostCommentsCountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostCommentsCountRequest proto.InternalMessageInfo

func (m *GetPostCommentsCountRequest) GetPostId() int64 {
	if m != nil {
		return m.PostId
	}
	return 0
}

type GetPostCommentsCountResponse struct {
	Count                int64    `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostCommentsCountResponse) Reset()         { *m = GetPostCommentsCountResponse{} }
func (m *GetPostCommentsCountResponse) String() string { return proto.CompactTextString(m) }
func (*GetPostCommentsCountResponse) ProtoMessage()    {}
func (*GetPostCommentsCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{5}
}

func (m *GetPostCommentsCountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostCommentsCountResponse.Unmarshal(m, b)
}
func (m *GetPostCommentsCountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostCommentsCountResponse.Marshal(b, m, deterministic)
}
func (m *GetPostCommentsCountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostCommentsCountResponse.Merge(m, src)
}
func (m *GetPostCommentsCountResponse) XXX_Size() int {
	return xxx_messageInfo_GetPostCommentsCountResponse.Size(m)
}
func (m *GetPostCommentsCountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostCommentsCountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostCommentsCountResponse proto.InternalMessageInfo

func (m *GetPostCommentsCountResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type SetPostStatusRequest struct {
	Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Public bool  `protobuf:"varint,2,opt,name=public,proto3" json:"public,omitempty"`
	// Whether to create_time and update_time
	Touch                bool     `protobuf:"varint,3,opt,name=touch,proto3" json:"touch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetPostStatusRequest) Reset()         { *m = SetPostStatusRequest{} }
func (m *SetPostStatusRequest) String() string { return proto.CompactTextString(m) }
func (*SetPostStatusRequest) ProtoMessage()    {}
func (*SetPostStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{6}
}

func (m *SetPostStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetPostStatusRequest.Unmarshal(m, b)
}
func (m *SetPostStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetPostStatusRequest.Marshal(b, m, deterministic)
}
func (m *SetPostStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetPostStatusRequest.Merge(m, src)
}
func (m *SetPostStatusRequest) XXX_Size() int {
	return xxx_messageInfo_SetPostStatusRequest.Size(m)
}
func (m *SetPostStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetPostStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetPostStatusRequest proto.InternalMessageInfo

func (m *SetPostStatusRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SetPostStatusRequest) GetPublic() bool {
	if m != nil {
		return m.Public
	}
	return false
}

func (m *SetPostStatusRequest) GetTouch() bool {
	if m != nil {
		return m.Touch
	}
	return false
}

type SetPostStatusResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetPostStatusResponse) Reset()         { *m = SetPostStatusResponse{} }
func (m *SetPostStatusResponse) String() string { return proto.CompactTextString(m) }
func (*SetPostStatusResponse) ProtoMessage()    {}
func (*SetPostStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{7}
}

func (m *SetPostStatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetPostStatusResponse.Unmarshal(m, b)
}
func (m *SetPostStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetPostStatusResponse.Marshal(b, m, deterministic)
}
func (m *SetPostStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetPostStatusResponse.Merge(m, src)
}
func (m *SetPostStatusResponse) XXX_Size() int {
	return xxx_messageInfo_SetPostStatusResponse.Size(m)
}
func (m *SetPostStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetPostStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetPostStatusResponse proto.InternalMessageInfo

type GetPostSourceRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostSourceRequest) Reset()         { *m = GetPostSourceRequest{} }
func (m *GetPostSourceRequest) String() string { return proto.CompactTextString(m) }
func (*GetPostSourceRequest) ProtoMessage()    {}
func (*GetPostSourceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{8}
}

func (m *GetPostSourceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostSourceRequest.Unmarshal(m, b)
}
func (m *GetPostSourceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostSourceRequest.Marshal(b, m, deterministic)
}
func (m *GetPostSourceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostSourceRequest.Merge(m, src)
}
func (m *GetPostSourceRequest) XXX_Size() int {
	return xxx_messageInfo_GetPostSourceRequest.Size(m)
}
func (m *GetPostSourceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostSourceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostSourceRequest proto.InternalMessageInfo

func (m *GetPostSourceRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetPostSourceResponse struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Content              string   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostSourceResponse) Reset()         { *m = GetPostSourceResponse{} }
func (m *GetPostSourceResponse) String() string { return proto.CompactTextString(m) }
func (*GetPostSourceResponse) ProtoMessage()    {}
func (*GetPostSourceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b419c78abee5f34, []int{9}
}

func (m *GetPostSourceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostSourceResponse.Unmarshal(m, b)
}
func (m *GetPostSourceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostSourceResponse.Marshal(b, m, deterministic)
}
func (m *GetPostSourceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostSourceResponse.Merge(m, src)
}
func (m *GetPostSourceResponse) XXX_Size() int {
	return xxx_messageInfo_GetPostSourceResponse.Size(m)
}
func (m *GetPostSourceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostSourceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostSourceResponse proto.InternalMessageInfo

func (m *GetPostSourceResponse) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *GetPostSourceResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*Post)(nil), "protocols.Post")
	proto.RegisterType((*GetPostRequest)(nil), "protocols.GetPostRequest")
	proto.RegisterType((*UpdatePostRequest)(nil), "protocols.UpdatePostRequest")
	proto.RegisterType((*DeletePostRequest)(nil), "protocols.DeletePostRequest")
	proto.RegisterType((*GetPostCommentsCountRequest)(nil), "protocols.GetPostCommentsCountRequest")
	proto.RegisterType((*GetPostCommentsCountResponse)(nil), "protocols.GetPostCommentsCountResponse")
	proto.RegisterType((*SetPostStatusRequest)(nil), "protocols.SetPostStatusRequest")
	proto.RegisterType((*SetPostStatusResponse)(nil), "protocols.SetPostStatusResponse")
	proto.RegisterType((*GetPostSourceRequest)(nil), "protocols.GetPostSourceRequest")
	proto.RegisterType((*GetPostSourceResponse)(nil), "protocols.GetPostSourceResponse")
}

func init() { proto.RegisterFile("protocols/post.proto", fileDescriptor_3b419c78abee5f34) }

var fileDescriptor_3b419c78abee5f34 = []byte{
	// 583 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0xf3, 0xd5, 0x64, 0xdc, 0xa6, 0x74, 0x95, 0xb6, 0xab, 0x04, 0xa9, 0xc6, 0x05, 0x94,
	0x93, 0x2d, 0x15, 0xc4, 0x85, 0x1b, 0xe1, 0x43, 0x1c, 0x90, 0x90, 0x13, 0x38, 0x70, 0xb1, 0xfc,
	0xb1, 0x71, 0xac, 0xda, 0x59, 0x93, 0x5d, 0x37, 0xea, 0x91, 0x1f, 0xc1, 0xff, 0x45, 0x3b, 0xbb,
	0xb1, 0x92, 0x00, 0xb7, 0x79, 0xcf, 0x33, 0x3b, 0x6f, 0xfc, 0x66, 0x60, 0x54, 0x6d, 0xb8, 0xe4,
	0x09, 0x2f, 0x84, 0x5f, 0x71, 0x21, 0x3d, 0x84, 0x64, 0xd0, 0xb0, 0x63, 0x27, 0xe3, 0x3c, 0x2b,
	0x98, 0x8f, 0x4c, 0x5c, 0x2f, 0xfd, 0x65, 0xce, 0x8a, 0x34, 0x2c, 0x23, 0x71, 0xaf, 0x93, 0xdd,
	0xdf, 0x6d, 0xe8, 0x7c, 0xe5, 0x42, 0x92, 0x21, 0xb4, 0xf2, 0x94, 0x5a, 0x8e, 0x35, 0x6d, 0x07,
	0xad, 0x3c, 0x25, 0x04, 0x3a, 0x69, 0x24, 0x19, 0x6d, 0x39, 0xd6, 0xb4, 0x1b, 0x60, 0x4c, 0xc6,
	0xd0, 0x2f, 0x79, 0x9a, 0x2f, 0x73, 0x96, 0xd2, 0x36, 0xf2, 0x0d, 0x26, 0x23, 0xe8, 0xca, 0x5c,
	0x16, 0x8c, 0x76, 0x1c, 0x6b, 0x3a, 0x08, 0x34, 0x20, 0x14, 0x4e, 0x12, 0xbe, 0x96, 0x6c, 0x2d,
	0x69, 0x17, 0xf9, 0x1d, 0x54, 0xef, 0x8b, 0xa2, 0xce, 0x68, 0x0f, 0x69, 0x8c, 0x15, 0x27, 0x1f,
	0x2b, 0x46, 0x4f, 0x34, 0xa7, 0x62, 0xd5, 0x33, 0x89, 0x24, 0xcb, 0xf8, 0xe6, 0x91, 0xf6, 0x51,
	0x5d, 0x83, 0xc9, 0x15, 0xf4, 0x84, 0x8c, 0x64, 0x2d, 0xe8, 0x00, 0x2b, 0x0c, 0x22, 0x13, 0x18,
	0x54, 0x51, 0xc6, 0xc2, 0x87, 0x9c, 0x6d, 0x29, 0xe8, 0x22, 0x45, 0x7c, 0xcf, 0xd9, 0x96, 0xbc,
	0x80, 0x61, 0xc2, 0xcb, 0x92, 0xad, 0x65, 0x68, 0x8a, 0x6d, 0xc7, 0x9a, 0xf6, 0x83, 0x33, 0xc3,
	0xce, 0xf5, 0x1b, 0xaa, 0xaf, 0x26, 0x04, 0x3d, 0x35, 0x7d, 0x0d, 0x56, 0xb3, 0x96, 0x4c, 0x46,
	0x82, 0x9e, 0xe9, 0x59, 0x11, 0xa0, 0x1a, 0x5e, 0x6f, 0x12, 0x46, 0x87, 0x46, 0x0d, 0x22, 0x72,
	0x03, 0xb6, 0x8e, 0x42, 0x1c, 0xee, 0x1c, 0x3f, 0x82, 0xa6, 0x16, 0x6a, 0x44, 0x35, 0x76, 0x94,
	0x09, 0xfa, 0xc4, 0x69, 0xe3, 0xd8, 0x51, 0x26, 0xdc, 0x5f, 0x16, 0x0c, 0x3f, 0x31, 0xa9, 0xac,
	0x09, 0xd8, 0xcf, 0x9a, 0x1d, 0x38, 0xd4, 0x45, 0x87, 0x6e, 0xc0, 0xde, 0xe6, 0x72, 0x15, 0x9a,
	0xa6, 0x2d, 0x9c, 0x02, 0x14, 0x35, 0xd7, 0x8d, 0x9f, 0xc1, 0x29, 0x26, 0xec, 0x1c, 0x68, 0x63,
	0x06, 0x16, 0xcd, 0x8c, 0x0b, 0x13, 0x18, 0x60, 0x0a, 0xf6, 0xef, 0xe0, 0xf7, 0xbe, 0x22, 0x16,
	0x4a, 0x43, 0x0d, 0x17, 0xdf, 0x2a, 0x65, 0xfc, 0xbe, 0x8a, 0x5b, 0xe8, 0xa8, 0x5d, 0x43, 0x1d,
	0xf6, 0xdd, 0xb9, 0xd7, 0x2c, 0x9b, 0x87, 0x59, 0xf8, 0x91, 0xbc, 0x05, 0xbb, 0xc6, 0x4a, 0x5c,
	0x35, 0x94, 0x66, 0xdf, 0x8d, 0x3d, 0xbd, 0x8d, 0xde, 0x6e, 0x1b, 0xbd, 0x8f, 0x6a, 0x1b, 0xbf,
	0x44, 0xe2, 0x3e, 0x00, 0x9d, 0xae, 0x62, 0xf7, 0x16, 0x2e, 0xde, 0xb3, 0x82, 0x1d, 0xb6, 0x3d,
	0x1a, 0xde, 0x7d, 0x03, 0x13, 0xf3, 0x7b, 0x66, 0xc6, 0x95, 0x19, 0xaf, 0xd7, 0x4d, 0xfa, 0x35,
	0x9c, 0x28, 0x21, 0x61, 0xb3, 0xd2, 0x3d, 0x05, 0x3f, 0xa7, 0xee, 0x6b, 0x78, 0xfa, 0xef, 0x3a,
	0x51, 0xf1, 0xb5, 0x60, 0xca, 0xda, 0x44, 0x11, 0xa6, 0x4c, 0x03, 0x77, 0x01, 0xa3, 0xb9, 0xae,
	0xd2, 0xdb, 0xf1, 0xb7, 0x2a, 0x7d, 0x34, 0x57, 0xd0, 0xab, 0xea, 0xb8, 0xc8, 0x13, 0xe3, 0x86,
	0x41, 0x78, 0x1c, 0xbc, 0x4e, 0x56, 0xc6, 0x02, 0x0d, 0xdc, 0x6b, 0xb8, 0x3c, 0x7a, 0x55, 0x8b,
	0x70, 0x5f, 0xc2, 0xc8, 0x88, 0xd4, 0x4e, 0xfe, 0xa7, 0x9d, 0xfb, 0x01, 0x2e, 0x8f, 0xf2, 0xcc,
	0x14, 0xbb, 0x43, 0xb2, 0xf6, 0x0e, 0x69, 0xef, 0x14, 0x5b, 0x07, 0xa7, 0xf8, 0xee, 0xf9, 0x0f,
	0x37, 0xcb, 0xe5, 0xaa, 0x8e, 0xbd, 0x84, 0x97, 0x7e, 0xc9, 0x1f, 0x44, 0xec, 0xcb, 0x88, 0xc7,
	0x05, 0xcf, 0xfc, 0xc6, 0xde, 0xb8, 0x87, 0xe1, 0xab, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4c,
	0x4e, 0x4e, 0x20, 0x75, 0x04, 0x00, 0x00,
}
