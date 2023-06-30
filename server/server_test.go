package server

import (
	"dd-log-proxy/logentry"
	"net"
	"testing"

	log "github.com/jlentink/yaglogger"
)

func Test_handleUDPMessage(t *testing.T) {
	channel := make(chan logentry.LogEntry)
	addrs, _ := net.InterfaceAddrs()

	addr := addrs[len(addrs)-1]
	buf := []byte("{\"service\":\"local-server\",\"hostname\":\"869c18ac2fd4\",\"ddsource\":\"monolog\",\"ddtags\":\"version:dev, env:dev\",\"message\":\"2023-06-26T18:25:43+02:00 This is log message 1\",\"level\":\"info\",\"traceid\":\"2247926187253978677\",\"spanid\":\"8024097423951482278\"}")
	go handleUDPMessage(addr, buf, channel)
	logEntry := <-channel
	if logEntry.Service != "local-server" {
		t.Errorf("error!")
	}
}

func Test_handleWrongUDPMessage(t *testing.T) {
	channel := make(chan logentry.LogEntry)
	addrs, _ := net.InterfaceAddrs()

	addr := addrs[len(addrs)-1]
	buf := []byte("this is not json!")
	handleUDPMessage(addr, buf, channel)
	if len(channel) != 0 {
		t.Errorf("error!")
	}
}

func Test_waitForUDPMessage(t *testing.T) {
	udpServer, err := net.ListenPacket("udp", "127.0.0.1:1337")
	defer udpServer.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := net.Dial("udp", "127.0.0.1:1337")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}

	channel := make(chan logentry.LogEntry)

	go waitForUDPMessage(channel, udpServer)

	buf := []byte("{\"service\":\"local-server\",\"hostname\":\"869c18ac2fd4\",\"ddsource\":\"monolog\",\"ddtags\":\"version:dev, env:dev\",\"message\":\"2023-06-26T18:25:43+02:00 This is log message 1\",\"level\":\"info\",\"traceid\":\"2247926187253978677\",\"spanid\":\"8024097423951482278\"}")
	conn.Write(buf)

	<-channel
}
func Test_waitForUDPMessageFailure(t *testing.T) {
	udpServer, err := net.ListenPacket("udp", "127.0.0.1:1337")
	if err != nil {
		log.Fatal(err.Error())
	}

	channel := make(chan logentry.LogEntry)

	go waitForUDPMessage(channel, udpServer)

	udpServer.Close()
}
