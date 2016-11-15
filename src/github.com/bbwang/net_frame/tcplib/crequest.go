/**
common request
**/
package tcplib

type CommHeader struct {
	uint32 length
	uint32 cmd
	uint32 seq
}

func (this *CommHeader) GetLen() uint32 {
	return this.length
}

func (this *CommHeader) GetCmd() uint32 {
	return this.cmd
}

func (this *CommHeader) GetSeq() uint32 {
	return this.seq
}
