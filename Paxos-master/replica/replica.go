package main


import (
	"os"
    "gopkg.in/ini.v1"
	"log"
	"paxos/pb"
	"google.golang.org/grpc"
	"fmt"
	"context"
	"command"
	"io/ioutil"
	"unsafe"
)

var filepath ="/home/happy/workplace/src/paxos/config/replica1.conf"
var log = "/home/happy/workplace/src/paxos/log"


const (
    Masteraddress = "localhost:50051"
    Num           = 3
)

type ClientAddress struct {
	Ip     	     string
	Port   	     int32
}


 
type Config struct {   //配置文件要通过tag来指定配置文件中的名称
  Ip string  `ini:"replica_ip"`
  Port int32  `ini:"replica_port"`
}


func ReadConfig(path string	)(*Config ,error) {
	fmt.Printf("load config\n")

	config :=new(Config)
	conf, err := ini.Load(path)   //加载配置文件
	
	if err != nil {
    	log.Println("load config file fail!")
    	return config, err
  	}
  
  	conf.BlockMode = false
  	err = conf.MapTo(config)   //解析成结构体
  
  	if err != nil {
    	
    	log.Println("mapto config file fail!")
    	return config, err
  	
  	}
	
	fmt.Printf("my address = %s\n",config.Ip)
 	
 	return config, nil

}




func RecordCommand(c *command.Command)error {
	
	writebuf := c.Serialize()

	fileObj,err := os.OpenFile(log,os.O_RDONLY|os.O_APPEND|os.O_RDWR,0644)
		

	contents,err := ioutil.ReadAll(fileObj)

    if err != nil {
    	return err
    }

    filelen := len(contents)

	fileObj.WriteAt(writebuf,int64(filelen))
	return nil
}


func GetLastCommand()(*command.Command,error) {
	var cmd command.Command

	
	fileObj,err := os.OpenFile(log,os.O_RDONLY,0644)

 	if err != nil {
    	return nil, err
    }

    defer fileObj.Close()

    contents ,err := ioutil.ReadAll(fileObj)

    if err != nil {
    	return nil, err
    }

    filelen := len(contents)

    readbuf := make([]byte, unsafe.Sizeof(cmd))
	
	lastcmdpos := filelen - int(unsafe.Sizeof(cmd))

	fileObj.ReadAt(readbuf,int64(lastcmdpos))
	
	c := &cmd

	c.Deserialize(readbuf)

	return c,nil
}





func ReadConfig(path string	)(*Config ,error) {
	fmt.Printf("load config\n")

	config :=new(Config)
	conf, err := ini.Load(path)   //加载配置文件
	
	if err != nil {
    	log.Println("load config file fail!")
    	return config, err
  	}
  
  	conf.BlockMode = false
  	err = conf.MapTo(config)   //解析成结构体
  
  	if err != nil {
    	
    	log.Println("mapto config file fail!")
    	return config, err
  	
  	}
	
	fmt.Printf("my address = %s\n",config.Ip)
 	
 	return config, nil

}


func main() {

    var peeraddress  [Num]ClientAddress
    var registerResponse *pb.RegisterResponse
    // Set up a connection to the gRPC server.


    config,err := ReadConfig(filepath)   //也可以通过os.arg或flag从命令行指定配置文件路径
	
	if err != nil {
    	log.Fatal(err)
    	return
  	}


	

    conn, err := grpc.Dial(Masteraddress, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    // Creates a new CustomerClient
    client := pb.NewMasterClient(conn)
    
    for {
		registerResponse, err = client.Register(context.Background(), &pb.RegisterRequest{
			Ip:         config.Ip,
			Port:       config.Port,
		})
		if len(registerResponse.ReplicaInfo) != 0{
			break
		}
	}

	for i:=0;i<len(registerResponse.ReplicaInfo);i++ {
				
		peeraddress[i].Ip = registerResponse.ReplicaInfo[i].Ip
		
		peeraddress[i].Port = registerResponse.ReplicaInfo[i].Port
		
		fmt.Printf("the peer address is%s : %d ", peeraddress[i].Ip,peeraddress[i].Port	)
	
	}

	for {

	}
   
}


func (r *ReplicaServer) Prepare(ctx context.Context, req *google_protobuf.Empty) (*pb.GetReplicaRequestResponse, error) {


}
