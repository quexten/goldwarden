package config

import "github.com/awnumar/memguard"

type LockedBuffer interface {
	Bytes() []byte
	Wipe()
}

func NewBuffer(size int, useMemguard bool) LockedBuffer {
	if useMemguard {
		return MemGuardLockedBuffer{memguard.NewBuffer(size)}
	}
	return MemoryLockedBuffer{make([]byte, size)}
}

func NewBufferFromBytes(bytes []byte, useMemguard bool) LockedBuffer {
	if useMemguard {
		return MemGuardLockedBuffer{memguard.NewBufferFromBytes(bytes)}
	}
	return MemoryLockedBuffer{bytes}
}

type MemGuardLockedBuffer struct {
	buffer *memguard.LockedBuffer
}

func (b MemGuardLockedBuffer) Wipe() {
	b.buffer.Destroy()
}

func (b MemGuardLockedBuffer) Bytes() []byte {
	return b.buffer.Bytes()
}

type MemoryLockedBuffer struct {
	buffer []byte
}

func (b MemoryLockedBuffer) Wipe() {
	for i := range b.buffer {
		b.buffer[i] = 0
	}
}

func (b MemoryLockedBuffer) Bytes() []byte {
	return b.buffer
}
