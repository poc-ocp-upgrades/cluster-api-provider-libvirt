package client

import (
	"io"
	libvirt "github.com/libvirt/libvirt-go"
)

type streamIO struct{ stream libvirt.Stream }

var _ io.Writer = &streamIO{}
var _ io.Reader = &streamIO{}
var _ io.Closer = &streamIO{}

func newStreamIO(s libvirt.Stream) *streamIO {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &streamIO{stream: s}
}
func (sio *streamIO) Read(p []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sio.stream.Recv(p)
}
func (sio *streamIO) Write(p []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sio.stream.Send(p)
}
func (sio *streamIO) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sio.stream.Finish()
}
