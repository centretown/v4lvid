// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package openmeteo_sdk

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type VariableWithValues struct {
	_tab flatbuffers.Table
}

func GetRootAsVariableWithValues(buf []byte, offset flatbuffers.UOffsetT) *VariableWithValues {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &VariableWithValues{}
	x.Init(buf, n+offset)
	return x
}

func FinishVariableWithValuesBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsVariableWithValues(buf []byte, offset flatbuffers.UOffsetT) *VariableWithValues {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &VariableWithValues{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedVariableWithValuesBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *VariableWithValues) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *VariableWithValues) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *VariableWithValues) Variable() Variable {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return Variable(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *VariableWithValues) MutateVariable(n Variable) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func (rcv *VariableWithValues) Unit() Unit {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return Unit(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *VariableWithValues) MutateUnit(n Unit) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *VariableWithValues) Value() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *VariableWithValues) MutateValue(n float32) bool {
	return rcv._tab.MutateFloat32Slot(8, n)
}

func (rcv *VariableWithValues) Values(j int) float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetFloat32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *VariableWithValues) ValuesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *VariableWithValues) MutateValues(j int, n float32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateFloat32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *VariableWithValues) ValuesInt64(j int) int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *VariableWithValues) ValuesInt64Length() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *VariableWithValues) MutateValuesInt64(j int, n int64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *VariableWithValues) Altitude() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutateAltitude(n int16) bool {
	return rcv._tab.MutateInt16Slot(14, n)
}

func (rcv *VariableWithValues) Aggregation() Aggregation {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return Aggregation(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *VariableWithValues) MutateAggregation(n Aggregation) bool {
	return rcv._tab.MutateByteSlot(16, byte(n))
}

func (rcv *VariableWithValues) PressureLevel() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutatePressureLevel(n int16) bool {
	return rcv._tab.MutateInt16Slot(18, n)
}

func (rcv *VariableWithValues) Depth() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutateDepth(n int16) bool {
	return rcv._tab.MutateInt16Slot(20, n)
}

func (rcv *VariableWithValues) DepthTo() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutateDepthTo(n int16) bool {
	return rcv._tab.MutateInt16Slot(22, n)
}

func (rcv *VariableWithValues) EnsembleMember() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutateEnsembleMember(n int16) bool {
	return rcv._tab.MutateInt16Slot(24, n)
}

func (rcv *VariableWithValues) PreviousDay() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *VariableWithValues) MutatePreviousDay(n int16) bool {
	return rcv._tab.MutateInt16Slot(26, n)
}

func VariableWithValuesStart(builder *flatbuffers.Builder) {
	builder.StartObject(12)
}
func VariableWithValuesAddVariable(builder *flatbuffers.Builder, variable Variable) {
	builder.PrependByteSlot(0, byte(variable), 0)
}
func VariableWithValuesAddUnit(builder *flatbuffers.Builder, unit Unit) {
	builder.PrependByteSlot(1, byte(unit), 0)
}
func VariableWithValuesAddValue(builder *flatbuffers.Builder, value float32) {
	builder.PrependFloat32Slot(2, value, 0.0)
}
func VariableWithValuesAddValues(builder *flatbuffers.Builder, values flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(values), 0)
}
func VariableWithValuesStartValuesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func VariableWithValuesAddValuesInt64(builder *flatbuffers.Builder, valuesInt64 flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(valuesInt64), 0)
}
func VariableWithValuesStartValuesInt64Vector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func VariableWithValuesAddAltitude(builder *flatbuffers.Builder, altitude int16) {
	builder.PrependInt16Slot(5, altitude, 0)
}
func VariableWithValuesAddAggregation(builder *flatbuffers.Builder, aggregation Aggregation) {
	builder.PrependByteSlot(6, byte(aggregation), 0)
}
func VariableWithValuesAddPressureLevel(builder *flatbuffers.Builder, pressureLevel int16) {
	builder.PrependInt16Slot(7, pressureLevel, 0)
}
func VariableWithValuesAddDepth(builder *flatbuffers.Builder, depth int16) {
	builder.PrependInt16Slot(8, depth, 0)
}
func VariableWithValuesAddDepthTo(builder *flatbuffers.Builder, depthTo int16) {
	builder.PrependInt16Slot(9, depthTo, 0)
}
func VariableWithValuesAddEnsembleMember(builder *flatbuffers.Builder, ensembleMember int16) {
	builder.PrependInt16Slot(10, ensembleMember, 0)
}
func VariableWithValuesAddPreviousDay(builder *flatbuffers.Builder, previousDay int16) {
	builder.PrependInt16Slot(11, previousDay, 0)
}
func VariableWithValuesEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}