/**
 * @brief XATMI main package
 * 	NOTES for finalizers! Note that if we pass from finalized object (typed ubf,
 *	expression tree or ATMI Context) pointer to C, the and the function call
 *	for the object is last in the object's go scope, the go might being to GC
 *	the go object, while the C function have received pointer in args. Thus
 *  C side in the middle of processing might get destructed object. e.g.
 *  c_str := C.BPrintStrC(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))
 *  after enter in C.BPrintStrC(), the GC might kill the u object. Thus to avoid
 *  this, we create a defered "no-op" call in the enter of go func. With "u.Buf.nop()"
 *	at the end of the functions.
 *
 * @file atmi.go
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2019, Mavimax, Ltd. All Rights Reserved.
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

#include <ndebug.h>
#include <string.h>
#include <stdlib.h>
#include <oatmi.h>
#include <onerror.h>

// Wrapper for TPNIT
static int go_tpinit(TPCONTEXT_T *p_ctxt) {
	return Otpinit(p_ctxt, NULL);
}

//ATMI library error code
static int go_tperrno(TPCONTEXT_T *p_ctxt) {
	return Otperrno(p_ctxt);
}

//Standard library error code
static int go_Nerror(TPCONTEXT_T *p_ctxt) {
	return ONerror(p_ctxt);
}


static void free_string(char* s) { free(s); }
static char * malloc_string(int size) { return malloc(size); }


static void go_param_to_tpqctl(
		TPQCTL *qc,
		long *ctl_flags,
		long *ctl_deq_time,
		long *ctl_priority,
		long *ctl_diagnostic,
		char *ctl_diagmsg,
		char *ctl_msgid,
		char *ctl_corrid,
		char *ctl_replyqueue,
		char *ctl_failurequeue,
		char *ctl_cltid,
		long *ctl_urcode,
		long *ctl_appkey,
		long *ctl_delivery_qos,
		long *ctl_reply_qos,
		long *ctl_exp_time)
{
	qc->flags = *ctl_flags;
	qc->deq_time = *ctl_deq_time;
	qc->priority = *ctl_priority;
	qc->diagnostic = *ctl_diagnostic;
	strcpy(qc->diagmsg, ctl_diagmsg);
	memcpy(qc->msgid, ctl_msgid, TMMSGIDLEN);
	memcpy(qc->corrid, ctl_corrid, TMCORRIDLEN);
	strcpy(qc->replyqueue, ctl_replyqueue);
	strcpy(qc->failurequeue, ctl_failurequeue);
	strcpy(qc->cltid.clientdata, ctl_cltid);
	qc->urcode = *ctl_urcode;
	qc->appkey = *ctl_appkey;
	qc->delivery_qos = *ctl_delivery_qos;
	qc->reply_qos = *ctl_reply_qos;
	qc->exp_time = *ctl_exp_time;
}

static void go_tpqctl_to_param(
		TPQCTL *qc,
		long *ctl_flags,
		long *ctl_deq_time,
		long *ctl_priority,
		long *ctl_diagnostic,
		char *ctl_diagmsg,
		char *ctl_msgid,
		char *ctl_corrid,
		char *ctl_replyqueue,
		char *ctl_failurequeue,
		char *ctl_cltid,
		long *ctl_urcode,
		long *ctl_appkey,
		long *ctl_delivery_qos,
		long *ctl_reply_qos,
		long *ctl_exp_time)
{
	*ctl_flags = qc->flags;
	*ctl_deq_time = qc->deq_time;
	*ctl_priority = qc->priority;
	*ctl_diagnostic = qc->diagnostic;
	strcpy(ctl_diagmsg, qc->diagmsg);
	memcpy(ctl_msgid, qc->msgid, TMMSGIDLEN);
	memcpy(ctl_corrid, qc->corrid, TMCORRIDLEN);
	strcpy(ctl_replyqueue, qc->replyqueue);
	strcpy(ctl_failurequeue, qc->failurequeue);
	strcpy(ctl_cltid, qc->cltid.clientdata);
	qc->urcode = *ctl_urcode;
	qc->appkey = *ctl_appkey;
	qc->delivery_qos = *ctl_delivery_qos;
	qc->reply_qos = *ctl_reply_qos;
	qc->exp_time = *ctl_exp_time;
}

static int go_tpenqueue (TPCONTEXT_T *p_ctx, char *qspace, char *qname, char *data, long len, long flags,
		long *ctl_flags,
		long *ctl_deq_time,
		long *ctl_priority,
		long *ctl_diagnostic,
		char *ctl_diagmsg,
		char *ctl_msgid,
		char *ctl_corrid,
		char *ctl_replyqueue,
		char *ctl_failurequeue,
		char *ctl_cltid,
		long *ctl_urcode,
		long *ctl_appkey,
		long *ctl_delivery_qos,
		long *ctl_reply_qos,
		long *ctl_exp_time
)
{
	int ret;
	TPQCTL qc;
	memset(&qc, 0, sizeof(qc));

	go_param_to_tpqctl(&qc,
			ctl_flags,
			ctl_deq_time,
			ctl_priority,
			ctl_diagnostic,
			ctl_diagmsg,
			ctl_msgid,
			ctl_corrid,
			ctl_replyqueue,
			ctl_failurequeue,
			ctl_cltid,
			ctl_urcode,
			ctl_appkey,
			ctl_delivery_qos,
			ctl_reply_qos,
			ctl_exp_time);

	ret = Otpenqueue (p_ctx, qspace, qname, &qc, data, len, flags);

	go_tpqctl_to_param(&qc,
			ctl_flags,
			ctl_deq_time,
			ctl_priority,
			ctl_diagnostic,
			ctl_diagmsg,
			ctl_msgid,
			ctl_corrid,
			ctl_replyqueue,
			ctl_failurequeue,
			ctl_cltid,
			ctl_urcode,
			ctl_appkey,
			ctl_delivery_qos,
			ctl_reply_qos,
			ctl_exp_time);

	return ret;
}

static int go_tpdequeue (TPCONTEXT_T *p_ctx,  char *qspace, char *qname, char **data, long *len, long flags,
		long *ctl_flags,
		long *ctl_deq_time,
		long *ctl_priority,
		long *ctl_diagnostic,
		char *ctl_diagmsg,
		char *ctl_msgid,
		char *ctl_corrid,
		char *ctl_replyqueue,
		char *ctl_failurequeue,
		char *ctl_cltid,
		long *ctl_urcode,
		long *ctl_appkey,
		long *ctl_delivery_qos,
		long *ctl_reply_qos,
		long *ctl_exp_time
)
{
	int ret;
	TPQCTL qc;
	memset(&qc, 0, sizeof(qc));

	go_param_to_tpqctl(&qc,
			ctl_flags,
			ctl_deq_time,
			ctl_priority,
			ctl_diagnostic,
			ctl_diagmsg,
			ctl_msgid,
			ctl_corrid,
			ctl_replyqueue,
			ctl_failurequeue,
			ctl_cltid,
			ctl_urcode,
			ctl_appkey,
			ctl_delivery_qos,
			ctl_reply_qos,
			ctl_exp_time);

	ret = Otpdequeue (p_ctx, qspace, qname, &qc, data, len, flags);

	go_tpqctl_to_param(&qc,
			ctl_flags,
			ctl_deq_time,
			ctl_priority,
			ctl_diagnostic,
			ctl_diagmsg,
			ctl_msgid,
			ctl_corrid,
			ctl_replyqueue,
			ctl_failurequeue,
			ctl_cltid,
			ctl_urcode,
			ctl_appkey,
			ctl_delivery_qos,
			ctl_reply_qos,
			ctl_exp_time);

	return ret;
}

//We need a tpfree version with NULL context
//So if we run in NULL context, then we must kill the new context appeared
//after the function call... (if any...)
//NOTE that tpfree will allocate auto-context if none currently present...
void go_tpfree(char *ptr)
{

    // Allocate new context + set it...
    TPCONTEXT_T c = tpnewctxt(0, 1);
	tpfree(ptr);
    tpfreectxt(c);

}

//Read the return code from current ATMI context
//@param p_ctx ATMI context
//@param c_err C error is set
//@return tpurcode (or 0 and error loaded)
static long go_tpurcode(TPCONTEXT_T *p_ctx, int *c_err)
{
	long ret=0;
	if (EXSUCCEED!=ndrx_tpsetctxt(*p_ctx, 0, CTXT_PRIV_ATMI|CTXT_PRIV_NSTD))
	{
		userlog("Failed to set ATMI context");
		*c_err=EXTRUE;
		goto out;
	}
	else
	{
		*c_err=EXFALSE;
	}

	ret=tpurcode;

	//Move process back to NULL context
	ndrx_tpsetctxt(TPNULLCONTEXT, 0L, CTXT_PRIV_ATMI|CTXT_PRIV_NSTD);

out:
	return ret;
}

*/
import "C"
import "unsafe"
import "fmt"
import "runtime"

