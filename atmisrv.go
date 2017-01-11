package atmi

/*
** XATMI main package, server API
**
** @file atmisrv.go
**
** -----------------------------------------------------------------------------
** Enduro/X Middleware Platform for Distributed Transaction Processing
** Copyright (C) 2015, ATR Baltic, SIA. All Rights Reserved.
** This software is released under one of the following licenses:
** GPL or ATR Baltic's license for commercial use.
** -----------------------------------------------------------------------------
** GPL license:
**
** This program is free software; you can redistribute it and/or modify it under
** the terms of the GNU General Public License as published by the Free Software
** Foundation; either version 2 of the License, or (at your option) any later
** version.
**
** This program is distributed in the hope that it will be useful, but WITHOUT ANY
** WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
** PARTICULAR PURPOSE. See the GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License along with
** this program; if not, write to the Free Software Foundation, Inc., 59 Temple
** Place, Suite 330, Boston, MA 02111-1307 USA
**
** -----------------------------------------------------------------------------
** A commercial use license is available from ATR Baltic, SIA
** contact@atrbaltic.com
** -----------------------------------------------------------------------------
 */

/*
#cgo pkg-config: atmisrvinteg

#include <string.h>
#include <stdlib.h>
#include <ndebug.h>
#include <oatmi.h>
#include <oatmisrv.h>

static void free_string(char* s) { free(s); }
static char * malloc_string(int size) { return malloc(size); }

extern int go_tpsrvinit();
extern void go_tpsrvdone();
extern void go_cb_dispatch_call(TPCONTEXT_T ctx, TPSVCINFO *p_svc, char *name, char *fname, char *cltid);
extern int go_periodcallback();
extern int go_b4pollcallback();
extern int go_pollevent(TPCONTEXT_T ctx, int fd, unsigned int events);

static int c_tpsrvinit(int argc, char **argv)
{
	int ret;
	//Pass the current context
	TPCONTEXT_T ctx;

	//Get the context
	tpgetctxt(&ctx, 0);

	ret = go_tpsrvinit(ctx);

	//Set back the context
	tpsetctxt(ctx, 0);

	return ret;
}

static void c_tpsrvdone(void)
{
	//Pass the current context
	TPCONTEXT_T ctx;

	//Get the context
	tpgetctxt(&ctx, 0);

	go_tpsrvdone(ctx);

	//Set back the context
	tpsetctxt(ctx, 0);

}

//Initialzie the callbacks
static void c_init(void)
{
	G_tpsvrinit__ = c_tpsrvinit;
	G_tpsvrdone__ = c_tpsrvdone;
}

static int run_serv(TPCONTEXT_T *p_ctx, int *argc, char **argv)
{
	return Ondrx_main(p_ctx, *argc, argv);
}

extern __thread void * G_atmi_tls;
//Proxy function for service call
static void _GO_SVC_ENTRY (TPSVCINFO *p_svc)
{
	//Pass the current context
	TPCONTEXT_T ctx;

        NDRX_LOG(log_debug, "_GO_SVC_ENTRY entry- getting context, current ATMI context: %p", G_atmi_tls);

	//Get the context
	tpgetctxt(&ctx, 0);

        NDRX_LOG(log_debug, "_GO_SVC_ENTRY got context %p, dispatching call to GO space", ctx);

	//Call the service entry
	go_cb_dispatch_call(ctx, p_svc, p_svc->name, p_svc->fname, p_svc->cltid.clientdata);

	//Set back the context
	tpsetctxt(ctx, 0);
}


//Wrapper for advertise
static int __run_advertise(TPCONTEXT_T *p_ctx, char *svcnm, char *fname)
{
	int ret;

	ret = Otpadvertise_full(p_ctx, svcnm, _GO_SVC_ENTRY, fname);

	return ret;
}

//Wrapper for doing for doing free of the service string
static void go_tpforward (TPCONTEXT_T *p_ctx, char *svc, char *data, long len, long flags)
{
	char svcnm[XATMI_SERVICE_NAME_LENGTH+1];
	//Pass the current context
	TPCONTEXT_T ctx;

	strncpy(svcnm, svc, XATMI_SERVICE_NAME_LENGTH+1);
	svcnm[XATMI_SERVICE_NAME_LENGTH] = '\0';
	free(svc);

	Otpforward (p_ctx, svc, data, len, flags);
}

//Wrapper for tpsubscribe() - to handle data types more accurately.
static long go_tpsubscribe (TPCONTEXT_T *p_ctx, char *eventexpr, char *filter,
			long ctl_flags, char *ctl_name1, char *ctl_name2, long flags)
{
	long ret;
	TPEVCTL ctl;
	strncpy(ctl.name1, ctl_name1, XATMI_SERVICE_NAME_LENGTH);
	ctl.name1[XATMI_SERVICE_NAME_LENGTH-1] = '\0';
	strncpy(ctl.name2, ctl_name2, XATMI_SERVICE_NAME_LENGTH);
	ctl.name2[XATMI_SERVICE_NAME_LENGTH-1] = '\0';
	ctl.flags = ctl_flags;

	ret = Otpsubscribe (p_ctx, eventexpr, filter, &ctl, flags);

	free(eventexpr);
	free(filter);
	free(ctl_name1);
	free(ctl_name2);

	return ret;

}

//Wrappers b4poll callbacks
static int c_b4pollcallback(void)
{
	int ret;
	//We run in server context, thus get the current handler
	//And pass it to go func
	//Pass the current context
	TPCONTEXT_T ctx;

	//Get the context
	tpgetctxt(&ctx, 0);

	ret = go_b4pollcallback(ctx);

	//Set back the context
	tpsetctxt(ctx, 0);

        return ret;

}

//Wrappers periodic callbacks
static int c_periodcallback(void)
{
	int ret;
	//We run in server context, thus get the current handler
	//And pass it to go func
	//Pass the current context
	TPCONTEXT_T ctx;

	//Get the context
	tpgetctxt(&ctx, 0);

	ret = go_periodcallback(ctx);

	//Set back the context
	tpsetctxt(ctx, 0);

        return ret;

}

static int c_tpext_addperiodcb(TPCONTEXT_T *p_ctx, int sec)
{
	return Otpext_addperiodcb(p_ctx, sec, c_periodcallback);
}

static int c_tpext_addb4pollcb(TPCONTEXT_T *p_ctx)
{
	return Otpext_addb4pollcb(p_ctx, c_b4pollcallback);
}


//The actual event callback, will proxy the even to go
static int c_pollevent(int fd, uint32_t events, void *ptr1)
{
	int ret;
	//Pass the current context
	TPCONTEXT_T ctx;

	//Get the context
	tpgetctxt(&ctx, 0);

	ret = go_pollevent(ctx, fd, (unsigned int)events);

	//Set back the context
	tpsetctxt(ctx, 0);

	return ret;
}

//Wrapper for FD poller
static int c_tpext_addpollerfd(TPCONTEXT_T *p_ctx, int fd, unsigned int events)
{
	return Otpext_addpollerfd(p_ctx, fd, events, NULL, c_pollevent);
}


*/
import "C"
import "os"
import "unsafe"

