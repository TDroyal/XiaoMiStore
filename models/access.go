package models

import "time"

// 权限表，给每重角色分配不同权限的
type Access struct {
	ID          int `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ModuleName  string //后台管理员管理平台左侧的一级菜单(例如：管理员管理、角色管理、权限管理)
	Type        int    // 1表示节点类型为模块（一级菜单），  2表示节点类型为菜单（二级菜单）  3表示操作，用于权限的判断
	ActionName  string //操作名称  （管理员列表，增加管理员）
	Url         string //操作对应的地址
	ModuleID    int    `gorm:"index"` // 自身是个1对多的关系  0表示它是一个一级菜单， 如果它和此表中的ID相等，就表示这个ID下对应二级菜单
	Sort        int
	Description string
	Status      int
	//  自身是个1对多的关系 左侧导航栏对应的二级导航栏
	AccessList []Access `gorm:"foreignKey:ModuleID; references:ID"`
}

func (Access) TableName() string {
	return "access"
}
