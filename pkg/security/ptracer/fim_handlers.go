// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux

// Package ptracer holds the start command of CWS injector
package ptracer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/DataDog/datadog-agent/pkg/security/proto/ebpfless"
	"github.com/DataDog/datadog-agent/pkg/util/native"
)

func isAcceptedRetval(retval int64) bool {
	return retval < 0 && retval != -int64(syscall.EACCES) && retval != -int64(syscall.EPERM)
}

func registerFIMHandlers(handlers map[int]syscallHandler) []string {
	fimHandlers := []syscallHandler{
		{
			IDs:        []syscallID{{Id: OpenNr, Name: "open"}},
			Func:       handleOpen,
			ShouldSend: func(ret int64) bool { return !isAcceptedRetval(ret) },
			SendIt:     true,
			RetFunc:    handleOpensRet,
		},
		{
			IDs:        []syscallID{{Id: OpenatNr, Name: "openat"}, {Id: Openat2Nr, Name: "openat2"}},
			Func:       handleOpenAt,
			ShouldSend: func(ret int64) bool { return !isAcceptedRetval(ret) },
			SendIt:     true,
			RetFunc:    handleOpensRet,
		},
		{
			IDs:        []syscallID{{Id: CreatNr, Name: "creat"}},
			Func:       handleCreat,
			ShouldSend: func(ret int64) bool { return !isAcceptedRetval(ret) },
			SendIt:     true,
			RetFunc:    handleOpensRet,
		},
		{
			IDs:        []syscallID{{Id: NameToHandleAtNr, Name: "name_to_handle_at"}},
			Func:       handleNameToHandleAt,
			ShouldSend: nil,
			SendIt:     false,
			RetFunc:    handleNameToHandleAtRet,
		},
		{
			IDs:        []syscallID{{Id: OpenByHandleAtNr, Name: "open_by_handle_at"}},
			Func:       handleOpenByHandleAt,
			ShouldSend: func(ret int64) bool { return !isAcceptedRetval(ret) },
			SendIt:     true,
			RetFunc:    handleOpensRet,
		},
		{
			IDs:        []syscallID{{Id: MemfdCreateNr, Name: "memfd_create"}},
			Func:       handleMemfdCreate,
			ShouldSend: func(ret int64) bool { return !isAcceptedRetval(ret) },
			SendIt:     true,
			RetFunc:    handleOpensRet,
		},
		{
			IDs:        []syscallID{{Id: FcntlNr, Name: "fcntl"}},
			Func:       handleFcntl,
			ShouldSend: nil,
			SendIt:     false,
			RetFunc:    handleFcntlRet,
		},
		{
			IDs:        []syscallID{{Id: DupNr, Name: "dup"}, {Id: Dup2Nr, Name: "dup2"}, {Id: Dup3Nr, Name: "dup3"}},
			Func:       handleDup,
			ShouldSend: nil,
			SendIt:     false,
			RetFunc:    handleDupRet,
		},
		{
			IDs:        []syscallID{{Id: CloseNr, Name: "close"}},
			Func:       handleClose,
			ShouldSend: nil,
			SendIt:     false,
			RetFunc:    nil,
		},
		{
			IDs:        []syscallID{{Id: UnlinkNr, Name: "unlink"}},
			Func:       handleUnlink,
			ShouldSend: func(ret int64) bool { return ret == 0 },
			SendIt:     true,
			RetFunc:    nil,
		},
		{
			IDs:        []syscallID{{Id: UnlinkatNr, Name: "unlinkat"}},
			Func:       handleUnlinkat,
			ShouldSend: func(ret int64) bool { return ret == 0 },
			SendIt:     true,
			RetFunc:    nil,
		},
		{
			IDs:        []syscallID{{Id: RmdirNr, Name: "rmdir"}},
			Func:       handleRmdir,
			ShouldSend: func(ret int64) bool { return ret == 0 },
			SendIt:     true,
			RetFunc:    nil,
		},
		{
			IDs:        []syscallID{{Id: RenameNr, Name: "rename"}},
			Func:       handleRename,
			ShouldSend: func(ret int64) bool { return ret == 0 },
			SendIt:     true,
			RetFunc:    handleRenamesRet,
		},
		{
			IDs:        []syscallID{{Id: RenameAtNr, Name: "renameat"}, {Id: RenameAt2Nr, Name: "renameat2"}},
			Func:       handleRenameAt,
			ShouldSend: func(ret int64) bool { return ret == 0 },
			SendIt:     true,
			RetFunc:    handleRenamesRet,
		},
	}

	syscallList := []string{}
	for _, h := range fimHandlers {
		for _, id := range h.IDs {
			if id.Id >= 0 { // insert only available syscalls
				handlers[id.Id] = h
				syscallList = append(syscallList, id.Name)
			}
		}
	}
	return syscallList
}

