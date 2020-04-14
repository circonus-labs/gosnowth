// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package noit

import (
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
)

type MetricValueUnionT struct {
	Type  MetricValueUnion
	Value interface{}
}

func MetricValueUnionPack(builder *flatbuffers.Builder, t *MetricValueUnionT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	switch t.Type {
	case MetricValueUnionIntValue:
		return IntValuePack(builder, t.Value.(*IntValueT))
	case MetricValueUnionUintValue:
		return UintValuePack(builder, t.Value.(*UintValueT))
	case MetricValueUnionLongValue:
		return LongValuePack(builder, t.Value.(*LongValueT))
	case MetricValueUnionUlongValue:
		return UlongValuePack(builder, t.Value.(*UlongValueT))
	case MetricValueUnionDoubleValue:
		return DoubleValuePack(builder, t.Value.(*DoubleValueT))
	case MetricValueUnionStringValue:
		return StringValuePack(builder, t.Value.(*StringValueT))
	case MetricValueUnionHistogram:
		return HistogramPack(builder, t.Value.(*HistogramT))
	case MetricValueUnionAbsentNumericValue:
		return AbsentNumericValuePack(builder, t.Value.(*AbsentNumericValueT))
	case MetricValueUnionAbsentStringValue:
		return AbsentStringValuePack(builder, t.Value.(*AbsentStringValueT))
	case MetricValueUnionAbsentHistogramValue:
		return AbsentHistogramValuePack(builder, t.Value.(*AbsentHistogramValueT))
	}
	return 0
}

func MetricValueUnionUnPack(t MetricValueUnion, table flatbuffers.Table) *MetricValueUnionT {
	switch t {
	case MetricValueUnionIntValue:
		x := IntValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionIntValue, Value: x.UnPack()}
	case MetricValueUnionUintValue:
		x := UintValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionUintValue, Value: x.UnPack()}
	case MetricValueUnionLongValue:
		x := LongValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionLongValue, Value: x.UnPack()}
	case MetricValueUnionUlongValue:
		x := UlongValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionUlongValue, Value: x.UnPack()}
	case MetricValueUnionDoubleValue:
		x := DoubleValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionDoubleValue, Value: x.UnPack()}
	case MetricValueUnionStringValue:
		x := StringValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionStringValue, Value: x.UnPack()}
	case MetricValueUnionHistogram:
		x := Histogram{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionHistogram, Value: x.UnPack()}
	case MetricValueUnionAbsentNumericValue:
		x := AbsentNumericValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionAbsentNumericValue, Value: x.UnPack()}
	case MetricValueUnionAbsentStringValue:
		x := AbsentStringValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionAbsentStringValue, Value: x.UnPack()}
	case MetricValueUnionAbsentHistogramValue:
		x := AbsentHistogramValue{_tab: table}
		return &MetricValueUnionT{Type: MetricValueUnionAbsentHistogramValue, Value: x.UnPack()}
	}
	return nil
}

type MetricValueUnion byte

const (
	MetricValueUnionNONE                 MetricValueUnion = 0
	MetricValueUnionIntValue             MetricValueUnion = 1
	MetricValueUnionUintValue            MetricValueUnion = 2
	MetricValueUnionLongValue            MetricValueUnion = 3
	MetricValueUnionUlongValue           MetricValueUnion = 4
	MetricValueUnionDoubleValue          MetricValueUnion = 5
	MetricValueUnionStringValue          MetricValueUnion = 6
	MetricValueUnionHistogram            MetricValueUnion = 7
	MetricValueUnionAbsentNumericValue   MetricValueUnion = 8
	MetricValueUnionAbsentStringValue    MetricValueUnion = 9
	MetricValueUnionAbsentHistogramValue MetricValueUnion = 10
)

var EnumNamesMetricValueUnion = map[MetricValueUnion]string{
	MetricValueUnionNONE:                 "NONE",
	MetricValueUnionIntValue:             "IntValue",
	MetricValueUnionUintValue:            "UintValue",
	MetricValueUnionLongValue:            "LongValue",
	MetricValueUnionUlongValue:           "UlongValue",
	MetricValueUnionDoubleValue:          "DoubleValue",
	MetricValueUnionStringValue:          "StringValue",
	MetricValueUnionHistogram:            "Histogram",
	MetricValueUnionAbsentNumericValue:   "AbsentNumericValue",
	MetricValueUnionAbsentStringValue:    "AbsentStringValue",
	MetricValueUnionAbsentHistogramValue: "AbsentHistogramValue",
}

