package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chefsalim/chameleon/worker"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	if len(os.Args) == 2 {
		fmt.Println("Connecting to ", os.Args[1])
		publisher.Connect(fmt.Sprintf("tcp://%s", os.Args[1]))
	} else {
		fmt.Println("Connecting to localhost:5567")
		publisher.Connect("tcp://localhost:5567")
	}

	//  Ensure subscriber connection has time to complete
	fmt.Println("Waiting for connection to complete")
	time.Sleep(time.Second)

	heartbeat := &worker.Heartbeat{}
	heartbeat.Endpoint = proto.String("12345@chameleon")
	hb_os := worker.Os(worker.Os_Linux)
	heartbeat.Os = &hb_os
	hb_state := worker.WorkerState(worker.WorkerState_Ready)
	heartbeat.State = &hb_state

	//  Send a heartbeat every 30 seconds
	for {
		fmt.Println("Sending heartbeat", heartbeat)
		hb, hb_err := proto.Marshal(heartbeat)
		if hb_err != nil {
			fmt.Println("Failed to marshall heartbeat", hb_err)
		}
		_, err := publisher.SendBytes(hb, 0)
		if err != nil {
			fmt.Println("Failed to send message", err)
		}
		time.Sleep(30 * time.Second)
	}
}
