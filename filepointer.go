package main

import "os"

// FilePointer method to move in log files
type FilePointer struct {
	Filename string
	File     *os.File
}

// NewFilePointer create new struc to locate line in file
func NewFilePointer(filename string) *FilePointer {
	f := &FilePointer{
		Filename: filename,
	}

	f.File, _ = os.Open(f.Filename)

	return f
}

// GetLine of the file
func (f *FilePointer) GetLine(pos int64, len int) []byte {

	// f.File.Seek(seek, len)
	buf := make([]byte, len)
	_, _ = f.File.ReadAt(buf, pos)

	return buf

}