//Servic call info
type TPSVCINFO struct {
	Name   string   /* Service name */
	Data   ATMIBuf  /* Buffer type */
	Flags  int64    /* Flags used for service invation */
	Cd     int      /* Call descriptor (generated by client) */
	Cltid  string   /* Client ID string - full client queue name */
	Appkey int64    /* RFU */
	Fname  string   /* Function name invoked (set at TpAdvertise second param) */
	Ctx    *ATMICtx /* ATMI Server Context */
}

//We need a list of functions and it's  parameter block
type fdpollcallback struct {
	cb   TPPollerFdCallback
	ptr1 interface{}
}

//Callback defintions:
type TPSrvInitFunc func(ctx *ATMICtx) int //TODO: Add parsed args after --
type TPSrvUninitFunc func(ctx *ATMICtx)
type TPServiceFunction func(ctx *ATMICtx, svc *TPSVCINFO)
type TPPeriodCallback func(ctx *ATMICtx) int //Periodic callback
type TPB4PollCallback func(ctx *ATMICtx) int //Before poll callback
type TPPollerFdCallback func(ctx *ATMICtx, fd int, events uint32, ptr1 interface{}) int

//Server init callbacks globals...
var cb_initf TPSrvInitFunc
var cb_uninitf TPSrvUninitFunc
var cb_priod TPPeriodCallback
var cb_b4poll TPB4PollCallback

