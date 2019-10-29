// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fetch

import (
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
)

type ReduceRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsReduceRequest(buf []byte, offset flatbuffers.UOffsetT) *ReduceRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ReduceRequest{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ReduceRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ReduceRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ReduceRequest) Label() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReduceRequest) Method() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReduceRequest) MethodParams(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *ReduceRequest) MethodParamsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ReduceRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ReduceRequestAddLabel(builder *flatbuffers.Builder, label flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(label), 0)
}
func ReduceRequestAddMethod(builder *flatbuffers.Builder, method flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(method), 0)
}
func ReduceRequestAddMethodParams(builder *flatbuffers.Builder, methodParams flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(methodParams), 0)
}
func ReduceRequestStartMethodParamsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ReduceRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type StreamRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsStreamRequest(buf []byte, offset flatbuffers.UOffsetT) *StreamRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StreamRequest{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *StreamRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StreamRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *StreamRequest) CheckUuid(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *StreamRequest) CheckUuidLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *StreamRequest) CheckUuidBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *StreamRequest) MutateCheckUuid(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *StreamRequest) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *StreamRequest) Kind() Kind {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return Kind(rcv._tab.GetInt8(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *StreamRequest) MutateKind(n Kind) bool {
	return rcv._tab.MutateInt8Slot(8, int8(n))
}

func (rcv *StreamRequest) Transform() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *StreamRequest) TransformParams(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *StreamRequest) TransformParamsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *StreamRequest) Label() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func StreamRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func StreamRequestAddCheckUuid(builder *flatbuffers.Builder, checkUuid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(checkUuid), 0)
}
func StreamRequestStartCheckUuidVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func StreamRequestAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(name), 0)
}
func StreamRequestAddKind(builder *flatbuffers.Builder, kind Kind) {
	builder.PrependInt8Slot(2, int8(kind), 0)
}
func StreamRequestAddTransform(builder *flatbuffers.Builder, transform flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(transform), 0)
}
func StreamRequestAddTransformParams(builder *flatbuffers.Builder, transformParams flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(transformParams), 0)
}
func StreamRequestStartTransformParamsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func StreamRequestAddLabel(builder *flatbuffers.Builder, label flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(label), 0)
}
func StreamRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Fetch struct {
	_tab flatbuffers.Table
}

func GetRootAsFetch(buf []byte, offset flatbuffers.UOffsetT) *Fetch {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Fetch{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Fetch) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Fetch) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Fetch) StartMs() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fetch) MutateStartMs(n uint64) bool {
	return rcv._tab.MutateUint64Slot(4, n)
}

func (rcv *Fetch) PeriodMs() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fetch) MutatePeriodMs(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *Fetch) Count() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fetch) MutateCount(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *Fetch) Streams(obj *StreamRequest, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Fetch) StreamsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Fetch) Reduce(obj *ReduceRequest, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Fetch) ReduceLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func FetchStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func FetchAddStartMs(builder *flatbuffers.Builder, startMs uint64) {
	builder.PrependUint64Slot(0, startMs, 0)
}
func FetchAddPeriodMs(builder *flatbuffers.Builder, periodMs uint32) {
	builder.PrependUint32Slot(1, periodMs, 0)
}
func FetchAddCount(builder *flatbuffers.Builder, count uint32) {
	builder.PrependUint32Slot(2, count, 0)
}
func FetchAddStreams(builder *flatbuffers.Builder, streams flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(streams), 0)
}
func FetchStartStreamsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FetchAddReduce(builder *flatbuffers.Builder, reduce flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(reduce), 0)
}
func FetchStartReduceVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FetchEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
