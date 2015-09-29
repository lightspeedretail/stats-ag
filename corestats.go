package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type Stats struct{}

func GetMemStats() string {
	v, err := mem.VirtualMemory()
	format := "total=%d available=%d used=%d used_percent=%f"
	if err == nil {
		return fmt.Sprintf(format, v.Total, v.Available, v.Used, v.UsedPercent)
	} else {
		return fmt.Sprintf(format, 0, 0, 0, 0)
	}
}

func GetLoadStats() string {
	l, err := load.LoadAvg()
	format := "last1=%f last5=%f last15=%f"
	if err == nil {
		return fmt.Sprintf(format, l.Load1, l.Load5, l.Load15)
	} else {
		return fmt.Sprintf(format, 0, 0, 0)
	}
}

func GetDiskStats() string {
	d, err := disk.DiskUsage("/")
	format := "fstype=%s total=%d free=%d used=%d used_percent=%f inodes_total=%d inodes_used=%d inodes_free=%d inodes_used_percent=%f"
	if err == nil {
		return fmt.Sprintf(format,
			d.Fstype, d.Total, d.Free, d.Used, d.UsedPercent, d.InodesTotal, d.InodesUsed, d.InodesFree, d.InodesUsedPercent)
	} else {
		return fmt.Sprintf(format, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	}
}

func GetCpuStats() string {
	v, err := cpu.CPUTimes(false)
	format := "user=%f system=%f idle=%f nice=%f iowait=%f irq=%f soft_irq=%f steal=%f guest=%f guest_nice=%f stolen=%f"
	if err == nil {
		return fmt.Sprintf(format,
			v[0].User, v[0].System, v[0].Idle, v[0].Nice, v[0].Iowait, v[0].Irq, v[0].Softirq, v[0].Steal,
			v[0].Guest, v[0].GuestNice, v[0].Stolen)
	} else {
		return fmt.Sprintf(format, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	}
}

func GetHostStats() string {
	v, err := host.HostInfo()
	format := "uptime=%d procs=%d os=%s platform=%s platform_family=%s platform_version=%s virtualization_system=%s virtualization_role=%s"
	if err == nil {
		return fmt.Sprintf(format, v.Uptime, v.Procs, v.OS, v.Platform, v.PlatformFamily, v.PlatformVersion, v.VirtualizationSystem, v.VirtualizationRole)
	} else {
		return fmt.Sprintf(format, 0, 0, 0, 0, 0, 0, 0, 0)
	}
}

func GetNetStats() string {
	v, err := net.NetIOCounters(false)
	format := "bytes_sent=%d bytes_recv=%d packets_sent=%d packets_recv=%d err_in=%d err_out=%d drop_in=%d drop_out=%d"
	if err == nil {
		return fmt.Sprintf(format, v[0].BytesSent, v[0].BytesRecv, v[0].PacketsSent, v[0].PacketsRecv, v[0].Errin, v[0].Errout, v[0].Dropin, v[0].Dropout)
	} else {
		return fmt.Sprintf(format, 0, 0, 0, 0, 0, 0, 0, 0)
	}
}
