package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// RandCID 随机下单时自定义的订单 ID
func RandCID() string {
	var id uuid.UUID
	var err error
	for {
		id, err = uuid.NewUUID()
		if err != nil {
			continue
		}
		break
	}

	ids := strings.Split(id.String(), "-")
	nid := strings.Join(ids[0:4], "")
	return nid
}

// IDMap ID 双向绑定
type IDMap struct {
	s2t map[string]string
	t2s map[string]string
}

// GetIDMap 获取双向绑定 MAP
func GetIDMap() *IDMap {
	return &IDMap{
		s2t: make(map[string]string),
		t2s: make(map[string]string),
	}
}

// Set 设置双向绑定
func (im *IDMap) Set(sid string, tid string) error {
	if _, ok := im.s2t[sid]; ok {
		return fmt.Errorf("sid exist")
	}
	if _, ok := im.t2s[tid]; ok {
		return fmt.Errorf("tid exist")
	}

	im.s2t[sid] = tid
	im.t2s[tid] = sid
	return nil
}

// GetSID 根据 tid 获取 sid
func (im *IDMap) GetSID(tid string) (string, error) {
	sid, ok := im.t2s[tid]
	if !ok {
		return "", fmt.Errorf("tid not exist")
	}
	return sid, nil
}

// GetTID 根据 sid 获取 tid
func (im *IDMap) GetTID(sid string) (string, error) {
	tid, ok := im.s2t[sid]
	if !ok {
		return "", fmt.Errorf("sid not exist")
	}
	return tid, nil
}

// DelSID 删除 sid 对应绑定
func (im *IDMap) DelSID(sid string) {
	tid, err := im.GetTID(sid)
	if err != nil {
		delete(im.t2s, tid)
	}

	delete(im.s2t, sid)
}

// DelTID 删除 tid 对应绑定
func (im *IDMap) DelTID(tid string) {
	sid, err := im.GetSID(tid)
	if err != nil {
		delete(im.s2t, sid)
	}

	delete(im.t2s, tid)
}