/*
 * SUCCEED/FAIL flags
 */
const (
	SUCCEED = 0
	FAIL    = -1
)

/*
 * List of ATMI Error codes (atmi.h/xatmi.h)
 */
const (
	TPMINVAL      = 0
	TPEABORT      = 1
	TPEBADDESC    = 2
	TPEBLOCK      = 3
	TPEINVAL      = 4
	TPELIMIT      = 5
	TPENOENT      = 6
	TPEOS         = 7
	TPEPERM       = 8
	TPEPROTO      = 9
	TPESVCERR     = 10
	TPESVCFAIL    = 11
	TPESYSTEM     = 12
	TPETIME       = 13
	TPETRAN       = 14
	TPGOTSIG      = 15
	TPERMERR      = 16
	TPEITYPE      = 17
	TPEOTYPE      = 18
	TPERELEASE    = 19
	TPEHAZARD     = 20
	TPEHEURISTIC  = 21
	TPEEVENT      = 22
	TPEMATCH      = 23
	TPEDIAGNOSTIC = 24
	TPEMIB        = 25
	TPINITFAIL    = 30
	TPMAXVAL      = 31
)

/*
 * flag bits for C language xatmi routines
 */
const (
	TPNOBLOCK      = 0x00000001
	TPSIGRSTRT     = 0x00000002
	TPNOREPLY      = 0x00000004
	TPNOTRAN       = 0x00000008
	TPTRAN         = 0x00000010
	TPNOTIME       = 0x00000020
	TPGETANY       = 0x00000080
	TPNOCHANGE     = 0x00000100
	TPCONV         = 0x00000400
	TPSENDONLY     = 0x00000800
	TPRECVONLY     = 0x00001000
	TPTRANSUSPEND  = 0x00040000 /* Suspend current transaction          */
	TPSOFTTIMEOUT  = 0x00080000 /* Software time-out, translated to XATMI timeout for caller */
	TPSOFTNOENT    = 0x00100000 /* No service entry                     */
	TPNOAUTBUF     = 0x00200000 /* Don't restore autbuf in srv context  */
	TPREGEXMATCH   = 0x00800000 /* Use regular expressoins for match    */
	TPNOCACHELOOK  = 0x01000000 /* Do not lookup cache                  */
	TPNOCACHEADD   = 0x02000000 /* Do not save data to cache            */
	TPNOCACHEDDATA = 0x04000000 /* Do not use cached data               */
)

/*
 * values for rval in tpreturn
 */
const (
	TPFAIL    = 0x0001
	TPSUCCESS = 0x0002
)

/*
 * events returned during conversational communication
 */
const (
	TPEV_DISCONIMM = 0x0001
	TPEV_SVCERR    = 0x0002
	TPEV_SVCFAIL   = 0x0004
	TPEV_SVCSUCC   = 0x0008
	TPEV_SENDONLY  = 0x0020
)

/*
 * Max message size (int bytes)

const (
	ATMI_MSG_MAX_SIZE = 65536
)

Replaced with atmi.ATMIMsgSizeMax()!
*/

/*
 * TPQCTL.flags flags
 */
