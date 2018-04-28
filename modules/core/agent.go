package core

import "net"

var AgentHandler Agent

type Agent struct {
	signal chan string
}

func InitAgent() {
	var agent Agent
	agent.startServ()
	AgentHandler = agent
}

func (agent Agent) startServ() {
	go func() {
		options := OptionsHandler.GetOption("Core")
		listener, err := net.Listen("tcp", "localhost:"+options["PortIPC"].(string))
		if err != nil {
			Logger("ListenerErrorIPC::"+err.Error(), "Agent", FATAL)
		}
		defer listener.Close()
		Logger("IPCServiceStarted", "Agent", MSG)
		for {
			conn, err := listener.Accept()
			if err != nil {
				Logger("AcceptErrorIPC::"+err.Error(), "Agent", FATAL)
			}
			go agent.handleRequest(conn)
		}
	}()
}

func (agent Agent) handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		Logger("ReadErrorIPC::"+err.Error(), "Agent", FATAL)
	}
	Logger("RecievedIPCSignal::"+string(buf), "Agent", MSG)
	agent.signal <- string(buf)
	conn.Write([]byte("Saviour:RecievedRequest"))
	conn.Close()
}

func (agent Agent) checkSignal() (bool, string) {
	var sig string
	exists := false
	select {
	case sig = <-agent.signal:
		exists = true
		Logger("RecievedIPCSignal::"+sig, "AGENT", MSG)
	default:
		// Do Nothing
	}
	return exists, sig
}
