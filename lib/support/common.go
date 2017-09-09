package support

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/fatih/color"
)

/////////////////////////////consts///////////////////////////////////////////

// VolByIDPrefix holds prefix for vol link names in /dev/disk/by-id dir.
const VolByIDPrefix string = "scsi-0DO_Volume_"

/////////////////////////////errors///////////////////////////////////////////

// ErrBadArgs used when no, or not all args provided.
var ErrBadArgs = errors.New("no, not enough, or wrong args provided")

// ErrUserSaysQuit used when user input signals stop exucting the program.
var ErrUserSaysQuit = errors.New("user says stop execution")

// ErrVolAttached used when vol shoul not be attached to a droplet.
type ErrVolAttached struct {
	Msg string
}

func (e ErrVolAttached) Error() string {
	return e.Msg
}

// YellowSp colors str.
var YellowSp = color.New(color.FgYellow).SprintFunc()

// YellowLn colors a line.
var YellowLn = color.New(color.FgYellow).PrintlnFunc()

// YellowPf color output with interpolation.
var YellowPf = color.New(color.FgYellow).PrintfFunc()

// GreenLn colors a line red.
var GreenLn = color.New(color.FgGreen).PrintlnFunc()

// RedPf output with interpolation colored red.
var RedPf = color.New(color.FgRed).PrintfFunc()

// RedLn colors a line red.
var RedLn = color.New(color.FgRed).PrintlnFunc()

// RedBold colors the output.
var RedBold = color.New(color.FgRed, color.Bold)

// Red colors the output.
var Red = color.New(color.FgRed)

/////////////// Common droplet properties//////////////////////////////////////

// MaxDroplets mas allowed droplets. Used in declaring slices of droplets.
var MaxDroplets = 8

// SSHKeys for the droplet.
var SSHKeys = godo.DropletCreateSSHKey{Fingerprint: "2f:2a:c4:eb:ec:38:35:cd:2a:d9:65:cf:59:12:df:44"}

// DropletSizes holds valid values.
var DropletSizes = []string{"512mb", "1gb", "2gb"}

// NameList to hold custom flag value for multiple names.
type NameList []string

func (n *NameList) String() string {
	return fmt.Sprintf("%v", *n)
}

// Set parsed the flag value
func (n *NameList) Set(value string) error {
	*n = strings.Split(value, ",")
	return nil
}

// Value interface for cusmom flag argument.
type Value interface {
	String() string
	Set(string) error
}

// ValidateRegions prints out error msg and exits for invalid regions.
func ValidateRegions(regPtr *string) error {
	regions := map[string]bool{"fra1": true}
	if _, valExist := regions[*regPtr]; !valExist {
		keys := make([]string, len(regions))
		i := 0
		for k := range regions {
			keys[i] = k
			i++
		}
		fmt.Printf("valid choices for region field are: %+v\n", keys)
		return ErrBadArgs
	}
	return nil
}

/////////////////////Supporting functions//////////////////////////////////////

// GetUserInput query for a user input and return int.
func GetUserInput(promt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promt)
	text, err := reader.ReadString('\n')
	// convert CRLF to LF
	if err != nil {
		panic("Cannot read user inut")
	}
	text = strings.Replace(text, "\n", "", -1)
	if text == "" {
		RedLn("Input is required!")
		os.Exit(1)
	}
	return strings.Replace(text, "\n", "", -1)
}

// UserConfirmDefaultN return false for any input but 'y|Y'
func UserConfirmDefaultN(prompt string) (conf bool, err error) {
	fmt.Println(prompt)
	fmt.Println("[y/N]")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return false, err
	}
	if char == 10 {
		return false, nil
	}
	if char != 'y' && char != 'Y' {
		return false, nil
	}
	return true, nil
}

// GetSSHKeyPath returns user provided path to ssh keys.
func GetSSHKeyPath() string {
	sshKeyPath := os.Getenv("DOSSHKeyPath")
	if sshKeyPath == "" {
		YellowLn("You can set env var DOSSHKeyPath!")
		sshKeyPath = GetUserInput("Type in a DOSSHKeyPath: ")
	}
	return sshKeyPath
}

// FetchSSHOutput executes ssh cmd and returns cmd output.
func FetchSSHOutput(userName, IP, sshKeyPath string, sshCmds []string) string {
	cmdArgs := append(
		[]string{
			"-o",
			"UserKnownHostsFile=/dev/null",
			"-o",
			"StrictHostKeyChecking=no",
			"-i",
			sshKeyPath,
			fmt.Sprintf("%s@%s", userName, IP),
		},
		sshCmds...,
	)
	// str := strings.Join(cmdArgs[:], " ")
	// fmt.Printf("str = %+v\n", str)
	// os.Exit(0)
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "ssh"
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		RedLn("There was an error running ssh command: ", err)
		os.Exit(1)
	}
	return string(cmdOut)
}

// ExecSSH execute ssh command interactively.
func ExecSSH(userName, IP string, cmds []string) {
	// sshCommander := SSHCommander{userName, IP, sshKeyPath}
	// sshOpt := fmt.Sprintf("-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i %s", sshKeyPath)
	sshKeyPath := GetSSHKeyPath()
	// cmds := []string{
	// 	"apt-get upgrade",
	// }
	arg := append(
		[]string{
			"-o",
			"UserKnownHostsFile=/dev/null",
			"-o",
			"StrictHostKeyChecking=no",
			"-i",
			sshKeyPath,
			fmt.Sprintf("%s@%s", userName, IP),
		},
		cmds...,
	)

	// cmdOut, err := exec.Command("ssh", arg...).Output()
	cmd := exec.Command("ssh", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	// cmdReader, err := cmd.StdoutPipe()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	// 	os.Exit(1)
	// }
	// scanner := bufio.NewScanner(cmdReader)
	// go func() {
	// 	for scanner.Scan() {
	// 		fmt.Printf("ssh output | %s\n", scanner.Text())
	// 	}
	// }()
	// err = cmd.Start()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	// 	os.Exit(1)
	// }

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
	// 	os.Exit(1)
	// }
}
