// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tchannel

import (
	"container/heap"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
)

var (
	// ErrInvalidConnectionState indicates that the connection is not in a valid state.
	ErrInvalidConnectionState = errors.New("connection is in an invalid state")

	// ErrNoPeers indicates that there are no peers.
	ErrNoPeers = errors.New("no peers available")

	// ErrPeerNotFound indicates that the specified peer was not found.
	ErrPeerNotFound = errors.New("peer not found")

	peerRng = NewRand(time.Now().UnixNano())
)

// Connectable is the interface used by peers to create connections.
type Connectable interface {
	// Connect tries to connect to the given hostPort.
	Connect(ctx context.Context, hostPort string) (*Connection, error)
	// Logger returns the logger to use.
	Logger() Logger
}

// PeerList maintains a list of Peers.
type PeerList struct {
	sync.RWMutex

	parent          *RootPeerList
	peersByHostPort map[string]*peerScore
	peerHeap        *peerHeap
	scoreCalculator ScoreCalculator
	lastSelected    uint64
}

func newPeerList(root *RootPeerList) *PeerList {
	return &PeerList{
		parent:          root,
		peersByHostPort: make(map[string]*peerScore),
		scoreCalculator: newPreferIncomingCalculator(),
		peerHeap:        newPeerHeap(),
	}
}

// SetStrategy sets customized peer selection stratedgy.
func (l *PeerList) SetStrategy(sc ScoreCalculator) {
	l.scoreCalculator = sc
}

// Siblings don't share peer lists (though they take care not to double-connect
// to the same hosts).
func (l *PeerList) newSibling() *PeerList {
	sib := newPeerList(l.parent)
	return sib
}

// Add adds a peer to the list if it does not exist, or returns any existing peer.
func (l *PeerList) Add(hostPort string) *Peer {
	if ps, _, ok := l.exists(hostPort); ok {
		return ps.Peer
	}
	l.Lock()
	defer l.Unlock()

	if p, ok := l.peersByHostPort[hostPort]; ok {
		return p.Peer
	}

	p := l.parent.Add(hostPort)
	p.addSC()
	ps := newPeerScore(p, l.scoreCalculator.GetScore(p))

	l.peersByHostPort[hostPort] = ps
	l.peerHeap.addPeer(ps)
	return p
}

// Get returns a peer from the peer list, or nil if none can be found.
func (l *PeerList) Get(prevSelected map[string]struct{}) (*Peer, error) {
	l.Lock()
	if l.peerHeap.Len() == 0 {
		l.Unlock()
		return nil, ErrNoPeers
	}

	// Select a peer, avoiding previously selected peers. If all peers have been previously
	// selected, then it's OK to repick them.
	peer := l.choosePeer(prevSelected, true /* avoidHost */)
	if peer == nil {
		peer = l.choosePeer(prevSelected, false /* avoidHost */)
	}
	if peer == nil {
		peer = l.choosePeer(nil, false /* avoidHost */)
	}
	l.Unlock()
	return peer, nil
}

// Remove removes a peer from the peer list. It returns an error if the peer cannot be found.
// Remove does not affect connections to the peer in any way.
func (l *PeerList) Remove(hostPort string) error {
	l.Lock()
	defer l.Unlock()

	p, ok := l.peersByHostPort[hostPort]
	if !ok {
		return ErrPeerNotFound
	}

	p.delSC()
	delete(l.peersByHostPort, hostPort)
	l.peerHeap.removePeer(p)

	return nil
}
func (l *PeerList) choosePeer(prevSelected map[string]struct{}, avoidHost bool) *Peer {
	var psPopList []*peerScore
	var ps *peerScore

	canChoosePeer := func(hostPort string) bool {
		if _, ok := prevSelected[hostPort]; ok {
			return false
		}
		if avoidHost {
			if _, ok := prevSelected[getHost(hostPort)]; ok {
				return false
			}
		}
		return true
	}

	size := l.peerHeap.Len()
	for i := 0; i < size; i++ {
		popped := l.peerHeap.popPeer()

		if canChoosePeer(popped.HostPort()) {
			ps = popped
			break
		}
		psPopList = append(psPopList, popped)
	}

	for _, p := range psPopList {
		heap.Push(l.peerHeap, p)
	}

	if ps == nil {
		return nil
	}

	l.peerHeap.pushPeer(ps)
	atomic.AddUint64(&ps.chosenCount, 1)
	return ps.Peer
}

