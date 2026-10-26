/**
 * @brief Native XA connection fixture
 * @file mock.c
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


#include <atmi.h>
#include <xa.h>
#include <tmenv.h>
#include <cconfig.h>

/** Opaque provider-owned connection marker. */
static long M_cookie[NDRX_CCONFIG_XA_MAX]={0};

/** @param rmid resource ID @return borrowed fixture connection */
static void *connection(int rmid)
{
    return &M_cookie[rmid];
}
/** @param info open/close string @param rmid RMID @param flags XA flags @return XA_OK */
static int info_entry(char *info, int rmid, long flags)
{
    (void)info; (void)rmid; (void)flags;
    return XA_OK;
}
/** @param xid transaction @param rmid RMID @param flags XA flags @return XA_OK */
static int xid_entry(XID *xid, int rmid, long flags)
{
    (void)xid; (void)rmid; (void)flags;
    return XA_OK;
}
/** @param xids output @param count capacity @param rmid RMID @param flags scan flags @return empty scan */
static int recover_entry(XID *xids, long count, int rmid, long flags)
{
    (void)xids; (void)count; (void)rmid; (void)flags;
    return 0;
}
/** @param handle async handle @param result result @param rmid RMID @param flags flags @return unsupported */
static int complete_entry(int *handle, int *result, int rmid, long flags)
{
    (void)handle; (void)result; (void)rmid; (void)flags;
    return XAER_INVAL;
}
/** Test XA switch; operations are synchronous and have no external effects. */
static struct xa_switch_t M_switch={
    .name="GOTEST", .flags=TMNOFLAGS, .version=0,
    .xa_open_entry=info_entry, .xa_close_entry=info_entry,
    .xa_start_entry=xid_entry, .xa_end_entry=xid_entry,
    .xa_rollback_entry=xid_entry, .xa_prepare_entry=xid_entry,
    .xa_commit_entry=xid_entry, .xa_recover_entry=recover_entry,
    .xa_forget_entry=xid_entry, .xa_complete_entry=complete_entry
};
/** @param rmid configured RMID @return registered fixture switch */
struct xa_switch_t *ndrx_get_xa_switch(int rmid)
{
    ndrx_xa_setgetconnn(rmid, connection);
    return &M_switch;
}
