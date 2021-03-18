package main

import (
	//"os"
	"github.com/google/uuid"
	//log "github.com/sirupsen/logrus"
)

func uuidgen()(string) {
	id := uuid.New()
	return id.String()
}

//func main(){
//	log.Errorf("[%s][modbuschannel][example_cc][Init] Error starting Simple chaincode",uuidgen())
//}
