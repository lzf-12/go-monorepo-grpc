// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: pb_schemas/inventory/v1/stock.proto

package inventoryv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrorCode int32

const (
	ErrorCode_UNDEFINED                        ErrorCode = 0
	ErrorCode_SKU_NOT_FOUND                    ErrorCode = 1
	ErrorCode_SKU_UOM_PAIR_NOT_MATCH           ErrorCode = 2
	ErrorCode_DB_UNAVAILABLE                   ErrorCode = 3
	ErrorCode_DB_ERROR_TRANSACTION             ErrorCode = 4
	ErrorCode_INTERNAL_ERROR                   ErrorCode = 5
	ErrorCode_INSUFFICIENT_QUANTITY_TO_RESERVE ErrorCode = 6
	ErrorCode_INSUFFICIENT_QUANTITY_TO_RELEASE ErrorCode = 7
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "UNDEFINED",
		1: "SKU_NOT_FOUND",
		2: "SKU_UOM_PAIR_NOT_MATCH",
		3: "DB_UNAVAILABLE",
		4: "DB_ERROR_TRANSACTION",
		5: "INTERNAL_ERROR",
		6: "INSUFFICIENT_QUANTITY_TO_RESERVE",
		7: "INSUFFICIENT_QUANTITY_TO_RELEASE",
	}
	ErrorCode_value = map[string]int32{
		"UNDEFINED":                        0,
		"SKU_NOT_FOUND":                    1,
		"SKU_UOM_PAIR_NOT_MATCH":           2,
		"DB_UNAVAILABLE":                   3,
		"DB_ERROR_TRANSACTION":             4,
		"INTERNAL_ERROR":                   5,
		"INSUFFICIENT_QUANTITY_TO_RESERVE": 6,
		"INSUFFICIENT_QUANTITY_TO_RELEASE": 7,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_schemas_inventory_v1_stock_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_pb_schemas_inventory_v1_stock_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{0}
}

// Inventory Item Definition
type InventoryItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sku           string                 `protobuf:"bytes,1,opt,name=sku,proto3" json:"sku,omitempty"`
	ReqQtyPerUom  float64                `protobuf:"fixed64,2,opt,name=req_qty_per_uom,json=reqQtyPerUom,proto3" json:"req_qty_per_uom,omitempty"`
	Uom           string                 `protobuf:"bytes,3,opt,name=uom,proto3" json:"uom,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InventoryItem) Reset() {
	*x = InventoryItem{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryItem) ProtoMessage() {}

func (x *InventoryItem) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryItem.ProtoReflect.Descriptor instead.
func (*InventoryItem) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{0}
}

func (x *InventoryItem) GetSku() string {
	if x != nil {
		return x.Sku
	}
	return ""
}

func (x *InventoryItem) GetReqQtyPerUom() float64 {
	if x != nil {
		return x.ReqQtyPerUom
	}
	return 0
}

func (x *InventoryItem) GetUom() string {
	if x != nil {
		return x.Uom
	}
	return ""
}

// Inventory Status for a single item
type InventoryStatus struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Sku               string                 `protobuf:"bytes,1,opt,name=sku,proto3" json:"sku,omitempty"`
	RequestedQuantity float64                `protobuf:"fixed64,2,opt,name=requested_quantity,json=requestedQuantity,proto3" json:"requested_quantity,omitempty"`
	AvailableQuantity float64                `protobuf:"fixed64,3,opt,name=available_quantity,json=availableQuantity,proto3" json:"available_quantity,omitempty"`
	ReservedQuantity  float64                `protobuf:"fixed64,4,opt,name=reserved_quantity,json=reservedQuantity,proto3" json:"reserved_quantity,omitempty"`
	TotalQuantity     float64                `protobuf:"fixed64,5,opt,name=total_quantity,json=totalQuantity,proto3" json:"total_quantity,omitempty"`
	SkuUom            string                 `protobuf:"bytes,6,opt,name=sku_uom,json=skuUom,proto3" json:"sku_uom,omitempty"`
	SkuPrice          float64                `protobuf:"fixed64,7,opt,name=sku_price,json=skuPrice,proto3" json:"sku_price,omitempty"`
	SkuCurrency       string                 `protobuf:"bytes,8,opt,name=sku_currency,json=skuCurrency,proto3" json:"sku_currency,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *InventoryStatus) Reset() {
	*x = InventoryStatus{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryStatus) ProtoMessage() {}

func (x *InventoryStatus) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryStatus.ProtoReflect.Descriptor instead.
func (*InventoryStatus) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{1}
}

