package kit

import "io"

type pipeConn struct {
	io.WriteCloser
	io.Reader
}

// PipeConn creates a pair of io.ReadWriteCloser that write and read to each other
func PipeConn() (io.ReadWriteCloser, io.ReadWriteCloser) {
	client, server := new(pipeConn), new(pipeConn)
	{
		r, w := io.Pipe()
		client.WriteCloser, server.Reader = w, r
	}
	{
		r, w := io.Pipe()
		server.WriteCloser, client.Reader = w, r
	}
	return client, server
}
