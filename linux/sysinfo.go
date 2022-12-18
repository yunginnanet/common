//go:build linux

package linux

import (
	"syscall"
	"time"
)

/*
	some interesting information on RAM, sysinfo, and /proc/meminfo

	- https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773
	- https://github.com/mmalecki/procps/blob/master/proc/sysinfo.c

	in the second link we see that even procps is parsing /proc/meminfo to get the RAM information
	i'd really like to avoid this, and I may end up estimating the available memory myself...
	the whole idea of this package is to focus on system calls and not on parsing files

	for now I'll just add a note to the RAMInfo struct definition.

	tldr; sysinfo is a bit incomplete due to it's lack of available memory calculation
*/

// from https://man7.org/linux/man-pages/man2/sysinfo.2.html
/*
	struct sysinfo {
		long uptime;             // Seconds since boot
		unsigned long loads[3];  // 1, 5, and 15 minute load averages
		unsigned long totalram;  // Total usable main memory size
		unsigned long freeram;   // Available memory size
		unsigned long sharedram; // Amount of shared memory
		unsigned long bufferram; // Memory used by buffers
		unsigned long totalswap; // Total swap space size
		unsigned long freeswap;  // Swap space still available
		unsigned short procs;    // Number of current processes
		unsigned long totalhigh; // Total high memory size
		unsigned long freehigh;  // Available high memory size
		unsigned int mem_unit;   // Memory unit size in bytes
		char _f[20-2*sizeof(long)-sizeof(int)]; // Padding to 64 bytes
	};
*/

type SystemInfo struct {
	// Uptime is the time since the system was booted.
	Uptime time.Duration
	// Loads is a 3 element array containing the 1, 5, and 15 minute load averages.
	Loads [3]uint64
	// RAM is a struct containing information about the system's RAM. See notes on this.
	RAM RAMInfo
	// Procs is the number of current processes.
	Procs uint16
}

// RAMInfo is a struct that contains information about the running system's memory.
// Please see important notes in the struct definition.
type RAMInfo struct {
	Total int64
	// Free is the amount of memory that is not being used.
	// This does not take into account memory that is being used for caching.
	// For more information, see the comments at the top of this file.
	Free      int64
	Used      int64
	Shared    int64
	Buffers   int64
	Cached    int64
	SwapTotal int64
	SwapFree  int64

	// unit is the memory unit size multiple in bytes.
	// It is used to calculate the above values when the struct is initialized.q
	unit uint32
}

// Sysinfo returns a SystemInfo struct containing information about the running system.
// For memory, please see important notes in the RAMInfo struct definition and the top of this file.
// Be sure to check err before using the returned value.
func Sysinfo() (systemInfo *SystemInfo, err error) {
	sysinf := syscall.Sysinfo_t{}
	err = syscall.Sysinfo(&sysinf)
	unit := uint64(sysinf.Unit)
	return &SystemInfo{
		Uptime: time.Duration(sysinf.Uptime) * time.Second,
		Loads:  [3]uint64{sysinf.Loads[0], sysinf.Loads[1], sysinf.Loads[2]},
		RAM: RAMInfo{
			Total:     int64(sysinf.Totalram * unit),
			Free:      int64(sysinf.Freeram * unit),
			Used:      int64((sysinf.Totalram - sysinf.Freeram) * unit),
			Shared:    int64(sysinf.Sharedram * unit),
			Buffers:   int64(sysinf.Bufferram * unit),
			Cached:    int64((sysinf.Totalram - sysinf.Freeram - sysinf.Bufferram) * unit),
			SwapTotal: int64(sysinf.Totalswap * unit),
			SwapFree:  int64(sysinf.Freeswap * unit),
			unit:      sysinf.Unit,
		},
		Procs: sysinf.Procs,
	}, nil
}

// Uptime returns the time since the system was booted.
// Be sure to check err before using the returned value.
func Uptime() (utime time.Duration, err error) {
	var sysinf *SystemInfo
	sysinf, err = Sysinfo()
	return sysinf.Uptime, err
}
