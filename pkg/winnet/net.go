package winnet

import (
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/user/win-tools/pkg/winproc"
)

const (
	afInet  = 2
	afInet6 = 23

	tcpTableOwnerPIDAll = 5
	udpTableOwnerPID    = 1

	stateClosed      = 1
	stateListen      = 2
	stateSynSent     = 3
	stateSynReceived = 4
	stateEstablished = 5
	stateFinWait1    = 6
	stateFinWait2    = 7
	stateCloseWait   = 8
	stateClosing     = 9
	stateLastAck     = 10
	stateTimeWait    = 11
	stateDeleteTCB   = 12
)

var (
	iphlpapi           = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTcp = iphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdp = iphlpapi.NewProc("GetExtendedUdpTable")
)

// Conn is a TCP or UDP endpoint with optional owning PID.
type Conn struct {
	Proto      string
	LocalIP    net.IP
	LocalPort  uint16
	RemoteIP   net.IP
	RemotePort uint16
	State      string
	StateCode  uint32
	PID        uint32
	ProcName   string
}

// Listening reports whether this is a listening socket.
func (c Conn) Listening() bool {
	return c.Proto == "udp" || c.Proto == "udp6" || c.StateCode == stateListen
}

// Options filters which connections to load.
type Options struct {
	TCP       bool
	UDP       bool
	Listening bool
	All       bool
	IPv4      bool
	IPv6      bool
	WithProc  bool
	LocalPort int
}

// List returns network connections according to opts.
func List(opts Options) ([]Conn, error) {
	if !opts.TCP && !opts.UDP {
		opts.TCP, opts.UDP = true, true
	}
	if !opts.IPv4 && !opts.IPv6 {
		opts.IPv4, opts.IPv6 = true, true
	}
	if !opts.Listening {
		opts.All = true
	}

	var out []Conn
	var err error
	var rows []Conn
	if opts.TCP && opts.IPv4 {
		rows, err = tcpRows(afInet)
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}
	if opts.TCP && opts.IPv6 {
		rows, err = tcpRows(afInet6)
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}
	if opts.UDP && opts.IPv4 {
		rows, err = udpRows(afInet)
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}
	if opts.UDP && opts.IPv6 {
		rows, err = udpRows(afInet6)
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}

	var nameByPID map[uint32]string
	if opts.WithProc {
		nameByPID = map[uint32]string{}
		if procs, perr := winproc.List(); perr == nil {
			for _, p := range procs {
				nameByPID[p.PID] = p.Name
			}
		}
	}

	filtered := make([]Conn, 0, len(out))
	for _, c := range out {
		if opts.Listening && !c.Listening() {
			continue
		}
		if opts.LocalPort > 0 && int(c.LocalPort) != opts.LocalPort {
			continue
		}
		if opts.WithProc && c.PID != 0 {
			c.ProcName = nameByPID[c.PID]
		}
		filtered = append(filtered, c)
	}
	return filtered, nil
}

func stateName(s uint32) string {
	switch s {
	case stateClosed:
		return "CLOSED"
	case stateListen:
		return "LISTEN"
	case stateSynSent:
		return "SYN_SENT"
	case stateSynReceived:
		return "SYN_RECV"
	case stateEstablished:
		return "ESTABLISHED"
	case stateFinWait1:
		return "FIN_WAIT1"
	case stateFinWait2:
		return "FIN_WAIT2"
	case stateCloseWait:
		return "CLOSE_WAIT"
	case stateClosing:
		return "CLOSING"
	case stateLastAck:
		return "LAST_ACK"
	case stateTimeWait:
		return "TIME_WAIT"
	case stateDeleteTCB:
		return "DELETE_TCB"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

func ntohs(p uint32) uint16 {
	return binary.BigEndian.Uint16([]byte{byte(p), byte(p >> 8)})
}

type mibTCPRowOwnerPID struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPID  uint32
}

type mibTCP6RowOwnerPID struct {
	LocalAddr     [16]byte
	LocalScopeId  uint32
	LocalPort     uint32
	RemoteAddr    [16]byte
	RemoteScopeId uint32
	RemotePort    uint32
	State         uint32
	OwningPID     uint32
}

type mibUDPRowOwnerPID struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPID uint32
}

type mibUDP6RowOwnerPID struct {
	LocalAddr    [16]byte
	LocalScopeId uint32
	LocalPort    uint32
	OwningPID    uint32
}