var EnumValuesMetricValueUnion = map[string]MetricValueUnion{
	"NONE":                 MetricValueUnionNONE,
	"IntValue":             MetricValueUnionIntValue,
	"UintValue":            MetricValueUnionUintValue,
	"LongValue":            MetricValueUnionLongValue,
	"UlongValue":           MetricValueUnionUlongValue,
	"DoubleValue":          MetricValueUnionDoubleValue,
	"StringValue":          MetricValueUnionStringValue,
	"Histogram":            MetricValueUnionHistogram,
	"AbsentNumericValue":   MetricValueUnionAbsentNumericValue,
	"AbsentStringValue":    MetricValueUnionAbsentStringValue,
	"AbsentHistogramValue": MetricValueUnionAbsentHistogramValue,
}

func (v MetricValueUnion) String() string {
	if s, ok := EnumNamesMetricValueUnion[v]; ok {
		return s
	}
	return "MetricValueUnion(" + strconv.FormatInt(int64(v), 10) + ")"
}

type IntValueT struct {
	Value int32
}

func IntValuePack(builder *flatbuffers.Builder, t *IntValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	IntValueStart(builder)
	IntValueAddValue(builder, t.Value)
	return IntValueEnd(builder)
}

func (rcv *IntValue) UnPack() *IntValueT {
	if rcv == nil {
		return nil
	}
	t := &IntValueT{}
	t.Value = rcv.Value()
	return t
}

type IntValue struct {
	_tab flatbuffers.Table
}

func GetRootAsIntValue(buf []byte, offset flatbuffers.UOffsetT) *IntValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &IntValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *IntValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IntValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *IntValue) Value() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IntValue) MutateValue(n int32) bool {
	return rcv._tab.MutateInt32Slot(4, n)
}

func IntValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func IntValueAddValue(builder *flatbuffers.Builder, value int32) {
	builder.PrependInt32Slot(0, value, 0)
}
func IntValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type UintValueT struct {
	Value uint32
}

func UintValuePack(builder *flatbuffers.Builder, t *UintValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	UintValueStart(builder)
	UintValueAddValue(builder, t.Value)
	return UintValueEnd(builder)
}

func (rcv *UintValue) UnPack() *UintValueT {
	if rcv == nil {
		return nil
	}
	t := &UintValueT{}
	t.Value = rcv.Value()
	return t
}

type UintValue struct {
	_tab flatbuffers.Table
}

func GetRootAsUintValue(buf []byte, offset flatbuffers.UOffsetT) *UintValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &UintValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *UintValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UintValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *UintValue) Value() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UintValue) MutateValue(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func UintValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func UintValueAddValue(builder *flatbuffers.Builder, value uint32) {
	builder.PrependUint32Slot(0, value, 0)
}
func UintValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type LongValueT struct {
	Value int64
}

func LongValuePack(builder *flatbuffers.Builder, t *LongValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	LongValueStart(builder)
	LongValueAddValue(builder, t.Value)
	return LongValueEnd(builder)
}

func (rcv *LongValue) UnPack() *LongValueT {
	if rcv == nil {
		return nil
	}
	t := &LongValueT{}
	t.Value = rcv.Value()
	return t
}

type LongValue struct {
	_tab flatbuffers.Table
}

func GetRootAsLongValue(buf []byte, offset flatbuffers.UOffsetT) *LongValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LongValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LongValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LongValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LongValue) Value() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LongValue) MutateValue(n int64) bool {
	return rcv._tab.MutateInt64Slot(4, n)
}

func LongValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func LongValueAddValue(builder *flatbuffers.Builder, value int64) {
	builder.PrependInt64Slot(0, value, 0)
}
func LongValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type UlongValueT struct {
	Value uint64
}

func UlongValuePack(builder *flatbuffers.Builder, t *UlongValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	UlongValueStart(builder)
	UlongValueAddValue(builder, t.Value)
	return UlongValueEnd(builder)
}

