package p2p

import "fmt"

type iPeer interface {
	instantiate() bool
	connect() bool
	transfer() bool
	close() bool
	dispose() bool
}

//Peer is the node
type Peer struct {
	name   string
	ip     string
	epoint int
}

//Instantiate returns the newly created peer
func (p *Peer) Instantiate(name string, ip string, epoint int) Peer {
	return Peer{name: name, ip: ip, epoint: epoint}
}

func (p *Peer) connect() bool {
	return true
}

func (p *Peer) transfer() bool {
	return true
}

func (p *Peer) close() bool {
	return true
}

func (p *Peer) dispose() bool {
	return true
}

func main() {
	fmt.Println(total)
}
