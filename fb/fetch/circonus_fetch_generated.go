// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fetch

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ReduceRequestT struct {
	Label        string
	Method       string
	MethodParams []string
}

func ReduceRequestPack(builder *flatbuffers.Builder, t *ReduceRequestT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	labelOffset := builder.CreateString(t.Label)
	methodOffset := builder.CreateString(t.Method)
	methodParamsOffset := flatbuffers.UOffsetT(0)
	if t.MethodParams != nil {
		methodParamsLength := len(t.MethodParams)
		methodParamsOffsets := make([]flatbuffers.UOffsetT, methodParamsLength)
		for j := 0; j < methodParamsLength; j++ {
			methodParamsOffsets[j] = builder.CreateString(t.MethodParams[j])
		}
		ReduceRequestStartMethodParamsVector(builder, methodParamsLength)
		for j := methodParamsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(methodParamsOffsets[j])
		}
		methodParamsOffset = builder.EndVector(methodParamsLength)
	}
	ReduceRequestStart(builder)
	ReduceRequestAddLabel(builder, labelOffset)
	ReduceRequestAddMethod(builder, methodOffset)
	ReduceRequestAddMethodParams(builder, methodParamsOffset)
	return ReduceRequestEnd(builder)
}

func (rcv *ReduceRequest) UnPack() *ReduceRequestT {
	if rcv == nil {
		return nil
	}
	t := &ReduceRequestT{}
	t.Label = string(rcv.Label())
	t.Method = string(rcv.Method())
	methodParamsLength := rcv.MethodParamsLength()
	t.MethodParams = make([]string, methodParamsLength)
	for j := 0; j < methodParamsLength; j++ {
		t.MethodParams[j] = string(rcv.MethodParams(j))
	}
	return t
}

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

type StreamRequestT struct {
	CheckUuid       []byte
	Name            string
	Kind            Kind
	Transform       string
	TransformParams []string
	Label           string
}

func StreamRequestPack(builder *flatbuffers.Builder, t *StreamRequestT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	checkUuidOffset := flatbuffers.UOffsetT(0)
	if t.CheckUuid != nil {
		checkUuidOffset = builder.CreateByteString(t.CheckUuid)
	}
	nameOffset := builder.CreateString(t.Name)
	transformOffset := builder.CreateString(t.Transform)
	transformParamsOffset := flatbuffers.UOffsetT(0)
	if t.TransformParams != nil {
		transformParamsLength := len(t.TransformParams)
		transformParamsOffsets := make([]flatbuffers.UOffsetT, transformParamsLength)
		for j := 0; j < transformParamsLength; j++ {
			transformParamsOffsets[j] = builder.CreateString(t.TransformParams[j])
		}
		StreamRequestStartTransformParamsVector(builder, transformParamsLength)
		for j := transformParamsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(transformParamsOffsets[j])
		}
		transformParamsOffset = builder.EndVector(transformParamsLength)
	}
	labelOffset := builder.CreateString(t.Label)
	StreamRequestStart(builder)
	StreamRequestAddCheckUuid(builder, checkUuidOffset)
	StreamRequestAddName(builder, nameOffset)
	StreamRequestAddKind(builder, t.Kind)
	StreamRequestAddTransform(builder, transformOffset)
	StreamRequestAddTransformParams(builder, transformParamsOffset)
	StreamRequestAddLabel(builder, labelOffset)
	return StreamRequestEnd(builder)
}

func (rcv *StreamRequest) UnPack() *StreamRequestT {
	if rcv == nil {
		return nil
	}
	t := &StreamRequestT{}
	t.CheckUuid = rcv.CheckUuidBytes()
	t.Name = string(rcv.Name())
	t.Kind = rcv.Kind()
	t.Transform = string(rcv.Transform())
	transformParamsLength := rcv.TransformParamsLength()
	t.TransformParams = make([]string, transformParamsLength)
	for j := 0; j < transformParamsLength; j++ {
		t.TransformParams[j] = string(rcv.TransformParams(j))
	}
	t.Label = string(rcv.Label())
	return t
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

type FetchT struct {
	StartMs  uint64
	PeriodMs uint32
	Count    uint32
	Streams  []*StreamRequestT
	Reduce   []*ReduceRequestT
}

func FetchPack(builder *flatbuffers.Builder, t *FetchT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	streamsOffset := flatbuffers.UOffsetT(0)
	if t.Streams != nil {
		streamsLength := len(t.Streams)
		streamsOffsets := make([]flatbuffers.UOffsetT, streamsLength)
		for j := 0; j < streamsLength; j++ {
			streamsOffsets[j] = StreamRequestPack(builder, t.Streams[j])
		}
		FetchStartStreamsVector(builder, streamsLength)
		for j := streamsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(streamsOffsets[j])
		}
		streamsOffset = builder.EndVector(streamsLength)
	}
	reduceOffset := flatbuffers.UOffsetT(0)
	if t.Reduce != nil {
		reduceLength := len(t.Reduce)
		reduceOffsets := make([]flatbuffers.UOffsetT, reduceLength)
		for j := 0; j < reduceLength; j++ {
			reduceOffsets[j] = ReduceRequestPack(builder, t.Reduce[j])
		}
		FetchStartReduceVector(builder, reduceLength)
		for j := reduceLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(reduceOffsets[j])
		}
		reduceOffset = builder.EndVector(reduceLength)
	}
	FetchStart(builder)
	FetchAddStartMs(builder, t.StartMs)
	FetchAddPeriodMs(builder, t.PeriodMs)
	FetchAddCount(builder, t.Count)
	FetchAddStreams(builder, streamsOffset)
	FetchAddReduce(builder, reduceOffset)
	return FetchEnd(builder)
}

func (rcv *Fetch) UnPack() *FetchT {
	if rcv == nil {
		return nil
	}
	t := &FetchT{}
	t.StartMs = rcv.StartMs()
	t.PeriodMs = rcv.PeriodMs()
	t.Count = rcv.Count()
	streamsLength := rcv.StreamsLength()
	t.Streams = make([]*StreamRequestT, streamsLength)
	for j := 0; j < streamsLength; j++ {
		x := StreamRequest{}
		rcv.Streams(&x, j)
		t.Streams[j] = x.UnPack()
	}
	reduceLength := rcv.ReduceLength()
	t.Reduce = make([]*ReduceRequestT, reduceLength)
	for j := 0; j < reduceLength; j++ {
		x := ReduceRequest{}
		rcv.Reduce(&x, j)
		t.Reduce[j] = x.UnPack()
	}
	return t
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
