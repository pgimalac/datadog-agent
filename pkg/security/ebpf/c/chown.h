#ifndef _CHOWN_H_
#define _CHOWN_H_

#include "syscalls.h"

struct chown_event_t {
    struct kevent_t event;
    struct process_context_t process;
    struct span_context_t span;
    struct container_context_t container;
    struct syscall_t syscall;
    struct file_t file;
    uid_t uid;
    gid_t gid;
};

int __attribute__((always_inline)) chown_approvers(struct syscall_cache_t *syscall) {
    return basename_approver(syscall, syscall->setattr.dentry, EVENT_CHOWN);
}

int __attribute__((always_inline)) trace__sys_chown(uid_t user, gid_t group) {
    // if we are already working on this syscall, we simply skip it
    struct syscall_cache_t *syscall_check = peek_syscall(EVENT_CHOWN);
    if (syscall_check) {
        return 0;
    }

    struct policy_t policy = fetch_policy(EVENT_CHOWN);
    if (is_discarded_by_process(policy.mode, EVENT_CHOWN)) {
        return 0;
    }

    struct syscall_cache_t syscall = {
        .type = EVENT_CHOWN,
        .policy = policy,
        .setattr = {
            .user = user,
            .group = group
        }
    };

    cache_syscall(&syscall);

    return 0;
}

SEC("kprobe/do_fchownat")
int kprobe_do_fchownat(struct pt_regs *ctx) {
    uid_t user = PT_REGS_PARM3(ctx);
    uid_t group = PT_REGS_PARM4(ctx);
    return trace__sys_chown(user, group);
}

SEC("kprobe/vfs_fchown")
int kprobe_vfs_fchown(struct pt_regs *ctx) {
    uid_t user = PT_REGS_PARM2(ctx);
    uid_t group = PT_REGS_PARM3(ctx);
    return trace__sys_chown(user, group);
}

SEC("kprobe/ksys_fchown")
int kprobe_ksys_fchown(struct pt_regs *ctx) {
    uid_t user = PT_REGS_PARM2(ctx);
    uid_t group = PT_REGS_PARM3(ctx);
    return trace__sys_chown(user, group);
}

SEC("kprobe/chown_common")
int kprobe_chown_common(struct pt_regs *ctx) {
    uid_t user = PT_REGS_PARM2(ctx);
    uid_t group = PT_REGS_PARM3(ctx);
    return trace__sys_chown(user, group);
}

int __attribute__((always_inline)) sys_chown_ret(void *ctx, int retval) {
    struct syscall_cache_t *syscall = pop_syscall(EVENT_CHOWN);
    if (!syscall) {
        return 0;
    }

    if (IS_UNHANDLED_ERROR(retval)) {
        return 0;
    }

    if (is_pipefs_mount_id(syscall->setattr.file.path_key.mount_id)) {
        return 0;
    }

    struct chown_event_t event = {
        .syscall.retval = retval,
        .file = syscall->setattr.file,
        .uid = syscall->setattr.user,
        .gid = syscall->setattr.group,
    };

    struct proc_cache_t *entry = fill_process_context(&event.process);
    fill_container_context(entry, &event.container);
    fill_span_context(&event.span);

    // dentry resolution in setattr.h

    send_event(ctx, EVENT_CHOWN, event);

    return 0;
}

int __attribute__((always_inline)) kprobe_sys_chown_ret(struct pt_regs *ctx) {
    int retval = PT_REGS_RC(ctx);
    return sys_chown_ret(ctx, retval);
}

SEC("kretprobe/do_fchownat")
int kretprobe_do_fchownat(struct pt_regs *ctx) {
    return kprobe_sys_chown_ret(ctx);
}

SEC("kretprobe/vfs_fchown")
int kretprobe_vfs_fchown(struct pt_regs *ctx) {
    return kprobe_sys_chown_ret(ctx);
}

SEC("kretprobe/ksys_fchown")
int kretprobe_ksys_fchown(struct pt_regs *ctx) {
    return kprobe_sys_chown_ret(ctx);
}

SEC("kretprobe/chown_common")
int kretprobe_chown_common(struct pt_regs *ctx) {
    return kprobe_sys_chown_ret(ctx);
}

SEC("tracepoint/handle_sys_chown_exit")
int tracepoint_handle_sys_chown_exit(struct tracepoint_raw_syscalls_sys_exit_t *args) {
    return sys_chown_ret(args, args->ret);
}

#endif
