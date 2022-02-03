//+build ignore

package ebpf

/*
#include "./c/tracer.h"
#include "./c/tcp_states.h"
#include "./c/prebuilt/offset-guess.h"
*/
import "C"

type ConnTuple C.conn_tuple_t
type TCPStats C.tcp_stats_t
type ConnStats C.conn_stats_ts_t
type Conn C.conn_t
type Batch C.batch_t
type Telemetry C.telemetry_t
type PortBinding C.port_binding_t
type PIDFD C.pid_fd_t
type UDPRecvSock C.udp_recv_sock_t
type BindSyscallArgs C.bind_syscall_args_t

// udp_recv_sock_t have *sock and *msghdr struct members, we make them opaque here
type _Ctype_struct_sock uint64
type _Ctype_struct_msghdr uint64

type TCPState uint8

const (
	Established TCPState = C.TCP_ESTABLISHED
	Close       TCPState = C.TCP_CLOSE
)

type ConnFlags uint32

const (
	LInit   ConnFlags = C.CONN_L_INIT
	RInit   ConnFlags = C.CONN_R_INIT
	Assured ConnFlags = C.CONN_ASSURED
)

const BatchSize = C.CONN_CLOSED_BATCH_SIZE
