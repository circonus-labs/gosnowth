// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package noit

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type MetricBatchT struct {
	Timestamp uint64
	CheckName string
	CheckUuid string
	AccountId int32
	Metrics []*MetricValueT
}

func MetricBatchPack(builder *flatbuffers.Builder, t *MetricBatchT) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	checkNameOffset := builder.CreateString(t.CheckName)
	checkUuidOffset := builder.CreateString(t.CheckUuid)
	metricsOffset := flatbuffers.UOffsetT(0)
	if t.Metrics != nil {
		metricsLength := len(t.Metrics)
		metricsOffsets := make([]flatbuffers.UOffsetT, metricsLength)
		for j := 0; j < metricsLength; j++ {
			metricsOffsets[j] = MetricValuePack(builder, t.Metrics[j])
		}
		MetricBatchStartMetricsVector(builder, metricsLength)
		for j := metricsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(metricsOffsets[j])
		}
		metricsOffset = builder.EndVector(metricsLength)
	}
	MetricBatchStart(builder)
	MetricBatchAddTimestamp(builder, t.Timestamp)
	MetricBatchAddCheckName(builder, checkNameOffset)
	MetricBatchAddCheckUuid(builder, checkUuidOffset)
	MetricBatchAddAccountId(builder, t.AccountId)
	MetricBatchAddMetrics(builder, metricsOffset)
	return MetricBatchEnd(builder)
}

func (rcv *MetricBatch) UnPack() *MetricBatchT {
	if rcv == nil { return nil }
	t := &MetricBatchT{}
	t.Timestamp = rcv.Timestamp()
	t.CheckName = string(rcv.CheckName())
	t.CheckUuid = string(rcv.CheckUuid())
	t.AccountId = rcv.AccountId()
	metricsLength := rcv.MetricsLength()
	t.Metrics = make([]*MetricValueT, metricsLength)
	for j := 0; j < metricsLength; j++ {
		x := MetricValue{}
		rcv.Metrics(&x, j)
		t.Metrics[j] = x.UnPack()
	}
	return t
}

type MetricBatch struct {
	_tab flatbuffers.Table
}

func GetRootAsMetricBatch(buf []byte, offset flatbuffers.UOffsetT) *MetricBatch {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MetricBatch{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *MetricBatch) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MetricBatch) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MetricBatch) Timestamp() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetricBatch) MutateTimestamp(n uint64) bool {
	return rcv._tab.MutateUint64Slot(4, n)
}

func (rcv *MetricBatch) CheckName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetricBatch) CheckUuid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetricBatch) AccountId() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetricBatch) MutateAccountId(n int32) bool {
	return rcv._tab.MutateInt32Slot(10, n)
}

func (rcv *MetricBatch) Metrics(obj *MetricValue, j int) bool {
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

func (rcv *MetricBatch) MetricsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func MetricBatchStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func MetricBatchAddTimestamp(builder *flatbuffers.Builder, timestamp uint64) {
	builder.PrependUint64Slot(0, timestamp, 0)
}
func MetricBatchAddCheckName(builder *flatbuffers.Builder, checkName flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(checkName), 0)
}
func MetricBatchAddCheckUuid(builder *flatbuffers.Builder, checkUuid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(checkUuid), 0)
}
func MetricBatchAddAccountId(builder *flatbuffers.Builder, accountId int32) {
	builder.PrependInt32Slot(3, accountId, 0)
}
func MetricBatchAddMetrics(builder *flatbuffers.Builder, metrics flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(metrics), 0)
}
func MetricBatchStartMetricsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MetricBatchEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}