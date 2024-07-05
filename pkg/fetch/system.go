package fetch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var CodeColor = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[1;31m",
	"green":  "\033[1;32m",
	"cyan":   "\033[1;33m",
	"yellow": "\033[1;34m",
	"purple": "\033[1;35m",
	"blue":   "\033[1;36m",
	"white":  "\033[1;37m",
}

var PlaceHolder = map[string]string{
	"${c0}": CodeColor["reset"],
	"${c1}": CodeColor["red"],
	"${c2}": CodeColor["green"],
	"${c3}": CodeColor["yellow"],
	"${c4}": CodeColor["blue"],
	"${c5}": CodeColor["purple"],
	"${c6}": CodeColor["cyan"],
	"${c7}": CodeColor["white"],
}

type SystemInfor struct {
	User     string
	Terminal string
	HostName HostNameInfor
	Cpu      CPUInfor
	Vm       VMInfor
	Disk     DiskInfo
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

func (si SystemInfor) GetUptime() string {
	uptime := si.HostName.UpTime
	days, hours, mins := uptimeToDaysHoursMins(uptime)

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d mins", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d mins", hours, mins)
	} else {
		return fmt.Sprintf("%d mins", mins)
	}
}

func (si SystemInfor) GetHost() string {
	return si.HostName.HostName
}

func (si SystemInfor) GetOS() string {
	return si.HostName.OS
}

func (si SystemInfor) GetKernelVersion() string {
	return si.HostName.KernelVersion
}

func (si SystemInfor) GetCpu() string {
	return si.Cpu.ModelName
}

func (si SystemInfor) GetMemmory() string {
	return fmt.Sprintf("%dMB / %dMB", si.Vm.Used/1024/1024, si.Vm.Total/1024/1024)
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (si SystemInfor) GetPackages() string {
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
	default:
		return ""
	}

	output, err := runCommand(cmd, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}

func (si SystemInfor) GetResolution() string {
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
	default:
		return ""
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		return ""
	}
	resolutions := strings.Split(strings.Trim(output, "\n"), "\n")
	return strings.Join(resolutions, ", ")
}

func (si SystemInfor) GetGpu() string {
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
	default:
		return ""
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}

func (si SystemInfor) GetShell() string {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = "sh"
		args = []string{"-c", "echo $SHELL"}
	case "windows":
		cmd = "powershell"
		args = []string{"-Command", "[System.Environment]::GetEnvironmentVariable('ComSpec')"}
	default:
		return ""
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}

func (si SystemInfor) GetTheme() string {
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
	default:
		return ""
	}
	output, err := runCommand(cmd, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}

// func (si SystemInfor) GetIcons() string {
// 	cmd, err := ExecLinuxCmd("gsettings get org.gnome.desktop.interface icon-theme")
// 	if err != nil {
// 		fmt.Printf("err: %s", err)
// 		return ""
// 	}
// 	return strings.Trim(cmd, "\n")
// }

func (si SystemInfor) formatInfo(label, info string) string {
	return fmt.Sprintf("%s%s: %s", label, PlaceHolder["${c0}"], info)
}

func (si SystemInfor) ListSysInfor(disable, seemore []string) []string {
	// We want to display by order
	listSysInform := []string{
		fmt.Sprintf(si.User + "@" + si.GetHost()),
		"-----------------------------------",
		si.formatInfo("OS", si.GetOS()),
		si.formatInfo("Host", si.GetHost()),
		si.formatInfo("Kernel", si.GetKernelVersion()),
		si.formatInfo("Uptime", si.GetUptime()),
		si.formatInfo("Packages", si.GetPackages()),
		si.formatInfo("Shell", si.GetShell()),
		si.formatInfo("Resolution", si.GetResolution()),
		si.formatInfo("Theme", si.GetTheme()),
		// si.formatInfo("Icons", si.GetIcons()),
		si.formatInfo("Terminal", si.GetUptime()),
		si.formatInfo("CPU", si.GetCpu()),
		si.formatInfo("GPU", si.GetGpu()),
		si.formatInfo("Memory", si.GetMemmory()),
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
	return SystemInfor{
		User:     getUser(),
		Terminal: getTerminal(),
		HostName: getHostName(),
		Cpu:      getCPU(),
		Vm:       getVM(),
		Disk:     getDisk(),
	}
}
