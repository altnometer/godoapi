package support

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/fatih/color"
)

/////////////////////////////files required///////////////////////////////////

// TSLArchSource is where you take letsencrypt ssl/tsl files.
var TSLArchSource = getTSLArchPath()

/////////////////////////////consts///////////////////////////////////////////

// VolByIDPrefix holds prefix for vol link names in /dev/disk/by-id dir.
const VolByIDPrefix string = "scsi-0DO_Volume_"

// Master1Name holds the name of k8s master node.
const Master1Name string = "master-1"

// NodeMemSize used when create k8s nodes as default value for -size arg.
const NodeMemSize = "2gb"

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

// GreenLn colors a line.
var GreenLn = color.New(color.FgGreen).PrintlnFunc()

// RedSp colors str.
var RedSp = color.New(color.FgRed).SprintFunc()

// RedSf color string with interpolation. RedSf("color is %s", myColor)
var RedSf = color.New(color.FgRed).SprintfFunc()

// RedPf output with interpolation colored red.
var RedPf = color.New(color.FgRed).PrintfFunc()

// RedLn colors a line.
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

////////////////////////functions//////////////////////////////////////////////

// StringInSlice checks if a given string is in a slice.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

func getTSLArchPath() string {
	homeDir, err := getHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := homeDir + "/letsencrypt.redmoo.gz"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(err)
	}
	return path
}

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

// UserConfirmDefaultY return true for input 'y|Y' or not input.
func UserConfirmDefaultY(prompt string) (conf bool, err error) {
	fmt.Println(prompt)
	fmt.Println("[Y/n]")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return false, err
	}
	if char == 10 {
		return true, nil
	}
	if char != 'y' && char != 'Y' {
		return false, nil
	}
	return true, nil
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

// PrintDropData prints out droplet data.
func PrintDropData(droplets []godo.Droplet) {
	if len(droplets) > 0 {
		for _, d := range droplets {
			ip, err := d.PublicIPv4()
			if err != nil {
				RedLn(err)
			}
			fmt.Printf("d.Name                %+v\n", d.Name)
			fmt.Printf("d.ID                  %+v\n", d.ID)
			fmt.Printf("d.Tags                %+v\n", d.Tags)
			// fmt.Printf("d.Size = %+v\n", d.Size)
			// fmt.Printf("d.Networks.V4 = %+v\n", d.Networks.V4)
			fmt.Printf("ip                    %+v\n", ip)
			fmt.Printf("d.Size.Slug           %+v\n", d.Size.Slug)
			fmt.Printf("d.Size.Memory         %+v\n", d.Size.Memory)
			fmt.Printf("d.Size.Vcpus          %+v\n", d.Size.Vcpus)
			fmt.Printf("d.Size.Disk           %+v\n", d.Size.Disk)
			fmt.Printf("d.Size.PriceMonthly   %+v\n", d.Size.PriceMonthly)
			// fmt.Printf("d = %+v\n", d)
			fmt.Println("***************************")
		}
	} else {
		GreenLn("No droplets exist.")
	}
}

// ValidateMemSize return true for a valid droplet memory size.
func ValidateMemSize(memSize string) (bool, error) {
	var sizeIsValid bool
	for _, size := range DropletSizes {
		if memSize == size {
			return true, nil
		}
	}
	if sizeIsValid == false {
		errMsg := RedSf("Valid -size values %+v\n", DropletSizes)
		return false, errors.New(errMsg)
	}
	return true, nil
}

/////////////////////execute cmd///////////////////////////////////////////////

// ExecCmd execute command interactively.
func ExecCmd(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

/////////////////////execute ssh///////////////////////////////////////////////

// GetSSHKeyPath returns user provided path to ssh keys.
func GetSSHKeyPath() string {
	sshKeyPath := os.Getenv("DOSSHKeyPath")
	if sshKeyPath == "" {
		YellowLn("You can set env var DOSSHKeyPath!")
		sshKeyPath = GetUserInput("Type in a DOSSHKeyPath: ")
	}
	return sshKeyPath
}

// GetSSHOutput executes ssh cmd and returns cmd output.
// sshKey must be added with "eval $(ssh-agent -s) && ssh-add ~/.ssh/privkey"
func GetSSHOutput(userName, ip string, sshCmds []string) string {
	cmdArgs := append(
		[]string{
			"-o",
			"UserKnownHostsFile=/dev/null",
			"-o",
			"StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", userName, ip),
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

// FetchSSHOutput executes ssh cmd and returns cmd output.
func FetchSSHOutput(userName, ip, sshKeyPath string, sshCmds []string) string {
	cmdArgs := append(
		[]string{
			"-o",
			"UserKnownHostsFile=/dev/null",
			"-o",
			"StrictHostKeyChecking=no",
			"-i",
			sshKeyPath,
			fmt.Sprintf("%s@%s", userName, ip),
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
func ExecSSH(userName, ip string, cmds []string) {
	// sshCommander := SSHCommander{userName, ip, sshKeyPath}
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
			fmt.Sprintf("%s@%s", userName, ip),
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
