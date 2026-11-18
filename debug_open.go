//go:build debug

package gokit

const IsDebug = true

/*** 用于死码消除, 仅指定 --tags=debug 时才会编译
if IsDebug {
	// debug代码
}
*/
