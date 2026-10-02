//go:build darwin

package ui

type KeyMode int

const (
	ModeEncryption KeyMode = iota
	ModeDecryption
	ModeBoth
)
