// go run Exercise4.go

package main

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
)

type Message struct {
    ID string
    Word string
   	CurrTime time.Time
    LocalIP string
    RemoteIP string
    RawWord string
}

func UDP_receive(port string, receiveCh chan Message)(err error) {
	
	baddr, err := net.ResolveUDPAddr("udp",":"+port)
	if err != nil {return err}

	localListenConn, err := net.ListenUDP("udp", baddr)
	if err != nil {return err}
	
	var receiveMessage Message
	
	for {
		buffer := make([]byte, 2048)
		n,addr,err := localListenConn.ReadFromUDP(buffer[0:])
		if err != nil {fmt.Println(err); return err}
		//fmt.Println(string(buffer))
		err = json.Unmarshal(buffer[:n], &receiveMessage)
		receiveMessage.LocalIP = addr.String()
		receiveMessage.CurrTime = time.Now()
		if err != nil {
			fmt.Println(err)
			return err
		}
		receiveMessage.RawWord = string(buffer)
		receiveCh <- receiveMessage
	}
}

func UDP_broadcast(baddr string, sendCh chan string) (error){

	tempConn, err := net.Dial("udp", baddr)
	if err != nil {return err}
	
	var msg Message
	msg.ID = "1"
	msg.Word = <- sendCh
	msg.RemoteIP = baddr
	
	buffer, err := json.Marshal(msg)
	if err != nil {return err}
	
	for{
		tempConn.Write([]byte(buffer))
		time.Sleep(100*time.Millisecond)
	}
}

func main() {
	receiveChannel := make(chan Message, 1024)
	sendChannel := make(chan string, 1024)
	//message := Message{}
	go UDP_broadcast("129.241.187.255:24568", sendChannel)
	go UDP_receive("24568", receiveChannel)
	
	time.Sleep(100*time.Millisecond)

	for {
		sendChannel <-"NOt generic"
		i := <- receiveChannel
		fmt.Println("\n\nMessage received on: ", i.CurrTime)
		fmt.Println("\nMessage ID was: ", i.ID)
		fmt.Println("\nMessage contents: ", i.Word)
		fmt.Println("\nLocal IP was: ", i.LocalIP)
		fmt.Println("\nRemote IP was: ", i.RemoteIP)
		fmt.Println("\nRaw contents: ", i.RawWord)
		fmt.Println("__________________________\n")
		time.Sleep(100*time.Millisecond)
	}
}