const (
	TPNOFLAGS         = 0x00000
	TPQCORRID         = 0x00001  /* set/get correlation id */
	TPQFAILUREQ       = 0x00002  /* set/get failure queue */
	TPQBEFOREMSGID    = 0x00004  /* RFU, enqueue before message id */
	TPQGETBYMSGIDOLD  = 0x00008  /* RFU, deprecated */
	TPQMSGID          = 0x00010  /* get msgid of enq/deq message */
	TPQPRIORITY       = 0x00020  /* set/get message priority */
	TPQTOP            = 0x00040  /* RFU, enqueue at queue top */
	TPQWAIT           = 0x00080  /* RFU, wait for dequeuing */
	TPQREPLYQ         = 0x00100  /* set/get reply queue */
	TPQTIME_ABS       = 0x00200  /* RFU, set absolute time */
	TPQTIME_REL       = 0x00400  /* RFU, set absolute time */
	TPQGETBYCORRIDOLD = 0x00800  /* deprecated */
	TPQPEEK           = 0x01000  /* peek */
	TPQDELIVERYQOS    = 0x02000  /* RFU, delivery quality of service */
	TPQREPLYQOS       = 0x04000  /* RFU, reply message quality of service */
	TPQEXPTIME_ABS    = 0x08000  /* RFU, absolute expiration time */
	TPQEXPTIME_REL    = 0x10000  /* RFU, relative expiration time */
	TPQEXPTIME_NONE   = 0x20000  /* RFU, never expire */
	TPQGETBYMSGID     = 0x40008  /* dequeue by msgid */
	TPQGETBYCORRID    = 0x80800  /* dequeue by corrid */
	TPQASYNC          = 0x100000 /* Async complete */
)

/*
 * Values for TQPCTL.diagnostic
 */
const (
	QMEINVAL     = -1
	QMEBADRMID   = -2
	QMENOTOPEN   = -3
	QMETRAN      = -4
	QMEBADMSGID  = -5
	QMESYSTEM    = -6
	QMEOS        = -7
	QMEABORTED   = -8
	QMENOTA      = -8 /* QMEABORTED */
	QMEPROTO     = -9
	QMEBADQUEUE  = -10
	QMENOMSG     = -11
	QMEINUSE     = -12
	QMENOSPACE   = -13
	QMERELEASE   = -14
	QMEINVHANDLE = -15
	QMESHARE     = -16
)

/*
 * Q constants
 */
const (
	TMMSGIDLEN       = 32
	TMCORRIDLEN      = 32
	TMQNAMELEN       = 15
	NDRX_MAX_ID_SIZE = 96
)

/*
 * Log levels for TPLOG (corresponding to ndebug.h)
 */
const (
	LOG_ALWAYS = 1
	LOG_ERROR  = 2
	LOG_WARN   = 3
	LOG_INFO   = 4
	LOG_DEBUG  = 5
	LOG_DUMP   = 6
)

/*
 * Logging facilites
 */
const (
	LOG_FACILITY_NDRX       = 0x00001 /* settings for ATMI logging             */
	LOG_FACILITY_UBF        = 0x00002 /* settings for UBF logging              */
	LOG_FACILITY_TP         = 0x00004 /* settings for TP logging               */
	LOG_FACILITY_TP_THREAD  = 0x00008 /* settings for TP, thread based logging */
	LOG_FACILITY_TP_REQUEST = 0x00010 /* Request logging, thread based         */
)

/*
 * Enduro/X standard library error codes
 */
const (
	NEINVALINI  = 1  /* Invalid INI file */
	NEMALLOC    = 2  /* Malloc failed */
	NEUNIX      = 3  /* Unix error occurred */
	NEINVAL     = 4  /* Invalid value passed to function */
	NESYSTEM    = 5  /* System failure */
	NEMANDATORY = 6  /* Mandatory field is missing */
	NEFORMAT    = 7  /* Format error */
	NETOUT      = 8  /* Time-out condition */
	NENOCONN    = 9  /* No connection */
	NELIMIT     = 10 /* Limit reached */
)

/**
 * Enduro/X extensions
 */
const (
	TPEX_NOCHANGE = 0x00000004 /**< Reject tpimport with error if 	*/
	TPEX_STRING   = 0x00000008 /**< Export buffer in base64 format 		*/
)

/**
 * Multi contexting defines
 */
const (
	TPNULLCONTEXT = 0 /**< NULL Context */
)

/*
 * Transaction ID type
 */
type TPTRANID struct {
	c_tptranid C.TPTRANID
}

/*
 * ATMI Context object
 */
type ATMICtx struct {
	gcoff int //dummy counter tricking the gc to suspend while using object in c
	c_ctx C.TPCONTEXT_T
}

/*
 * Server context data (used for server's main thread
 * switching taks to worker thread)
 */
type TPSRVCTXDATA struct {
	c_ptr *C.char
}

/*
 * Event controll struct
 */
type TPEVCTL struct {
	flags int64
	name1 string
	name2 string
}

/*
 * Queue control structure
 */
type TPQCTL struct {
	flags        int64             /* indicates which of the values are set */
	deq_time     int64             /* absolute/relative  time for dequeuing */
	priority     int64             /* enqueue priority */
	diagnostic   int64             /* indicates reason for failure */
	diagmsg      string            /* diagnostic message */
	msgid        [TMMSGIDLEN]byte  /* id of message before which to queue */
	corrid       [TMCORRIDLEN]byte /* correlation id used to identify message */
	replyqueue   string            /* queue name for reply message */
	failurequeue string            /* queue name for failure message */
	cltid        string            /* client identifier for originating client */
	urcode       int64             /* application user-return code */
	appkey       int64             /* application authentication client key */
	delivery_qos int64             /* delivery quality of service  */
	reply_qos    int64             /* reply message quality of service  */
	exp_time     int64             /* expiration time  */
}

///////////////////////////////////////////////////////////////////////////////////
// ATMI Buffers section
///////////////////////////////////////////////////////////////////////////////////

//ATMI buffer
type ATMIBuf struct {
	gcoff int

	C_ptr *C.char
	//We will need some API for length & buffer setting
	//Probably we need a wrapper for lenght function
	C_len C.long

	//have finalizer installed
	HaveFinalizer bool

	//Have some context, just a reference to, for ATMI buffer operations
	Ctx *ATMICtx
}

//Base interface for typed buffer
type TypedBuffer interface {
	GetBuf() *ATMIBuf
}

//Have inteface to base ATMI buffer
func (u *ATMIBuf) GetBuf() *ATMIBuf {
	return u
}

//Do nothing, to trick the GC
func (u *ATMIBuf) nop() int {
	u.gcoff++
	return u.gcoff
}

