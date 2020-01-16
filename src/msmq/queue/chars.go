/*
 * @Description: BUG无处不在。
 * @Author: 晴天
 * @Date: 2020-01-16 17:17:10
 * @LastEditTime : 2020-01-16 17:20:31
 * @LastEditors  : 晴天
 */
package queue

import (
	"bytes"
)

type (
	Char1  [1]byte
	Char3  [3]byte
	Char5  [5]byte
	Char7  [7]byte
	Char8  [8]byte
	Char9  [9]byte
	Char10 [10]byte
	Char11 [11]byte
	Char12 [12]byte
	Char16 [16]byte
	Char20 [20]byte
	Char24 [24]byte
)

func (c Char1) String() string  { return String(c[:]) }
func (c Char3) String() string  { return String(c[:]) }
func (c Char5) String() string  { return String(c[:]) }
func (c Char7) String() string  { return String(c[:]) }
func (c Char8) String() string  { return String(c[:]) }
func (c Char9) String() string  { return String(c[:]) }
func (c Char10) String() string { return String(c[:]) }
func (c Char11) String() string { return String(c[:]) }
func (c Char12) String() string { return String(c[:]) }
func (c Char16) String() string { return String(c[:]) }
func (c Char20) String() string { return String(c[:]) }
func (c Char24) String() string { return String(c[:]) }

func (c *Char1) Set(s string)  { CString(c[:], s) }
func (c *Char3) Set(s string)  { CString(c[:], s) }
func (c *Char5) Set(s string)  { CString(c[:], s) }
func (c *Char7) Set(s string)  { CString(c[:], s) }
func (c *Char8) Set(s string)  { CString(c[:], s) }
func (c *Char9) Set(s string)  { CString(c[:], s) }
func (c *Char10) Set(s string) { CString(c[:], s) }
func (c *Char11) Set(s string) { CString(c[:], s) }
func (c *Char12) Set(s string) { CString(c[:], s) }
func (c *Char16) Set(s string) { CString(c[:], s) }
func (c *Char20) Set(s string) { CString(c[:], s) }
func (c *Char24) Set(s string) { CString(c[:], s) }

func String(bs []byte) string {

	if i := bytes.Index(bs, []byte{0}); i != -1 {
		return string(bs[:i])
	}

	return ""
}

// Bytes copies string characters into a byte array, and is null terminateds it.
func CString(out []byte, in string) {

	// Copy bytes from string into buffer
	n := copy(out[:len(out)-1], in[:])

	// If n is less than the length of the output buffer
	if n < len(out) {
		// Set the next byte to zero
		out[n] = 0
	} else {
		// Set the last byte to zero
		out[n-1] = 0
	}
}
