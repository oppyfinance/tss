package storage

import (
	"golang.org/x/crypto/sha3"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-peerstore/addr"
	tnet "github.com/libp2p/go-libp2p-testing/net"
	maddr "github.com/multiformats/go-multiaddr"
	. "gopkg.in/check.v1"

	"github.com/oppyfinance/tss/conversion"
)

type FileStateMgrTestSuite struct{}

var _ = Suite(&FileStateMgrTestSuite{})

func TestPackage(t *testing.T) { TestingT(t) }

func (s *FileStateMgrTestSuite) SetUpTest(c *C) {
	conversion.SetupBech32Prefix()
}

func (s *FileStateMgrTestSuite) TestNewFileStateMgr(c *C) {
	folder := os.TempDir()
	f := filepath.Join(folder, "test", "test1", "test2")
	password := "my password!"
	h := sha3.New256()
	h.Write([]byte(password))
	sk := h.Sum(nil)
	defer func() {
		err := os.RemoveAll(f)
		c.Assert(err, IsNil)
	}()
	fsm, err := NewFileStateMgr(f, sk)
	c.Assert(err, IsNil)
	c.Assert(fsm, NotNil)
	_, err = os.Stat(f)
	c.Assert(err, IsNil)
	fileName, err := fsm.getFilePathName("whatever")
	c.Assert(err, NotNil)

	fileName, err = fsm.getFilePathName("oppypub1addwnpepq2dwek9hkrlxjxadrlmy9fr42gqyq6029q0hked46l3u6a9fxqel6v0rcq9")
	c.Assert(err, IsNil)
	c.Assert(fileName, Equals, filepath.Join(f, "localstate-oppypub1addwnpepq2dwek9hkrlxjxadrlmy9fr42gqyq6029q0hked46l3u6a9fxqel6v0rcq9.json"))
}

func (s *FileStateMgrTestSuite) TestSaveLocalState(c *C) {

	password := "my password!"
	h := sha3.New256()
	h.Write([]byte(password))
	sk := h.Sum(nil)

	stateItem := KeygenLocalState{
		PubKey:    "wasdfasdfasdfasdfasdfasdf",
		LocalData: keygen.NewLocalPartySaveData(5),
		ParticipantKeys: []string{
			"A", "B", "C",
		},
		LocalPartyKey: "A",
	}
	folder := os.TempDir()
	f := filepath.Join(folder, "test", "test1", "test2")
	defer func() {
		err := os.RemoveAll(f)
		c.Assert(err, IsNil)
	}()
	fsm, err := NewFileStateMgr(f, sk)
	c.Assert(err, IsNil)
	c.Assert(fsm, NotNil)
	stateItem.PubKey = "oppypub1addwnpepq2dwek9hkrlxjxadrlmy9fr42gqyq6029q0hked46l3u6a9fxqel6v0rcq9"
	c.Assert(fsm.SaveLocalState(stateItem), IsNil)
	filePathName := filepath.Join(f, "localstate-"+stateItem.PubKey+".json")
	_, err = os.Stat(filePathName)
	c.Assert(err, IsNil)
	item, err := fsm.GetLocalState(stateItem.PubKey)
	c.Assert(err, IsNil)
	c.Assert(reflect.DeepEqual(stateItem, item), Equals, true)
}

func (s *FileStateMgrTestSuite) TestSaveAddressBook(c *C) {
	testAddresses := make(map[peer.ID]addr.AddrList)
	password := "my password!"
	h := sha3.New256()
	h.Write([]byte(password))
	sk := h.Sum(nil)

	var t *testing.T
	id1 := tnet.RandIdentityOrFatal(t)
	id2 := tnet.RandIdentityOrFatal(t)
	id3 := tnet.RandIdentityOrFatal(t)
	mockAddr, err := maddr.NewMultiaddr("/ip4/192.168.3.5/tcp/6668")
	c.Assert(err, IsNil)
	peers := []peer.ID{id1.ID(), id2.ID(), id3.ID()}
	for _, each := range peers {
		testAddresses[each] = []maddr.Multiaddr{mockAddr}
	}
	folder := os.TempDir()
	f := filepath.Join(folder, "test")
	defer func() {
		err := os.RemoveAll(f)
		c.Assert(err, IsNil)
	}()
	fsm, err := NewFileStateMgr(f, sk)
	c.Assert(err, IsNil)
	c.Assert(fsm, NotNil)
	c.Assert(fsm.SaveAddressBook(testAddresses), IsNil)
	filePathName := filepath.Join(f, "address_book.seed")
	_, err = os.Stat(filePathName)
	c.Assert(err, IsNil)
	item, err := fsm.RetrieveP2PAddresses()
	c.Assert(err, IsNil)
	c.Assert(item, HasLen, 3)
}