func (rcv *UlongValue) UnPack() *UlongValueT {
	if rcv == nil {
		return nil
	}
	t := &UlongValueT{}
	t.Value = rcv.Value()
	return t
}

type UlongValue struct {
	_tab flatbuffers.Table
}

func GetRootAsUlongValue(buf []byte, offset flatbuffers.UOffsetT) *UlongValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &UlongValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *UlongValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UlongValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *UlongValue) Value() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UlongValue) MutateValue(n uint64) bool {
	return rcv._tab.MutateUint64Slot(4, n)
}

func UlongValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func UlongValueAddValue(builder *flatbuffers.Builder, value uint64) {
	builder.PrependUint64Slot(0, value, 0)
}
func UlongValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type DoubleValueT struct {
	Value float64
}

func DoubleValuePack(builder *flatbuffers.Builder, t *DoubleValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	DoubleValueStart(builder)
	DoubleValueAddValue(builder, t.Value)
	return DoubleValueEnd(builder)
}

func (rcv *DoubleValue) UnPack() *DoubleValueT {
	if rcv == nil {
		return nil
	}
	t := &DoubleValueT{}
	t.Value = rcv.Value()
	return t
}

type DoubleValue struct {
	_tab flatbuffers.Table
}

func GetRootAsDoubleValue(buf []byte, offset flatbuffers.UOffsetT) *DoubleValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DoubleValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *DoubleValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DoubleValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DoubleValue) Value() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *DoubleValue) MutateValue(n float64) bool {
	return rcv._tab.MutateFloat64Slot(4, n)
}

func DoubleValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func DoubleValueAddValue(builder *flatbuffers.Builder, value float64) {
	builder.PrependFloat64Slot(0, value, 0.0)
}
func DoubleValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type StringValueT struct {
	Value string
}

func StringValuePack(builder *flatbuffers.Builder, t *StringValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	valueOffset := builder.CreateString(t.Value)
	StringValueStart(builder)
	StringValueAddValue(builder, valueOffset)
	return StringValueEnd(builder)
}

func (rcv *StringValue) UnPack() *StringValueT {
	if rcv == nil {
		return nil
	}
	t := &StringValueT{}
	t.Value = string(rcv.Value())
	return t
}

type StringValue struct {
	_tab flatbuffers.Table
}

func GetRootAsStringValue(buf []byte, offset flatbuffers.UOffsetT) *StringValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StringValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *StringValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StringValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *StringValue) Value() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func StringValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func StringValueAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(value), 0)
}
func StringValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type AbsentNumericValueT struct {
}

func AbsentNumericValuePack(builder *flatbuffers.Builder, t *AbsentNumericValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	AbsentNumericValueStart(builder)
	return AbsentNumericValueEnd(builder)
}

func (rcv *AbsentNumericValue) UnPack() *AbsentNumericValueT {
	if rcv == nil {
		return nil
	}
	t := &AbsentNumericValueT{}
	return t
}

type AbsentNumericValue struct {
	_tab flatbuffers.Table
}

func GetRootAsAbsentNumericValue(buf []byte, offset flatbuffers.UOffsetT) *AbsentNumericValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AbsentNumericValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AbsentNumericValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AbsentNumericValue) Table() flatbuffers.Table {
	return rcv._tab
}

func AbsentNumericValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func AbsentNumericValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type AbsentStringValueT struct {
}

func AbsentStringValuePack(builder *flatbuffers.Builder, t *AbsentStringValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	AbsentStringValueStart(builder)
	return AbsentStringValueEnd(builder)
}

func (rcv *AbsentStringValue) UnPack() *AbsentStringValueT {
	if rcv == nil {
		return nil
	}
	t := &AbsentStringValueT{}
	return t
}

type AbsentStringValue struct {
	_tab flatbuffers.Table
}

func GetRootAsAbsentStringValue(buf []byte, offset flatbuffers.UOffsetT) *AbsentStringValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AbsentStringValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AbsentStringValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AbsentStringValue) Table() flatbuffers.Table {
	return rcv._tab
}

func AbsentStringValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func AbsentStringValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type AbsentHistogramValueT struct {
}

