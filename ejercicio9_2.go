package main

import "sync"

var initTableOnce sync.Once


var pc [256]byte

func initTable() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}


func PopCount(x uint64) int {
	initTableOnce.Do(initTable)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}