//
// handlers called on syscall entrance
//

func handleOpenAt(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	fd := tracer.ReadArgInt32(regs, 0)

	filename, err := tracer.ReadArgString(process.Pid, regs, 1)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFd(process, filename, fd)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: filename,
		Flags:    uint32(tracer.ReadArgUint64(regs, 2)),
		Mode:     uint32(tracer.ReadArgUint64(regs, 3)),
	}

	return fillFileMetadata(filename, msg.Open, disableStats)
}

func handleOpen(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	filename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFilename(process, filename)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: filename,
		Flags:    uint32(tracer.ReadArgUint64(regs, 1)),
		Mode:     uint32(tracer.ReadArgUint64(regs, 2)),
	}

	return fillFileMetadata(filename, msg.Open, disableStats)
}

func handleCreat(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	filename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFilename(process, filename)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: filename,
		Flags:    unix.O_CREAT | unix.O_WRONLY | unix.O_TRUNC,
		Mode:     uint32(tracer.ReadArgUint64(regs, 1)),
	}

	return fillFileMetadata(filename, msg.Open, disableStats)
}

func handleMemfdCreate(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	filename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}
	filename = "memfd:" + filename

	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: filename,
		Flags:    uint32(tracer.ReadArgUint64(regs, 1)),
	}
	return nil
}

func handleNameToHandleAt(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	fd := tracer.ReadArgInt32(regs, 0)

	filename, err := tracer.ReadArgString(process.Pid, regs, 1)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFd(process, filename, fd)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: filename,
	}
	return nil
}

func handleOpenByHandleAt(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	pFileHandleData, err := tracer.ReadArgData(process.Pid, regs, 1, 8 /*sizeof uint32 + sizeof int32*/)
	if err != nil {
		return err
	}
	var handleBytes uint32
	var handleType int32
	buf := bytes.NewReader(pFileHandleData[:4])
	err = binary.Read(buf, native.Endian, &handleBytes)
	if err != nil {
		return err
	}
	buf = bytes.NewReader(pFileHandleData[4:8])
	err = binary.Read(buf, native.Endian, &handleType)
	if err != nil {
		return err
	}

	key := fileHandleKey{
		handleBytes: handleBytes,
		handleType:  handleType,
	}
	val, ok := process.Res.FileHandleCache[key]
	if !ok {
		return errors.New("didn't find correspondance in the file handle cache")
	}
	msg.Type = ebpfless.SyscallTypeOpen
	msg.Open = &ebpfless.OpenSyscallMsg{
		Filename: val.pathName,
		Flags:    uint32(tracer.ReadArgUint64(regs, 2)),
	}
	return fillFileMetadata(val.pathName, msg.Open, disableStats)
}

func handleUnlinkat(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	fd := tracer.ReadArgInt32(regs, 0)

	filename, err := tracer.ReadArgString(process.Pid, regs, 1)
	if err != nil {
		return err
	}

	flags := tracer.ReadArgInt32(regs, 2)

	filename, err = getFullPathFromFd(process, filename, fd)
	if err != nil {
		return err
	}

	if flags == unix.AT_REMOVEDIR {
		msg.Type = ebpfless.SyscallTypeRmdir
		msg.Rmdir = &ebpfless.RmdirSyscallMsg{
			File: ebpfless.OpenSyscallMsg{
				Filename: filename,
			},
		}
		err = fillFileMetadata(filename, &msg.Rmdir.File, disableStats)
	} else {
		msg.Type = ebpfless.SyscallTypeUnlink
		msg.Unlink = &ebpfless.UnlinkSyscallMsg{
			File: ebpfless.OpenSyscallMsg{
				Filename: filename,
			},
		}
		err = fillFileMetadata(filename, &msg.Unlink.File, disableStats)
	}
	return err
}

func handleUnlink(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	filename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFilename(process, filename)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeUnlink
	msg.Unlink = &ebpfless.UnlinkSyscallMsg{
		File: ebpfless.OpenSyscallMsg{
			Filename: filename,
		},
	}
	return fillFileMetadata(filename, &msg.Unlink.File, disableStats)
}

func handleRmdir(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	filename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}

	filename, err = getFullPathFromFilename(process, filename)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeRmdir
	msg.Rmdir = &ebpfless.RmdirSyscallMsg{
		File: ebpfless.OpenSyscallMsg{
			Filename: filename,
		},
	}
	return fillFileMetadata(filename, &msg.Rmdir.File, disableStats)
}