//Function maps
var funcmaps map[string]TPServiceFunction
var funcpollers map[int]fdpollcallback

//export go_tpsrvinit
func go_tpsrvinit(ctx C.TPCONTEXT_T) C.int {

	var ret int
	ret = FAIL

	ac := MakeATMICtx(ctx)

	if nil != cb_initf {
		ret = cb_initf(ac)
	}

	return C.int(ret)
}

//export go_periodcallback
func go_periodcallback(ctx C.TPCONTEXT_T) C.int {
	var ret int
	ret = FAIL

	ac := MakeATMICtx(ctx)

	if nil != cb_initf {
		ret = cb_priod(ac)
	}

	return C.int(ret)

}

//export go_b4pollcallback
func go_b4pollcallback(ctx C.TPCONTEXT_T) C.int {
	var ret int
	ret = FAIL

	ac := MakeATMICtx(ctx)

	if nil != cb_initf {
		ret = cb_b4poll(ac)
	}

	return C.int(ret)

}

//export go_tpsrvdone
func go_tpsrvdone(ctx C.TPCONTEXT_T) {

	ac := MakeATMICtx(ctx)
	if nil != cb_uninitf {
		cb_uninitf(ac)
	}
}

//export go_cb_dispatch_call
func go_cb_dispatch_call(ctx C.TPCONTEXT_T, p_svc *C.TPSVCINFO, name *C.char, fname *C.char, cltid *C.char) {

	var svc TPSVCINFO

	ac := MakeATMICtx(ctx)
	//Conver the svc info
	svc.Cd = int(p_svc.cd)
	svc.Flags = int64(p_svc.flags)
	svc.Appkey = int64(p_svc.appkey)
	svc.Name = C.GoString(name)
	svc.Fname = C.GoString(fname)
	svc.Cltid = C.GoString(cltid)
	svc.Ctx = MakeATMICtx(ctx)

	//Set the data buffer...
	//TODO: Probably we want to cast it to some typed buffer...
	svc.Data.C_ptr = p_svc.data
	svc.Data.C_len = p_svc.len
	svc.Data.Ctx = svc.Ctx

	//Finalizer not needed here - auto-buffer (will be automatically free by endurox)
	//runtime.SetFinalizer(&svc.Data, nil)

	//Dispatch the call to target function...
	funcmaps[svc.Fname](ac, &svc)
}

//Continue main thread processing (go back to server polling)
func (ac *ATMICtx) TpContinue() {
	/*  This does nothing - no need to call 
        C.Otpcontinue(&ac.c_ctx)
        */
}

//We should pass here init & un-init functions...
//So that we can start the processing
//@param initf	callback to init function
//@param uninitf	callback to un-init function
//@return Enduro/X service exit code, ATMI Error
func (ac *ATMICtx) TpRun(initf TPSrvInitFunc, uninitf TPSrvUninitFunc) ATMIError {
	var err ATMIError
	C.c_init()

	//make the map of function hash
	funcmaps = make(map[string]TPServiceFunction)
	funcpollers = make(map[int]fdpollcallback)

	if nil == initf {
		/* invalid params.. */
		err = NewCustomATMIError(TPEINVAL, "init function cannot be null!")
	}
	cb_initf = initf
	cb_uninitf = uninitf

	argc := C.int(len(os.Args))
	argv := make([]*C.char, argc)
	for i, arg := range os.Args {
		argv[i] = C.CString(arg)
	}

	c_ret := C.run_serv(&ac.c_ctx, &argc, &argv[0]) // Run the Enduro/X server process

	/* Generate error, if server failed */
	if 0 != c_ret {
		err = NewCustomATMIError(TPESYSTEM, "ATMI Server failed")
	}

	for _, arg := range argv {
		C.free(unsafe.Pointer(arg))
	}

	return err
}

