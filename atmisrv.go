package atmi

/*
#cgo LDFLAGS: -latmisrvinteg -latmi -lrt -lm -lubf -lnstd -ldl

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>

static void free_string(char* s) { free(s); }
static char * malloc_string(int size) { return malloc(size); }

extern int go_tpsrvinit();
extern void go_tpsrvdone();
extern void go_cb_dispatch_call(TPSVCINFO *p_svc, char *name, char *fname, char *cltid);
extern int go_periodcallback();
extern int go_pollevent(int fd, unsigned int events);

static int c_tpsrvinit(int argc, char **argv)
{
	return go_tpsrvinit();
}

static void c_tpsrvdone(void)
{
	go_tpsrvdone();
}

//Initialzie the callbacks
static void c_init(void)
{
	G_tpsvrinit__ = c_tpsrvinit;
	G_tpsvrdone__ = c_tpsrvdone;
}

static int run_serv(int *argc, char **argv)
{
	return ndrx_main(*argc, argv);
}

//Proxy function for service call
static void _GO_SVC_ENTRY (TPSVCINFO *p_svc)
{
	//Call the service entry
	go_cb_dispatch_call(p_svc, p_svc->name, p_svc->fname, p_svc->cltid.clientdata);
}

//Wrapper for advertise
static int __run_advertise(char *svcnm, char *fname)
{
	int ret;

	ret = tpadvertise_full(svcnm, _GO_SVC_ENTRY, fname);

	return ret;
}

//Wrapper for doing for doing free of the service string
static void go_tpforward (char *svc, char *data, long len, long flags)
{
	char svcnm[XATMI_SERVICE_NAME_LENGTH+1];

	strncpy(svcnm, svc, XATMI_SERVICE_NAME_LENGTH+1);
	svcnm[XATMI_SERVICE_NAME_LENGTH] = '\0';
	free(svc);

	tpforward (svc, data, len, flags);
}

//Wrapper for tpsubscribe() - to handle data types more accurately.
static long go_tpsubscribe (char *eventexpr, char *filter,
			long ctl_flags, char *ctl_name1, char *ctl_name2, long flags)
{
	long ret;
	TPEVCTL ctl;
	strncpy(ctl.name1, ctl_name1, XATMI_SERVICE_NAME_LENGTH);
	ctl.name1[XATMI_SERVICE_NAME_LENGTH] = '\0';
	strncpy(ctl.name2, ctl_name2, XATMI_SERVICE_NAME_LENGTH);
	ctl.name2[XATMI_SERVICE_NAME_LENGTH] = '\0';
	ctl.flags = ctl_flags;

	ret = tpsubscribe (eventexpr, filter, &ctl, flags);

	free(eventexpr);
	free(filter);
	free(ctl_name1);
	free(ctl_name2);

	return ret;

}

//Wrappers periodic callbacks
static int c_periodcallback(void)
{
	return go_periodcallback();
}

static int c_tpext_addperiodcb(int sec)
{
	tpext_addperiodcb(sec, c_periodcallback);
}


//The actual event callback, will proxy the even to go
static int c_pollevent(int fd, uint32_t events, void *ptr1)
{
	return go_pollevent(fd, (unsigned int)events);
}

//Wrapper for FD poller
static int c_tpext_addpollerfd(int fd, unsigned int events)
{
	return tpext_addpollerfd(fd, events, NULL, c_pollevent);
}


*/
import "C"
import "os"
import "unsafe"

//import "runtime"

//Servic call info
type TPSVCINFO struct {
	Name   string
	Data   ATMIBuf
	Flags  int64
	Cd     int
	Cltid  string
	Appkey int64
	Fname  string
}

//We need a list of functions and it's  parameter block
type fdpollcallback struct {
	cb   TPPollerFdCallback
	ptr1 interface{}
}

//Callback defintions:
type TPSrvInitFunc func() int //TODO: Add parsed args after --
type TPSrvUninitFunc func()
type TPServiceFunction func(svc *TPSVCINFO)
type TPPeriodCallback func() int
type TPPollerFdCallback func(fd int, events uint32, ptr1 interface{}) int

//Server init callbacks globals...
var cb_initf TPSrvInitFunc
var cb_uninitf TPSrvUninitFunc
var cb_priod TPPeriodCallback

