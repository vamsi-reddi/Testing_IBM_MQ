package ibmq

import (
	"log"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)


type IBMQ struct {
	IBMQManager ibmmq.MQQueueManager
	iBMQObject  ibmmq.MQObject
	iBMQMOD     *ibmmq.MQOD
}

func (ibmq *IBMQ)ConnectToQueueManager() bool {
	goQMgrName := "QM1"
	channel := "DEV.APP.SVRCONN"
    connName := "localhost(1414)"
    user := "app"
    password := "passw0rd"

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
	ibmq.IBMQManager, err = ibmmq.Connx(goQMgrName,cno)

	if err != nil {
		log.Println("error connecting to IBMQ: ", err)
		return false
	}

	return ibmq.ConnectToQueue("queue")
}


func (ibmq *IBMQ) ConnectToQueue(queue string) bool {
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


func (ibmq *IBMQ) PutMessageIntoQueue(message []byte) error{

	gopmo := ibmmq.NewMQPMO()
	gopmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

	return ibmq.iBMQObject.Put(nil, gopmo, message)
}


func (ibmq *IBMQ) GetMessageFromQueue() (string, error){

	gmo := ibmmq.NewMQGMO()
	gmo.Options = ibmmq.MQGMO_WAIT | ibmmq.MQGMO_FAIL_IF_QUIESCING
	gmo.WaitInterval = 3 * 1000

	buffer := make([]byte, 1024)

	datalen, err := ibmq.iBMQObject.Get(nil, gmo, buffer)

	if err != nil {
		log.Println("error fetching msgs form queue: ", err)
		return "", err
	}

	return string(buffer[:datalen]), nil
}