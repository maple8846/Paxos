package main

import (
	"flag"
	"log"
	"sync"
	"paxos/pb"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"context"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

var numNodes *int = flag.Int("N", 3, "Number of replicas. Defaults to 3.")
var portnum *int = flag.Int("port", 50051, "Port # to listen on. Defaults to 50051")

type Replica struct {
	Ip   string
	Port int32
}



type MasterServer struct {
	N           int
	ReplicaList []Replica
	lock        *sync.Mutex
	nodes       []*grpc.ClientConn
	alive       []bool
	ready       bool
}

func main() {
	flag.Parse()

	log.Printf("Master starting on port %d\n", *portnum)
	log.Printf("...waiting for %d replicas\n", *numNodes)

	master := &MasterServer{*numNodes,
		make([]Replica, 0, *numNodes),
		new(sync.Mutex),
		make([]*grpc.ClientConn, *numNodes),
		make([]bool, *numNodes),
		false}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServer(s, master)

	if err := s.Serve(lis); err != nil {
		fmt.Printf("gRPC server error", "err", err)
	}
	// Register reflection service on gRPC server.
}

func (s *MasterServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	
	ReplicaNode := &Replica{
		Ip:   req.Ip,
		Port: req.Port,
	}


	nlen := len(s.ReplicaList)
	index := nlen

	for i, peer := range s.ReplicaList {
		if peer.Ip == ReplicaNode.Ip && peer.Port == ReplicaNode.Port {
			index = i
			break
		}
	}

	if index == nlen {
		s.ReplicaList = s.ReplicaList[0 : nlen+1]
		s.ReplicaList[index].Ip = ReplicaNode.Ip
		s.ReplicaList[index].Port = ReplicaNode.Port
		nlen++
	}

	if nlen == s.N {
		s.ready = true

	}

	if s.ready == true {

		var replicaInfo *pb.SingleReplicaInfo
		registerResponse := new(pb.RegisterResponse)

		for i := 0; i < s.N; i++ {
			replicaInfo = new(pb.SingleReplicaInfo)
			replicaInfo.Ip = s.ReplicaList[i].Ip
			replicaInfo.Port = s.ReplicaList[i].Port
			registerResponse.ReplicaInfo = append(registerResponse.ReplicaInfo, replicaInfo)
		}
		return registerResponse, nil
	}

	return &pb.RegisterResponse{},nil
}

func (s *MasterServer) GetReplicaList(ctx context.Context, req *google_protobuf.Empty) (*pb.GetReplicaRequestResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Printf("enter getReplicaList procedure")
	getResponse := new(pb.GetReplicaRequestResponse)

	fmt.Printf("s.N",s.N)

	for i := 0; i < s.N; i++ {
		replicaInfo := new(pb.SingleReplicaInfo)
		replicaInfo.Ip = s.ReplicaList[i].Ip
		replicaInfo.Port = s.ReplicaList[i].Port
		fmt.Printf("replicaInfo",replicaInfo.Ip)
		getResponse.ReplicaInfo = append(getResponse.ReplicaInfo, replicaInfo)
	}

	return getResponse, nil

}
