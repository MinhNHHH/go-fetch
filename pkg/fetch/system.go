package fetch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	asci "github.com/minhnh/fetch/internal/ascii"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemInfor struct {
	User       string
	Terminal   string
	HostName   HostNameInfor
	Cpu        CPUInfor
	Vm         VMInfor
	Disk       DiskInfo
	Packages   string
	GPU        string
	Theme      string
	Resolution string
	Shell      string
	Icon       string
}

type HostNameInfor struct {
	HostName        string
	UpTime          uint64
	BootTime        uint64
	Procs           uint64
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

type CPUInfor struct {
	VendorId  string
	Model     string
	ModelName string
	Mhz       float64
	CacheSize int32
}

type VMInfor struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
	Free        uint64
	Active      uint64
	Inactive    uint64
}

type DiskInfo struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func uptimeToDaysHoursMins(uptimeSeconds uint64) (days, hours, mins uint64) {
	// Calculate days
	days = uptimeSeconds / (24 * 3600)

	// Calculate remaining seconds after extracting days
	remainingSeconds := uptimeSeconds % (24 * 3600)

	// Calculate hours
	hours = remainingSeconds / 3600

	// Calculate remaining seconds after extracting hours
	remainingSeconds %= 3600

	// Calculate minutes
	mins = remainingSeconds / 60

	return days, hours, mins
}

func getUser() string {
	return os.Getenv("USER")
}

func getTerminal() string {
	return os.Getenv("TERM_PROGRAM")
}

func getCPU() CPUInfor {
	cpuStat, err := cpu.Info()
	if err != nil {
		log.Fatalf("error when getting cpu information: %s", err.Error())
	}
	if len(cpuStat) == 0 {
		log.Fatalf("can not get cpu information")
	}
	cpuInfor := CPUInfor{
		VendorId:  cpuStat[0].VendorID,
		Model:     cpuStat[0].Model,
		ModelName: cpuStat[0].ModelName,
		Mhz:       cpuStat[0].Mhz,
		CacheSize: cpuStat[0].CacheSize,
	}
	return cpuInfor
}

func getVM() VMInfor {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("error when getting vm information: %s", err.Error())
	}
	vmInfor := VMInfor{
		Total:       vmStat.Total,
		Available:   vmStat.Available,
		Used:        vmStat.Used,
		UsedPercent: vmStat.UsedPercent,
		Free:        vmStat.Free,
		Active:      vmStat.Active,
		Inactive:    vmStat.Inactive,
	}
	return vmInfor
}

func getDisk() DiskInfo {
	diskStat, err := disk.Usage("/") // If you're in Unix change this "\\" for "/"
	if err != nil {
		log.Fatalf("error when getting disk information: %s", err.Error())
	}
	diskInfor := DiskInfo{
		Total:       diskStat.Total,
		Used:        diskStat.Used,
		UsedPercent: diskStat.UsedPercent,
		Free:        diskStat.Free,
	}

	return diskInfor
}

func getHostName() HostNameInfor {
	hostStat, err := host.Info()
	if err != nil {
		log.Fatalf("error when getting hostname information: %s", err.Error())
	}
	hostName := HostNameInfor{
		HostName:        hostStat.Hostname,
		UpTime:          hostStat.Uptime,
		BootTime:        hostStat.BootTime,
		Procs:           hostStat.Procs,
		OS:              hostStat.OS,
		Platform:        hostStat.Platform,
		PlatformFamily:  hostStat.PlatformFamily,
		PlatformVersion: hostStat.PlatformVersion,
		KernelVersion:   hostStat.KernelVersion,
		KernelArch:      hostStat.KernelArch,
	}
	return hostName
}

func GetUptime(uptime uint64) string {
	days, hours, mins := uptimeToDaysHoursMins(uptime)

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d mins", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d mins", hours, mins)
	} else {
		return fmt.Sprintf("%d mins", mins)
	}
}

func GetMemmory(vmUsed, vmTotal uint64) string {
	return fmt.Sprintf("%dMB / %dMB", vmUsed/1024/1024, vmTotal/1024/1024)
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func GetPackages(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "sh"
		args = []string{"-c", "dpkg --list | grep '^ii' | wc -l"}
	case "darwin":
		cmd = "sh"
		args = []string{"-c", "brew list | wc -l"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "Get-Package | Measure-Object | Select-Object -ExpandProperty Count"}
	}

	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get packages: %s", err)
		ch <- ""
		return
	}
	ch <- strings.TrimSpace(output)
}

