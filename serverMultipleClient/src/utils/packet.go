package utils

//global variable declaration
//Usage:
//Peer:
//	Ptype: Peer/backup
//	Pcontent: ""
//	Psize: size of Ptype
//Master:
//	Ptype: Response
//	Pcontent:""
//	Psize: sizeof Ptype
//Client:
//	Ptype: Fetch/Store
//	Pcontent: FileName/Data chunk
//	Psize:""/chunk size

type Packet struct {
	Ptype int
	Pcontent string
	Psize uintptr
}


/*
Usage:
	Master:
		Ptype: Reponse
		Backup: True/False
		NetAddress: Backup peer addr
 */
type Response struct {
	Ptype int
	Backup bool
	NetAddress string
}

/*
Usage:
	Master:
		Ptype: Response
		PrimaryNetAddr: Primary Peer Network
						 address
		BackupNetAddr: Backup Peer Network
						address
 */
type clientResponse struct {
	Ptype int
	PrimaryNetAddr string
	BackupNetAddr string
}


//constants which identify the
//the packet type
const(
	PEER = 1
	FETCH = 2
	STORE = 3
	BACKUP = 4
	RESPONSE = 5
	UPDATE = 6
)

/*
Function creates the packet based on
the input and returns the instace of the
packet
Params:
	p_type: Packet Type(PEER/
						FETCH/
						STORE/
						BACKUP/
						RESPONSE
	content: Content the packet will carry
	p_size: size of the packet

Returns: Instance of packet struct
 */
func CreatePacket(p_type int, content string, p_size uintptr) Packet {
	//create a packet with defined parameters
	packet_t := Packet{Ptype:p_type, Pcontent:content, Psize:p_size}

	//return the created packet
	return packet_t
}

/*
Function creates the Response packet based
on the input and returns the instance of the
response packet struct
Params:
	p_type: Packet Type(PEER/
						FETCH/
						STORE/
						BACKUP/
						RESPONSE
	content: Content the packet will carry
	p_size: size of the packet

Returns: Instance of packet struct
 */
func CreateResponse(p_type int, backup bool, parent string) Response {
	//create a packet with defined parameters
	//and return the created packet
	return Response{Ptype:p_type, Backup:backup, NetAddress:parent}
}
