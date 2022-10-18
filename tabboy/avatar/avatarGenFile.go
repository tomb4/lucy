// Generated by github.com/davyxu/tabtoy
// DO NOT EDIT!!
// Version:
package main

import "errors"

type TableEnumValue struct {
	Name  string
	Index int32
}

type AvatarData struct {
	Id            int32    `tb_name:"ID"`
	FirstClass    string   `tb_name:"一级分类"`
	SecondClass   string   `tb_name:"二级分类"`
	Name          string   `tb_name:"名称"`
	ColorName     string   `tb_name:"颜色名称"`
	ColorList     []string `tb_name:"颜色列表"`
	IsShow        string   `tb_name:"是否展示"`
	LangIndonesia string   `tb_name:"印尼语"`
	OpType        int32    `tb_name:"操作类型"`
	SelectIdx     int32    `tb_name:"选中项"`
}

type AvatarElement struct {
	Id          int32  `tb_name:"ID"`
	FirstClass  string `tb_name:"一级分类"`
	SecondClass string `tb_name:"二级分类"`
	Name        string `tb_name:"名称"`
	UnityUri    string `tb_name:"U3D资源"`
	AppUri      string `tb_name:"APP资源"`
	Params      string `tb_name:"参数项"`
}

// Combine struct
type Table struct {
	AvatarData    []*AvatarData    // table: AvatarData
	AvatarElement []*AvatarElement // table: AvatarElement

	// Indices
	AvatarDataById    map[int32]*AvatarData    `json:"-"` // table: AvatarData
	AvatarElementById map[int32]*AvatarElement `json:"-"` // table: AvatarElement

	// Handlers
	postHandlers []func(*Table) error `json:"-"`
	preHandlers  []func(*Table) error `json:"-"`

	indexHandler map[string]func() `json:"-"`
	resetHandler map[string]func() `json:"-"`
}

// 注册加载后回调(用于构建数据)
func (self *Table) RegisterPostEntry(h func(*Table) error) {

	if h == nil {
		panic("empty postload handler")
	}

	self.postHandlers = append(self.postHandlers, h)
}

// 注册加载前回调(用于清除数据)
func (self *Table) RegisterPreEntry(h func(*Table) error) {

	if h == nil {
		panic("empty preload handler")
	}

	self.preHandlers = append(self.preHandlers, h)
}

// 清除索引和数据
func (self *Table) ResetData() error {

	err := self.InvokePreHandler()
	if err != nil {
		return err
	}

	return self.ResetTable("")
}

// 全局表构建索引及通知回调
func (self *Table) BuildData() error {

	err := self.IndexTable("")
	if err != nil {
		return err
	}

	return self.InvokePostHandler()
}

// 调用加载前回调
func (self *Table) InvokePreHandler() error {
	for _, h := range self.preHandlers {
		if err := h(self); err != nil {
			return err
		}
	}

	return nil
}

// 调用加载后回调
func (self *Table) InvokePostHandler() error {
	for _, h := range self.postHandlers {
		if err := h(self); err != nil {
			return err
		}
	}

	return nil
}

// 为表建立索引. 表名为空时, 构建所有表索引
func (self *Table) IndexTable(tableName string) error {

	if tableName == "" {

		for _, h := range self.indexHandler {
			h()
		}
		return nil

	} else {
		if h, ok := self.indexHandler[tableName]; ok {
			h()
		}

		return nil
	}
}

// 重置表格数据
func (self *Table) ResetTable(tableName string) error {
	if tableName == "" {
		for _, h := range self.resetHandler {
			h()
		}

		return nil
	} else {
		if h, ok := self.resetHandler[tableName]; ok {
			h()
			return nil
		}

		return errors.New("reset table failed, table not found: " + tableName)
	}
}

// 初始化表实例
func NewTable() *Table {

	self := &Table{
		indexHandler: make(map[string]func()),
		resetHandler: make(map[string]func()),
	}

	self.indexHandler["AvatarData"] = func() {

		for _, v := range self.AvatarData {
			self.AvatarDataById[v.Id] = v
		}
	}

	self.indexHandler["AvatarElement"] = func() {

		for _, v := range self.AvatarElement {
			self.AvatarElementById[v.Id] = v
		}
	}

	self.resetHandler["AvatarData"] = func() {
		self.AvatarData = nil

		self.AvatarDataById = map[int32]*AvatarData{}
	}
	self.resetHandler["AvatarElement"] = func() {
		self.AvatarElement = nil

		self.AvatarElementById = map[int32]*AvatarElement{}
	}

	self.ResetData()

	return self
}

func init() {

}
