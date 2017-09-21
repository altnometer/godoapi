package support

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func getSSHConfig(userName string) (*ssh.ClientConfig, error) {
	// key, err := getSSHKey(sshKeyPath)
	authMeth, err := getAuthMeth()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("key = %+v\n", key)
	confPtr := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			// ssh.PublicKeys(key),
			authMeth,
		},
		HostKeyCallback: keyPrint,
	}
	confPtr.SetDefaults()
	return confPtr, nil

}

func getSSHKey(sshKeyPath string) (ssh.Signer, error) {
	buffer, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return key, nil
}
func getAuthMeth() (ssh.AuthMethod, error) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		fmt.Printf("Did you run %s?\n",
			RedSp("eval $(ssh-agent -s) && ssh-add ~/.ssh/privkey"),
		)
		return nil, err
	}
	sshAuthMeth := ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	return sshAuthMeth, nil
}

func keyPrint(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	// fmt.Printf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}

// GetSSHClient return *ssh.Client
func GetSSHClient(userName, host string) (*ssh.Client, error) {
	sshConfig, err := getSSHConfig(userName)
	if err != nil {
		return nil, err
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, 22), sshConfig)
	// fmt.Printf("client.User() = %+v\n", client.User())
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}
	return client, nil
}

// GetSSHSession return ssh session. Don't forget to session.Close().
// Run command with session.Run("/usr/bin/whoami").
// Read output with var b bytes.Buffer, session.Stdout = &b,
// fmt.Println(b.String())
func GetSSHSession(sshClient *ssh.Client) (*ssh.Session, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}
	return session, nil
}

// GetSSHInterSession return interactive ssh session.
// Close() session after use.
// Set env var with err = session.Setenv("LC_USR_DIR", "/usr").
// Run commands with err = session.Run("ls -l $LC_USR_DIR").
func GetSSHInterSession(sshClient *ssh.Client) (*ssh.Session, error) {
	session, err := GetSSHSession(sshClient)
	if err != nil {
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)
	return session, nil
}
