package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func checkAllowDomain(name string) (bool, []string) {
	for i := 0; i < len(Data); i++ {
		for x := 0; x < len(Data[i].Domain); x++ {
			if strings.HasSuffix(name, Data[i].Domain[x]) {
				return true, Data[i].Remote
			}
		}
	}

	return false, nil
}

func startHTTP() {
	for {
		log.Printf("[APP][HTTP] %v", beginHTTP())

		time.Sleep(time.Second * 4)
	}
}

func beginHTTP() error {
	ln, err := net.Listen("tcp", "127.0.0.1:60080")
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Println("[APP][HTTP] Started")

	for {
		client, err := ln.Accept()
		if err != nil {
			if errno, ok := err.(net.Error); ok {
				if errno.Temporary() {
					continue
				}
			}

			return err
		}

		go handleHTTP(client)
	}
}

func handleHTTP(client net.Conn) {
	defer client.Close()

	data := make([]byte, 1400)
	size, err := client.Read(data)
	if err != nil {
		return
	}
	data = data[:size]

	offset := bytes.Index(data, []byte{0x0d, 0x0a, 0x0d, 0x0a})
	if offset == -1 {
		return
	}

	list := make(map[string]string)

	{
		hdr := bytes.Split(data[:offset], []byte{0x0d, 0x0a})
		for i := 0; i < len(hdr); i++ {
			if i == 0 {
				continue
			}

			SPL := strings.SplitN(string(hdr[i]), ":", 2)
			if len(SPL) < 2 {
				continue
			}

			list[strings.ToUpper(strings.TrimSpace(SPL[0]))] = strings.TrimSpace(SPL[1])
		}
	}

	domain, ok := list["HOST"]
	if !ok {
		return
	}

	var remote net.Conn
	contains, remotes := checkAllowDomain(domain)
	if !contains {
		if remote, err = net.Dial("tcp", net.JoinHostPort(domain, "80")); err != nil {
			return
		}
	} else {
		if remote, err = net.Dial("tcp", remotes[0]); err != nil {
			return
		}

		log.Printf("[APP][HTTP] %s <-> %s", client.RemoteAddr(), domain)
	}
	defer remote.Close()

	if _, err := remote.Write(data); err != nil {
		return
	}
	data = nil

	Pipe(client, remote)
}

func startTLS() {
	for {
		log.Printf("[APP][TLS] %v", beginTLS())

		time.Sleep(time.Second * 4)
	}
}

func beginTLS() error {
	ln, err := net.Listen("tcp", "127.0.0.1:60443")
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Println("[APP][TLS] Started")

	for {
		client, err := ln.Accept()
		if err != nil {
			if errno, ok := err.(net.Error); ok {
				if errno.Temporary() {
					continue
				}
			}

			return err
		}

		go handleTLS(client)
	}
}

func handleTLS(client net.Conn) {
	defer client.Close()

	data := make([]byte, 1400)
	size, err := client.Read(data)
	if err != nil || size <= 44 {
		return
	}
	data = data[:size]

	if data[0] != 0x16 {
		return
	}

	offset := 0
	offset += 1 // Content Type
	offset += 2 // Version
	offset += 2 // Length

	// Handshake Type
	if data[offset] != 0x01 {
		log.Printf("[APP][TLS][%s] Not Client Hello", client.RemoteAddr())
		return
	}
	offset += 1

	offset += 3  // Length
	offset += 2  // Version
	offset += 32 // Random

	// Session ID
	length := int(data[offset])
	offset += 1
	offset += length
	if size <= offset+1 {
		return
	}

	// Cipher Suites
	length = (int(data[offset]) << 8) + int(data[offset+1])
	offset += 2
	offset += length
	if size <= offset {
		return
	}

	// Compression Methods
	length = int(data[offset])
	offset += 1
	offset += length

	// Extension Length
	offset += 2
	if size <= offset+1 {
		return
	}

	domain := ""
	for size > offset+2 && domain == "" {
		// Extension Type
		name := (int(data[offset]) << 8) + int(data[offset+1])
		offset += 2
		if size <= offset+1 {
			return
		}

		// Extension Length
		length = (int(data[offset]) << 8) + int(data[offset+1])
		offset += 2

		// Extension: Server Name
		if name == 0 {
			// Server Name List Length
			offset += 2
			if size <= offset {
				return
			}

			// Server Name Type
			if data[offset] != 0x00 {
				log.Printf("[APP][TLS][%s] Not Host Name", client.RemoteAddr())
				return
			}
			offset += 1
			if size <= offset+1 {
				return
			}

			// Server Name Length
			length = (int(data[offset]) << 8) + int(data[offset+1])
			offset += 2
			if size <= offset+length {
				return
			}

			// Server Name
			domain = string(data[offset : offset+length])

			// Get Out
			break
		}

		// Extension Data
		offset += length
	}

	var remote net.Conn
	contains, remotes := checkAllowDomain(domain)
	if !contains {
		if remote, err = net.Dial("tcp", net.JoinHostPort(domain, "443")); err != nil {
			return
		}
	} else {
		if remote, err = net.Dial("tcp", remotes[1]); err != nil {
			return
		}

		log.Printf("[APP][TLS] %s <-> %s", client.RemoteAddr(), domain)
	}
	defer remote.Close()

	if _, err := remote.Write(data); err != nil {
		return
	}
	data = nil

	Pipe(client, remote)
}

func Pipe(client net.Conn, remote net.Conn) {
	go func() {
		_, _ = io.CopyBuffer(remote, client, make([]byte, 1400))
		_ = client.SetDeadline(time.Now())
		_ = remote.SetDeadline(time.Now())
	}()

	_, _ = io.CopyBuffer(client, remote, make([]byte, 1400))
	_ = client.SetDeadline(time.Now())
	_ = remote.SetDeadline(time.Now())
}
