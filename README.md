# Application Server for GO (ASG)

Version: v8.0.0 - Using Object-mode API.

ASG is application server for Golang, which makes it possible to process 
distributed transactions in Golang. It is possible to reload the application 
components without service interruption. Basically system is service oriented 
where server components advertises services (service is just a literal name like 
"GETBALANCE", "TRXREQ", etc.), then later these services are called by client 
binaries. Clients and servers can be located on different physical machines and 
they can call each other with out knowledge of their psychical location. As 
these server binaries are stateless, they can be started in multiple copies. 
This ensures fault tolerant processing. Clients and Server Services comunicate 
via middleware which supports three kind of buffers for request/response data: 
Arbitrary string, Byte array, Unified Buffer Format (UBF), JSON buffer.


## Build & test status

| OS   |      Status      | OS       |      Status   |OS       |      Status   |
|----------|:-------------:|----------|:-------------:|----------|:-------------:|
|RHEL/Oracle Linux 8| [![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-ol8)](http://www.silodev.com:9090/jenkins/job/endurox-go-ol8/) | Centos 6|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-centos6)](http://www.silodev.com:9090/jenkins/job/endurox-go-centos6/)|FreeBSD 11|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-freebsd11)](http://www.silodev.com:9090/jenkins/job/endurox-go-freebsd11/)|
|Oracle Linux 7|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-ol7)](http://www.silodev.com:9090/jenkins/job/endurox-go-ol7/)|OSX 11.4|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-osx11_4)](http://www.silodev.com:9090/jenkins/job/endurox-go-osx11_4/)|raspbian10_arv7l|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-raspbian10_arv7l)](http://www.silodev.com:9090/jenkins/job/endurox-go-raspbian10_arv7l/)|
|SLES 12|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-sles12)](http://www.silodev.com:9090/jenkins/job/endurox-go-sles12/)|SLES 15|[![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-sles15)](http://www.silodev.com:9090/jenkins/job/endurox-go-sles15/)|Ubuntu 14.04| [![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-ubuntu14)](http://www.silodev.com:9090/jenkins/job/endurox-go-ubuntu14/)|
|Ubuntu 18.04| [![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-ubuntu18)](http://www.silodev.com:9090/jenkins/job/endurox-go-ubuntu18/)|AIX 7.2| [![Build Status](http://www.silodev.com:9090/jenkins/buildStatus/icon?job=endurox-go-aix7_2)](http://www.silodev.com:9090/jenkins/job/endurox-go-aix7_2/)|


## Documentation

[The Enduro-GO API Document](doc/endurox-go-book.adoc)

Enduro/X documentation is located here: http://www.endurox.org/dokuwiki

Basic ASG application layouts can be checked out from this repository "tests" folder.

## Foundation

ASG is built on Enduro/X middleware framework, which by itself implements 
extended XATMI specification. For distributed transaction processing XA API is 
used. XA must be supported by underlaying SQL (or any other resource) driver. 
The platform is build on GNU/Linux technology and it utilizes Posix kernel 
queues for gaining high IPC throughput.

Before try to build ASG, you need to install Enduro/X Middleware Platform, 
https://github.com/endurox-dev/endurox
Binary packages for Enduro/X are available here: http://www.endurox.org/projects/endurox/files

## Persistent message queue

Enduro/X provides a queuing subsystem called TMQ (Transactional Message Queue). 
This facility provides persistent queues that allows applications to explicitly 
enqueue and dequeue messages from named queues. Queues can be ordered by message 
en-queue time in LIFO or FIFO order. Queues are managed by an XA compliant 
resource manager allowing queue operations to participate in distributed 
transactions. An automated queue forwarding feature is provided that will remove 
entries from a queue and invoke an associated Enduro/X ATMI services, placing 
the reply message on an associated reply queue and failed messages to failure 
queue. The basic usage of persistent queues can be checked in 
tests/07_basic_durable_queue folder. TMQ API consists of TpEnqueue() and 
TpDequeue() calls or automated dequeue and message forwarding to destination 
service.

## Buffer Types

Typed buffers are used for data transport between services.

### String buffer

It is possible to send to services arbitrary strings. These could be JSON, XML 
or whatever data. The service might respond with the same buffer format, with 
changed contents. 

### Byte array (Carray)

It is possible to send to services byte arrays. The data could include binary 
zeros.

### Unified Buffer Format buffer

UBF buffer basically is hash-list where for each value there could be array of 
elements (e.g. one or more). The buffer is typed. Fields are predefined in field 
definition tables, later with Enduro/X's 'mkfldhdr -m1' can be generated field 
constant tables which provides go format.

### JSON buffer

JSON Text format buffer is supported. This can be used to call Enduro/X server 
or receive JSON format calls from the system.

## XA SQL Drivers

Currently Enduro/X supports Oracle DB OCI driver. The patched version for XA 
processing is available here: https://github.com/endurox-dev/go-oci8. When doing 
processing in XA mode, the connection string must be empty ("").

## Contact

Forums: http://www.endurox.org/projects/endurox-go/boards

# Releases

- Version 8.0.0 released on 09/01/2022 (stable) Support #754, Support #780

