package toold

import (
	"fmt"
	"io"
	"os"

	//"fmt"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

//ObjectID id
type ObjectID string

//objectIDCounter 计数
var objectIDCounter = uint32(0)

//machineID id
var machineID = readMachineID()

func readMachineID() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {

			//  panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	// fmt.Println("readMachineId:" + string(id))
	return id
}

/*
NewObjectID 创建id
*/
func NewObjectID() ObjectID {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineID[0]
	b[5] = machineID[1]
	b[6] = machineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectID(b[:])
}

/*
Hex hex
*/
func (id ObjectID) Hex() string {
	return hex.EncodeToString([]byte(id))
}

/*
GuuID guuid
*/
func GuuID() string {
	objID := NewObjectID()
	return objID.Hex()
}

//UUID UUID
func UUID() string {
	ul, _ := uuid.NewV4()
	return fmt.Sprintf("%v", ul)
}
