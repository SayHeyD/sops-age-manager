//go:build darwin || windows

package ui

type KeyMode int

const (
	ModeEncryption KeyMode = iota
	ModeDecryption
	ModeBoth
)
