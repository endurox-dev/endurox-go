/**
 * @brief Context and handle adapter for Go scripting callbacks
 * @file script_bridge.c
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
#include "script_bridge.h"
#include <stdlib.h>
#include <string.h>
#include "_cgo_export.h"

/** Go frame handles are valid only during synchronous callbacks. */
static __thread uintptr_t M_script_frame=0;

/** @return active Go callback-frame handle on this thread, or zero */
uintptr_t go_script_frame_get(void)
{
    return M_script_frame;
}

/**
 * Set the callback frame while retaining only an integer token in C.
 * @param handle new callback-frame handle
 * @return previous callback-frame handle
 */
uintptr_t go_script_frame_swap(uintptr_t handle)
{
    uintptr_t previous=M_script_frame;
    M_script_frame=handle;
    return previous;
}

/**
 * Dispatch a script callback through Go, restoring the caller's context.
 * @param vm scripting VM
 * @param name callback name
 * @param param parameter buffer slot, borrowed for the call
 * @param param_len parameter length slot
 * @param idata input buffer slot, borrowed for the call
 * @param ilen input length slot
 * @param odata output buffer slot, ownership returned to the native caller
 * @param olen output length slot
 * @param user C-allocated uintptr_t containing the Go callback handle
 * @param flags callback flags
 * @return callback result or EXFAIL on a context-switch failure
 */
static int go_script_dispatch(tpscr_vm_t *vm, const char *name,
        char **param, long *param_len, char **idata, long *ilen,
        char **odata, long *olen, void *user, long flags)
{
    int ret=EXFAIL;
    TPCONTEXT_T context=TPNULLCONTEXT;
    uintptr_t handle=*(uintptr_t *)user;

    if (TPMULTICONTEXTS!=tpgetctxt(&context, 0))
    {
        return EXFAIL;
    }

    ret=goScriptCallback(context, vm, (char *)name, param, param_len,
            idata, ilen, odata, olen, handle, flags);

    if (EXSUCCEED!=tpsetctxt(context, 0))
    {
        ret=EXFAIL;
    }

    return ret;
}

/**
 * Copy failure diagnostics while still on the failing thread. The ATMI error
 * is read first, as the scripting error queries reset it.
 * @param context object-API context
 * @param err destination
 */
static void go_script_err_capture(TPCONTEXT_T *context, go_script_err_t *err)
{
    const char *text;

    err->atmi_code=Otperrno(context);
    text=Otpstrerror(context, err->atmi_code);
    err->atmi_text=(NULL!=text ? strdup(text) : NULL);
    err->script_code=Otpscrerrno(context, NULL);
    text=Otpscrerror(context, NULL);
    err->script_text=(NULL!=text ? strdup(text) : NULL);
}

tpscr_vm_t *go_tpscrinit(TPCONTEXT_T *context, const tpscr_cfg_t *cfg,
        long flags, go_script_err_t *err)
{
    tpscr_vm_t *vm=Otpscrinit(context, cfg, flags);

    if (NULL==vm)
    {
        go_script_err_capture(context, err);
    }

    return vm;
}

int go_tpscruninit(TPCONTEXT_T *context, tpscr_vm_t *vm, long flags,
        go_script_err_t *err)
{
    int ret=Otpscruninit(context, vm, flags);

    if (EXSUCCEED!=ret)
    {
        go_script_err_capture(context, err);
    }

    return ret;
}

int go_tpscrregcb(TPCONTEXT_T *context, tpscr_vm_t *vm,
        const char *name, void *token, long flags)
{
    return Otpscrregcb(context, vm, name, go_script_dispatch, token, flags);
}
/* vim: set ts=4 sw=4 et smartindent: */