//Function maps
var funcmaps map[string]TPServiceFunction
var funcpollers map[int]fdpollcallback

//export go_tpsrvinit
func go_tpsrvinit() C.int {

	var ret int

	ret = FAIL

	if nil != cb_initf {
		ret = cb_initf()
	}

	return C.int(ret)
}

//export go_periodcallback
func go_periodcallback() C.int {
	return C.int(cb_priod())
}

//export go_tpsrvdone
func go_tpsrvdone() {

	if nil != cb_uninitf {
		cb_uninitf()
	}
}

//export go_cb_dispatch_call
func go_cb_dispatch_call(p_svc *C.TPSVCINFO, name *C.char, fname *C.char, cltid *C.char) {

	var svc TPSVCINFO

	//Conver the svc info
	svc.Cd = int(p_svc.cd)
	svc.Flags = int64(p_svc.flags)
	svc.Appkey = int64(p_svc.appkey)
	svc.Name = C.GoString(name)
	svc.Fname = C.GoString(fname)
	svc.Cltid = C.GoString(cltid)

	//Set the data buffer...
	//TODO: Probably we want to cast it to some typed buffer...
	svc.Data.C_ptr = p_svc.data
	svc.Data.C_len = p_svc.len

	//Finalizer not needed here - auto-buffer (will be automatically free by endurox)
	//runtime.SetFinalizer(&svc.Data, nil)

	//Dispatch the call to target function...
	funcmaps[svc.Fname](&svc)

}

//Continue main thread processing (go back to server polling)
func TpContinue() {
	C.tpcontinue()
}

