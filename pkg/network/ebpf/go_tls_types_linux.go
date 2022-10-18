// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs -- -fsigned-char -v http_types.go

package ebpf

type Location struct {
	Stack_offset int64
	X_register   int64
	In_register  uint8
	Exists       uint8
	Pad_cgo_0    [6]byte
}
type SliceLocation struct {
	Ptr Location
	Len Location
	Cap Location
}
type GoroutineIDMetadata struct {
	Runtime_g_tls_addr_offset uint64
	Goroutine_id_offset       uint64
	Runtime_g_register        int64
	Runtime_g_in_register     uint8
	Pad_cgo_0                 [7]byte
}
type TlsConnLayout struct {
	Tls_conn_inner_conn_offset uint64
	Tcp_conn_inner_conn_offset uint64
	Conn_fd_offset             uint64
	Net_fd_pfd_offset          uint64
	Fd_sysfd_offset            uint64
}

type TlsOffsetsData struct {
	Goroutine_id       GoroutineIDMetadata
	Conn_layout        TlsConnLayout
	Read_conn_pointer  Location
	Read_buffer        SliceLocation
	Read_return_bytes  Location
	Write_conn_pointer Location
	Write_buffer       SliceLocation
	Write_return_bytes Location
	Write_return_error Location
	Close_conn_pointer Location
}
