package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner/parameters"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/clients"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2os"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2params"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2vm"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

////////////////////////////////////////

func PulumiVMConnectionToSSHConnectionInfo(outputs auto.OutputMap, vmname string) (*sshConnectionInfo, error) {
	identityFile, _ := runner.GetProfile().ParamStore().GetWithDefault(parameters.PublicKeyPath, "")

	outputName := fmt.Sprintf("%s-connection", vmname)
	out, ok := outputs[outputName]
	if !ok {
		return nil, fmt.Errorf("output '%s' not found", outputName)
	}
	vmconn := out.Value.(map[string]interface{})
	// Convert stack outputs into SSH connection info
	conn := &sshConnectionInfo{
		Host:         vmconn["host"].(string),
		User:         vmconn["user"].(string),
		Port:         22,
		IdentityFile: identityFile,
	}
	return conn, nil
}

func PulumiCreateEnv(stackName string, envFactory func(ctx *pulumi.Context) error, upResultHandler func(upResult auto.UpResult) error) error {
	// Create or Get a stack,
	stackManager := infra.GetStackManager()
	config := runner.ConfigMap{}
	failIfNotExist := false
	ctx := context.Background()
	_, upResult, err := stackManager.GetStackNoDeleteOnFailure(ctx, stackName, config, func(ctx *pulumi.Context) error {
		// Create the resources in the stack
		if err := envFactory(ctx); err != nil {
			return fmt.Errorf("setup vms in remote instance: %w", err)
		}
		return nil
	}, failIfNotExist)
	if err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	return upResultHandler(upResult)
}

////////////////////////////////////////

// devenv constants
const (
	containerName  = "build"
	dockerImage    = "datadog/agent-buildimages-windows_x64"
	dockerImageTag = "ltsc2022"
	agentRepoPath  = `$HOME\projects\datadog-agent`
	goModPath      = `C:\dev\go\pkg\mod`
)

func DisableWindowsDefender(conn Connection) error {
	res, err := conn.Exec(`(Get-MpComputerStatus).RealTimeProtectionEnabled`)
	if err != nil || res != "False" {
		fmt.Println("Disabling Windows Defender...")
		destPath := `C:\Windows\Temp\Configure-Antivirus.ps1`
		res, err := conn.Exec(fmt.Sprintf(`Invoke-WebRequest -UseBasicParsing "https://raw.githubusercontent.com/actions/runner-images/060ad1383a7231762cbd2d61a04116a2dc7905e0/images/win/scripts/Installers/Configure-Antivirus.ps1" -o %s`, destPath))
		if err != nil {
			return fmt.Errorf("failed to download script: %w\n%s", err, res)
		}
		res, err = conn.Exec(destPath)
		if err != nil {
			return fmt.Errorf("failed to disable defender: %w\n%s", err, res)
		}
		fmt.Println("Disabled Windows Defender!")
	}
	return nil
}

// Install docker reboots the host
func InstallDocker(conn Connection) error {
	// Install docker
	res, err := conn.Exec(`docker --version`)
	if err != nil {
		fmt.Println("Installing docker...")
		destPath := `C:\Windows\Temp\install-docker-ce.ps1`
		res, err = conn.Exec(fmt.Sprintf(`Invoke-WebRequest -UseBasicParsing "https://raw.githubusercontent.com/microsoft/Windows-Containers/Main/helpful_tools/Install-DockerCE/install-docker-ce.ps1" -o %s`, destPath))
		if err != nil {
			return fmt.Errorf("failed to download script: %w\n%s", err, res)
		}
		res, err = conn.Exec(`C:\Windows\Temp\install-docker-ce.ps1`)
		if err != nil {
			return fmt.Errorf("failed to install docker-ce: %w\n%s", err, res)
		}
		fmt.Println("Waiting for host to reboot...")
		// TODO: How to wait for host to reboot? If we try to reconnect too early we will reconnect
		//       before the reboot
		time.Sleep(30 * time.Second)
		fmt.Println("Reconnecting to host after reboot...")
		err = conn.Reconnect()
		if err != nil {
			return err
		}
		fmt.Println("connected!")
		// TODO: install-docker-ce registers a scheduled task but I'm not sure that it runs on SSH logon
		//       so for now we run it manually
		res, err = conn.Exec(`C:\Windows\Temp\install-docker-ce.ps1`)
		if err != nil {
			return fmt.Errorf("failed to install docker-ce: %w\n%s", err, res)
		}
		_, err := conn.Exec(`docker --version`)
		if err != nil {
			return fmt.Errorf("failed to install docker-ce: %w", err)
		}
		fmt.Println("docker installed!")
	}
	return nil
}

func InstallChocolatey(conn Connection) error {
	// Install choco
	res, err := conn.Exec(`choco --version`)
	if err != nil {
		fmt.Println("Installing chocolatey...")
		res, err = conn.Exec(`Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`)
		if err != nil {
			return fmt.Errorf("failed to install chocolatey: %w\n%s", err, res)
		}
		// Restart sshd and reconnect to reload env vars
		res, err = conn.Exec("Restart-Service sshd")
		if err != nil {
			return fmt.Errorf("Failed to restart sshd: %w\n%s", err, res)
		}
		err = conn.Reconnect()
		if err != nil {
			return err
		}
		fmt.Println("choco installed!")
	}
	return nil
}