// GetOrAdd returns a peer for the given hostPort, creating one if it doesn't yet exist.
func (l *PeerList) GetOrAdd(hostPort string) *Peer {
	if ps, _, ok := l.exists(hostPort); ok {
		return ps.Peer
	}
	return l.Add(hostPort)
}

// Copy returns a map of the peer list. This method should only be used for testing.
func (l *PeerList) Copy() map[string]*Peer {
	l.RLock()
	defer l.RUnlock()

	listCopy := make(map[string]*Peer)
	for k, v := range l.peersByHostPort {
		listCopy[k] = v.Peer
	}
	return listCopy
}

// exists checks if a hostport exists in the peer list.
func (l *PeerList) exists(hostPort string) (*peerScore, uint64, bool) {
	var score uint64

	l.RLock()
	ps, ok := l.peersByHostPort[hostPort]
	if ok {
		score = ps.score
	}
	l.RUnlock()

	return ps, score, ok
}

// updatePeer is called when there is a change that may cause the peer's score to change.
// The new score is calculated, and the peer heap is updated with the new score if the score changes.
func (l *PeerList) updatePeer(p *Peer) {
	ps, psScore, ok := l.exists(p.hostPort)
	if !ok {
		return
	}

	newScore := l.scoreCalculator.GetScore(p)
	if newScore == psScore {
		return
	}

	l.Lock()
	ps.score = newScore
	l.peerHeap.updatePeer(ps)
	l.Unlock()
}

// peerScore represents a peer and scoring for the peer heap.
// It is not safe for concurrent access, it should only be used through the PeerList.
type peerScore struct {
	*Peer

	score uint64
	index int
	order uint64
}

func newPeerScore(p *Peer, score uint64) *peerScore {
	return &peerScore{
		Peer:  p,
		score: score,
		index: -1,
	}
}

// Peer represents a single autobahn service or client with a unique host:port.
type Peer struct {
	sync.RWMutex

	channel      Connectable
	hostPort     string
	onConnChange func(*Peer)

	// scCount is the number of subchannels that this peer is added to.
	scCount uint32

	// connections are mutable, and are protected by the mutex.
	inboundConnections  []*Connection
	outboundConnections []*Connection
	chosenCount         uint64

	// onUpdate is a test-only hook.
	onUpdate func(*Peer)
}

func newPeer(channel Connectable, hostPort string, onConnChange func(*Peer)) *Peer {
	if hostPort == "" {
		panic("Cannot create peer with blank hostPort")
	}
	return &Peer{
		channel:      channel,
		hostPort:     hostPort,
		onConnChange: onConnChange,
	}
}

// HostPort returns the host:port used to connect to this peer.
func (p *Peer) HostPort() string {
	return p.hostPort
}

// getActive returns a list of active connections.
// TODO(prashant): Should we clear inactive connections?
func (p *Peer) getActive() []*Connection {
	var active []*Connection
	p.runWithConnections(func(c *Connection) {
		if c.IsActive() {
			active = append(active, c)
		}
	})
	return active
}

func randConn(conns []*Connection) *Connection {
	return conns[peerRng.Intn(len(conns))]
}

