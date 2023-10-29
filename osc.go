package osc

import (
	"net"
	"time"
)

//func Dial(localPort int, remoteAddr string) (conn net.Conn, err error) {
func Dial(localAddr string, remoteAddr string) (conn net.Conn, err error) {
	// Takes a local and remote address and returns a net.Conn
	//     addresses should be provided in the form: "ip:port"
	//     LocalAddr may be in the form ":port"
	laddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return conn, err
	}
	dialer := &net.Dialer{
		//LocalAddr: &net.UDPAddr{
		//Port: localPort,
		//},
		LocalAddr: laddr,
		Timeout:   10 * time.Second,
	}
	conn, err = dialer.Dial("udp", remoteAddr)
	return conn, err
}

func Inquire(conn net.Conn, msg Message, timeout time.Duration) (reply Message, err error) {
	// Takes a Conn and an osc Message
	//   Sends the message to a server, and listens for a response
	//   Returns the responding Message

	// The given timeout will define the deadlines for both sending and listening
	callTime := time.Now()

	// Send message
	err = Send(conn, msg, timeout)
	if err != nil {
		return reply, err
	}

	// Redefine timeout for listening
	timeout = timeout - time.Since(callTime)

	// Wait for reply
	reply, err = Listen(conn, timeout)
	if err != nil {
		return reply, err
	}

	return reply, nil
}

func Send(conn net.Conn, msg Message, timeout time.Duration) error {
	// Send an OSC message of type Message to the Conn connection

	// Set deadline from the given timeout
	conn.SetWriteDeadline(time.Now().Add(timeout))

	// Make the packet from the components if it doesn't already exist
	if msg.Packet.Len() == 0 {
		err := msg.MakePacket()
		if err != nil {
			return err
		}
	}

	// Write the bytes to the connection
	_, err := conn.Write(msg.Packet.Bytes())
	return err
}

func Listen(conn net.Conn, timeout time.Duration) (msg Message, err error) {
	// Act as a server and listen for an incoming OSC message

	// Set deadline from the given timeout
	conn.SetReadDeadline(time.Now().Add(timeout))

	// Make a []byte of length 512 and read into it
	byt := make([]byte, 512)
	_, err = conn.Read(byt)
	if err != nil {
		return msg, err
	}

	// Write bytes to packet
	msg.Packet.Write(byt)

	// Parse the []byte in msg.Packet and populate the properties of msg
	//     The incoming arguments will be decoded in ParseMessage() method
	err = msg.ParseMessage()

	// Return msg
	return msg, err
}
