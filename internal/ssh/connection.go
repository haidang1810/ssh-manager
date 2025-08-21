package ssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"ssh-manager/internal/models"
	"ssh-manager/internal/utils"
)

// Connect establishes an interactive SSH session to a remote server.
func Connect(conn *models.Connection) error {
	var authMethods []ssh.AuthMethod

	// Decrypt password if it exists

decryptedPassword := conn.Password
	if decryptedPassword != "" {
		d, err := utils.Decrypt(decryptedPassword)
		if err == nil {
			decryptedPassword = d
		} else {
			// If decryption fails, assume it's a plaintext password (for backward compatibility)
			fmt.Fprintf(os.Stderr, "Warning: Failed to decrypt password for %s. Assuming plaintext. Error: %v\n", conn.Name, err)
		}
	}

	if conn.KeyPath != "" {
		key, err := ioutil.ReadFile(conn.KeyPath)
		if err != nil {
			return fmt.Errorf("unable to read private key: %w", err)
		}

		var signer ssh.Signer
		var parseErr error

		// Try parsing without a passphrase first
		signer, parseErr = ssh.ParsePrivateKey(key)

		// If parsing without passphrase fails and it's due to passphrase protection, prompt for it
		if parseErr != nil && strings.Contains(parseErr.Error(), "passphrase protected") {
			fmt.Print("Enter passphrase for private key: ")
			bytePassphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println() // Newline after password input
			if err != nil {
				return fmt.Errorf("failed to read passphrase: %w", err)
			}
			passphrase := string(bytePassphrase)

			signer, parseErr = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
		}

		// If there's still an error after trying with/without passphrase, return it
		if parseErr != nil {
			return fmt.Errorf("unable to parse private key: %w", parseErr)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if decryptedPassword != "" {
		authMethods = append(authMethods, ssh.Password(decryptedPassword))
	}

	sshConfig := &ssh.ClientConfig{
		User: conn.User,
		Auth: authMethods,
		// IMPORTANT: In a real-world application, you should not ignore host key validation.
		// This is a security risk. For this project, we will add it for convenience,
		// but it should be replaced with a proper host key callback.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conn.Host, conn.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up terminal modes
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to make terminal raw: %w", err)
	}
	defer terminal.Restore(fd, oldState)

	// Set up standard I/O
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Request PTY
	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	if err := session.RequestPty("xterm-256color", termHeight, termWidth, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("failed to request pty: %w", err)
	}

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Wait for session to finish
	return session.Wait()
}