package main


import (
	
	"log"
	"paxos/pb"
	"google.golang.org/grpc"
	"fmt"
	"context"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)


const (
    Masteraddress = "localhost:50051"
    Num           = 3
)


type ClientAddress struct {
	Ip     	     string
	Port   	     int32
}




func main() {

    var peeraddress  [Num]ClientAddress


    // Set up a connection to the gRPC server.

    conn, err := grpc.Dial(Masteraddress, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    // Creates a new CustomerClient
    client := pb.NewMasterClient(conn)
        
    para := new(google_protobuf.Empty)
	
	getReplicaRequestResponse, err := client.GetReplicaList(context.Background(),para)

  	if err != nil {
  		log.Fatalf("did not connect: %v", err)
  		return
  	}

	for i:=0;i<len(getReplicaRequestResponse.ReplicaInfo);i++ {
				
		peeraddress[i].Ip = getReplicaRequestResponse.ReplicaInfo[i].Ip
		
		peeraddress[i].Port = getReplicaRequestResponse.ReplicaInfo[i].Port
		
		fmt.Printf("the peer address is%s : %d ", peeraddress[i].Ip,peeraddress[i].Port	)
	
	}
   
}