//Advertise service
//@param svcname		Service Name
//@param funcname	Function Name
//@return ATMI Error
func (ac *ATMICtx) TpAdvertise(svcname string, funcname string, fptr TPServiceFunction) ATMIError {
	var err ATMIError

	if nil == fptr {
		return NewCustomATMIError(TPEINVAL, "Service function must not be nil!")
	}

	c_svcname := C.CString(svcname)
	c_funcname := C.CString(funcname)

	ret := C.__run_advertise(&ac.c_ctx, c_svcname, c_funcname)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	} else {
		/* Add the function to the map */
		funcmaps[funcname] = fptr
	}

	return err
}

//Return the ATMI call and go to Q poller
//@param rvel 	Return value (TPFAIL or TPSUCCESS)
//@param rcode	Return code (used for custom purposes)
//@param tb	ATMI buffer
//@param flags	Flags
func (ac *ATMICtx) TpReturn(rval int, rcode int64, tb TypedBuffer, flags int64) {

	data := tb.GetBuf()
	C.Otpreturn(&ac.c_ctx, C.int(rval), C.long(rcode), data.C_ptr, data.C_len, C.long(flags))
}

//Forward the call to specified poller and return to Q poller
//@param svc 	Service name to forward the call to
//@param data	ATMI buffer
//@param flags	Flags
func (ac *ATMICtx) TpForward(svc string, tb TypedBuffer, flags int64) {

	data := tb.GetBuf()
	C.go_tpforward(&ac.c_ctx, C.CString(svc), data.C_ptr, data.C_len, C.long(flags))
}

