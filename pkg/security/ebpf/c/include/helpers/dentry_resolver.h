#ifndef _HELPERS_DENTRY_RESOLVER_H_
#define _HELPERS_DENTRY_RESOLVER_H_

#include "constants/custom.h"
#include "maps.h"

#include "buffer_selector.h"

// int __attribute__((always_inline)) resolve_dentry(void *ctx, int dr_type) {
//     if (dr_type == DR_KPROBE) {
//         bpf_tail_call_compat(ctx, &dentry_resolver_kprobe_progs, DR_KPROBE_AD_FILTER_KEY);
//     } else if (dr_type == DR_TRACEPOINT) {
//         bpf_tail_call_compat(ctx, &dentry_resolver_tracepoint_progs, DR_TRACEPOINT_AD_FILTER_KEY);
//     }
//     return 0;
// }

int __attribute__((always_inline)) monitor_resolution_err(u32 resolution_err) {
    if (resolution_err > 0) {
        struct bpf_map_def *erpc_stats = select_buffer(&dr_erpc_stats_fb, &dr_erpc_stats_bb, ERPC_MONITOR_KEY);
        if (erpc_stats == NULL) {
            return 0;
        }

        struct dr_erpc_stats_t *stats = bpf_map_lookup_elem(erpc_stats, &resolution_err);
        if (stats == NULL) {
            return 0;
        }
        __sync_fetch_and_add(&stats->count, 1);
    }
    return 0;
}

// u32 __attribute__((always_inline)) parse_erpc_request(struct dr_erpc_state_t *state, void *data) {
//     u32 err = 0;
//     int ret = bpf_probe_read(&state->key, sizeof(state->key), data);
//     if (ret < 0) {
//         err = DR_ERPC_READ_PAGE_FAULT;
//         goto exit;
//     }
//     ret = bpf_probe_read(&state->userspace_buffer, sizeof(state->userspace_buffer), data + sizeof(state->key));
//     if (ret < 0) {
//         err = DR_ERPC_READ_PAGE_FAULT;
//         goto exit;
//     }
//     ret = bpf_probe_read(&state->buffer_size, sizeof(state->buffer_size), data + sizeof(state->key) + sizeof(state->userspace_buffer));
//     if (ret < 0) {
//         err = DR_ERPC_READ_PAGE_FAULT;
//         goto exit;
//     }
//     ret = bpf_probe_read(&state->challenge, sizeof(state->challenge), data + sizeof(state->key) + sizeof(state->userspace_buffer) + sizeof(state->buffer_size));
//     if (ret < 0) {
//         err = DR_ERPC_READ_PAGE_FAULT;
//         goto exit;
//     }

//     state->iteration = 0;
//     state->ret = 0;
//     state->cursor = 0;

// exit:
//     return err;
// }

// int __attribute__((always_inline)) handle_dr_request(struct pt_regs *ctx, void *data, u32 dr_erpc_key) {
//     u32 resolution_err = parse_erpc_request(state, data);
//     if (resolution_err > 0) {
//         goto exit;
//     }

//     bpf_tail_call_compat(ctx, &dentry_resolver_kprobe_progs, dr_erpc_key);

// exit:
//     monitor_resolution_err(resolution_err);
//     return 0;
// }

int __attribute__((always_inline)) handle_resolve_parent_dentry(void *ctx, void *data) {
    u32 zero = 0;
    struct dr_erpc_state_t *state = bpf_map_lookup_elem(&dr_erpc_state, &zero);
    if (!state) {
        return 0;
    }

    u32 err = 0;
    int ret = bpf_probe_read(&state->key, sizeof(state->key), data);
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->userspace_buffer, sizeof(state->userspace_buffer), data + sizeof(state->key));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->buffer_size, sizeof(state->buffer_size), data + sizeof(state->key) + sizeof(state->userspace_buffer));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->challenge, sizeof(state->challenge), data + sizeof(state->key) + sizeof(state->userspace_buffer) + sizeof(state->buffer_size));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }

    state->iteration = 0;
    state->ret = 0;
    state->cursor = 0;

    bpf_tail_call_compat(ctx, &erpc_progs, ERPC_RESOLVE_PARENT_DENTRY_KEY);
    err = DR_ERPC_TAIL_CALL_ERROR;

exit:
    monitor_resolution_err(err);
    return 0;
}

int __attribute__((always_inline)) handle_resolve_pathsegment(void *ctx, void *data) {
    u32 zero = 0;
    struct dr_erpc_state_t *state = bpf_map_lookup_elem(&dr_erpc_state, &zero);
    if (!state) {
        return 0;
    }

    u32 err = 0;
    int ret = bpf_probe_read(&state->userspace_buffer, sizeof(state->userspace_buffer), data);
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->buffer_size, sizeof(state->buffer_size), data + sizeof(state->userspace_buffer));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->path_ref.cpu, sizeof(state->path_ref.cpu), data + sizeof(state->userspace_buffer) + sizeof(state->buffer_size));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->path_ref.read_cursor, sizeof(state->path_ref.read_cursor), data + sizeof(state->userspace_buffer) + sizeof(state->buffer_size) + sizeof(state->path_ref.cpu));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->path_ref.len, sizeof(state->path_ref.len), data + sizeof(state->userspace_buffer) + sizeof(state->buffer_size) + sizeof(state->path_ref.cpu) + sizeof(state->path_ref.read_cursor));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }
    ret = bpf_probe_read(&state->challenge, sizeof(state->challenge), data + sizeof(state->userspace_buffer) + sizeof(state->buffer_size) + sizeof(state->path_ref.cpu) + sizeof(state->path_ref.read_cursor) + sizeof(state->path_ref.len));
    if (ret < 0) {
        err = DR_ERPC_READ_PAGE_FAULT;
        goto exit;
    }


    u32 total_len = sizeof(state->challenge) + sizeof(state->path_ref.watermark) * 2 + state->path_ref.len;
    if (total_len > state->buffer_size) {
        err = DR_ERPC_BUFFER_SIZE;
        goto exit;
    }

    if (state->path_ref.read_cursor >= PR_RING_BUFFER_SIZE || total_len  > PR_RING_BUFFER_SIZE) {
        err = DR_ERPC_CACHE_MISS; // TODO: use a specific error type for malformed request
        goto exit;
    }

    ret = bpf_probe_write_user((void *)state->userspace_buffer, &state->challenge, sizeof(state->challenge));
    if (ret < 0) {
        err = ret == -14 ? DR_ERPC_WRITE_PAGE_FAULT : DR_ERPC_UNKNOWN_ERROR;
        goto exit;
    }

    state->iteration = 0;
    state->ret = 0;
    state->cursor = sizeof(state->challenge);

    bpf_tail_call_compat(ctx, &erpc_progs, ERPC_RESOLVE_PATHSEGMENT_KEY);
    err = DR_ERPC_TAIL_CALL_ERROR;

exit:
    monitor_resolution_err(err);
    return 0;
}

#endif
