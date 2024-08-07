package session

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"time"
)

// Session session对象
type Session struct {
	mSessionID        string
	mUserID           interface{}
	mLastTimeAccessed time.Time
	mValue            map[string]interface{} //持久化数据
	mTempValue        map[string]interface{}
	mOnSave           func(sid, value string)
	mManager          *Manager
}

// HasData 查找数据
func (sess *Session) HasData(key string) bool {
	_, ok := sess.mTempValue[key]
	if !ok {
		_, ok = sess.mValue[key]
	}
	return ok
}

// GetData 获取数据
func (sess *Session) GetData(key string) interface{} {
	data := sess.mTempValue[key]
	if data == nil {
		data = sess.mValue[key]
	}
	return data
}

// PutTempData 存储临时数据
func (sess *Session) PutTempData(key string, value interface{}) {
	sess.mTempValue[key] = value
}

// RemoveTempData 移除临时数据
func (sess *Session) RemoveTempData(key string) {
	delete(sess.mTempValue, key)
}

// PutData 存储数据
func (sess *Session) PutData(key string, value interface{}) {
	//log.Println("put data: ", this, key, value)
	sess.mValue[key] = value
	if sess.mOnSave != nil {
		sess.save()
	}
	if sess.mManager != nil {
		sess.saveGob()
	}
}

// RemoveData 移除数据
func (sess *Session) RemoveData(key string) {
	delete(sess.mValue, key)
	if sess.mOnSave != nil {
		sess.save()
	}
	if sess.mManager != nil {
		sess.saveGob()
	}
}

// Close 关闭会话
func (sess *Session) Close() {
	mgr := sess.mManager
	mgr.EndSession(sess.mSessionID)
}

// SessID sid
func (sess *Session) SessID() string {
	return sess.mSessionID
}

// UserID 用户ID
func (sess *Session) UserID() interface{} {
	return sess.mUserID
}

// save 数据持久化
func (sess *Session) save() {
	if sess.mOnSave == nil {
		return
	}

	bs, err := json.Marshal(sess.mValue)
	if err != nil {
		log.Println("[session] save err", err)
		return
	}

	sess.mOnSave(sess.mSessionID, string(bs))
}

func (sess *Session) saveGob() {
	if sess.mManager.mOption.OnSave == nil {
		return
	}

	var result bytes.Buffer
	err := gob.NewEncoder(&result).Encode(sess.mValue)
	if err != nil {
		log.Println("[session] save gob err", err)
		return
	}

	if err := sess.mManager.mOption.OnSave(
		sess.mSessionID,
		result.Bytes()); err != nil {
		log.Println("[session]OnSave failed, err:", err)
	}
}