func InstallGit(conn Connection) error {
	// Install git
	res, err := conn.Exec(`git --version`)
	if err != nil {
		fmt.Println("Installing git...")
		res, err = conn.Exec(`choco install -y git.install --params "'/GitAndUnixToolsOnPath /WindowsTerminal /NoAutoCrlf'"`)
		if err != nil {
			return fmt.Errorf("failed to install git: %w\n%s", err, res)
		}
		// choco refreshenv to put git on PATH
		res, err = conn.Exec(`import-module C:\ProgramData\chocolatey\helpers\chocolateyProfile.psm1; refreshenv`)
		if err != nil {
			return fmt.Errorf("failed to choco refreshenv: %w\n%s", err, res)
		}
		// Restart sshd and reconnect to reset env vars
		res, err = conn.Exec(`Restart-Service sshd`)
		if err != nil {
			return fmt.Errorf("Failed to restart sshd: %w\n%s", err, res)
		}
		err = conn.Reconnect()
		if err != nil {
			return err
		}
		fmt.Println("git installed!")
	}
	return nil
}

func CreateBuildContainer(conn Connection) error {
	// create dir for go deps, will be mounted into container so new containers can share the deps
	res, err := conn.Exec(fmt.Sprintf(`mkdir -Force %s`, goModPath))
	if err != nil {
		return fmt.Errorf("failed to create go mod path: %w\n%s", err, res)
	}

	res, err = conn.Exec(fmt.Sprintf("docker container inspect %s", containerName))
	if err != nil {
		imagePath := fmt.Sprintf("%s:%s", dockerImage, dockerImageTag)
		fmt.Printf("Creating container: %s from %s\n", containerName, imagePath)
		fmt.Println("The image is ~25GB so the download will take some time...")
		res, err = conn.Exec(fmt.Sprintf(`docker run --quiet --name %s --detach -v "%s:C:\mnt" -v "%s:C:\dev\go\pkg\mod" --storage-opt size=50G %s ping -t 127.0.0.1`,
			containerName,
			agentRepoPath,
			goModPath,
			imagePath))
		if err != nil {
			return fmt.Errorf("failed to create container: %w\n%s", err, res)
		}
	}
	return nil
}

func CloneAgentRepo(conn Connection) error {
	// Clone repo
	// TODO: Doesn't vscode support an option where it rsync's the local code to the remote host to build?
	//       agent-sandbox gets wiped (weekly?), so code shouldn't live on the VM, or instance should live in a different account
	res, err := conn.Exec(fmt.Sprintf(`mkdir %s`, agentRepoPath))
	if err != nil {
		// path already exists
		return nil
	}
	fmt.Println("Cloning datadogagent...")
	res, err = conn.Exec(fmt.Sprintf(`git clone https://github.com/datadog/datadog-agent %s`, agentRepoPath))
	if err != nil {
		return fmt.Errorf("failed to clone agent repo: %w\n%s", err, res)
	}
	return nil
}

////////////////////////////////////////

func InstallDevEnv(conn Connection) error {
	var err error
	err = DisableWindowsDefender(conn)
	if err != nil {
		return err
	}

	err = InstallDocker(conn)
	if err != nil {
		return err
	}

	err = InstallChocolatey(conn)
	if err != nil {
		return err
	}

	err = InstallGit(conn)
	if err != nil {
		return err
	}

	err = CloneAgentRepo(conn)
	if err != nil {
		return err
	}

	err = CreateBuildContainer(conn)
	if err != nil {
		return err
	}

	return nil
}

func DisplaySSHConnectionInfo(info *sshConnectionInfo) {
	fmt.Printf(`Add the following to ~/.ssh/config
# Added by E2E
Host %s %s
  HostName %s
  IdentityFile %s
  User %s
`,
		info.Host, "winbuild",
		info.Host,
		info.IdentityFile,
		info.User,
	)
	fmt.Printf("Then connect with `ssh %s` or `ssh winbuild`\n", info.Host)
}

////////////////////////////////////////

// Pulumi VM constants
const (
	// stack name will be prefixed based on test infra config
	defaultPulumiStackName = "devenv"
	defaultPulumiVMName    = "winbuild"
	// 4 vCPU, 16GB RAM
	ec2InstanceType = "t2.xlarge"
	// Name of pre-built image with devenv already setup. The docker build image download takes ~20 minutes so this saves a lot of time.
	// TODO: automate updating this image
	//       Example command: aws ec2 create-image --region us-east-1 --instance-id i-0f41f055d75adf722 --name datadog-agent-winbuild-devenv --no-reboot
	prebuiltAMIName = "datadog-agent-winbuild-devenv"
)