func AbsentHistogramValuePack(builder *flatbuffers.Builder, t *AbsentHistogramValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	AbsentHistogramValueStart(builder)
	return AbsentHistogramValueEnd(builder)
}

func (rcv *AbsentHistogramValue) UnPack() *AbsentHistogramValueT {
	if rcv == nil {
		return nil
	}
	t := &AbsentHistogramValueT{}
	return t
}

type AbsentHistogramValue struct {
	_tab flatbuffers.Table
}

func GetRootAsAbsentHistogramValue(buf []byte, offset flatbuffers.UOffsetT) *AbsentHistogramValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AbsentHistogramValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AbsentHistogramValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AbsentHistogramValue) Table() flatbuffers.Table {
	return rcv._tab
}

func AbsentHistogramValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func AbsentHistogramValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type HistogramBucketT struct {
	Val   int8
	Exp   int8
	Count uint64
}

func HistogramBucketPack(builder *flatbuffers.Builder, t *HistogramBucketT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	HistogramBucketStart(builder)
	HistogramBucketAddVal(builder, t.Val)
	HistogramBucketAddExp(builder, t.Exp)
	HistogramBucketAddCount(builder, t.Count)
	return HistogramBucketEnd(builder)
}

func (rcv *HistogramBucket) UnPack() *HistogramBucketT {
	if rcv == nil {
		return nil
	}
	t := &HistogramBucketT{}
	t.Val = rcv.Val()
	t.Exp = rcv.Exp()
	t.Count = rcv.Count()
	return t
}

type HistogramBucket struct {
	_tab flatbuffers.Table
}

func GetRootAsHistogramBucket(buf []byte, offset flatbuffers.UOffsetT) *HistogramBucket {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &HistogramBucket{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *HistogramBucket) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *HistogramBucket) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *HistogramBucket) Val() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *HistogramBucket) MutateVal(n int8) bool {
	return rcv._tab.MutateInt8Slot(4, n)
}

func (rcv *HistogramBucket) Exp() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *HistogramBucket) MutateExp(n int8) bool {
	return rcv._tab.MutateInt8Slot(6, n)
}

func (rcv *HistogramBucket) Count() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *HistogramBucket) MutateCount(n uint64) bool {
	return rcv._tab.MutateUint64Slot(8, n)
}

func HistogramBucketStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func HistogramBucketAddVal(builder *flatbuffers.Builder, val int8) {
	builder.PrependInt8Slot(0, val, 0)
}
func HistogramBucketAddExp(builder *flatbuffers.Builder, exp int8) {
	builder.PrependInt8Slot(1, exp, 0)
}
func HistogramBucketAddCount(builder *flatbuffers.Builder, count uint64) {
	builder.PrependUint64Slot(2, count, 0)
}
func HistogramBucketEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type HistogramT struct {
	Buckets    []*HistogramBucketT
	Cumulative bool
}

func HistogramPack(builder *flatbuffers.Builder, t *HistogramT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	bucketsOffset := flatbuffers.UOffsetT(0)
	if t.Buckets != nil {
		bucketsLength := len(t.Buckets)
		bucketsOffsets := make([]flatbuffers.UOffsetT, bucketsLength)
		for j := 0; j < bucketsLength; j++ {
			bucketsOffsets[j] = HistogramBucketPack(builder, t.Buckets[j])
		}
		HistogramStartBucketsVector(builder, bucketsLength)
		for j := bucketsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(bucketsOffsets[j])
		}
		bucketsOffset = builder.EndVector(bucketsLength)
	}
	HistogramStart(builder)
	HistogramAddBuckets(builder, bucketsOffset)
	HistogramAddCumulative(builder, t.Cumulative)
	return HistogramEnd(builder)
}

func (rcv *Histogram) UnPack() *HistogramT {
	if rcv == nil {
		return nil
	}
	t := &HistogramT{}
	bucketsLength := rcv.BucketsLength()
	t.Buckets = make([]*HistogramBucketT, bucketsLength)
	for j := 0; j < bucketsLength; j++ {
		x := HistogramBucket{}
		rcv.Buckets(&x, j)
		t.Buckets[j] = x.UnPack()
	}
	t.Cumulative = rcv.Cumulative()
	return t
}

type Histogram struct {
	_tab flatbuffers.Table
}

