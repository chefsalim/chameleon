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

	receiver, _ := zmq.NewSocket(zmq.DEALER)
	receiver.SetIdentity("12345@chameleon")

	if len(os.Args) == 2 {
		fmt.Println("Connecting to ", os.Args[1])
		publisher.Connect(fmt.Sprintf("tcp://%s:5567", os.Args[1]))
		receiver.Connect(fmt.Sprintf("tcp://%s:5566", os.Args[1]))
	} else {
		fmt.Println("Connecting to localhost")
		publisher.Connect("tcp://localhost:5567")
		receiver.Connect("tcp://localhost:5566")
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

	go heartbeater(heartbeat, publisher)

	// Wait to receive a job command
	for {
		// Receive empty message
		_, err := receiver.RecvBytes(0)
		if err != nil {
			fmt.Println("Unable to receive first message!")
			continue
		}

		// Receive command
		cmd_b, err := receiver.RecvBytes(0)
		if err != nil {
			fmt.Println("Unable to receive command!")
			continue
		}
		//fmt.Printf("Got cmd bytes (len: %d): %v\n", len(cmd_b), cmd_b)

		cmd := &worker.WorkerCommand{}
		proto.Unmarshal(cmd_b, cmd)
		fmt.Printf("Got cmd: %v\n", cmd)

		// Receive job
		job_b, err := receiver.RecvBytes(0)
		if err != nil {
			fmt.Println("Unable to receive job!")
			continue
		}
		//fmt.Printf("Got job bytes (len: %d): %v\n", len(job_b), job_b)

		job := &worker.Job{}
		proto.Unmarshal(job_b, job)
		fmt.Printf("Got job: %v\n", job)
	}
}

func heartbeater(heartbeat *worker.Heartbeat, publisher *zmq.Socket) {
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
