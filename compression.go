package diskv

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

// Compression is an interface that Diskv uses to implement compression of data.
// Writer takes a destination io.Writer and returns a WriteCloser that
// compresses all data written through it. Reader takes a source io.Reader and
// returns a ReadCloser that decompresses all data read through it. You may
// define these methods on your own type, or use one of the NewCompression
// helpers.
type Compression interface {
	Writer(dst io.Writer) (io.WriteCloser, error)
	Reader(src io.Reader) (io.ReadCloser, error)
}

type genericCompression struct {
	wf func(w io.Writer) (io.WriteCloser, error)
	rf func(r io.Reader) (io.ReadCloser, error)
}

func (g *genericCompression) Writer(dst io.Writer) (io.WriteCloser, error) {
	return g.wf(dst)
}

func (g *genericCompression) Reader(src io.Reader) (io.ReadCloser, error) {
	return g.rf(src)
}

//
//
//

func NewGzipCompression() Compression {
	return NewGzipCompressionLevel(flate.DefaultCompression)
}

func NewGzipCompressionLevel(level int) Compression {
	return &genericCompression{
		wf: func(w io.Writer) (io.WriteCloser, error) { return gzip.NewWriterLevel(w, level) },
		rf: func(r io.Reader) (io.ReadCloser, error) { return gzip.NewReader(r) },
	}
}

func NewZlibCompression() Compression {
	return NewZlibCompressionLevel(flate.DefaultCompression)
}

func NewZlibCompressionLevel(level int) Compression {
	return NewZlibCompressionLevelDict(level, nil)
}

func NewZlibCompressionLevelDict(level int, dict []byte) Compression {
	return &genericCompression{
		func(w io.Writer) (io.WriteCloser, error) { return zlib.NewWriterLevelDict(w, level, dict) },
		func(r io.Reader) (io.ReadCloser, error) { return zlib.NewReaderDict(r, dict) },
	}
}
