/* ibyte_pool.go - interface */
/*
modification history
--------------------
2016/10/12, by zhangjiyang01@baidu.com
*/
/*
DESCRIPTION
*/

package byte_pool

type IBytePool interface {
	Set(int32, []byte) error
	Get(int32) []byte

	MaxElemSize() int
}
