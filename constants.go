/**
 * @brief Enduro/X application constants
 *
 * @file constants.go
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
package atmi

/*
#cgo pkg-config: atmisrvinteg

#include <oatmi.h>
#include <nerror.h>
*/
import "C"

// Constants for XATMI APIs and shared error results.
const (
	NDRX_TPSCR_PACKAGE                 = C.NDRX_TPSCR_PACKAGE
	NDRX_TPSCR_REPLACE                 = C.NDRX_TPSCR_REPLACE
	NDRX_TPSCR_FLAT                    = C.NDRX_TPSCR_FLAT
	TPCONVMAXSTR                       = C.TPCONVMAXSTR
	TPTOSTRING                         = C.TPTOSTRING
	TPCONVCLTID                        = C.TPCONVCLTID
	TPCONVTRANID                       = C.TPCONVTRANID
	TPCONVXID                          = C.TPCONVXID
	NDRX_SAGA_NAME_MAX                 = C.NDRX_SAGA_NAME_MAX
	NDRX_SAGA_PAYLOAD_MAX              = C.NDRX_SAGA_PAYLOAD_MAX
	NDRX_SAGA_MIME_MAX                 = C.NDRX_SAGA_MIME_MAX
	NDRX_SAGA_URL_MAX                  = C.NDRX_SAGA_URL_MAX
	NDRX_SAGA_KEY_MAX                  = C.NDRX_SAGA_KEY_MAX
	NDRX_SAGA_ACTION_COMMIT            = C.NDRX_SAGA_ACTION_COMMIT
	NDRX_SAGA_ACTION_ROLLBACK          = C.NDRX_SAGA_ACTION_ROLLBACK
	NDRX_SAGA_ACTION_FORGET            = C.NDRX_SAGA_ACTION_FORGET
	NDRX_SAGA_ACTION_AFTER             = C.NDRX_SAGA_ACTION_AFTER
	NDRX_SAGA_ACTION_STATUS            = C.NDRX_SAGA_ACTION_STATUS
	NDRX_SAGA_URC_ACCEPTED             = C.NDRX_SAGA_URC_ACCEPTED
	NDRX_SAGA_URC_UNKNOWN              = C.NDRX_SAGA_URC_UNKNOWN
	NDRX_SAGA_OUTCOME_CLOSED           = C.NDRX_SAGA_OUTCOME_CLOSED
	NDRX_SAGA_OUTCOME_CANCELLED        = C.NDRX_SAGA_OUTCOME_CANCELLED
	NDRX_SAGA_OUTCOME_FAILED_TO_CLOSE  = C.NDRX_SAGA_OUTCOME_FAILED_TO_CLOSE
	NDRX_SAGA_OUTCOME_FAILED_TO_CANCEL = C.NDRX_SAGA_OUTCOME_FAILED_TO_CANCEL
	NDRX_SAGA_BRXID_HDR                = C.NDRX_SAGA_BRXID_HDR
	NDRX_SAGA_BRXID_MAX                = C.NDRX_SAGA_BRXID_MAX
	TPSAGA_NOLOCK                      = C.TPSAGA_NOLOCK
	TPSAGA_JOINXID                     = C.TPSAGA_JOINXID
	TPSAGA_UPDATE                      = C.TPSAGA_UPDATE
	TPSAGA_URLSONLY                    = C.TPSAGA_URLSONLY
	TPTXS_XID                          = C.TPTXS_XID
	TPTXS_BRXID                        = C.TPTXS_BRXID
	TPTXS_TRANID                       = C.TPTXS_TRANID
	TPTXS_STATUS_NOTFOUND              = C.TPTXS_STATUS_NOTFOUND
	TPTXS_STATUS_ACTIVE                = C.TPTXS_STATUS_ACTIVE
	TPTXS_STATUS_PREPARING             = C.TPTXS_STATUS_PREPARING
	TPTXS_STATUS_COMMITTING            = C.TPTXS_STATUS_COMMITTING
	TPTXS_STATUS_ABORTING              = C.TPTXS_STATUS_ABORTING
	TPTXS_STATUS_COMMIT_HEURISTIC      = C.TPTXS_STATUS_COMMIT_HEURISTIC
	TPTXS_STATUS_ABORT_HEURISTIC       = C.TPTXS_STATUS_ABORT_HEURISTIC
	TPTXS_STATUS_COMMITTED             = C.TPTXS_STATUS_COMMITTED
	TPTXS_STATUS_ABORTED               = C.TPTXS_STATUS_ABORTED
	NDRX_TRACEPARENT_LEN               = C.NDRX_TRACEPARENT_LEN
	TPQKEEPORIG                        = C.TPQKEEPORIG
	TPTXOPENALL                        = C.TPTXOPENALL
	TPSAGALRA                          = C.TPSAGALRA
	NDRX_B58_CIDMAX                    = C.NDRX_B58_CIDMAX
	NDRX_DLM_LM_NULL                   = C.NDRX_DLM_LM_NULL
	NDRX_DLM_LM_CONCURRENT_READ        = C.NDRX_DLM_LM_CONCURRENT_READ
	NDRX_DLM_LM_CONCURRENT_WRITE       = C.NDRX_DLM_LM_CONCURRENT_WRITE
	NDRX_DLM_LM_PROTECTED_READ         = C.NDRX_DLM_LM_PROTECTED_READ
	NDRX_DLM_LM_PROTECTED_WRITE        = C.NDRX_DLM_LM_PROTECTED_WRITE
	NDRX_DLM_LM_EXCLUSIVE              = C.NDRX_DLM_LM_EXCLUSIVE
	NDRX_DLM_ISLOCKED_NOLOCK           = C.NDRX_DLM_ISLOCKED_NOLOCK
	NDRX_DLM_ISLOCKED_LKDATA           = C.NDRX_DLM_ISLOCKED_LKDATA
	NDRX_DLM_ISLOCKED_EPRECOND         = C.NDRX_DLM_ISLOCKED_EPRECOND
	NDRX_DLM_ISLOCKED_LKMATCH          = C.NDRX_DLM_ISLOCKED_LKMATCH
	NDRX_DLM_ISLOCKED_LKNL             = C.NDRX_DLM_ISLOCKED_LKNL
	NDRX_DLM_ISLOCKED_LKCR             = C.NDRX_DLM_ISLOCKED_LKCR
	NDRX_DLM_ISLOCKED_LKCW             = C.NDRX_DLM_ISLOCKED_LKCW
	NDRX_DLM_ISLOCKED_LKPR             = C.NDRX_DLM_ISLOCKED_LKPR
	NDRX_DLM_ISLOCKED_LKPW             = C.NDRX_DLM_ISLOCKED_LKPW
	NDRX_DLM_ISLOCKED_LKEX             = C.NDRX_DLM_ISLOCKED_LKEX
	NDRX_DLM_ISLOCKED_LKONL            = C.NDRX_DLM_ISLOCKED_LKONL
	NDRX_DLM_ISLOCKED_LKOCR            = C.NDRX_DLM_ISLOCKED_LKOCR
	NDRX_DLM_ISLOCKED_LKOCW            = C.NDRX_DLM_ISLOCKED_LKOCW
	NDRX_DLM_ISLOCKED_LKOPR            = C.NDRX_DLM_ISLOCKED_LKOPR
	NDRX_DLM_ISLOCKED_LKOPW            = C.NDRX_DLM_ISLOCKED_LKOPW
	NDRX_DLM_ISLOCKED_LKOEX            = C.NDRX_DLM_ISLOCKED_LKOEX
	NDRX_DLM_FL_FLATUBF                = C.NDRX_DLM_FL_FLATUBF
	NDRX_DLM_FL_FLATBLOB               = C.NDRX_DLM_FL_FLATBLOB
	NDRX_DLM_FL_CONVERT                = C.NDRX_DLM_FL_CONVERT
	NDRX_DLM_FL_MIGRATER               = C.NDRX_DLM_FL_MIGRATER
	NDRX_DLM_FL_NOBLOCK                = C.NDRX_DLM_FL_NOBLOCK
	NDRX_DLM_FL_REGEX                  = C.NDRX_DLM_FL_REGEX
	NDRX_DLM_FL_NODATA                 = C.NDRX_DLM_FL_NODATA
	NDRX_DLM_FL_EXPIRY                 = C.NDRX_DLM_FL_EXPIRY
	NDRX_DLM_FL_EXPFENCE               = C.NDRX_DLM_FL_EXPFENCE
	NDRX_DLM_FL_EXPEVENT               = C.NDRX_DLM_FL_EXPEVENT
	NDRX_DLM_FL_FENCEEVENT             = C.NDRX_DLM_FL_FENCEEVENT
	NDRX_DLM_FL_EXPPIDFENCE            = C.NDRX_DLM_FL_EXPPIDFENCE
	NDRX_DLM_FL_LOCALINFO              = C.NDRX_DLM_FL_LOCALINFO
	NDRX_DLM_FL_KEYASCID               = C.NDRX_DLM_FL_KEYASCID
	NDRX_DLM_FL_BYPID                  = C.NDRX_DLM_FL_BYPID
	NDRX_DLM_FL_NOPIDCHK               = C.NDRX_DLM_FL_NOPIDCHK
	NDRX_DLM_FL_IGNOWNER               = C.NDRX_DLM_FL_IGNOWNER
	NDRX_DLM_FL_PRECOND                = C.NDRX_DLM_FL_PRECOND
	NDRX_DLM_FL_LINREAD                = C.NDRX_DLM_FL_LINREAD
	NDRX_DLM_FL_FORCE                  = C.NDRX_DLM_FL_FORCE
	NDRX_DLM_FL_EVFLAT                 = C.NDRX_DLM_FL_EVFLAT
	NDRX_DLM_CMD_LOCK                  = C.NDRX_DLM_CMD_LOCK
	NDRX_DLM_CMD_UNLOCK                = C.NDRX_DLM_CMD_UNLOCK
	NDRX_DLM_CMD_ISLOCKED              = C.NDRX_DLM_CMD_ISLOCKED
	NDRX_DLM_CMD_SET                   = C.NDRX_DLM_CMD_SET
	NDRX_DLM_CMD_CAS                   = C.NDRX_DLM_CMD_CAS
	NDRX_DLM_CMD_GET                   = C.NDRX_DLM_CMD_GET
	NDRX_DLM_CMD_DEL                   = C.NDRX_DLM_CMD_DEL
	NDRX_DLM_CMD_RENAME                = C.NDRX_DLM_CMD_RENAME
	NDRX_DLM_CMD_EXISTS                = C.NDRX_DLM_CMD_EXISTS
	NDRX_DLM_CMD_DROPDB                = C.NDRX_DLM_CMD_DROPDB
	NDRX_DLM_CMD_EXPIRELOCK            = C.NDRX_DLM_CMD_EXPIRELOCK
	NDRX_DLM_CMD_FENCELOCK             = C.NDRX_DLM_CMD_FENCELOCK
	NDRX_DLM_CMD_EXPIREDATA            = C.NDRX_DLM_CMD_EXPIREDATA
	NDRX_DLM_OP_DLMCMD                 = C.NDRX_DLM_OP_DLMCMD
	NDRX_DLM_OP_CANCEL                 = C.NDRX_DLM_OP_CANCEL
	NDRX_DLM_OP_LSSTORE                = C.NDRX_DLM_OP_LSSTORE
	NDRX_DLM_OP_ADDNODE                = C.NDRX_DLM_OP_ADDNODE
	NDRX_DLM_OP_RMNODE                 = C.NDRX_DLM_OP_RMNODE
	NDRX_DLM_OP_XFER                   = C.NDRX_DLM_OP_XFER
	NDRX_DLM_OP_FENCEOK                = C.NDRX_DLM_OP_FENCEOK
	NDRX_DLM_OP_DLMINIT                = C.NDRX_DLM_OP_DLMINIT
	NDRX_DLM_OP_COMPACT                = C.NDRX_DLM_OP_COMPACT
	NDRX_DLM_OP_T_DLM                  = C.NDRX_DLM_OP_T_DLM
	NDRX_DLM_FENCE_CODE_BUFSZ          = C.NDRX_DLM_FENCE_CODE_BUFSZ
	NDRX_DLM_FENCE_ST_SCH              = C.NDRX_DLM_FENCE_ST_SCH
	NDRX_DLM_FENCE_ST_ATT              = C.NDRX_DLM_FENCE_ST_ATT
	NDRX_DLM_FENCE_ST_EXE              = C.NDRX_DLM_FENCE_ST_EXE
	NDRX_DLM_FENCE_ST_EXD              = C.NDRX_DLM_FENCE_ST_EXD
	NDRX_DLM_FENCE_ST_FAI              = C.NDRX_DLM_FENCE_ST_FAI
	NDRX_DLM_FENCE_ST_STL              = C.NDRX_DLM_FENCE_ST_STL
	NDRX_DLM_FENCE_ST_FOK              = C.NDRX_DLM_FENCE_ST_FOK
	NDRX_DLM_FENCE_TY_INF              = C.NDRX_DLM_FENCE_TY_INF
	NDRX_DLM_FENCE_TY_LOC              = C.NDRX_DLM_FENCE_TY_LOC
	NDRX_DLM_FENCE_TY_GLB              = C.NDRX_DLM_FENCE_TY_GLB
	NDRX_DLM_NODE_FENCE_EV             = C.NDRX_DLM_NODE_FENCE_EV
	NDRX_DLM_MEMBERSHIP_EV             = C.NDRX_DLM_MEMBERSHIP_EV
	NDRX_DLM_STATE_EV                  = C.NDRX_DLM_STATE_EV
	NDRX_DLM_XFER_EV                   = C.NDRX_DLM_XFER_EV
	TPCTXTAG_PROCESS                   = C.TPCTXTAG_PROCESS
	TPCTXTAG_THREAD                    = C.TPCTXTAG_THREAD
	TPCTXTAG_NEXT                      = C.TPCTXTAG_NEXT
	TPCTXTAG_LOG                       = C.TPCTXTAG_LOG
	TPCTXTAG_TRACE                     = C.TPCTXTAG_TRACE
	TPCTXTAG_BOTH                      = C.TPCTXTAG_BOTH
	TPTRACE_F_ACTIVE                   = C.TPTRACE_F_ACTIVE
	TPTRACE_F_NEW_ROOT                 = C.TPTRACE_F_NEW_ROOT
	TPTRACE_STATUS_UNSET               = C.TPTRACE_STATUS_UNSET
	TPTRACE_STATUS_OK                  = C.TPTRACE_STATUS_OK
	TPTRACE_STATUS_ERROR               = C.TPTRACE_STATUS_ERROR
	NDRX_TRACE_STATE_MAX               = C.NDRX_TRACE_STATE_MAX
	NDRX_SCR_CFG_VERSION_1             = C.NDRX_SCR_CFG_VERSION_1
	NDRX_SCR_INITSTRING_BUFSZ          = C.NDRX_SCR_INITSTRING_BUFSZ
	NDRX_SCR_INITSTRING_PYTHON         = C.NDRX_SCR_INITSTRING_PYTHON
	NMINVAL                            = C.NMINVAL
	NEPLUGIN                           = C.NEPLUGIN
	NENOSPACE                          = C.NENOSPACE
	NEINVALKEY                         = C.NEINVALKEY
	NENOENT                            = C.NENOENT
	NEWRITE                            = C.NEWRITE
	NEEXEC                             = C.NEEXEC
	NESUPPORT                          = C.NESUPPORT
	NEEXISTS                           = C.NEEXISTS
	NEVERSION                          = C.NEVERSION
	NEBUSY                             = C.NEBUSY
	NESTALE                            = C.NESTALE
	NEEOF                              = C.NEEOF
	NESYNC                             = C.NESYNC
	NECONSTRAINT                       = C.NECONSTRAINT
	NESTATE                            = C.NESTATE
	NEUNKNOWN                          = C.NEUNKNOWN
	NESTMT                             = C.NESTMT
	NEACCESS                           = C.NEACCESS
	NECONTEXT                          = C.NECONTEXT
	NEMATCH                            = C.NEMATCH
	NEPROTO                            = C.NEPROTO
	NEUBF                              = C.NEUBF
	NEATMI                             = C.NEATMI
	NEMSGSIZE                          = C.NEMSGSIZE
	NECANCELED                         = C.NECANCELED
	NEINPROGRESS                       = C.NEINPROGRESS
	NEDEADLK                           = C.NEDEADLK
	NEPRECOND                          = C.NEPRECOND
	NEUNAVAILABLE                      = C.NEUNAVAILABLE
	NMAXVAL                            = C.NMAXVAL
)

/* vim: set ts=4 sw=4 et smartindent: */