// GetConnection returns an active connection to this peer. If no active connections
// are found, it will create a new outbound connection and return it.
func (p *Peer) GetConnection(ctx context.Context) (*Connection, error) {
	// TODO(prashant): Use some sort of scoring to pick a connection.
	if activeConns := p.getActive(); len(activeConns) > 0 {
		return randConn(activeConns), nil
	}

	// No active connections, make a new outgoing connection.
	c, err := p.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// AddInboundConnection adds an active inbound connection to the peer's connection list.
// If a connection is not active, ErrInvalidConnectionState will be returned.
func (p *Peer) AddInboundConnection(c *Connection) error {
	switch c.readState() {
	case connectionActive, connectionStartClose:
		// TODO(prashantv): Block inbound connections when the connection is not active.
		break
	default:
		return ErrInvalidConnectionState
	}

	p.Lock()
	p.inboundConnections = append(p.inboundConnections, c)
	p.Unlock()

	p.connectionStateChanged(c)
	return nil
}

// addSC adds a reference to a peer from a subchannel (e.g. peer list).
func (p *Peer) addSC() {
	p.Lock()
	p.scCount++
	p.Unlock()
}

// delSC removes a reference to a peer from a subchannel (e.g. peer list).
func (p *Peer) delSC() {
	p.Lock()
	p.scCount--
	p.Unlock()
}

// canRemove returns whether this peer can be safely removed from the root peer list.
func (p *Peer) canRemove() bool {
	p.RLock()
	count := len(p.inboundConnections) + len(p.outboundConnections) + int(p.scCount)
	p.RUnlock()
	return count == 0
}

// AddOutboundConnection adds an active outbound connection to the peer's connection list.
// If a connection is not active, ErrInvalidConnectionState will be returned.
func (p *Peer) AddOutboundConnection(c *Connection) error {
	switch c.readState() {
	case connectionActive, connectionStartClose:
		break
	default:
		return ErrInvalidConnectionState
	}

	p.Lock()
	p.outboundConnections = append(p.outboundConnections, c)
	p.Unlock()

	p.connectionStateChanged(c)
	return nil
}

// checkInboundConnection will check whether the changed connection is an inbound
// connection, and will remove any closed connections.
func (p *Peer) checkInboundConnection(changed *Connection) (updated bool, isInbound bool) {
	newConns := p.inboundConnections[:0]
	for _, c := range p.inboundConnections {
		if c == changed {
			isInbound = true
		}

		if c.readState() != connectionClosed {
			newConns = append(newConns, c)
		} else {
			updated = true
		}
	}
	if updated {
		p.inboundConnections = newConns
	}

	return updated, isInbound
}

// connectionStateChanged is called when one of the peers' connections states changes.
func (p *Peer) connectionStateChanged(changed *Connection) {
	p.Lock()
	updated, _ := p.checkInboundConnection(changed)
	p.Unlock()

	if updated {
		p.onConnChange(p)
	}
}

// Connect adds a new outbound connection to the peer.
func (p *Peer) Connect(ctx context.Context) (*Connection, error) {
	c, err := p.channel.Connect(ctx, p.hostPort)
	if err != nil {
		return nil, err
	}

	if err := p.AddOutboundConnection(c); err != nil {
		return nil, err
	}

	return c, nil
}

// BeginCall starts a new call to this specific peer, returning an OutboundCall that can
// be used to write the arguments of the call.
func (p *Peer) BeginCall(ctx context.Context, serviceName, methodName string, callOptions *CallOptions) (*OutboundCall, error) {
	if callOptions == nil {
		callOptions = defaultCallOptions
	}
	callOptions.RequestState.AddSelectedPeer(p.HostPort())

	conn, err := p.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	if callOptions == nil {
		callOptions = defaultCallOptions
	}
	call, err := conn.beginCall(ctx, serviceName, methodName, callOptions)
	if err != nil {
		return nil, err
	}

	return call, err
}

// NumConnections returns the number of inbound and outbound connections for this peer.
func (p *Peer) NumConnections() (inbound int, outbound int) {
	p.RLock()
	inbound = len(p.inboundConnections)
	outbound = len(p.outboundConnections)
	p.RUnlock()
	return inbound, outbound
}

// NumPendingOutbound returns the number of pending outbound calls.
func (p *Peer) NumPendingOutbound() int {
	count := 0
	p.RLock()
	for _, c := range p.outboundConnections {
		count += c.outbound.count()
	}

	for _, c := range p.inboundConnections {
		count += c.outbound.count()
	}
	p.RUnlock()
	return count
}

func (p *Peer) runWithConnections(f func(*Connection)) {
	p.RLock()
	for _, c := range p.inboundConnections {
		f(c)
	}

	for _, c := range p.outboundConnections {
		f(c)
	}
	p.RUnlock()
}

func (p *Peer) callOnUpdateComplete() {
	p.RLock()
	f := p.onUpdate
	p.RUnlock()

	if f != nil {
		f(p)
	}
}
