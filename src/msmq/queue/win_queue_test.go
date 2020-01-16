/*
 * @Description: BUG无处不在。
 * @Author: 晴天
 * @Date: 2020-01-16 16:41:54
 * @LastEditTime : 2020-01-16 17:22:01
 * @LastEditors  : 晴天
 */

package queue

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var (
		_ip    = "."
		_email = "cc"
		label  = "test"
	)
	InitQueue(_ip, _email, label)

	var _c Char24
	_c.Set("SLANMK")

	_msg_int := Msg{FCode: 400, Tbl: 100, Cstring: _c}
	err := SendEntry(_ip, _email, &_msg_int)
	if err != nil {
		t.Fatal(err)
	}
	_msg_out := Msg{}
	err = ReadEntry(_ip, _email, &_msg_out)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("读取数据：FCOde:%d,Tbl:%d,char:%s\n",
		_msg_out.FCode,
		_msg_out.Tbl,
		_msg_out.Cstring.String())
}