func GetRootAsHistogram(buf []byte, offset flatbuffers.UOffsetT) *Histogram {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Histogram{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Histogram) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Histogram) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Histogram) Buckets(obj *HistogramBucket, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Histogram) BucketsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Histogram) Cumulative() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Histogram) MutateCumulative(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func HistogramStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func HistogramAddBuckets(builder *flatbuffers.Builder, buckets flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(buckets), 0)
}
func HistogramStartBucketsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func HistogramAddCumulative(builder *flatbuffers.Builder, cumulative bool) {
	builder.PrependBoolSlot(1, cumulative, false)
}
func HistogramEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type MetricValueT struct {
	Name       string
	Timestamp  uint64
	Value      *MetricValueUnionT
	Generation int16
	StreamTags []string
}

func MetricValuePack(builder *flatbuffers.Builder, t *MetricValueT) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	nameOffset := builder.CreateString(t.Name)
	valueOffset := MetricValueUnionPack(builder, t.Value)

	streamTagsOffset := flatbuffers.UOffsetT(0)
	if t.StreamTags != nil {
		streamTagsLength := len(t.StreamTags)
		streamTagsOffsets := make([]flatbuffers.UOffsetT, streamTagsLength)
		for j := 0; j < streamTagsLength; j++ {
			streamTagsOffsets[j] = builder.CreateString(t.StreamTags[j])
		}
		MetricValueStartStreamTagsVector(builder, streamTagsLength)
		for j := streamTagsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(streamTagsOffsets[j])
		}
		streamTagsOffset = builder.EndVector(streamTagsLength)
	}
	MetricValueStart(builder)
	MetricValueAddName(builder, nameOffset)
	MetricValueAddTimestamp(builder, t.Timestamp)
	if t.Value != nil {
		MetricValueAddValueType(builder, t.Value.Type)
	}
	MetricValueAddValue(builder, valueOffset)
	MetricValueAddGeneration(builder, t.Generation)
	MetricValueAddStreamTags(builder, streamTagsOffset)
	return MetricValueEnd(builder)
}

func (rcv *MetricValue) UnPack() *MetricValueT {
	if rcv == nil {
		return nil
	}
	t := &MetricValueT{}
	t.Name = string(rcv.Name())
	t.Timestamp = rcv.Timestamp()
	valueTable := flatbuffers.Table{}
	if rcv.Value(&valueTable) {
		t.Value = MetricValueUnionUnPack(rcv.ValueType(), valueTable)
	}
	t.Generation = rcv.Generation()
	streamTagsLength := rcv.StreamTagsLength()
	t.StreamTags = make([]string, streamTagsLength)
	for j := 0; j < streamTagsLength; j++ {
		t.StreamTags[j] = string(rcv.StreamTags(j))
	}
	return t
}

type MetricValue struct {
	_tab flatbuffers.Table
}

func GetRootAsMetricValue(buf []byte, offset flatbuffers.UOffsetT) *MetricValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MetricValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *MetricValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MetricValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MetricValue) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetricValue) Timestamp() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetricValue) MutateTimestamp(n uint64) bool {
	return rcv._tab.MutateUint64Slot(6, n)
}

func (rcv *MetricValue) ValueType() MetricValueUnion {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return MetricValueUnion(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *MetricValue) MutateValueType(n MetricValueUnion) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *MetricValue) Value(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func (rcv *MetricValue) Generation() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetricValue) MutateGeneration(n int16) bool {
	return rcv._tab.MutateInt16Slot(12, n)
}

func (rcv *MetricValue) StreamTags(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *MetricValue) StreamTagsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func MetricValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func MetricValueAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func MetricValueAddTimestamp(builder *flatbuffers.Builder, timestamp uint64) {
	builder.PrependUint64Slot(1, timestamp, 0)
}
func MetricValueAddValueType(builder *flatbuffers.Builder, valueType MetricValueUnion) {
	builder.PrependByteSlot(2, byte(valueType), 0)
}
func MetricValueAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(value), 0)
}
func MetricValueAddGeneration(builder *flatbuffers.Builder, generation int16) {
	builder.PrependInt16Slot(4, generation, 0)
}
func MetricValueAddStreamTags(builder *flatbuffers.Builder, streamTags flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(streamTags), 0)
}
func MetricValueStartStreamTagsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MetricValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
