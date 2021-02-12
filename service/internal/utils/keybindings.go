package utils

import "github.com/micmonay/keybd_event"

type KeyItem struct {
	KeyCode int
	Shift   bool
	Altgr   bool
}

var KeyEnglish = []KeyItem{
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	// Space
	{keybd_event.VK_SPACE, false, false},
	// !
	{keybd_event.VK_1, true, false},
	// "
	{keybd_event.VK_SP7, true, false},
	// #
	{keybd_event.VK_3, true, false},
	// $
	{keybd_event.VK_4, true, false},
	// %
	{keybd_event.VK_5, true, false},
	// &
	{keybd_event.VK_7, true, false},
	// '
	{keybd_event.VK_SP7, false, false},
	// (
	{keybd_event.VK_9, true, false},
	// )
	{keybd_event.VK_0, true, false},
	// *
	{keybd_event.VK_8, true, false},
	// +
	{keybd_event.VK_SP3, true, false},
	// ,
	{keybd_event.VK_SP9, false, false},
	// -
	{keybd_event.VK_2, false, false},
	// .
	{keybd_event.VK_SP10, false, false},
	// /
	{keybd_event.VK_SP11, false, false},
	// 0
	{keybd_event.VK_0, false, false},
	// 1
	{keybd_event.VK_1, false, false},
	// 2
	{keybd_event.VK_2, false, false},
	// 3
	{keybd_event.VK_3, false, false},
	// 4
	{keybd_event.VK_4, false, false},
	// 5
	{keybd_event.VK_5, false, false},
	// 6
	{keybd_event.VK_6, false, false},
	// 7
	{keybd_event.VK_7, false, false},
	// 8
	{keybd_event.VK_8, false, false},
	// 9
	{keybd_event.VK_9, false, false},
	// :
	{keybd_event.VK_SP6, true, false},
	// ;
	{keybd_event.VK_SP6, false, false},
	// <
	{keybd_event.VK_SP9, true, false},
	// =
	{keybd_event.VK_SP3, false, false},
	// >
	{keybd_event.VK_SP10, true, false},
	// ?
	{keybd_event.VK_SP11, true, false},
	// @
	{keybd_event.VK_2, true, false},
	// A
	{keybd_event.VK_A, true, false},
	// B
	{keybd_event.VK_B, true, false},
	// C
	{keybd_event.VK_C, true, false},
	// D
	{keybd_event.VK_D, true, false},
	// E
	{keybd_event.VK_E, true, false},
	// F
	{keybd_event.VK_F, true, false},
	// G
	{keybd_event.VK_G, true, false},
	// H
	{keybd_event.VK_H, true, false},
	// I
	{keybd_event.VK_I, true, false},
	// J
	{keybd_event.VK_J, true, false},
	// K
	{keybd_event.VK_K, true, false},
	// L
	{keybd_event.VK_L, true, false},
	// M
	{keybd_event.VK_M, true, false},
	// N
	{keybd_event.VK_N, true, false},
	// O
	{keybd_event.VK_O, true, false},
	// P
	{keybd_event.VK_P, true, false},
	// Q
	{keybd_event.VK_Q, true, false},
	// R
	{keybd_event.VK_R, true, false},
	// S
	{keybd_event.VK_S, true, false},
	// T
	{keybd_event.VK_T, true, false},
	// U
	{keybd_event.VK_U, true, false},
	// V
	{keybd_event.VK_V, true, false},
	// W
	{keybd_event.VK_W, true, false},
	// X
	{keybd_event.VK_X, true, false},
	// Y
	{keybd_event.VK_Y, true, false},
	// Z
	{keybd_event.VK_Z, true, false},
	// [
	{keybd_event.VK_SP4, false, false},
	// \
	{keybd_event.VK_SP8, false, false},
	// ]
	{keybd_event.VK_SP5, false, false},
	// ^
	{keybd_event.VK_6, true, false},
	// _
	{keybd_event.VK_SP2, true, false},
	// `
	{keybd_event.VK_SP12, false, false},
	// a
	{keybd_event.VK_A, false, false},
	// b
	{keybd_event.VK_B, false, false},
	// c
	{keybd_event.VK_C, false, false},
	// d
	{keybd_event.VK_D, false, false},
	// e
	{keybd_event.VK_E, false, false},
	// f
	{keybd_event.VK_F, false, false},
	// g
	{keybd_event.VK_G, false, false},
	// h
	{keybd_event.VK_H, false, false},
	// i
	{keybd_event.VK_I, false, false},
	// j
	{keybd_event.VK_J, false, false},
	// k
	{keybd_event.VK_K, false, false},
	// l
	{keybd_event.VK_L, false, false},
	// m
	{keybd_event.VK_M, false, false},
	// n
	{keybd_event.VK_N, false, false},
	// o
	{keybd_event.VK_O, false, false},
	// p
	{keybd_event.VK_P, false, false},
	// q
	{keybd_event.VK_Q, false, false},
	// r
	{keybd_event.VK_R, false, false},
	// s
	{keybd_event.VK_S, false, false},
	// t
	{keybd_event.VK_T, false, false},
	// u
	{keybd_event.VK_U, false, false},
	// v
	{keybd_event.VK_V, false, false},
	// w
	{keybd_event.VK_W, false, false},
	// x
	{keybd_event.VK_X, false, false},
	// y
	{keybd_event.VK_Y, false, false},
	// z
	{keybd_event.VK_Z, false, false},
	// {
	{keybd_event.VK_SP4, true, false},
	// |
	{keybd_event.VK_SP8, true, false},
	// }
	{keybd_event.VK_SP5, true, false},
	// ~
	{keybd_event.VK_SP12, true, false},
}

