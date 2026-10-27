/**
 * @brief C adapter for Go scripting callbacks
 * @file script_bridge.h
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2026, Mavimax, Ltd. All Rights Reserved.
 * This software is released under one of the following licenses:
 * LGPL or Mavimax's license for commercial use.
 * See LICENSE file for full text.
 *
 * C (as designed by Dennis Ritchie and later authors) language code is licensed
 * under Enduro/X Modified GNU Affero General Public License, version 3.
 * See LICENSE_C file for full text.
 * -----------------------------------------------------------------------------
 * LGPL license:
 *
 * This program is free software; you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License, version 3 as published
 * by the Free Software Foundation;
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
 * PARTICULAR PURPOSE. See the GNU Lesser General Public License, version 3
 * for more details.
 *
 * You should have received a copy of the Lesser General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA
 *
 * -----------------------------------------------------------------------------
 * A commercial use license is available from Mavimax, Ltd
 * contact@mavimax.com
 * -----------------------------------------------------------------------------
 */
/** Header guard for the Go scripting callback adapter. */
#ifndef ENDUROX_GO_SCRIPT_BRIDGE_H
#define ENDUROX_GO_SCRIPT_BRIDGE_H

#include <stdint.h>
#include <oatmi.h>

/**
 * Register the Go callback adapter with the selected VM.
 * @param context object-API context
 * @param vm scripting VM
 * @param name callback registration name
 * @param token C-allocated uintptr_t containing a runtime/cgo.Handle
 * @param flags native callback flags
 * @return native tpscrregcb result
 */
int go_tpscrregcb(TPCONTEXT_T *context, tpscr_vm_t *vm,
        const char *name, void *token, long flags);

/**
 * Diagnostics of a failed VM initialization or shutdown, read on the thread of
 * the failing call: the provider keeps them thread-local when no VM is usable.
 */
typedef struct
{
    int atmi_code;      /**< tperrno of the context after the failure */
    char *atmi_text;    /**< malloc'ed tpstrerror() text, NULL if unavailable */
    int script_code;    /**< provider error code without a VM */
    char *script_text;  /**< malloc'ed provider message, NULL if unavailable */
} go_script_err_t;

/**
 * Create a VM; on failure capture its diagnostics in the same native call.
 * @param context object-API context
 * @param cfg provider configuration, or NULL
 * @param flags native initialization flags
 * @param err diagnostics on failure; free its texts with free()
 * @return native VM, or NULL with err filled
 */
tpscr_vm_t *go_tpscrinit(TPCONTEXT_T *context, const tpscr_cfg_t *cfg,
        long flags, go_script_err_t *err);

/**
 * Destroy a VM; on failure capture its diagnostics in the same native call.
 * @param context object-API context
 * @param vm native VM, not used again after this call
 * @param flags native shutdown flags
 * @param err diagnostics on failure; free its texts with free()
 * @return native tpscruninit result
 */
int go_tpscruninit(TPCONTEXT_T *context, tpscr_vm_t *vm, long flags,
        go_script_err_t *err);

/**
 * Free compiler output; capture failure diagnostics on the same native thread.
 * @param context object-API context
 * @param bytes provider-owned compiler allocation
 * @param err diagnostics on failure; free its texts with free()
 * @return native tpscrfree result
 */
int go_tpscrfree(TPCONTEXT_T *context, char *bytes, go_script_err_t *err);

/**
 * Read the active callback-frame handle on this native thread.
 * @return runtime/cgo.Handle token, or zero outside a callback
 */
uintptr_t go_script_frame_get(void);

/**
 * Replace the active callback-frame handle on this native thread.
 * @param handle runtime/cgo.Handle token, or zero when restoring no callback
 * @return previous frame handle
 */
uintptr_t go_script_frame_swap(uintptr_t handle);

#endif
/* vim: set ts=4 sw=4 et smartindent: */