func (x *InventoryStatus) GetSku() string {
	if x != nil {
		return x.Sku
	}
	return ""
}

func (x *InventoryStatus) GetRequestedQuantity() float64 {
	if x != nil {
		return x.RequestedQuantity
	}
	return 0
}

func (x *InventoryStatus) GetAvailableQuantity() float64 {
	if x != nil {
		return x.AvailableQuantity
	}
	return 0
}

func (x *InventoryStatus) GetReservedQuantity() float64 {
	if x != nil {
		return x.ReservedQuantity
	}
	return 0
}

func (x *InventoryStatus) GetTotalQuantity() float64 {
	if x != nil {
		return x.TotalQuantity
	}
	return 0
}

func (x *InventoryStatus) GetSkuUom() string {
	if x != nil {
		return x.SkuUom
	}
	return ""
}

func (x *InventoryStatus) GetSkuPrice() float64 {
	if x != nil {
		return x.SkuPrice
	}
	return 0
}

func (x *InventoryStatus) GetSkuCurrency() string {
	if x != nil {
		return x.SkuCurrency
	}
	return ""
}

type ReservedItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OrderId       string                 `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReservedItem) Reset() {
	*x = ReservedItem{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReservedItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReservedItem) ProtoMessage() {}

func (x *ReservedItem) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReservedItem.ProtoReflect.Descriptor instead.
func (*ReservedItem) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{2}
}

func (x *ReservedItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ReservedItem) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

// Request to check inventory
type StandardInventoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderId       string                 `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Items         []*InventoryItem       `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StandardInventoryRequest) Reset() {
	*x = StandardInventoryRequest{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StandardInventoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StandardInventoryRequest) ProtoMessage() {}

func (x *StandardInventoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StandardInventoryRequest.ProtoReflect.Descriptor instead.
func (*StandardInventoryRequest) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{3}
}

func (x *StandardInventoryRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *StandardInventoryRequest) GetItems() []*InventoryItem {
	if x != nil {
		return x.Items
	}
	return nil
}

// Successful response
type InventoryStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*InventoryStatus     `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InventoryStatusResponse) Reset() {
	*x = InventoryStatusResponse{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryStatusResponse) ProtoMessage() {}

func (x *InventoryStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryStatusResponse.ProtoReflect.Descriptor instead.
func (*InventoryStatusResponse) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{4}
}

func (x *InventoryStatusResponse) GetItems() []*InventoryStatus {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *InventoryStatusResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type InventoryReservationResponse struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	OrderId               string                 `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	SuccessProcessedItems *SuccessProcessedItems `protobuf:"bytes,2,opt,name=success_processed_items,json=successProcessedItems,proto3" json:"success_processed_items,omitempty"`
	FailedProcessedItems  *FailedProcessedItems  `protobuf:"bytes,3,opt,name=failed_processed_items,json=failedProcessedItems,proto3" json:"failed_processed_items,omitempty"`
	Timestamp             *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *InventoryReservationResponse) Reset() {
	*x = InventoryReservationResponse{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryReservationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryReservationResponse) ProtoMessage() {}

func (x *InventoryReservationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryReservationResponse.ProtoReflect.Descriptor instead.
func (*InventoryReservationResponse) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{5}
}

func (x *InventoryReservationResponse) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *InventoryReservationResponse) GetSuccessProcessedItems() *SuccessProcessedItems {
	if x != nil {
		return x.SuccessProcessedItems
	}
	return nil
}

func (x *InventoryReservationResponse) GetFailedProcessedItems() *FailedProcessedItems {
	if x != nil {
		return x.FailedProcessedItems
	}
	return nil
}

func (x *InventoryReservationResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type ReservationHistory struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OrderId       string                 `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Sku           string                 `protobuf:"bytes,3,opt,name=sku,proto3" json:"sku,omitempty"`
	Quantity      float64                `protobuf:"fixed64,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Uom           string                 `protobuf:"bytes,5,opt,name=uom,proto3" json:"uom,omitempty"`
	Status        string                 `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	ReservedAt    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=reserved_at,json=reservedAt,proto3" json:"reserved_at,omitempty"`
	ReleasedAt    *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=released_at,json=releasedAt,proto3" json:"released_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReservationHistory) Reset() {
	*x = ReservationHistory{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReservationHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReservationHistory) ProtoMessage() {}

func (x *ReservationHistory) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReservationHistory.ProtoReflect.Descriptor instead.
func (*ReservationHistory) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{6}
}

func (x *ReservationHistory) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ReservationHistory) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *ReservationHistory) GetSku() string {
	if x != nil {
		return x.Sku
	}
	return ""
}

