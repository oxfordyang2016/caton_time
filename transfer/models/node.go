package models

import (
	"time"
)
//node 's function is download and upload
type Node struct {
	Id             uint64    `xorm:"pk autoincr"`
	MachineCode    string    `xorm:"varchar(255) not null unique"`
	Nid            string    `xorm:"varchar(64) not null unique"`
	Name           string    `xorm:"varchar(255)"`
	Zid            string    `xomr:"varchar(255)"`
	RegisterTime   time.Time `xorm:"-"`
	LastLoginTime  time.Time `xorm:"DateTime"`
	LastLogoutTime time.Time `xorm:"DateTime"`
	CreateAt       time.Time `xorm:"DateTime created"`
	UpdateAt       time.Time `xorm:"DateTime updated"`
}

func CreateNode(machine_code, nid string) (*Node, error) {
	n := &Node{
		MachineCode: machine_code,
		Nid:         nid,
	}
	if _, err := DB().Insert(n); err != nil {
		return nil, err
	}
	return n, nil
}

func GetNode(nid string) (n *Node, err error) {
	n = new(Node)
	var existed bool
	if existed, err = DB().Where("nid=?", nid).Get(n); err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	n.RegisterTime = n.CreateAt
	return n, nil
}

func GetNodeByMachineCode(mc string) (n *Node, err error) {
	n = new(Node)
	var existed bool
	if existed, err = DB().Where("machine_code=?", mc).Get(n); err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	n.RegisterTime = n.CreateAt
	return n, nil
}

func (n *Node) Setup(zid string, name string, omitempty bool) error {
	var err error
	new_n := new(Node)
	new_n.Zid = zid
	new_n.Name = name
	if !omitempty {
		_, err = DB().Where("nid=?", n.Nid).Cols("zid", "name").Update(new_n)
	} else {
		_, err = DB().Where("nid=?", n.Nid).Update(new_n)
	}

	return err
}

func (n *Node) UpdateLoginTime(t time.Time) error {
	n.LastLoginTime = t
	if _, err := DB().Where("nid=?", n.Nid).Update(n); err != nil {
		return err
	}
	return nil
}

func (n *Node) UpdateLogoutTime(t time.Time) error {
	n.LastLogoutTime = t
	if _, err := DB().Where("nid=?", n.Nid).Update(n); err != nil {
		return err
	}
	return nil
}