//Max message size
//@return buffer size configured by Enduro/X, min 64K
func ATMIMsgSizeMax() int64 {
	return int64(C.ndrx_msgsizemax())
}

///////////////////////////////////////////////////////////////////////////////////
// Error Handlers, ATMI level
///////////////////////////////////////////////////////////////////////////////////

//ATMI Error type
type atmiError struct {
	code    int
	message string
}

//ATMI error interface
type ATMIError interface {
	Error() string
	Code() int
	Message() string
}

//Do nothing, to trick the GC
func (ac *ATMICtx) nop() int {
	ac.gcoff++
	return ac.gcoff
}

//Generate ATMI error, read the codes
func (ac *ATMICtx) NewATMIError() ATMIError {
	var err atmiError
	err.code = int(C.go_tperrno(&ac.c_ctx))
	err.message = C.GoString(C.Otpstrerror(&ac.c_ctx, C.go_tperrno(&ac.c_ctx)))
	return err
}

//Build a custom error
//@param err		Error buffer to build
//@param code	Error code to setup
//@param msg		Error message
func NewCustomATMIError(code int, msg string) ATMIError {
	var err atmiError
	err.code = code
	err.message = msg
	return err
}

//Standard error interface
func (e atmiError) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

//code getter
func (e atmiError) Code() int {
	return e.code
}

//message getter
func (e atmiError) Message() string {
	return e.message
}

///////////////////////////////////////////////////////////////////////////////////
// Error Handlers, NSTD - Enduro/X Standard library
///////////////////////////////////////////////////////////////////////////////////

//NSTD Error type
type nstdError struct {
	code    int
	message string
}

//NSTD error interface
type NSTDError interface {
	Error() string
	Code() int
	Message() string
}

//Generate NSTD error, read the codes
func (ac *ATMICtx) NewNstdError() NSTDError {
	var err nstdError
	err.code = int(C.go_Nerror(&ac.c_ctx))
	err.message = C.GoString(C.ONstrerror(&ac.c_ctx, C.go_Nerror(&ac.c_ctx)))
	return err
}

//Build a custom error. Can be used at Go level sources
//To simulate standard error
//@param err		Error buffer to build
//@param code	Error code to setup
//@param msg		Error message
func NewCustomNstdError(code int, msg string) NSTDError {
	var err nstdError
	err.code = code
	err.message = msg
	return err
}

//Standard error interface
func (e nstdError) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

//Error code getter
func (e nstdError) Code() int {
	return e.code
}

//Error message getter
func (e nstdError) Message() string {
	return e.message
}

///////////////////////////////////////////////////////////////////////////////////
// API Section
// TODO: Think about persistent association with thread. So that
//       in XA case it would be simpler to manipulate with DB + XATMI...
///////////////////////////////////////////////////////////////////////////////////

//Allocate new ATMI context. This is the context with most of the XATMI operations
//are made. Single go routine can have multiple contexts at the same time.
//The function does not open queues or init XATMI sub-system unless the dependant
//operation is called. For example you may allocat the context and use it for logging
//that will not make overhead for system queues.
//@return ATMI Error, Pointer to ATMI Context object
func NewATMICtx() (*ATMICtx, ATMIError) {
	var ret ATMICtx
	ret.c_ctx = C.tpnewctxt(0, 0)
	if nil == ret.c_ctx {
		return nil, NewCustomATMIError(TPESYSTEM, "Failed to allocate "+
			"new context - see ULOG for details")
	}

	runtime.SetFinalizer(&ret, freeATMICtx)

	return &ret, nil
}

//Free up the ATMI Context
//Internally this will call the TpTerm too to termiante any XATMI client
//session in progress.
func (ac *ATMICtx) FreeATMICtx() {
	ac.TpTerm() //This extra, but let it be
	C.Otpfreectxt(&ac.c_ctx, ac.c_ctx)
}

//Associate current OS thread with context
//This might be needed for global transaction processing
//Which uses underlaying OS threads for transaction association
func (ac *ATMICtx) AssocThreadWithCtx() ATMIError {

	if ret := C.tpsetctxt(ac.c_ctx, 0); SUCCEED != ret {
		return ac.NewATMIError()
	}

	return nil
}

//Disassocate current os thread from context
//This might be needed for global transaction processing
//Which uses underlaying OS threads for transaction association
func (ac *ATMICtx) DisassocThreadFromCtx() ATMIError {

	if ret := C.tpgetctxt(&ac.c_ctx, 0); SUCCEED != ret {
		return ac.NewATMIError()
	}
	return nil
}

//Kill the ATMI context (internal version for finalizer)
func freeATMICtx(ac *ATMICtx) {
	if nil != ac.c_ctx {
		//ac.TpTerm() //This extra, but let it be - not needed, free will do.
		C.Otpfreectxt(&ac.c_ctx, ac.c_ctx)
	}
}

//Make context object from C pointer. Function can be used in case
//If doing any direct XATMI operations and you have a C context handler.
//Which can be promoted to Go level ATMI Context.
//@param c_ctx Context ATMI object
//@return ATMI Context Object
func MakeATMICtx(c_ctx C.TPCONTEXT_T) *ATMICtx {
	var ret ATMICtx
	ret.c_ctx = c_ctx
	return &ret
}

//TODO, maybe we need to use error deligates, so that user can override the error handling object?

//Allocate buffer
//Accepts the standard ATMI values
//We should add error handling here
//@param	 b_type 		Buffer type
//@param	 b_subtype 	Buffer sub-type
//@param	 size		Buffer size request
//@return 			ATMI Buffer, atmiError
func (ac *ATMICtx) TpAlloc(b_type string, b_subtype string, size int64) (*ATMIBuf, ATMIError) {

	var buf ATMIBuf
	var err ATMIError

	c_type := C.CString(b_type)
	c_subtype := C.CString(b_subtype)

	size_l := C.long(size)

	buf.C_ptr = C.Otpalloc(&ac.c_ctx, c_type, c_subtype, size_l)

	//Check the error
	if nil == buf.C_ptr {
		err = ac.NewATMIError()
	}

	C.free(unsafe.Pointer(c_type))
	C.free(unsafe.Pointer(c_subtype))

	runtime.SetFinalizer(&buf, tpfree)
	buf.HaveFinalizer = true

	ac.nop() //keep context until the end of the func, and only then allow gc
	return &buf, err
}

