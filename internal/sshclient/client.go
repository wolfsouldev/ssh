package sshclient

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// ConnectWithPassword connects to an SSH server using password authentication.
func ConnectWithPassword(host string, port int, user, password string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return connect(host, port, config)
}

// ConnectWithKey connects to an SSH server using private key authentication.
func ConnectWithKey(host string, port int, user, privateKey, keyPassphrase string) error {
	var signer ssh.Signer
	var err error

	if keyPassphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(keyPassphrase))
	} else {
		signer, err = ssh.ParsePrivateKey([]byte(privateKey))
	}
	if err != nil {
		return fmt.Errorf("parsing private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return connect(host, port, config)
}

func connect(host string, port int, config *ssh.ClientConfig) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	defer session.Close()

	return startInteractiveSession(session)
}

func startInteractiveSession(session *ssh.Session) error {
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Put terminal in raw mode
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("setting raw terminal: %w", err)
	}
	defer term.Restore(fd, oldState)

	// Get terminal size
	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24
	}

	// Request PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return fmt.Errorf("requesting PTY: %w", err)
	}

	// Handle window resize
	go handleResize(session, fd)

	if err := session.Shell(); err != nil {
		return fmt.Errorf("starting shell: %w", err)
	}

	return session.Wait()
}

func handleResize(session *ssh.Session, fd int) {
	prevW, prevH, _ := term.GetSize(fd)
	for {
		time.Sleep(1 * time.Second)
		w, h, err := term.GetSize(fd)
		if err != nil {
			continue
		}
		if w != prevW || h != prevH {
			_ = session.WindowChange(h, w)
			prevW, prevH = w, h
		}
	}
}

// TestConnection tests if an SSH connection can be established.
func TestConnection(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 5*1000000000) // 5 seconds
	if err != nil {
		return fmt.Errorf("cannot reach %s: %w", addr, err)
	}
	conn.Close()
	return nil
}