//We should pass here init & un-init functions...
//So that we can start the processing
//@param initf	callback to init function
//@param uninitf	callback to un-init function
//@return Enduro/X service exit code, ATMI Error
func TpRun(initf TPSrvInitFunc, uninitf TPSrvUninitFunc) ATMIError {
	var err ATMIError
	C.c_init()

	//make the map of function hash
	funcmaps = make(map[string]TPServiceFunction)
	funcpollers = make(map[int]fdpollcallback)

	if nil == initf {
		/* invalid params.. */
		err = NewCustomAtmiError(TPEINVAL, "init function cannot be null!")
	}
	cb_initf = initf
	cb_uninitf = uninitf

	argc := C.int(len(os.Args))
	argv := make([]*C.char, argc)
	for i, arg := range os.Args {
		argv[i] = C.CString(arg)
	}

	c_ret := C.run_serv(&argc, &argv[0]) // Run the Enduro/X server process

	/* Generate error, if server failed */
	if 0 != c_ret {
		err = NewCustomAtmiError(TPESYSTEM, "ATMI Server failed")
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
func TpAdvertise(svcname string, funcname string, fptr TPServiceFunction) ATMIError {
	var err ATMIError

	if nil == fptr {
		return NewCustomAtmiError(TPEINVAL, "Service function must not be nil!")
	}

	c_svcname := C.CString(svcname)
	c_funcname := C.CString(funcname)

	ret := C.__run_advertise(c_svcname, c_funcname)

	if SUCCEED != ret {
		err = NewAtmiError()
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
func TpReturn(rval int, rcode int64, tb TypedBuffer, flags int64) {

	data := tb.GetBuf()
	C.tpreturn(C.int(rval), C.long(rcode), data.C_ptr, data.C_len, C.long(flags))
}

//Forward the call to specified poller and return to Q poller
//@param svc 	Service name to forward the call to
//@param data	ATMI buffer
//@param flags	Flags
func TpForward(svc string, tb TypedBuffer, flags int64) {

	data := tb.GetBuf()
	C.go_tpforward(C.CString(svc), data.C_ptr, data.C_len, C.long(flags))
}

//Unadvertise service dynamically
//@param	svcname	Service Name
//@return ATMI Error
func TpUnadvertise(svcname string) ATMIError {
	var err ATMIError
	c_svcname := C.CString(svcname)

	ret := C.tpunadvertise(c_svcname)

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Unsubscribe from event broker
//@param	subscription	Subscription ID (retruned by TPSubscribe())
//@param flags	Flags
//@return Number of subscriptions deleted, ATMI Error
func TpUnsubscribe(subscription int64, flags int64) (int, ATMIError) {
	var err ATMIError
	ret := C.tpunsubscribe(C.long(subscription), C.long(flags))
	if FAIL == ret {
		err = NewAtmiError()
	}

	return int(ret), err
}

//Subscribe service to some specified event
//@param	eventexpr	Subscription ID (retruned by TPSubscribe())
//@param filter	Event filter expression (regex)
//@param ctl Control struct
//@param flags	Flags
//@return Subscription id, ATMI Error
func TpSubscribe(eventexpr string, filter string, ctl *TPEVCTL, flags int64) (int64, ATMIError) {
	var err ATMIError
	ret := C.go_tpsubscribe(C.CString(eventexpr), C.CString(filter),
		C.long(ctl.flags), C.CString(ctl.name1), C.CString(ctl.name2), C.long(flags))

	if FAIL == ret {
		err = NewAtmiError()
	}

	return int64(ret), err
}

//Get Server Call thread context data (free of *TPSRVCTXDATA must be done by user)
//@return contect data, ATMI Error
func TpSrvGetCtxData() (*TPSRVCTXDATA, ATMIError) {
	var err ATMIError
	var data *TPSRVCTXDATA
	c_ptr := C.tpsrvgetctxdata()

	if nil == c_ptr {
		err = NewAtmiError()
	} else {
		data = new(TPSRVCTXDATA)
		data.c_ptr = c_ptr
	}

	return data, err
}

//Restore thread context data
//@return ATMI Error
func TpSrvSetCtxData(data *TPSRVCTXDATA, flags int64) ATMIError {
	var err ATMIError
	var ret C.int
	if nil == data || nil == data.c_ptr {
		/* Set Error */
		err = NewCustomAtmiError(TPEINVAL, "Tpsrvsetctxdata - data is nil, but mandatory!")
		goto out
	}

	ret = C.tpsrvsetctxdata(data.c_ptr, C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

out:
	return err
}

//Free the server context data
//@param data	Context data block
func TpSrvFreeCtxData(data *TPSRVCTXDATA) {
	if nil != data && nil != data.c_ptr {
		C.free(unsafe.Pointer(data.c_ptr))
	}
}

//Remove the polling file descriptor
//@param fd 		FD to poll on
//@return ATMI Error
func TpExtDelPollerfd(fd int) ATMIError {
	var err ATMIError
	ret := C.tpext_delpollerfd(C.int(fd))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Delet del periodic callback
//@return ATMI Error
func TpExtDelPeriodCB() ATMIError {
	var err ATMIError
	ret := C.tpext_delperiodcb()

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Delete before-doing-poll callback
//@return ATMI Error
func TpExtDelB4PollCB() ATMIError {
	var err ATMIError
	ret := C.tpext_delb4pollcb()

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Set periodic before poll callback func
//@return ATMI Error
func TpExtAddPeriodCB(secs int, cb TPPeriodCallback) ATMIError {
	var err ATMIError

	if nil == cb {
		/* Set Error */
		err = NewCustomAtmiError(TPEINVAL, "Tpext_addperiodcb - cb is nil, but mandatory!")
		return err /* <<<< RETURN! */
	}

	cb_priod = cb
	ret := C.c_tpext_addperiodcb(C.int(secs))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//export go_pollevent
func go_pollevent(fd C.int, events C.uint) C.int {

	poller := funcpollers[int(fd)]
	ret := poller.cb(int(fd), uint32(events), poller.ptr1)

	return C.int(ret)
}

//Add custom File Descriptor (FD) to Q poller
//@param events	Epoll events
//@param ptr1	Custom data block to be passed to callback func
//@param cb 	Callback func
//@return ATMI Error
func TpExtAddPollerFD(fd int, events uint32, ptr1 interface{}, cb TPPollerFdCallback) ATMIError {
	var err ATMIError

	if nil == cb {
		/* Set Error */
		err = NewCustomAtmiError(TPEINVAL, "Tpext_addpollerfd - cb is nil, but mandatory!")
		return err /* <<<< RETURN! */
	}

	ret := C.c_tpext_addpollerfd(C.int(fd), C.uint(events))

	if SUCCEED != ret {
		err = NewAtmiError()
	} else {
		var cbblock fdpollcallback
		cbblock.cb = cb
		cbblock.ptr1 = ptr1
		funcpollers[fd] = cbblock
	}

	return err
}