func tcpRows(af int) ([]Conn, error) {
	var size uint32
	r1, _, err := procGetExtendedTcp.Call(0, uintptr(unsafe.Pointer(&size)), 1, uintptr(af), tcpTableOwnerPIDAll, 0)
	if size == 0 && r1 != 0 && r1 != uintptr(syscall.ERROR_INSUFFICIENT_BUFFER) {
		return nil, fmt.Errorf("GetExtendedTcpTable size: %v", err)
	}
	buf := make([]byte, size)
	r1, _, err = procGetExtendedTcp.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 1, uintptr(af), tcpTableOwnerPIDAll, 0)
	if r1 != 0 {
		return nil, fmt.Errorf("GetExtendedTcpTable: %v", err)
	}
	if len(buf) < 4 {
		return nil, nil
	}
	n := binary.LittleEndian.Uint32(buf[:4])
	out := make([]Conn, 0, n)
	offset := 4
	if af == afInet {
		rowSize := int(unsafe.Sizeof(mibTCPRowOwnerPID{}))
		for i := uint32(0); i < n; i++ {
			if offset+rowSize > len(buf) {
				break
			}
			row := (*mibTCPRowOwnerPID)(unsafe.Pointer(&buf[offset]))
			offset += rowSize
			ip := make(net.IP, 4)
			binary.LittleEndian.PutUint32(ip, row.LocalAddr)
			rip := make(net.IP, 4)
			binary.LittleEndian.PutUint32(rip, row.RemoteAddr)
			out = append(out, Conn{
				Proto:      "tcp",
				LocalIP:    ip,
				LocalPort:  ntohs(row.LocalPort),
				RemoteIP:   rip,
				RemotePort: ntohs(row.RemotePort),
				State:      stateName(row.State),
				StateCode:  row.State,
				PID:        row.OwningPID,
			})
		}
	} else {
		rowSize := int(unsafe.Sizeof(mibTCP6RowOwnerPID{}))
		for i := uint32(0); i < n; i++ {
			if offset+rowSize > len(buf) {
				break
			}
			row := (*mibTCP6RowOwnerPID)(unsafe.Pointer(&buf[offset]))
			offset += rowSize
			lip := make(net.IP, 16)
			copy(lip, row.LocalAddr[:])
			rip := make(net.IP, 16)
			copy(rip, row.RemoteAddr[:])
			out = append(out, Conn{
				Proto:      "tcp6",
				LocalIP:    lip,
				LocalPort:  ntohs(row.LocalPort),
				RemoteIP:   rip,
				RemotePort: ntohs(row.RemotePort),
				State:      stateName(row.State),
				StateCode:  row.State,
				PID:        row.OwningPID,
			})
		}
	}
	return out, nil
}

func udpRows(af int) ([]Conn, error) {
	var size uint32
	r1, _, err := procGetExtendedUdp.Call(0, uintptr(unsafe.Pointer(&size)), 1, uintptr(af), udpTableOwnerPID, 0)
	if size == 0 && r1 != 0 && r1 != uintptr(syscall.ERROR_INSUFFICIENT_BUFFER) {
		return nil, fmt.Errorf("GetExtendedUdpTable size: %v", err)
	}
	buf := make([]byte, size)
	r1, _, err = procGetExtendedUdp.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 1, uintptr(af), udpTableOwnerPID, 0)
	if r1 != 0 {
		return nil, fmt.Errorf("GetExtendedUdpTable: %v", err)
	}
	if len(buf) < 4 {
		return nil, nil
	}
	n := binary.LittleEndian.Uint32(buf[:4])
	out := make([]Conn, 0, n)
	offset := 4
	if af == afInet {
		rowSize := int(unsafe.Sizeof(mibUDPRowOwnerPID{}))
		for i := uint32(0); i < n; i++ {
			if offset+rowSize > len(buf) {
				break
			}
			row := (*mibUDPRowOwnerPID)(unsafe.Pointer(&buf[offset]))
			offset += rowSize
			ip := make(net.IP, 4)
			binary.LittleEndian.PutUint32(ip, row.LocalAddr)
			out = append(out, Conn{
				Proto:     "udp",
				LocalIP:   ip,
				LocalPort: ntohs(row.LocalPort),
				RemoteIP:  net.IPv4zero,
				State:     "UNCONN",
				PID:       row.OwningPID,
			})
		}
	} else {
		rowSize := int(unsafe.Sizeof(mibUDP6RowOwnerPID{}))
		for i := uint32(0); i < n; i++ {
			if offset+rowSize > len(buf) {
				break
			}
			row := (*mibUDP6RowOwnerPID)(unsafe.Pointer(&buf[offset]))
			offset += rowSize
			lip := make(net.IP, 16)
			copy(lip, row.LocalAddr[:])
			out = append(out, Conn{
				Proto:     "udp6",
				LocalIP:   lip,
				LocalPort: ntohs(row.LocalPort),
				RemoteIP:  net.IPv6zero,
				State:     "UNCONN",
				PID:       row.OwningPID,
			})
		}
	}
	return out, nil
}

// FormatAddr formats ip:port for display.
func FormatAddr(ip net.IP, port uint16) string {
	if ip == nil || ip.IsUnspecified() {
		return fmt.Sprintf("*:%d", port)
	}
	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%s:%d", v4.String(), port)
	}
	return fmt.Sprintf("[%s]:%d", ip.String(), port)
}