var KeyGerman = []KeyItem{
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	// Space
	{keybd_event.VK_SPACE, false, false},
	// !
	{keybd_event.VK_1, true, false},
	// "
	{keybd_event.VK_2, true, false},
	// #
	{keybd_event.VK_SP8, false, false},
	// $
	{keybd_event.VK_4, true, false},
	// %
	{keybd_event.VK_5, true, false},
	// &
	{keybd_event.VK_6, true, false},
	// '
	{keybd_event.VK_SP8, true, false},
	// (
	{keybd_event.VK_8, true, false},
	// )
	{keybd_event.VK_9, true, false},
	// *
	{keybd_event.VK_SP5, true, false},
	// +
	{keybd_event.VK_SP5, false, false},
	// ,
	{keybd_event.VK_SP9, false, false},
	// -
	{keybd_event.VK_SP11, false, false},
	// .
	{keybd_event.VK_SP10, false, false},
	// /
	{keybd_event.VK_7, true, false},
	// 0
	{keybd_event.VK_0, false, false},
	// 1
	{keybd_event.VK_1, false, false},
	// 2
	{keybd_event.VK_2, false, false},
	// 3
	{keybd_event.VK_3, false, false},
	// 4
	{keybd_event.VK_4, false, false},
	// 5
	{keybd_event.VK_5, false, false},
	// 6
	{keybd_event.VK_6, false, false},
	// 7
	{keybd_event.VK_7, false, false},
	// 8
	{keybd_event.VK_8, false, false},
	// 9
	{keybd_event.VK_9, false, false},
	// :
	{keybd_event.VK_SP10, true, false},
	// ;
	{keybd_event.VK_SP9, true, false},
	// <
	{keybd_event.VK_SP12, false, false},
	// =
	{keybd_event.VK_0, true, false},
	// >
	{keybd_event.VK_SP12, true, false},
	// ?
	{keybd_event.VK_SP2, true, false},
	// @ todo
	{keybd_event.VK_Q, false, true},
	// A
	{keybd_event.VK_A, true, false},
	// B
	{keybd_event.VK_B, true, false},
	// C
	{keybd_event.VK_C, true, false},
	// D
	{keybd_event.VK_D, true, false},
	// E
	{keybd_event.VK_E, true, false},
	// F
	{keybd_event.VK_F, true, false},
	// G
	{keybd_event.VK_G, true, false},
	// H
	{keybd_event.VK_H, true, false},
	// I
	{keybd_event.VK_I, true, false},
	// J
	{keybd_event.VK_J, true, false},
	// K
	{keybd_event.VK_K, true, false},
	// L
	{keybd_event.VK_L, true, false},
	// M
	{keybd_event.VK_M, true, false},
	// N
	{keybd_event.VK_N, true, false},
	// O
	{keybd_event.VK_O, true, false},
	// P
	{keybd_event.VK_P, true, false},
	// Q
	{keybd_event.VK_Q, true, false},
	// R
	{keybd_event.VK_R, true, false},
	// S
	{keybd_event.VK_S, true, false},
	// T
	{keybd_event.VK_T, true, false},
	// U
	{keybd_event.VK_U, true, false},
	// V
	{keybd_event.VK_V, true, false},
	// W
	{keybd_event.VK_W, true, false},
	// X
	{keybd_event.VK_X, true, false},
	// Y
	{keybd_event.VK_Z, true, false},
	// Z
	{keybd_event.VK_Y, true, false},
	// [ todo
	{keybd_event.VK_8, false, true},
	// \ todo
	{keybd_event.VK_SP2, false, true},
	// ] todo
	{keybd_event.VK_9, false, true},
	// ^
	{keybd_event.VK_SP1, true, false},
	// _
	{keybd_event.VK_SP11, true, false},
	// `
	{keybd_event.VK_SP3, true, false},
	// a
	{keybd_event.VK_A, false, false},
	// b
	{keybd_event.VK_B, false, false},
	// c
	{keybd_event.VK_C, false, false},
	// d
	{keybd_event.VK_D, false, false},
	// e
	{keybd_event.VK_E, false, false},
	// f
	{keybd_event.VK_F, false, false},
	// g
	{keybd_event.VK_G, false, false},
	// h
	{keybd_event.VK_H, false, false},
	// i
	{keybd_event.VK_I, false, false},
	// j
	{keybd_event.VK_J, false, false},
	// k
	{keybd_event.VK_K, false, false},
	// l
	{keybd_event.VK_L, false, false},
	// m
	{keybd_event.VK_M, false, false},
	// n
	{keybd_event.VK_N, false, false},
	// o
	{keybd_event.VK_O, false, false},
	// p
	{keybd_event.VK_P, false, false},
	// q
	{keybd_event.VK_Q, false, false},
	// r
	{keybd_event.VK_R, false, false},
	// s
	{keybd_event.VK_S, false, false},
	// t
	{keybd_event.VK_T, false, false},
	// u
	{keybd_event.VK_U, false, false},
	// v
	{keybd_event.VK_V, false, false},
	// w
	{keybd_event.VK_W, false, false},
	// x
	{keybd_event.VK_X, false, false},
	// y
	{keybd_event.VK_Z, false, false},
	// z
	{keybd_event.VK_Y, false, false},
	// { todo
	{keybd_event.VK_7, false, true},
	// | todo
	{keybd_event.VK_SP12, false, true},
	// } todo
	{keybd_event.VK_0, false, true},
	// ~ todo
	{keybd_event.VK_SP5, false, true},
}