func GetResolution(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "sh"
		args = []string{"-c", "xrandr | grep '*' | awk '{print $1}'"}
	case "darwin":
		cmd = "sh"
		args = []string{"-c", "system_profiler SPDisplaysDataType | grep Resolution"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "Get-WmiObject -Class Win32_VideoController | Select-Object -ExpandProperty VideoModeDescription"}

	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get resolution: %s", err)
		ch <- ""
		return
	}
	resolutions := strings.Split(strings.Trim(output, "\n"), "\n")
	ch <- strings.TrimSpace(strings.Join(resolutions, ", "))
}

func GetGpu(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "sh"
		args = []string{"-c", "lspci | grep -i 'vga\\|3d\\|2d' | awk -F: '{print $3}'"}
	case "darwin":
		cmd = "sh"
		args = []string{"-c", "system_profiler SPDisplaysDataType | grep 'Chipset Model:'"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "Get-WmiObject -Class Win32_VideoController | Select-Object -ExpandProperty Name"}
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get Gpu: %s", err)
		ch <- ""
		return
	}
	ch <- strings.TrimSpace(output)
}

func GetShell(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = "sh"
		args = []string{"-c", "echo $SHELL"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "[System.Environment]::GetEnvironmentVariable('ComSpec')"}
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get theme %s", err)
		ch <- ""
		return
	}
	ch <- strings.TrimSpace(output)
}

func GetTheme(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "sh"
		args = []string{"-c", "gsettings get org.gnome.desktop.interface gtk-theme"}
	case "darwin":
		cmd = "sh"
		args = []string{"-c", "defaults read -g AppleInterfaceStyle"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "(Get-ItemProperty -Path HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize).AppsUseLightThemes"}
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get theme: %s", err)
		ch <- ""
		return
	}
	ch <- strings.TrimSpace(output)
}

func GetIcons(ch chan<- string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "sh"
		args = []string{"-c", "gsettings get org.gnome.desktop.interface icon-theme"}
	case "darwin":
		// macOS does not have a system-wide icon theme configuration like Linux.
		ch <- "Apple"
		return
	case "windows":
		// Windows does not have a system-wide icon theme configuration like Linux.
		ch <- "Windows"
		return
	}

	output, err := runCommand(cmd, args...)
	if err != nil {
		fmt.Printf("Error get Icons: %s", err)
		ch <- ""
		return
	}
	ch <- strings.TrimSpace(output)
}

func (si SystemInfor) formatInfo(label, info string) string {
	return fmt.Sprintf("%s%s: %s", label, asci.PlaceHolder["${c0}"], info)
}

func (si SystemInfor) ListSysInfor(disable, seemore []string) []string {
	// We want to display by order
	listSysInform := []string{
		fmt.Sprintf(si.User + "@" + si.HostName.HostName),
		"-----------------------------------",
		si.formatInfo("OS", si.HostName.OS),
		si.formatInfo("Host", si.HostName.HostName),
		si.formatInfo("Kernel", si.HostName.KernelVersion),
		si.formatInfo("Uptime", GetUptime(si.HostName.UpTime)),
		si.formatInfo("Packages", si.Packages),
		si.formatInfo("Shell", si.Shell),
		si.formatInfo("Resolution", si.Resolution),
		si.formatInfo("Theme", si.Theme),
		si.formatInfo("Icons", si.Icon),
		si.formatInfo("CPU", si.Cpu.ModelName),
		si.formatInfo("GPU", si.GPU),
		si.formatInfo("Memory", GetMemmory(si.Vm.Used, si.Vm.Total)),
	}

	if len(disable) > 0 {
		for _, typeInfo := range disable {
			for index, str := range listSysInform {
				findDisable := strings.Contains(strings.ToLower(str), typeInfo)
				if findDisable {
					listSysInform = append(listSysInform[:index], listSysInform[index+1:]...)
				}
			}
		}
	}

	return listSysInform
}

func NewSysInfor() SystemInfor {
	packagesChan := make(chan string)
	resolutionChan := make(chan string)
	themChan := make(chan string)
	gpuChan := make(chan string)
	shellChan := make(chan string)
	iconChan := make(chan string)

	go GetIcons(iconChan)
	go GetResolution(resolutionChan)
	go GetPackages(packagesChan)
	go GetTheme(themChan)
	go GetGpu(gpuChan)
	go GetShell(shellChan)

	return SystemInfor{
		User:       getUser(),
		Terminal:   getTerminal(),
		HostName:   getHostName(),
		Cpu:        getCPU(),
		Vm:         getVM(),
		Disk:       getDisk(),
		Packages:   <-packagesChan,
		Resolution: <-resolutionChan,
		GPU:        <-gpuChan,
		Theme:      <-themChan,
		Shell:      <-shellChan,
		Icon:       <-iconChan,
	}
}