func PulumiVM(ctx *pulumi.Context, newImage bool) error {

	// create an EC2 instance with windows amd64
	// TODO: How to specify OS? This defaults to Windows Server 2022
	opts := []ec2params.Option{
		ec2params.WithName(defaultPulumiVMName),
	}
	if newImage {
		opts = append(opts, ec2params.WithArch(ec2os.WindowsOS, "x86_64"))
	} else {
		ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Filters: []ec2.GetAmiFilter{
				{
					Name: "name",
					Values: []string{
						prebuiltAMIName,
					},
				},
			},
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to lookup pre-built AMI ID: %w", err)
		}
		opts = append(opts, ec2params.WithImageName(ami.Id, "x86_64", ec2os.WindowsOS))
	}

	// WithInstanceType must come after WithImageName b/c WithImageName overwrites the instance type
	opts = append(opts, ec2params.WithInstanceType(ec2InstanceType))

	_, err := ec2vm.NewEc2VM(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to make ec2vm: %w", err)
	}
	return nil
}

func LoadConnInfoFromFile(path string) (*sshConnectionInfo, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening json file '%s': %w", path, err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading json file '%s': %w", path, err)
	}

	var c sshConnectionInfo
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	return &c, nil
}

func main() {
	var connInfo *sshConnectionInfo
	var err error

	var usePulumi bool
	var pulumiNewImage bool
	var pulumiStackName string
	var sshConnFile string

	flag.BoolVar(&usePulumi, "pulumi", false, "")
	flag.BoolVar(&pulumiNewImage, "pulumi-new-image", false, "Create a new image from scratch, instead of using the pre-built image")
	flag.StringVar(&pulumiStackName, "pulumi-stack-name", defaultPulumiStackName, "Pulumi stack name with vm-connection to connect to")
	flag.StringVar(&sshConnFile, "file", "", "JSON file containing connection info")

	flag.Parse()

	if usePulumi {
		err = PulumiCreateEnv(pulumiStackName,
			func(ctx *pulumi.Context) error {
				return PulumiVM(ctx, pulumiNewImage)
			},
			func(upResult auto.UpResult) error {
				// Map Pulumi VM connection object to generic SSH connection info object
				var err error
				outputs := upResult.Outputs

				connInfo, err = PulumiVMConnectionToSSHConnectionInfo(outputs, defaultPulumiVMName)
				return err
			})
		if err != nil {
			fmt.Printf("pulumi stack failure: %v", err)
			return
		}
	} else if len(sshConnFile) > 0 {
		connInfo, err = LoadConnInfoFromFile(sshConnFile)
		if err != nil {
			fmt.Printf("failed reading connection info from file: %v", err)
			return
		}
	} else {
		flag.Usage()
		return
	}

	DisplaySSHConnectionInfo(connInfo)

	conn, err := NewSSHConnection(connInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	out, err := conn.Exec("whoami")
	fmt.Println(out)
	out, err = conn.Exec("ipconfig")
	fmt.Println(out)

	if !usePulumi || pulumiNewImage {
		err = InstallDevEnv(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("Enter the build container with `docker exec -it %s powershell`\n", containerName)
}

////////////////////////////////////////

type Connection interface {
	Exec(command string) (string, error)
	Reconnect() error
	Wait() error
	Close()
}

////////////////////////////////////////

type sshConnection struct {
	info   *sshConnectionInfo
	client *ssh.Client
}

func NewSSHConnection(info *sshConnectionInfo) (*sshConnection, error) {
	c := &sshConnection{info: info}
	err := c.Reconnect()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (conn *sshConnection) Exec(command string) (string, error) {
	return PsExec(conn.client, command)
}

func (conn *sshConnection) Reconnect() error {
	// close current connection
	if conn.client != nil {
		conn.client.Close()
		conn.client = nil
	}
	// connect
	client, err := ConnectSSH(conn.info)
	if err != nil {
		return err
	}
	conn.client = client
	return nil
}

func (conn *sshConnection) Wait() error {
	var err error
	if conn.client != nil {
		err = conn.client.Wait()
		conn.client.Close()
		conn.client = nil
	}
	return err
}

func (conn *sshConnection) Close() {
	if conn.client != nil {
		conn.client.Close()
		conn.client = nil
	}
}

type sshConnectionInfo struct {
	Host         string `json:"Host"` // ip or hostnme
	User         string `json:"User"`
	Port         uint   `json:"Port"`
	IdentityFile string `json:"IdentityFile"`
}

func ConnectSSH(connInfo *sshConnectionInfo) (*ssh.Client, error) {
	client, _, err := clients.GetSSHClient(
		connInfo.User,
		fmt.Sprintf("%s:%d", connInfo.Host, connInfo.Port),
		nil,
		2*time.Second, 5)
	return client, err
}

// //////////////////////////////////////

func PsExec(client *ssh.Client, command string) (string, error) {
	s, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer s.Close()

	var outstr string
	out, err := s.CombinedOutput(command)
	if out != nil {
		outstr = strings.TrimSuffix(string(out), "\r\n")
	}
	if err != nil {
		return outstr, err
	}
	return outstr, nil
}

////////////////////////////////////////