//Change the context of the buffers (needed for error handling)
func (buf *ATMIBuf) TpSetCtxt(ac *ATMICtx) {
	buf.Ctx = ac
}

//Reallocate the buffer
//@param buf		ATMI buffer
//@return 		ATMI Error
func (buf *ATMIBuf) TpRealloc(size int64) ATMIError {

	var err ATMIError

	buf.C_ptr = C.Otprealloc(&buf.Ctx.c_ctx, buf.C_ptr, C.long(size))

	if nil == buf.C_ptr {
		err = buf.Ctx.NewATMIError()
	}

	buf.nop()

	return err
}

//Initialize client
//@return		ATMI Error
func (ac *ATMICtx) TpInit() ATMIError {

	var err ATMIError

	if SUCCEED != C.go_tpinit(&ac.c_ctx) {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

// Do the service call, assume using the same buffer
// for return value.
// This works for self describing buffers. Otherwise we need a buffer size in
// ATMIBuf.
// @param svc	service name
// @param buf	ATMI buffer
// @param flags 	Flags to be used
// @return atmiError
func (ac *ATMICtx) TpCall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {

	var err ATMIError
	c_svc := C.CString(svc)

	buf := tb.GetBuf()

	ret := C.Otpcall(&ac.c_ctx, c_svc, buf.C_ptr, buf.C_len, &buf.C_ptr, &buf.C_len, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	//Check the types
	switch tb.(type) {
	case *TypedVIEW:
		if v, ok := tb.(*TypedVIEW); ok {
			itype := ""
			_, errA := ac.TpTypes(buf, &itype, &v.view)

			if nil != err && nil != errA {
				err = errA
			}
		}
		break
	}

	C.free(unsafe.Pointer(c_svc))

	ac.nop() //keep context until the end of the func, and only then allow gc
	return int(ret), err
}

//TP Async call
//@param svc		Service Name to call
//@param buf		ATMI buffer
//@param flags	Flags to be used for call (see flags section)
//@return		Call Descriptor (cd), ATMI Error
func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	var err ATMIError
	c_svc := C.CString(svc)

	buf := tb.GetBuf()

	ret := C.Otpacall(&ac.c_ctx, c_svc, buf.C_ptr, buf.C_len, C.long(flags))

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	C.free(unsafe.Pointer(c_svc))

	buf.nop() //keep context until the end of the func, and only then allow gc
	ac.nop()  //keep context until the end of the func, and only then allow gc
	return int(ret), err
}

//Get async call reply
//@param cd	call
//@param buf	ATMI buffer
//@param flags call flags
func (ac *ATMICtx) TpGetRply(cd *int, tb TypedBuffer, flags int64) (int, ATMIError) {

	var err ATMIError
	var c_cd C.int

	buf := tb.GetBuf()

	ret := C.Otpgetrply(&ac.c_ctx, &c_cd, &buf.C_ptr, &buf.C_len, C.long(flags))
	*cd = int(c_cd)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop()  //keep context until the end of the func, and only then allow gc
	buf.nop() //keep context until the end of the func, and only then allow gc

	return int(ret), err
}

//Cancel async call
//@param cd		Call descriptor
//@return ATMI error
func (ac *ATMICtx) TpCancel(cd int) ATMIError {

	var err ATMIError

	ret := C.Otpcancel(&ac.c_ctx, C.int(cd))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Connect to service in conversational mode
//@param svc		Service name
//@param data	ATMI buffers
//@param flags	Flags
//@return		call descriptor (cd), ATMI error
func (ac *ATMICtx) TpConnect(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {

	var err ATMIError
	c_svc := C.CString(svc)

	data := tb.GetBuf()

	ret := C.Otpconnect(&ac.c_ctx, c_svc, data.C_ptr, data.C_len, C.long(flags))

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	C.free(unsafe.Pointer(c_svc))

	data.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc

	return int(ret), err
}

//Disconnect from conversation
//@param cd		Call Descriptor
//@return ATMI Error
func (ac *ATMICtx) TpDiscon(cd int) ATMIError {

	var err ATMIError

	ret := C.Otpdiscon(&ac.c_ctx, C.int(cd))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Receive data from conversation
//@param cd			call descriptor
//@param	 data		ATMI buffer
//@param revent		Return Event
//@return			ATMI Error
func (ac *ATMICtx) TpRecv(cd int, tb TypedBuffer, flags int64, revent *int64) ATMIError {

	var err ATMIError

	c_revent := C.long(*revent)

	data := tb.GetBuf()

	ret := C.Otprecv(&ac.c_ctx, C.int(cd), &data.C_ptr, &data.C_len, C.long(flags), &c_revent)

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	*revent = int64(c_revent)

	data.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

//Receive data from conversation
//@param cd			call descriptor
//@param	 data		ATMI buffer
//@param revent		Return Event
//@return			ATMI Error
func (ac *ATMICtx) TpSend(cd int, tb TypedBuffer, flags int64, revent *int64) ATMIError {

	var err ATMIError

	c_revent := C.long(*revent)

	data := tb.GetBuf()

	ret := C.Otpsend(&ac.c_ctx, C.int(cd), data.C_ptr, data.C_len, C.long(flags), &c_revent)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	*revent = int64(c_revent)

	data.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

//Free the ATMI buffer
//@param buf		ATMI buffer
func (ac *ATMICtx) TpFree(buf *ATMIBuf) {

	C.Otpfree(&ac.c_ctx, buf.C_ptr)
	buf.C_ptr = nil

	ac.nop() //keep context until the end of the func, and only then allow gc

}

//Free the ATMI buffer (internal version, for finalizer)
//Context less operation
//@param buf		ATMI buffer
func tpfree(buf *ATMIBuf) {
	//Kill any context is appeared.
	//Protect us from garbadge collector
	if buf.C_ptr != nil {
		C.go_tpfree(buf.C_ptr)
		buf.C_ptr = nil
	}
}

//Commit global transaction
//@param	 flags		flags for abort operation
func (ac *ATMICtx) TpCommit(flags int64) ATMIError {

	var err ATMIError

	ret := C.Otpcommit(&ac.c_ctx, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Abort global transaction
//@param	 flags		flags for abort operation (must be 0)
//@return ATMI Error
func (ac *ATMICtx) TpAbort(flags int64) ATMIError {

	var err ATMIError

	ret := C.Otpabort(&ac.c_ctx, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Open XA Sub-system
//@return ATMI Error
func (ac *ATMICtx) TpOpen() ATMIError {

	var err ATMIError

	ret := C.Otpopen(&ac.c_ctx)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

// Close XA Sub-system
//@return ATMI Error
func (ac *ATMICtx) TpClose() ATMIError {

	var err ATMIError

	ret := C.Otpclose(&ac.c_ctx)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

//Check are we in globa transaction?
//@return 	0 - not in global Tx, 1 - in global Tx
func (ac *ATMICtx) TpGetLev() int {

	ret := C.Otpgetlev(&ac.c_ctx)

	ac.nop() //keep context until the end of the func, and only then allow gc
	return int(ret)
}

//Begin transaction
//@param timeout		Transaction Timeout
//@param flags		Transaction flags
//@return	ATMI Error
func (ac *ATMICtx) TpBegin(timeout uint64, flags int64) ATMIError {

	var err ATMIError

	ret := C.Otpbegin(&ac.c_ctx, C.ulong(timeout), C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Suspend transaction
//@param tranid	Transaction Id reference
//@param flags	Flags for suspend (must be 0)
//@return 	ATMI Error
func (ac *ATMICtx) TpSuspend(tranid *TPTRANID, flags int64) ATMIError {

	var err ATMIError

	ret := C.Otpsuspend(&ac.c_ctx, &tranid.c_tptranid, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

//Resume transaction
//@param tranid	Transaction Id reference
//@param flags	Flags for tran resume (must be 0)
//@return 	ATMI Error
func (ac *ATMICtx) TpResume(tranid *TPTRANID, flags int64) ATMIError {

	var err ATMIError

	ret := C.Otpresume(&ac.c_ctx, &tranid.c_tptranid, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return err
}

//Get cluster node id
//@return		Node Id
func (ac *ATMICtx) TpGetnodeId() int64 {

	ret := C.Otpgetnodeid(&ac.c_ctx)

	ac.nop() //keep context until the end of the func, and only then allow gc
	return int64(ret)
}

//Post the event to subscribers
//@param eventname	Name of the event to post
//@param data		ATMI buffer
//@param flags		flags
//@return		Number Of events posted, ATMI error
func (ac *ATMICtx) TpPost(eventname string, tb TypedBuffer, len int64, flags int64) (int, ATMIError) {
	var err ATMIError
	c_eventname := C.CString(eventname)

	data := tb.GetBuf()

	ret := C.Otppost(&ac.c_ctx, c_eventname, data.C_ptr, data.C_len, C.long(flags))

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	C.free(unsafe.Pointer(c_eventname))

	data.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc
	return int(ret), err
}

//Return ATMI buffer info
//@param ptr 	Pointer to ATMI buffer
//@param itype	ptr to string to return the buffer type  (can be nil), if set
//then on output value will be UBF, CARRAY, STRING or JSON other buffers currently
//are not supported.
//@param subtype ptr to string to return sub-type (can be nil)
//@return	Buffer lenght if no error or -1 if error, ATMI error
func (ac *ATMICtx) TpTypes(ptr *ATMIBuf, itype *string, subtype *string) (int64, ATMIError) {

	var err ATMIError

	/* we should allocat the fields there...  */

	var c_type *C.char
	var c_subtype *C.char

	c_type = C.malloc_string(16)
	c_subtype = C.malloc_string(16)

	ret := C.Otptypes(&ac.c_ctx, ptr.C_ptr, c_type, c_subtype)

	if FAIL == ret {
		err = ac.NewATMIError()
	} else {
		if nil != itype && nil != c_type {
			*itype = C.GoString(c_type)
		}

		if nil != subtype && nil != c_subtype {
			*subtype = C.GoString(c_subtype)
		}
	}

	if nil != c_type {
		C.free_string(c_type)
	}

	if nil != c_subtype {
		C.free_string(c_subtype)
	}

	ptr.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc
	return int64(ret), err
}

//Return ATMI buffer info
//@param itype	ptr to string to return the buffer type  (can be nil), if set
//then on output value will be UBF, CARRAY, STRING or JSON other buffers currently
//are not supported.
//@param subtype ptr to string to return sub-type (can be nil)
//@return	Buffer lenght if no error or -1 if error, ATMI error
func (ptr *ATMIBuf) TpTypes(itype *string, subtype *string) (int64, ATMIError) {
	return ptr.Ctx.TpTypes(ptr, itype, subtype)
}

//Terminate the client
//@return ATMI error
func (ac *ATMICtx) TpTerm() ATMIError {

	ret := C.Otpterm(&ac.c_ctx)
	if SUCCEED != ret {
		return ac.NewATMIError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return nil
}

//Glue function for tpenqueue and tpdequeue
//@param qspace	Name of the event to post
//@param qname		ATMI buffer
//@param ctl		Control structure
//@param tb		Typed buffer
//@param flags		ATMI call flags
//@param is_enq		Is Enqueue? If not then dequeue
//@return		ATMI error
func (ac *ATMICtx) tp_enq_deq(qspace string, qname string, ctl *TPQCTL, tb TypedBuffer, flags int64, is_enq bool) ATMIError {

	var err ATMIError

	c_qspace := C.CString(qspace)
	defer C.free(unsafe.Pointer(c_qspace))

	c_qname := C.CString(qname)
	defer C.free(unsafe.Pointer(c_qname))

	c_ctl_flags := C.long(ctl.flags)
	c_ctl_deq_time := C.long(ctl.deq_time)
	c_ctl_priority := C.long(ctl.priority)
	c_ctl_diagnostic := C.long(ctl.diagnostic)
	c_ctl_diagmsg := C.calloc(1, 256)
	c_ctl_diagmsg_ptr := (*C.char)(unsafe.Pointer(c_ctl_diagmsg))
	defer C.free(unsafe.Pointer(c_ctl_diagmsg))

	c_ctl_msgid := C.malloc(TMMSGIDLEN)
	c_ctl_msgid_ptr := (*C.char)(unsafe.Pointer(c_ctl_msgid))
	defer C.free(unsafe.Pointer(c_ctl_msgid))
	for i := 0; i < TMMSGIDLEN; i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_ctl_msgid) + uintptr(i))) = C.char(ctl.msgid[i])
	}

	c_ctl_corrid := C.malloc(TMCORRIDLEN)
	c_ctl_corrid_ptr := (*C.char)(unsafe.Pointer(c_ctl_corrid))
	defer C.free(unsafe.Pointer(c_ctl_corrid))
	for i := 0; i < TMCORRIDLEN; i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_ctl_corrid) + uintptr(i))) = C.char(ctl.corrid[i])
	}

	/* Allocate the buffer for reply q, because we might receive this on
	   dequeue.
	*/
	c_ctl_replyqueue_tmp := C.CString(ctl.replyqueue)
	defer C.free(unsafe.Pointer(c_ctl_replyqueue_tmp))
	c_ctl_replyqueue := C.malloc(TMQNAMELEN + 1)
	c_ctl_replyqueue_ptr := (*C.char)(unsafe.Pointer(c_ctl_corrid))
	defer C.free(unsafe.Pointer(c_ctl_replyqueue))

	if C.strlen(c_ctl_replyqueue_tmp) > TMQNAMELEN {
		return NewCustomATMIError(TPEINVAL,
			fmt.Sprintf("Invalid reply queue len, max: %d", TMQNAMELEN))
	}
	C.strcpy(c_ctl_replyqueue_ptr, c_ctl_replyqueue_tmp)

	/* Allocate the buffer for failure q, because we might receive this on
	   dequeue.
	*/
	c_ctl_failurequeue_tmp := C.CString(ctl.failurequeue)
	defer C.free(unsafe.Pointer(c_ctl_failurequeue_tmp))
	c_ctl_failurequeue := C.malloc(TMQNAMELEN + 1)
	c_ctl_failurequeue_ptr := (*C.char)(unsafe.Pointer(c_ctl_corrid))
	defer C.free(unsafe.Pointer(c_ctl_failurequeue))

	if C.strlen(c_ctl_failurequeue_tmp) > TMQNAMELEN {
		return NewCustomATMIError(TPEINVAL,
			fmt.Sprintf("Invalid failure queue len, max: %d", TMQNAMELEN))
	}
	C.strcpy(c_ctl_failurequeue_ptr, c_ctl_failurequeue_tmp)

	/* The same goes with client id... we might return it on dequeue */
	c_ctl_cltid_tmp := C.CString(ctl.cltid)
	defer C.free(unsafe.Pointer(c_ctl_cltid_tmp))
	c_ctl_cltid := C.malloc(TMQNAMELEN + 1)
	c_ctl_cltid_ptr := (*C.char)(unsafe.Pointer(c_ctl_corrid))
	defer C.free(unsafe.Pointer(c_ctl_cltid))

	if C.strlen(c_ctl_cltid_tmp) > NDRX_MAX_ID_SIZE {
		return NewCustomATMIError(TPEINVAL,
			fmt.Sprintf("Invalid client id len, max: %d", NDRX_MAX_ID_SIZE))
	}
	C.strcpy(c_ctl_cltid_ptr, c_ctl_cltid_tmp)

	c_ctl_urcode := C.long(ctl.urcode)
	c_ctl_appkey := C.long(ctl.appkey)
	c_ctl_delivery_qos := C.long(ctl.delivery_qos)
	c_ctl_reply_qos := C.long(ctl.reply_qos)
	c_ctl_exp_time := C.long(ctl.exp_time)

	buf := tb.GetBuf()

	var ret C.int
	if is_enq {
		ret = C.go_tpenqueue(&ac.c_ctx, c_qspace, c_qname, buf.C_ptr, buf.C_len, C.long(flags),
			&c_ctl_flags,
			&c_ctl_deq_time,
			&c_ctl_priority,
			&c_ctl_diagnostic,
			c_ctl_diagmsg_ptr,
			c_ctl_msgid_ptr,
			c_ctl_corrid_ptr,
			c_ctl_replyqueue_ptr,
			c_ctl_failurequeue_ptr,
			c_ctl_cltid_ptr,
			&c_ctl_urcode,
			&c_ctl_appkey,
			&c_ctl_delivery_qos,
			&c_ctl_reply_qos,
			&c_ctl_exp_time)
	} else {
		ret = C.go_tpdequeue(&ac.c_ctx, c_qspace, c_qname, &buf.C_ptr, &buf.C_len, C.long(flags),
			&c_ctl_flags,
			&c_ctl_deq_time,
			&c_ctl_priority,
			&c_ctl_diagnostic,
			c_ctl_diagmsg_ptr,
			c_ctl_msgid_ptr,
			c_ctl_corrid_ptr,
			c_ctl_replyqueue_ptr,
			c_ctl_failurequeue_ptr,
			c_ctl_cltid_ptr,
			&c_ctl_urcode,
			&c_ctl_appkey,
			&c_ctl_delivery_qos,
			&c_ctl_reply_qos,
			&c_ctl_exp_time)
	}

	/* transfer back to structure values we got... */
	ctl.flags = int64(c_ctl_flags)
	ctl.deq_time = int64(c_ctl_deq_time)
	ctl.priority = int64(c_ctl_priority)
	ctl.diagnostic = int64(c_ctl_diagnostic)

	ctl.diagmsg = C.GoString(c_ctl_diagmsg_ptr)

	for i := 0; i < TMMSGIDLEN; i++ {
		ctl.msgid[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_ctl_msgid) + uintptr(i))))
	}

	for i := 0; i < TMCORRIDLEN; i++ {
		ctl.corrid[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_ctl_corrid) + uintptr(i))))
	}

	ctl.replyqueue = C.GoString(c_ctl_replyqueue_ptr)
	ctl.failurequeue = C.GoString(c_ctl_failurequeue_ptr)
	ctl.cltid = C.GoString(c_ctl_cltid_ptr)

	ctl.urcode = int64(c_ctl_urcode)
	ctl.appkey = int64(c_ctl_appkey)
	ctl.delivery_qos = int64(c_ctl_delivery_qos)
	ctl.reply_qos = int64(c_ctl_reply_qos)
	ctl.exp_time = int64(c_ctl_exp_time)

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	buf.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc

	return err
}

//Enqueue message to Q
//@param qspace	Name of the event to post
//@param qname		ATMI buffer
//@param ctl		Control structure
//@param tb		Typed buffer
//@param flags		ATMI call flags
//@return		ATMI error
func (ac *ATMICtx) TpEnqueue(qspace string, qname string, ctl *TPQCTL, tb TypedBuffer, flags int64) ATMIError {
	return ac.tp_enq_deq(qspace, qname, ctl, tb, flags, true)
}

//Dequeue message from Q
//@param qspace	Name of the event to post
//@param qname		ATMI buffer
//@param ctl		Control structure
//@param tb		Typed buffer
//@param flags		ATMI call flags
//@return		ATMI error
func (ac *ATMICtx) TpDequeue(qspace string, qname string, ctl *TPQCTL, tb TypedBuffer, flags int64) ATMIError {
	return ac.tp_enq_deq(qspace, qname, ctl, tb, flags, false)
}

//Have some common testing functinos (used by Enduro/X modules)
//@param a Paramter a to test
//@param b Paramter b to test (compare with a)
//@param message additional error message
func (ac *ATMICtx) TpAssertEqualPanic(a interface{}, b interface{}, message string) {
	aa := fmt.Sprintf("%v", a)
	bb := fmt.Sprintf("%v", b)

	if aa == bb {
		return
	}
	msg2 := fmt.Sprintf("%v != %v", a, b)

	panic(fmt.Sprintf("TpPanicAssertEqual: %s: %s", message, msg2))

}

//Have some common testing functinos (used by Enduro/X modules)
//@param a Paramter a to test
//@param b Paramter b to test (compare with a)
//@param message additional error message
func (ac *ATMICtx) TpAssertNotEqualPanic(a interface{}, b interface{}, message string) {
	aa := fmt.Sprintf("%v", a)
	bb := fmt.Sprintf("%v", b)

	if aa != bb {
		return
	}
	msg2 := fmt.Sprintf("%v != %v", a, b)

	panic(fmt.Sprintf("TpPanicAssertNotEqual: %s: %s", message, msg2))

}

//Have some common testing functinos (used by Enduro/X modules)
//@param a Paramter a to test
//@param b Paramter b to test (compare with a)
//@param message additional error message
func (ac *ATMICtx) TpAssertEqualError(a interface{}, b interface{}, message string) error {
	aa := fmt.Sprintf("%v", a)
	bb := fmt.Sprintf("%v", b)

	if aa == bb {
		return nil
	}
	msg2 := fmt.Sprintf("%v != %v", a, b)

	return fmt.Errorf("TpPanicAssertEqual: %s: %s", message, msg2)

}

//Have some common testing functinos (used by Enduro/X modules)
//@param a Paramter a to test
//@param b Paramter b to test (compare with a)
//@param message additional error message
func (ac *ATMICtx) TpAssertNotEqualError(a interface{}, b interface{}, message string) error {
	aa := fmt.Sprintf("%v", a)
	bb := fmt.Sprintf("%v", b)

	if aa != bb {
		return nil
	}
	msg2 := fmt.Sprintf("%v != %v", a, b)

	return fmt.Errorf("TpAssertNotEqualError: %s: %s", message, msg2)

}

//Export the buffer to JSON format.
//@param tb TypecdTyped buffer
//@param flags 0 (JSON text) or TPEX_STRING (export in base64 format)
func (ac *ATMICtx) TpExport(tb TypedBuffer, flags int64) (string, ATMIError) {

	var err ATMIError
	buf := tb.GetBuf()

	c_str_buf := C.malloc(C.size_t(ATMIMsgSizeMax() * 2))
	c_str_buf_ptr := (*C.char)(unsafe.Pointer(c_str_buf))
	defer C.free(unsafe.Pointer(c_str_buf))
	c_len := C.long(ATMIMsgSizeMax() * 2)

	//WELL! go might be too agressive, if it loads in the function call
	//the buf object and for the caller tb/buf is last used in the block
	//the GC might start to kill the buf, thus meaning it will corrupt the
	//UBF buffer object while C is processing !!!!!!!!
	//Thus needs to lock buf somehow while we are in the C
	if ret := C.Otpexport(&ac.c_ctx, buf.C_ptr, buf.C_len, c_str_buf_ptr, &c_len,
		C.long(flags)); ret != SUCCEED {
		err = ac.NewATMIError()
	}

	if nil != err {
		return "", err
	}

	//Have buffer usage after C, avoid GC during the C call, if this is last
	//buffer use
	buf.nop()

	ac.nop() //keep context until the end of the func, and only then allow gc
	return C.GoString(c_str_buf_ptr), nil

}

//Import the UBF buffer from the json string which is exported by TpExport
//The tb TypedBuffer will be updated according to incoming data
//@param jsondata json string data according to texport(3)
//@param tb typed buffer where to install the incoming buffer
//@param flags TPEX_STRING if decode as base64, TPEX_NOCHANGE do not change tb format
//	if buffer type is different
func (ac *ATMICtx) TpImport(jsondata string, tb TypedBuffer, flags int64) ATMIError {
	var err ATMIError
	buf := tb.GetBuf()

	c_jsondata := C.CString(jsondata)
	defer C.free(unsafe.Pointer(c_jsondata))

	if ret := C.Otpimport(&ac.c_ctx, c_jsondata, C.long(0), &buf.C_ptr, &buf.C_len,
		C.long(flags)); ret != SUCCEED {
		err = ac.NewATMIError()
	}

	buf.nop()
	ac.nop() //keep context until the end of the func, and only then allow gc
	//Have buffer usage after C, avoid GC during the C call

	return err
}

//Return size of long. Used to detect running platform
//Basically all int64 APIs used by ATMI package will
//resolve from Enduro/X point of view this count of bytes.
//@return number of bytes in "long" data type
func ExSizeOfLong() int {

	return C.EX_SIZEOF_LONG
}

//Return user return code from last service call
//(i.e. value from second argument of tpreturn - rcode)
//@return tpurcode from ATMI context
func (ac *ATMICtx) TpURCode() (int64, ATMIError) {

	var c_err C.int
	ret := int64(C.go_tpurcode(&ac.c_ctx, &c_err))

	//mark the scope of ac
	ac.nop()

	if 1 == c_err {
		return 0, ac.NewATMIError()
	}

	return ret, nil
}

/* vim: set ts=4 sw=4 et smartindent: */
