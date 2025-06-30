package ibmq

import (
	"fmt"
	"log"
	"sync"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

type IBMQ struct {
	IBMQManager ibmmq.MQQueueManager
	iBMQObject  ibmmq.MQObject
	iBMQMOD     *ibmmq.MQOD
	mu          sync.Mutex
	isConnected bool
}

// Close properly closes all IBM MQ connections and resources
func (ibmq *IBMQ) Close() error {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	var errors []error

	// Close the queue object if it exists
	if ibmq.isConnected && ibmq.iBMQObject != (ibmmq.MQObject{}) {
		if err := ibmq.iBMQObject.Close(0); err != nil {
			log.Printf("Error closing queue object: %v", err)
			errors = append(errors, err)
		} else {
			log.Println("Queue object closed successfully")
		}
		ibmq.iBMQObject = ibmmq.MQObject{}
	}

	// Disconnect from the queue manager if connected
	if ibmq.isConnected && ibmq.IBMQManager != (ibmmq.MQQueueManager{}) {
		if err := ibmq.IBMQManager.Disc(); err != nil {
			log.Printf("Error disconnecting from queue manager: %v", err)
			errors = append(errors, err)
		} else {
			log.Println("Disconnected from queue manager successfully")
		}
		ibmq.IBMQManager = ibmmq.MQQueueManager{}
		ibmq.isConnected = false
	}

	// Clear the queue descriptor
	ibmq.iBMQMOD = nil

	if len(errors) > 0 {
		return errors[0] // Return the first error for simplicity
	}
	return nil
}

// IsConnected returns whether the IBM MQ connection is active
func (ibmq *IBMQ) IsConnected() bool {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()
	return ibmq.isConnected && ibmq.IBMQManager != (ibmmq.MQQueueManager{})
}

func (ibmq *IBMQ) ConnectToQueueManager() bool {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	// If already connected, return true
	if ibmq.isConnected && ibmq.IBMQManager != (ibmmq.MQQueueManager{}) {
		return true
	}

	goQMgrName := ""
	channel := ""
	connName := ""
	user := ""
	password := ""

	cno := ibmmq.NewMQCNO()

	cd := ibmmq.NewMQCD()
	cd.ChannelName = channel
	cd.ConnectionName = connName
	cno.ClientConn = cd

	csp := ibmmq.NewMQCSP()
	csp.UserId = user
	csp.Password = password
	cno.SecurityParms = csp

	var err error
	ibmq.IBMQManager, err = ibmmq.Connx(goQMgrName, cno)

	if err != nil {
		log.Println("error connecting to IBMQ: ", err)
		return false
	}

	ibmq.isConnected = true
	return ibmq.ConnectToQueue("queue")
}

func (ibmq *IBMQ) ConnectToQueue(queue string) bool {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	// Close existing queue object if any
	if ibmq.isConnected && ibmq.iBMQObject != (ibmmq.MQObject{}) {
		if err := ibmq.iBMQObject.Close(0); err != nil {
			log.Printf("Error closing existing queue object: %v", err)
		}
		ibmq.iBMQObject = ibmmq.MQObject{}
	}

	ibmq.iBMQMOD = ibmmq.NewMQOD()
	ibmq.iBMQMOD.ObjectType = ibmmq.MQOT_Q
	ibmq.iBMQMOD.ObjectName = queue

	openOptions := ibmmq.MQOO_OUTPUT

	var err error
	ibmq.iBMQObject, err = ibmq.IBMQManager.Open(ibmq.iBMQMOD, openOptions)

	if err != nil {
		log.Println("failed to connect to queue: ", err)
		return false
	}

	return true
}

func (ibmq *IBMQ) PutMessageIntoQueue(message []byte) error {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	if !ibmq.isConnected || ibmq.iBMQObject == (ibmmq.MQObject{}) {
		return fmt.Errorf("not connected to queue")
	}

	gopmo := ibmmq.NewMQPMO()
	gopmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

	return ibmq.iBMQObject.Put(nil, gopmo, message)
}

func (ibmq *IBMQ) GetMessageFromQueue() (string, error) {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	if !ibmq.isConnected || ibmq.iBMQObject == (ibmmq.MQObject{}) {
		return "", fmt.Errorf("not connected to queue")
	}

	gmo := ibmmq.NewMQGMO()
	gmo.Options = ibmmq.MQGMO_WAIT | ibmmq.MQGMO_FAIL_IF_QUIESCING
	gmo.WaitInterval = 3 * 1000

	buffer := make([]byte, 1024)

	datalen, err := ibmq.iBMQObject.Get(nil, gmo, buffer)

	if err != nil {
		log.Println("error fetching msgs from queue: ", err)
		return "", err
	}

	return string(buffer[:datalen]), nil
}

// Reconnect attempts to reconnect to the queue manager
func (ibmq *IBMQ) Reconnect() error {
	ibmq.mu.Lock()
	defer ibmq.mu.Unlock()

	// Close existing connections
	if err := ibmq.closeInternal(); err != nil {
		log.Printf("Error during cleanup before reconnect: %v", err)
	}

	// Attempt to reconnect
	if ibmq.ConnectToQueueManager() {
		return nil
	}
	return fmt.Errorf("failed to reconnect to queue manager")
}

// closeInternal is an internal method for cleanup without mutex (used by Reconnect)
func (ibmq *IBMQ) closeInternal() error {
	var errors []error

	if ibmq.isConnected && ibmq.iBMQObject != (ibmmq.MQObject{}) {
		if err := ibmq.iBMQObject.Close(0); err != nil {
			errors = append(errors, err)
		}
		ibmq.iBMQObject = ibmmq.MQObject{}
	}

	if ibmq.isConnected && ibmq.IBMQManager != (ibmmq.MQQueueManager{}) {
		if err := ibmq.IBMQManager.Disc(); err != nil {
			errors = append(errors, err)
		}
		ibmq.IBMQManager = ibmmq.MQQueueManager{}
		ibmq.isConnected = false
	}

	ibmq.iBMQMOD = nil

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}
