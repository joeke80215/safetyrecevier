package safetyrecevier

import (
	"bytes"
	"io/ioutil"
	"os"
)

type readHandler func([]byte) (int, error)
type seekHandler func(int64, int) (int64, error)
type appendHandler func([]byte, int)

type handler struct {
	read   readHandler
	seek   seekHandler
	append appendHandler
}

type tmpFileReader interface {
	Read(p []byte) (n int, err error)
	Close() error
	Seek(offset int64, whence int) (int64, error)
}

// SafeReceive safe receive handler
type SafeReceive struct {
	maxMemorySize int
	b             []byte
	byteSize      int
	isLargeFile   bool
	tmpFile       *os.File
	tmpFileRead   tmpFileReader
	reader        *bytes.Reader
	handle        handler
}

// New receiver
func New(maxSize int) *SafeReceive {
	sr := &SafeReceive{
		maxMemorySize: maxSize,
	}
	sr.handle = handler{
		read:   sr.read,
		seek:   sr.seek,
		append: sr.append,
	}
	return sr
}

// Read implement read interface
func (s *SafeReceive) Read(p []byte) (n int, err error) {
	return s.handle.read(p)
}

// Seek implement seek interface
func (s *SafeReceive) Seek(offset int64, whence int) (int64, error) {
	return s.handle.seek(offset, whence)
}

// CloseReceive close receive stream
func (s *SafeReceive) CloseReceive() error {
	return s.closeReceive()
}

// CloseReader close reader and remove tmp file
func (s *SafeReceive) CloseReader() error {
	return s.closeReader()
}

// Append chunk,n from expect from Read()
func (s *SafeReceive) Append(chunk []byte, n int) {
	s.handle.append(chunk, n)
}

func (s *SafeReceive) append(chunk []byte, n int) {
	s.b = append(s.b, chunk[:n]...)
	var chunkSize int
	if chunkSize = len(chunk); chunkSize >= s.maxMemorySize || chunkSize+s.byteSize >= s.maxMemorySize {
		s.isLargeFile = true
		s.setLargeFileHandler()
		s.initTmpFile()
		s.flush()
		return
	}
	s.byteSize += chunkSize
}

// GetIsLargeFile is the large file
func (s *SafeReceive) GetIsLargeFile() bool {
	return s.isLargeFile
}

// GetBytes get bytes if the large file return nil
func (s *SafeReceive) GetBytes() []byte {
	return s.b
}

// GetTmpFilePath get tmp file path
func (s *SafeReceive) GetTmpFilePath() string {
	if s.tmpFile == nil {
		return ""
	}
	return s.tmpFile.Name()
}

func (s *SafeReceive) largeFileAppend(chunk []byte, n int) {
	if _, err := s.tmpFile.Write(chunk[:n]); err != nil {
		panic(err)
	}
}

func (s *SafeReceive) initTmpFile() {
	pwdDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	s.tmpFile, err = ioutil.TempFile(pwdDir, "")
	if err != nil {
		panic(err)
	}
}

func (s *SafeReceive) flush() {
	if _, err := s.tmpFile.Write(s.b); err != nil {
		panic(err)
	}
	s.b = nil
	s.byteSize = 0
}

func (s *SafeReceive) setLargeFileHandler() {
	s.handle.read = s.largeFileRead
	s.handle.seek = s.largeSeek
	s.handle.append = s.largeFileAppend
}

func (s *SafeReceive) largeFileRead(p []byte) (n int, err error) {
	if s.tmpFileRead == nil {
		s.tmpFileRead, err = os.Open(s.tmpFile.Name())
		if err != nil {
			return
		}
	}
	return s.tmpFileRead.Read(p)
}

func (s *SafeReceive) read(p []byte) (n int, err error) {
	if s.reader == nil {
		s.reader = bytes.NewReader(s.b)
	}
	return s.reader.Read(p)
}

func (s *SafeReceive) largeSeek(offset int64, whence int) (int64, error) {
	return s.tmpFile.Seek(offset, whence)
}

func (s *SafeReceive) seek(offset int64, whence int) (int64, error) {
	if s.reader == nil {
		s.reader = bytes.NewReader(s.b)
	}
	return s.reader.Seek(offset, whence)
}

func (s *SafeReceive) closeReceive() (err error) {
	if s.tmpFile != nil {
		err = s.tmpFile.Close()
	}
	return
}

func (s *SafeReceive) closeReader() (err error) {
	if s.tmpFileRead != nil {
		err = s.tmpFileRead.Close()
	}
	if s.tmpFile != nil {
		os.Remove(s.tmpFile.Name())
	}
	return
}
