package chainbase

import (
	"reflect"
)

const (
	//ReadOnly ..
	ReadOnly = iota
	//ReadWrite ..
	ReadWrite = iota
)

//StrcmpLess ..
type StrcmpLess struct {
}

//Oid ..
type Oid struct {
}

//Object ..
type Object struct {
}

//UndoState ..
type UndoState struct {
}

//IntIncrementer ..
type IntIncrementer struct {
}

//AbstractSession ..
type AbstractSession interface {
	Push()
	Squash()
	Undo()
	Revision() uint64
}

//SessionImpl ..
type SessionImpl struct {
	session AbstractSession
}

//Push ..
func (s *SessionImpl) Push() {}

//Squash ..
func (s *SessionImpl) Squash() {}

//Undo ..
func (s *SessionImpl) Undo() {}

//Revision ..
func (s *SessionImpl) Revision() uint64 {
	return s.session.Revision()
}

//New ..
func (s *SessionImpl) New(a AbstractSession) {
	s.session = a
}

//AbstractIndex ..
type AbstractIndex interface {
	SetRevision(revision uint64)
	StartUndoSession(enabled bool) *AbstractSession
	Revision() int64
	Undo()
	Squash()
	Commit(revision int64)
	UndoAll()
	TypeID() uint32
	RowCount() uint64
	TypeName() string
	UndoStackRevisionRange() map[int64]int64
	RemoveObject(ID int64)
}

//template<typename MultiIndexType> 这个type 是传多索引容器进来，但是这里是传map进来

//GenericIndex ..
type GenericIndex struct {
	Session
	revision        int64
	nextID          uint32
	indices         map[interface{}]interface{}
	sizeOfValueType uint32
	sizeOfThis      uint32
}

//Revision ..
func (g GenericIndex) Revision() int64 {
	return g.revision
}

//UndoStackRevisionRange ..
func (g GenericIndex) UndoStackRevisionRange() map[int64]int64 {
	begin := g.revision
	end := g.revision

	var a map[int64]int64
	a[begin] = end

	// if( _stack.size() > 0 ) {
	//    begin = _stack.front().revision - 1;
	//    end   = _stack.back().revision;
	// }

	return a
}

//RemoveObject ..
func (g *GenericIndex) RemoveObject(ID int64) {

}

//IndexImpl ..
type IndexImpl struct {
	idxPtr *interface{}
	base   AbstractIndex
}

//SetRevision ..
func (i *IndexImpl) SetRevision(revision uint64) {
	i.base.SetRevision(revision)
}

//StartUndoSession ..
func (i *IndexImpl) StartUndoSession(enabled bool) AbstractSession {
	s := new(SessionImpl)
	a := i.base.StartUndoSession(enabled)
	s.New(*a)
	return s
}

//Revision ..
func (i *IndexImpl) Revision() int64 {
	return i.base.Revision()
}

//Undo ..
func (i *IndexImpl) Undo() {
	i.base.Undo()
}

//Squash ..
func (i *IndexImpl) Squash() {
	i.base.Squash()
}

//Commit ..
func (i *IndexImpl) Commit(revision int64) {
	i.base.Commit(revision)
}

//UndoAll ..
func (i *IndexImpl) UndoAll() {
	i.base.UndoAll()
}

//TypeID ..
func (i *IndexImpl) TypeID() uint32 {
	return uint32(reflect.TypeOf(i.base).Elem().Len())
}

//RowCount ..
func (i *IndexImpl) RowCount() uint64 {
	return uint64(i.base.Revision())
}

//TypeName ..
func (i *IndexImpl) TypeName() string {
	return reflect.TypeOf(i.base).String()
}

//UndoStackRevisionRange ..
func (i *IndexImpl) UndoStackRevisionRange() map[int64]int64 {
	return i.base.UndoStackRevisionRange()
}

//RemoveObject ..
func (i *IndexImpl) RemoveObject(ID int64) {

}

//Session ..
type Session struct {
	rversion uint64
}

//DataBase ..
type DataBase struct {
	Session
	sizeOfValueType uint32
	sizeOfThis      uint32
}

//GetRversion ..
func (c DataBase) GetRversion() uint64 {
	return c.rversion
}

//SetRversion ..
func (c *DataBase) SetRversion(u uint64) {
	c.rversion = u
}

//Undo ..
func (c DataBase) Undo() {
	println("DataBase Undo")
}