func handleRename(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	oldFilename, err := tracer.ReadArgString(process.Pid, regs, 0)
	if err != nil {
		return err
	}

	oldFilename, err = getFullPathFromFilename(process, oldFilename)
	if err != nil {
		return err
	}

	newFilename, err := tracer.ReadArgString(process.Pid, regs, 1)
	if err != nil {
		return err
	}

	newFilename, err = getFullPathFromFilename(process, newFilename)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeRename
	msg.Rename = &ebpfless.RenameSyscallMsg{
		OldFile: ebpfless.OpenSyscallMsg{
			Filename: oldFilename,
		},
		NewFile: ebpfless.OpenSyscallMsg{
			Filename: newFilename,
		},
	}
	return fillFileMetadata(oldFilename, &msg.Rename.OldFile, disableStats)
}

func handleRenameAt(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	oldFD := tracer.ReadArgInt32(regs, 0)

	oldFilename, err := tracer.ReadArgString(process.Pid, regs, 1)
	if err != nil {
		return err
	}

	oldFilename, err = getFullPathFromFd(process, oldFilename, oldFD)
	if err != nil {
		return err
	}

	newFD := tracer.ReadArgInt32(regs, 2)

	newFilename, err := tracer.ReadArgString(process.Pid, regs, 3)
	if err != nil {
		return err
	}

	newFilename, err = getFullPathFromFd(process, newFilename, newFD)
	if err != nil {
		return err
	}

	msg.Type = ebpfless.SyscallTypeRename
	msg.Rename = &ebpfless.RenameSyscallMsg{
		OldFile: ebpfless.OpenSyscallMsg{
			Filename: oldFilename,
		},
		NewFile: ebpfless.OpenSyscallMsg{
			Filename: newFilename,
		},
	}
	return fillFileMetadata(oldFilename, &msg.Rename.OldFile, disableStats)
}

func handleFcntl(tracer *Tracer, _ *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	msg.Type = ebpfless.SyscallTypeFcntl
	msg.Fcntl = &ebpfless.FcntlSyscallMsg{
		Fd:  tracer.ReadArgUint32(regs, 0),
		Cmd: tracer.ReadArgUint32(regs, 1),
	}
	return nil
}

func handleDup(tracer *Tracer, _ *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	// using msg to temporary store arg0, as it will be erased by the return value on ARM64
	msg.Dup = &ebpfless.DupSyscallFakeMsg{
		OldFd: tracer.ReadArgInt32(regs, 0),
	}
	return nil
}

func handleClose(tracer *Tracer, process *Process, _ *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	fd := tracer.ReadArgInt32(regs, 0)
	delete(process.Res.Fd, fd)
	return nil
}

//
// handlers called on syscall return
//

func handleNameToHandleAtRet(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	if msg.Open == nil {
		return errors.New("msg empty")
	}

	if ret := tracer.ReadRet(regs); ret < 0 {
		return errors.New("syscall failed")
	}

	pFileHandleData, err := tracer.ReadArgData(process.Pid, regs, 2, 8 /*sizeof uint32 + sizeof int32*/)
	if err != nil {
		return err
	}
	var handleBytes uint32
	var handleType int32
	buf := bytes.NewReader(pFileHandleData[:4])
	err = binary.Read(buf, native.Endian, &handleBytes)
	if err != nil {
		return err
	}
	buf = bytes.NewReader(pFileHandleData[4:8])
	err = binary.Read(buf, native.Endian, &handleType)
	if err != nil {
		return err
	}

	key := fileHandleKey{
		handleBytes: handleBytes,
		handleType:  handleType,
	}
	process.Res.FileHandleCache[key] = &fileHandleVal{
		pathName: msg.Open.Filename,
	}
	return nil
}

func handleOpensRet(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	if ret := tracer.ReadRet(regs); ret > 0 {
		process.Res.Fd[int32(ret)] = msg.Open.Filename
	}
	return nil
}

func handleFcntlRet(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	if ret := tracer.ReadRet(regs); ret >= 0 {
		// maintain fd/path mapping
		if msg.Fcntl.Cmd == unix.F_DUPFD || msg.Fcntl.Cmd == unix.F_DUPFD_CLOEXEC {
			if path, exists := process.Res.Fd[int32(msg.Fcntl.Fd)]; exists {
				process.Res.Fd[int32(ret)] = path
			}
		}
	}
	return nil
}

func handleDupRet(tracer *Tracer, process *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, _ bool) error {
	if ret := tracer.ReadRet(regs); ret >= 0 {
		if path, ok := process.Res.Fd[msg.Dup.OldFd]; ok {
			// maintain fd/path in case of dups
			process.Res.Fd[int32(ret)] = path
		}
	}
	return nil
}

func handleRenamesRet(tracer *Tracer, _ *Process, msg *ebpfless.SyscallMsg, regs syscall.PtraceRegs, disableStats bool) error {
	if ret := tracer.ReadRet(regs); ret == 0 {
		return fillFileMetadata(msg.Rename.NewFile.Filename, &msg.Rename.NewFile, disableStats)
	}
	return nil
}
