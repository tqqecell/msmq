/*
 * @Description: BUG无处不在。
 * @Author: 晴天
 * @Date: 2020-01-16 15:27:30
 * @LastEditTime : 2020-01-16 17:30:12
 * @LastEditors  : 晴天
 */
package queue

import (
	"bytes"
	"encoding/binary"
	"fmt"

	// 使用共享库而不是cgo绑定Windows COM。
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const (
	MQ_DENY_NONE = 0

	MQ_SEND_ACCESS    = 2
	MQ_RECEIVE_ACCESS = 1
)

var (
	vtMissing = ole.NewVariant(ole.VT_ERROR, 0)
)

type Msg struct {
	FCode   int32  `json:"fCode"`
	Tbl     int32  `json:"tbl"`
	Cstring Char24 `json:"slabNo"`
}

/**
 * @description:初始化
 * @param {type}
 * @return:
 */
func InitQueue(ip, queueName, label string) error {
	ole.CoInitialize(0)
	_, err := openQueue(ip, queueName, MQ_SEND_ACCESS, MQ_DENY_NONE)
	if err != nil {
		return creatQuery(queueName, label)
	}
	return nil
}

/**
 * @description:发送
 * @param {type}
 * @return:
 */
func SendEntry(ip, queueName string, msg *Msg) error {
	_que, err := openQueue(ip, queueName, MQ_SEND_ACCESS, MQ_DENY_NONE)
	if err != nil {
		return err
	}
	defer _que.Release()
	_msg, err := creatMessage()
	if err != nil {
		return err
	}
	defer _msg.Release()
	defer oleutil.MustCallMethod(_que, "Close")

	//数据缓存区
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.LittleEndian, msg)
	if err != nil {
		return err
	}
	oleutil.MustPutProperty(_msg, "Body", buf.Bytes())
	oleutil.MustCallMethod(_msg, "Send", _que)
	return nil
}

/**
 * @description:读取一条信息
 * @param {type}
 * @return:
 */
func ReadEntry(ip, queueName string, msg *Msg) error {
	_que, err := openQueue(ip, queueName, MQ_RECEIVE_ACCESS, MQ_DENY_NONE)
	if err != nil {
		return err
	}
	defer _que.Release()
	// MSMQQueue.Receive 检索队列中的第一条消息，从队列中删除消息
	// Transaction
	// 	MQ_NO_TRANSACTION 非事务
	// 	MQ_MTS_TRANSACTION: 默认值
	// 	MQ_XA_TRANSACTION: 外部XA 事务
	// WantDestinationQueue
	// WantBody 默认为true 是否索引邮件正文
	// ReceiveTimeout 超时时间
	_msg := oleutil.MustCallMethod(_que, "Receive", &vtMissing, &vtMissing, true, -1).ToIDispatch()
	if _msg == nil {
		return fmt.Errorf("未收到消息")
	}
	defer _msg.Release()
	defer oleutil.MustCallMethod(_que, "Close")
	// 索引消息长度
	_msg_size := oleutil.MustGetProperty(_msg, "BodyLength").Val
	if _msg_size == 0 {
		return fmt.Errorf("消息长度为0")
	}
	// MSMQMessage.Body
	_body := oleutil.MustGetProperty(_msg, "Body").ToArray().ToByteArray()

	buf := bytes.NewBuffer(_body)
	err = binary.Read(buf, binary.LittleEndian, msg)
	if err != nil {
		return err
	}
	return nil
}

/**
 * @description: 打开
 * @param {type}
 * @return:
 */
func openQueue(ip, queueName string, access, share int) (*ole.IDispatch, error) {
	_que_info, err := creatQueue()
	if err != nil {
		return nil, err
	}
	defer _que_info.Release()

	// 连接FormatName
	// Direct=tcp:xxx.xxx.xxx.xxx \private$\yourqname
	// Direct=os:.\private$\yourqname
	formatName := fmt.Sprintf("DIRECT=os:%s\\private$\\%s", ip, queueName)
	// 更新MSMQQueueInfo 属性FormatName
	oleutil.MustPutProperty(_que_info, "FormatName", formatName)

	// MSMQQueueInfo.Open
	// 参数 access 制定应用程序如何访问队列
	// 	MQ_PEEK_ACCESS 只能查看邮件。它们不能从队列中删除。
	// 	MQ_SEND_ACCESS 消息只能发送到队列
	// 	MQ_RECEIVE_ACCESS 读取、删除、查看、清除
	// 参数 share 谁可以访问队列
	// 	MQ_DENY_NONE 队列对Everyone组的所有成员都可用，如果Access 设置为MQ_PEEK_ACCESS或MQ_SEND_ACCESS 则必须使用此设置
	// 	MQ_DENY_RECEIVE_SHARE 限定从队列中检索消息的进程，如果队列已由另一个进程打开，其他进程调用时失败
	_tmp_queue, err := oleutil.CallMethod(_que_info, "Open", access, share)
	if err != nil {
		return nil, err
	}
	// 将变量转换成对象
	_que := _tmp_queue.ToIDispatch()
	if _que == nil {
		return nil, fmt.Errorf("队列不存在")
	}
	return _que, nil
}

/**
 * @description: 创建
 * @param {type}
 * @return:
 */
func creatQuery(pathName, label string) error {
	_que_info, err := creatQueue()
	if err != nil {
		return err
	}
	defer _que_info.Release()
	pathName = fmt.Sprintf(".\\Private$\\%s", pathName)
	oleutil.MustPutProperty(_que_info, "PathName", pathName)
	oleutil.MustPutProperty(_que_info, "Label", label)
	_, err = oleutil.CallMethod(_que_info, "Create", false, true)
	if err != nil {
		return err
	}
	return nil
}

/**
 * @description:创建
 * 	MSMQQueueInfo COM 接口
 * 		Create 创建queue 返回 MSMQQueueInfo 对象
 * 		Delete 移除queue 邮箱
 *		Open 打开queue 返回 MSMQQueue 对象
 *  	Refresh 刷新消息队列
 *  	Update 将当前MSMQQueueInfo对象属性更新到服务器
 * @param {type}
 * @return:
 */
func creatQueue() (*ole.IDispatch, error) {
	// 根据接口类型从programID创建对象
	_unknown, err := oleutil.CreateObject("MSMQ.MSMQQueueInfo")
	if err != nil {
		return nil, err
	}
	defer _unknown.Release()
	return _unknown.MustQueryInterface(ole.IID_IDispatch), nil
}

/**
 * @description:消息
 * @param {type}
 * @return:
 */
func creatMessage() (*ole.IDispatch, error) {
	_unknown, err := oleutil.CreateObject("MSMQ.MSMQMessage")
	if err != nil {
		return nil, err
	}
	defer _unknown.Release()
	return _unknown.MustQueryInterface(ole.IID_IDispatch), nil
}