func (x *ReservationHistory) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *ReservationHistory) GetUom() string {
	if x != nil {
		return x.Uom
	}
	return ""
}

func (x *ReservationHistory) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ReservationHistory) GetReservedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ReservedAt
	}
	return nil
}

func (x *ReservationHistory) GetReleasedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ReleasedAt
	}
	return nil
}

type SuccessProcessedItems struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*ReservationHistory  `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SuccessProcessedItems) Reset() {
	*x = SuccessProcessedItems{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SuccessProcessedItems) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuccessProcessedItems) ProtoMessage() {}

func (x *SuccessProcessedItems) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuccessProcessedItems.ProtoReflect.Descriptor instead.
func (*SuccessProcessedItems) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{7}
}

func (x *SuccessProcessedItems) GetItems() []*ReservationHistory {
	if x != nil {
		return x.Items
	}
	return nil
}

type FailedProcessedItems struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*InventoryStatus     `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FailedProcessedItems) Reset() {
	*x = FailedProcessedItems{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FailedProcessedItems) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailedProcessedItems) ProtoMessage() {}

func (x *FailedProcessedItems) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailedProcessedItems.ProtoReflect.Descriptor instead.
func (*FailedProcessedItems) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{8}
}

func (x *FailedProcessedItems) GetItems() []*InventoryStatus {
	if x != nil {
		return x.Items
	}
	return nil
}

