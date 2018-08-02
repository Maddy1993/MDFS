package utils

//global variable declaration
type Packet struct {
	Ptype int
	Pcontent string
	Psize uintptr
}

type Response struct {
	Ptype int
	Backup bool
	NetAddress string
}

//constants which identify the
//the packet type
const(
	PEER = 1
	FETCH = 2
	STORE = 3
	BACKUP = 4
	RESPONSE = 5
)

func CreatePacket(p_type int, content string, p_size uintptr) Packet {
	//create a packet with defined parameters
	packet_t := Packet{Ptype:p_type, Pcontent:content, Psize:p_size}

	//return the created packet
	return packet_t
}