//Unadvertise service dynamically
//@param	svcname	Service Name
//@return ATMI Error
func (ac *ATMICtx) TpUnadvertise(svcname string) ATMIError {
	var err ATMIError
	c_svcname := C.CString(svcname)

	ret := C.Otpunadvertise(&ac.c_ctx, c_svcname)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//Unsubscribe from event broker
//@param	subscription	Subscription ID (retruned by TPSubscribe())
//@param flags	Flags
//@return Number of subscriptions deleted, ATMI Error
func (ac *ATMICtx) TpUnsubscribe(subscription int64, flags int64) (int, ATMIError) {
	var err ATMIError
	ret := C.Otpunsubscribe(&ac.c_ctx, C.long(subscription), C.long(flags))
	if FAIL == ret {
		err = ac.NewATMIError()
	}

	return int(ret), err
}

//Subscribe service to some specified event
//@param	eventexpr	Subscription ID (retruned by TPSubscribe())
//@param filter	Event filter expression (regex)
//@param ctl Control struct
//@param flags	Flags
//@return Subscription id, ATMI Error
func (ac *ATMICtx) TpSubscribe(eventexpr string, filter string, ctl *TPEVCTL, flags int64) (int64, ATMIError) {
	var err ATMIError
	ret := C.go_tpsubscribe(&ac.c_ctx, C.CString(eventexpr), C.CString(filter),
		C.long(ctl.flags), C.CString(ctl.name1), C.CString(ctl.name2), C.long(flags))

	if FAIL == ret {
		err = ac.NewATMIError()
	}

	return int64(ret), err
}

//Get Server Call thread context data (free of *TPSRVCTXDATA must be done by user)
//@return contect data, ATMI Error
func (ac *ATMICtx) TpSrvGetCtxData() (*TPSRVCTXDATA, ATMIError) {
	var err ATMIError
	var data *TPSRVCTXDATA
	c_ptr := C.Otpsrvgetctxdata(&ac.c_ctx)

	if nil == c_ptr {
		err = ac.NewATMIError()
	} else {
		data = new(TPSRVCTXDATA)
		data.c_ptr = c_ptr
	}

	return data, err
}

//Restore thread context data
//@return ATMI Error
func (ac *ATMICtx) TpSrvSetCtxData(data *TPSRVCTXDATA, flags int64) ATMIError {
	var err ATMIError
	var ret C.int
	if nil == data || nil == data.c_ptr {
		/* Set Error */
		err = NewCustomATMIError(TPEINVAL, "Tpsrvsetctxdata - data is nil, but mandatory!")
		goto out
	}

	ret = C.Otpsrvsetctxdata(&ac.c_ctx, data.c_ptr, C.long(flags))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

out:
	return err
}

//Free the server context data
//@param data	Context data block
func (ac *ATMICtx) TpSrvFreeCtxData(data *TPSRVCTXDATA) {
	if nil != data && nil != data.c_ptr {
		C.free(unsafe.Pointer(data.c_ptr))
	}
}

//Remove the polling file descriptor
//@param fd 		FD to poll on
//@return ATMI Error
func (ac *ATMICtx) TpExtDelPollerFD(fd int) ATMIError {
	var err ATMIError
	ret := C.Otpext_delpollerfd(&ac.c_ctx, C.int(fd))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//Set periodic before poll callback func
//@param cb	Callback function with "func(ctx *ATMICtx) int" signature
//@return ATMI Error
func (ac *ATMICtx) TpExtAddB4PollCB(cb TPB4PollCallback) ATMIError {
	var err ATMIError

	if nil == cb {
		/* Set Error */
		err = NewCustomATMIError(TPEINVAL, "TpExtAddB4PollCB - cb is nil,"+
			" but mandatory!")
		return err /* <<<< RETURN! */
	}

	cb_b4poll = cb
	ret := C.c_tpext_addb4pollcb(&ac.c_ctx)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//Delete before-doing-poll callback
//@return ATMI Error
func (ac *ATMICtx) TpExtDelB4PollCB() ATMIError {
	var err ATMIError
	ret := C.Otpext_delb4pollcb(&ac.c_ctx)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//Set periodic poll callback function.
//Function is called from main service dispatcher in case if given number of seconds
//are elapsed. If the service is doing some work currenlty then it will not be interrupted.
//If the service workload was longer than period, then given period will be lost and
//will be serviced and next sleep period or after receiving next service call.
//@param secs 	Interval in secods between calls. This basically is number of seconds in
//which service will sleep and wake up.
//@param cb 	Callback function with signature: "func(ctx *ATMICtx) int".
//@return ATMI Error
func (ac *ATMICtx) TpExtAddPeriodCB(secs int, cb TPPeriodCallback) ATMIError {
	var err ATMIError

	if nil == cb {
		/* Set Error */
		err = NewCustomATMIError(TPEINVAL, "Tpext_addperiodcb - cb is nil, but mandatory!")
		return err /* <<<< RETURN! */
	}

	cb_priod = cb
	ret := C.c_tpext_addperiodcb(&ac.c_ctx, C.int(secs))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//Delete del periodic callback
//@return ATMI Error
func (ac *ATMICtx) TpExtDelPeriodCB() ATMIError {
	var err ATMIError
	ret := C.Otpext_delperiodcb(&ac.c_ctx)

	if SUCCEED != ret {
		err = ac.NewATMIError()
	}

	return err
}

//export go_pollevent
func go_pollevent(ctx C.TPCONTEXT_T, fd C.int, events C.uint) C.int {

	ac := MakeATMICtx(ctx)

	poller := funcpollers[int(fd)]
	ret := poller.cb(ac, int(fd), uint32(events), poller.ptr1)

	return C.int(ret)
}

//Add custom File Descriptor (FD) to Q poller
//@param events	Epoll events
//@param ptr1	Custom data block to be passed to callback func
//@param cb 	Callback func
//@return ATMI Error
func (ac *ATMICtx) TpExtAddPollerFD(fd int, events uint32, ptr1 interface{}, cb TPPollerFdCallback) ATMIError {
	var err ATMIError

	if nil == cb {
		/* Set Error */
		err = NewCustomATMIError(TPEINVAL, "Tpext_addpollerfd - cb is nil, but mandatory!")
		return err /* <<<< RETURN! */
	}

	ret := C.c_tpext_addpollerfd(&ac.c_ctx, C.int(fd), C.uint(events))

	if SUCCEED != ret {
		err = ac.NewATMIError()
	} else {
		var cbblock fdpollcallback
		cbblock.cb = cb
		cbblock.ptr1 = ptr1
		funcpollers[fd] = cbblock
	}

	return err
}

//Return server id
//@return server_id
func (ac *ATMICtx) TpGetSrvId() int {
	return int(C.Otpgetsrvid(&ac.c_ctx))
}