type ErrorDetails struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ErrorCode     ErrorCode              `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3,enum=pb_schemas.inventory.v1.ErrorCode" json:"error_code,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ErrorDetails) Reset() {
	*x = ErrorDetails{}
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ErrorDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorDetails) ProtoMessage() {}

func (x *ErrorDetails) ProtoReflect() protoreflect.Message {
	mi := &file_pb_schemas_inventory_v1_stock_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorDetails.ProtoReflect.Descriptor instead.
func (*ErrorDetails) Descriptor() ([]byte, []int) {
	return file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP(), []int{9}
}

func (x *ErrorDetails) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_UNDEFINED
}

func (x *ErrorDetails) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

var File_pb_schemas_inventory_v1_stock_proto protoreflect.FileDescriptor

const file_pb_schemas_inventory_v1_stock_proto_rawDesc = "" +
	"\n" +
	"#pb_schemas/inventory/v1/stock.proto\x12\x17pb_schemas.inventory.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"Z\n" +
	"\rInventoryItem\x12\x10\n" +
	"\x03sku\x18\x01 \x01(\tR\x03sku\x12%\n" +
	"\x0freq_qty_per_uom\x18\x02 \x01(\x01R\freqQtyPerUom\x12\x10\n" +
	"\x03uom\x18\x03 \x01(\tR\x03uom\"\xae\x02\n" +
	"\x0fInventoryStatus\x12\x10\n" +
	"\x03sku\x18\x01 \x01(\tR\x03sku\x12-\n" +
	"\x12requested_quantity\x18\x02 \x01(\x01R\x11requestedQuantity\x12-\n" +
	"\x12available_quantity\x18\x03 \x01(\x01R\x11availableQuantity\x12+\n" +
	"\x11reserved_quantity\x18\x04 \x01(\x01R\x10reservedQuantity\x12%\n" +
	"\x0etotal_quantity\x18\x05 \x01(\x01R\rtotalQuantity\x12\x17\n" +
	"\asku_uom\x18\x06 \x01(\tR\x06skuUom\x12\x1b\n" +
	"\tsku_price\x18\a \x01(\x01R\bskuPrice\x12!\n" +
	"\fsku_currency\x18\b \x01(\tR\vskuCurrency\"9\n" +
	"\fReservedItem\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\border_id\x18\x02 \x01(\tR\aorderId\"s\n" +
	"\x18StandardInventoryRequest\x12\x19\n" +
	"\border_id\x18\x01 \x01(\tR\aorderId\x12<\n" +
	"\x05items\x18\x02 \x03(\v2&.pb_schemas.inventory.v1.InventoryItemR\x05items\"\x93\x01\n" +
	"\x17InventoryStatusResponse\x12>\n" +
	"\x05items\x18\x01 \x03(\v2(.pb_schemas.inventory.v1.InventoryStatusR\x05items\x128\n" +
	"\ttimestamp\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\"\xc0\x02\n" +
	"\x1cInventoryReservationResponse\x12\x19\n" +
	"\border_id\x18\x01 \x01(\tR\aorderId\x12f\n" +
	"\x17success_processed_items\x18\x02 \x01(\v2..pb_schemas.inventory.v1.SuccessProcessedItemsR\x15successProcessedItems\x12c\n" +
	"\x16failed_processed_items\x18\x03 \x01(\v2-.pb_schemas.inventory.v1.FailedProcessedItemsR\x14failedProcessedItems\x128\n" +
	"\ttimestamp\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\"\x91\x02\n" +
	"\x12ReservationHistory\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\border_id\x18\x02 \x01(\tR\aorderId\x12\x10\n" +
	"\x03sku\x18\x03 \x01(\tR\x03sku\x12\x1a\n" +
	"\bquantity\x18\x04 \x01(\x01R\bquantity\x12\x10\n" +
	"\x03uom\x18\x05 \x01(\tR\x03uom\x12\x16\n" +
	"\x06status\x18\x06 \x01(\tR\x06status\x12;\n" +
	"\vreserved_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"reservedAt\x12;\n" +
	"\vreleased_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"releasedAt\"Z\n" +
	"\x15SuccessProcessedItems\x12A\n" +
	"\x05items\x18\x01 \x03(\v2+.pb_schemas.inventory.v1.ReservationHistoryR\x05items\"V\n" +
	"\x14FailedProcessedItems\x12>\n" +
	"\x05items\x18\x01 \x03(\v2(.pb_schemas.inventory.v1.InventoryStatusR\x05items\"v\n" +
	"\fErrorDetails\x12A\n" +
	"\n" +
	"error_code\x18\x01 \x01(\x0e2\".pb_schemas.inventory.v1.ErrorCodeR\terrorCode\x12#\n" +
	"\rerror_message\x18\x02 \x01(\tR\ferrorMessage*\xd7\x01\n" +
	"\tErrorCode\x12\r\n" +
	"\tUNDEFINED\x10\x00\x12\x11\n" +
	"\rSKU_NOT_FOUND\x10\x01\x12\x1a\n" +
	"\x16SKU_UOM_PAIR_NOT_MATCH\x10\x02\x12\x12\n" +
	"\x0eDB_UNAVAILABLE\x10\x03\x12\x18\n" +
	"\x14DB_ERROR_TRANSACTION\x10\x04\x12\x12\n" +
	"\x0eINTERNAL_ERROR\x10\x05\x12$\n" +
	" INSUFFICIENT_QUANTITY_TO_RESERVE\x10\x06\x12$\n" +
	" INSUFFICIENT_QUANTITY_TO_RELEASE\x10\a2\xff\x02\n" +
	"\x10InventoryService\x12s\n" +
	"\n" +
	"CheckStock\x121.pb_schemas.inventory.v1.StandardInventoryRequest\x1a0.pb_schemas.inventory.v1.InventoryStatusResponse\"\x00\x12z\n" +
	"\fReserveStock\x121.pb_schemas.inventory.v1.StandardInventoryRequest\x1a5.pb_schemas.inventory.v1.InventoryReservationResponse\"\x00\x12z\n" +
	"\fReleaseStock\x121.pb_schemas.inventory.v1.StandardInventoryRequest\x1a5.pb_schemas.inventory.v1.InventoryReservationResponse\"\x00B3Z1ops-monorepo/protogen/go/inventory/v1;inventoryv1b\x06proto3"

var (
	file_pb_schemas_inventory_v1_stock_proto_rawDescOnce sync.Once
	file_pb_schemas_inventory_v1_stock_proto_rawDescData []byte
)

func file_pb_schemas_inventory_v1_stock_proto_rawDescGZIP() []byte {
	file_pb_schemas_inventory_v1_stock_proto_rawDescOnce.Do(func() {
		file_pb_schemas_inventory_v1_stock_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pb_schemas_inventory_v1_stock_proto_rawDesc), len(file_pb_schemas_inventory_v1_stock_proto_rawDesc)))
	})
	return file_pb_schemas_inventory_v1_stock_proto_rawDescData
}

var file_pb_schemas_inventory_v1_stock_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_schemas_inventory_v1_stock_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_pb_schemas_inventory_v1_stock_proto_goTypes = []any{
	(ErrorCode)(0),                       // 0: pb_schemas.inventory.v1.ErrorCode
	(*InventoryItem)(nil),                // 1: pb_schemas.inventory.v1.InventoryItem
	(*InventoryStatus)(nil),              // 2: pb_schemas.inventory.v1.InventoryStatus
	(*ReservedItem)(nil),                 // 3: pb_schemas.inventory.v1.ReservedItem
	(*StandardInventoryRequest)(nil),     // 4: pb_schemas.inventory.v1.StandardInventoryRequest
	(*InventoryStatusResponse)(nil),      // 5: pb_schemas.inventory.v1.InventoryStatusResponse
	(*InventoryReservationResponse)(nil), // 6: pb_schemas.inventory.v1.InventoryReservationResponse
	(*ReservationHistory)(nil),           // 7: pb_schemas.inventory.v1.ReservationHistory
	(*SuccessProcessedItems)(nil),        // 8: pb_schemas.inventory.v1.SuccessProcessedItems
	(*FailedProcessedItems)(nil),         // 9: pb_schemas.inventory.v1.FailedProcessedItems
	(*ErrorDetails)(nil),                 // 10: pb_schemas.inventory.v1.ErrorDetails
	(*timestamppb.Timestamp)(nil),        // 11: google.protobuf.Timestamp
}
var file_pb_schemas_inventory_v1_stock_proto_depIdxs = []int32{
	1,  // 0: pb_schemas.inventory.v1.StandardInventoryRequest.items:type_name -> pb_schemas.inventory.v1.InventoryItem
	2,  // 1: pb_schemas.inventory.v1.InventoryStatusResponse.items:type_name -> pb_schemas.inventory.v1.InventoryStatus
	11, // 2: pb_schemas.inventory.v1.InventoryStatusResponse.timestamp:type_name -> google.protobuf.Timestamp
	8,  // 3: pb_schemas.inventory.v1.InventoryReservationResponse.success_processed_items:type_name -> pb_schemas.inventory.v1.SuccessProcessedItems
	9,  // 4: pb_schemas.inventory.v1.InventoryReservationResponse.failed_processed_items:type_name -> pb_schemas.inventory.v1.FailedProcessedItems
	11, // 5: pb_schemas.inventory.v1.InventoryReservationResponse.timestamp:type_name -> google.protobuf.Timestamp
	11, // 6: pb_schemas.inventory.v1.ReservationHistory.reserved_at:type_name -> google.protobuf.Timestamp
	11, // 7: pb_schemas.inventory.v1.ReservationHistory.released_at:type_name -> google.protobuf.Timestamp
	7,  // 8: pb_schemas.inventory.v1.SuccessProcessedItems.items:type_name -> pb_schemas.inventory.v1.ReservationHistory
	2,  // 9: pb_schemas.inventory.v1.FailedProcessedItems.items:type_name -> pb_schemas.inventory.v1.InventoryStatus
	0,  // 10: pb_schemas.inventory.v1.ErrorDetails.error_code:type_name -> pb_schemas.inventory.v1.ErrorCode
	4,  // 11: pb_schemas.inventory.v1.InventoryService.CheckStock:input_type -> pb_schemas.inventory.v1.StandardInventoryRequest
	4,  // 12: pb_schemas.inventory.v1.InventoryService.ReserveStock:input_type -> pb_schemas.inventory.v1.StandardInventoryRequest
	4,  // 13: pb_schemas.inventory.v1.InventoryService.ReleaseStock:input_type -> pb_schemas.inventory.v1.StandardInventoryRequest
	5,  // 14: pb_schemas.inventory.v1.InventoryService.CheckStock:output_type -> pb_schemas.inventory.v1.InventoryStatusResponse
	6,  // 15: pb_schemas.inventory.v1.InventoryService.ReserveStock:output_type -> pb_schemas.inventory.v1.InventoryReservationResponse
	6,  // 16: pb_schemas.inventory.v1.InventoryService.ReleaseStock:output_type -> pb_schemas.inventory.v1.InventoryReservationResponse
	14, // [14:17] is the sub-list for method output_type
	11, // [11:14] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pb_schemas_inventory_v1_stock_proto_init() }
func file_pb_schemas_inventory_v1_stock_proto_init() {
	if File_pb_schemas_inventory_v1_stock_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pb_schemas_inventory_v1_stock_proto_rawDesc), len(file_pb_schemas_inventory_v1_stock_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_schemas_inventory_v1_stock_proto_goTypes,
		DependencyIndexes: file_pb_schemas_inventory_v1_stock_proto_depIdxs,
		EnumInfos:         file_pb_schemas_inventory_v1_stock_proto_enumTypes,
		MessageInfos:      file_pb_schemas_inventory_v1_stock_proto_msgTypes,
	}.Build()
	File_pb_schemas_inventory_v1_stock_proto = out.File
	file_pb_schemas_inventory_v1_stock_proto_goTypes = nil
	file_pb_schemas_inventory_v1_stock_proto_depIdxs = nil
}
