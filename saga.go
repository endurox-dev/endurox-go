/**
 * @brief Saga registration, recovery and transaction identifiers
 *
 * @file saga.go
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
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"encoding/base64"
	"runtime"
	"strings"
	"unsafe"
)

// XID holds a native global transaction identifier. Obtain it from a Saga
// registration or TpXIDFromString; it is distinct from a branch XID string.
type XID struct{ c_xid C.XID }

// TpSagaCtl owns one native Saga control and its payload. Empty callback URLs omit actions.
// The zero value is ready to use. Copies of an initialized value share the native record; use a new
// value for independent state. Access a record serially.
type TpSagaCtl struct {
	c_ctl *C.TPSAGACTL
}

// native returns the single backing C record, allocating it on first use.
// @return Native record shared by this wrapper and its copies
func (ctl *TpSagaCtl) native() *C.TPSAGACTL {
	if nil == ctl.c_ctl {
		ctl.c_ctl = new(C.TPSAGACTL)
		runtime.SetFinalizer(ctl.c_ctl, freeSagaControl)
	}
	return ctl.c_ctl
}

// GetFlags reads registration mode flags from the native record.
// @return Current registration mode flags
func (ctl *TpSagaCtl) GetFlags() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().flags)
}

// SetFlags writes registration mode flags to the native record.
// @param value New registration mode flags
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetFlags(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Flags exceeds native long")
	}
	ctl.native().flags = C.long(value)
	return nil
}

// GetTimeout reads the lifetime requested at first registration.
// @return Seconds from first registration, or zero for no participant limit
func (ctl *TpSagaCtl) GetTimeout() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().timeout)
}

// SetTimeout sets the optional registration lifetime. It can shorten the global
// deadline and bounds later renewal. Metadata updates retain the original limit.
// @param value Seconds from first registration; zero adds no limit
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetTimeout(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if value < 0 || int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Invalid Saga timeout")
	}
	ctl.native().timeout = C.long(value)
	return nil
}

// GetSagaName reads Saga profile name from the native record.
// @return Current Saga profile name
func (ctl *TpSagaCtl) GetSagaName() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().saga_name[0])
}

// SetSagaName writes Saga profile name to the native record.
// @param value New Saga profile name
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetSagaName(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().saga_name[:], value, "SagaName")
}

// GetSagaKey reads idempotency key from the native record.
// @return Current idempotency key
func (ctl *TpSagaCtl) GetSagaKey() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().saga_key[0])
}

// SetSagaKey writes idempotency key to the native record.
// @param value New idempotency key
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetSagaKey(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().saga_key[:], value, "SagaKey")
}

// GetContentType reads payload MIME type from the native record.
// @return Current payload MIME type
func (ctl *TpSagaCtl) GetContentType() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().content_type[0])
}

// SetContentType writes payload MIME type to the native record.
// @param value New payload MIME type
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetContentType(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().content_type[:], value, "ContentType")
}

// GetBranchXID reads participant branch identifier from the native record.
// @return Current participant branch identifier
func (ctl *TpSagaCtl) GetBranchXID() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().branch_xid[0])
}

// SetBranchXID writes participant branch identifier to the native record.
// @param value New participant branch identifier
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetBranchXID(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().branch_xid[:], value, "BranchXID")
}

// GetCommitURL reads commit callback URL from the native record.
// @return Current commit callback URL
func (ctl *TpSagaCtl) GetCommitURL() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().commit_url[0])
}

// SetCommitURL writes commit callback URL to the native record.
// @param value New commit callback URL
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetCommitURL(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().commit_url[:], value, "CommitURL")
}

// GetRollbackURL reads rollback callback URL from the native record.
// @return Current rollback callback URL
func (ctl *TpSagaCtl) GetRollbackURL() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().rollback_url[0])
}

// SetRollbackURL writes rollback callback URL to the native record.
// @param value New rollback callback URL
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetRollbackURL(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().rollback_url[:], value, "RollbackURL")
}

// GetStatusURL reads status callback URL from the native record.
// @return Current status callback URL
func (ctl *TpSagaCtl) GetStatusURL() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().status_url[0])
}

// SetStatusURL writes status callback URL to the native record.
// @param value New status callback URL
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetStatusURL(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().status_url[:], value, "StatusURL")
}

// GetForgetURL reads forget callback URL from the native record.
// @return Current forget callback URL
func (ctl *TpSagaCtl) GetForgetURL() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().forget_url[0])
}

// SetForgetURL writes forget callback URL to the native record.
// @param value New forget callback URL
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetForgetURL(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().forget_url[:], value, "ForgetURL")
}

// GetAfterURL reads after callback URL from the native record.
// @return Current after callback URL
func (ctl *TpSagaCtl) GetAfterURL() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().after_url[0])
}

// SetAfterURL writes after callback URL to the native record.
// @param value New after callback URL
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpSagaCtl) SetAfterURL(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	return copyCText(ctl.native().after_url[:], value, "AfterURL")
}

// freeSagaControl releases payload storage after every wrapper copy is unreachable.
// @param ctl Native control whose payload is no longer in use
func freeSagaControl(ctl *C.TPSAGACTL) {
	C.free(unsafe.Pointer(ctl.payload))
}

// GetPayload copies the bytes selected by the native payload length.
// @return Independent payload bytes; nil when the payload is empty
func (ctl *TpSagaCtl) GetPayload() []byte {
	defer runtime.KeepAlive(ctl)
	if nil == ctl || nil == ctl.c_ctl || 0 == ctl.c_ctl.payload_len {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(ctl.c_ctl.payload), C.int(ctl.c_ctl.payload_len))
}

// SetPayload copies bytes into control-owned C memory and sets payload_len.
// Passing nil or an empty slice releases the payload immediately.
// @param value Payload bytes, at most NDRX_SAGA_PAYLOAD_MAX
// @return Nil on success, otherwise an error with the previous payload intact
func (ctl *TpSagaCtl) SetPayload(value []byte) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl || len(value) > NDRX_SAGA_PAYLOAD_MAX {
		return NewCustomATMIError(TPEINVAL, "invalid Saga control or payload length")
	}
	native := ctl.native()
	if 0 == len(value) {
		C.free(unsafe.Pointer(native.payload))
		native.payload = nil
		native.payload_len = 0
		return nil
	}
	if nil == native.payload {
		data := C.calloc(C.size_t(NDRX_SAGA_PAYLOAD_MAX), 1)
		if nil == data {
			return NewCustomATMIError(TPEOS, "failed to allocate Saga payload")
		}
		native.payload = (*C.char)(data)
	}
	C.memset(unsafe.Pointer(native.payload), 0, C.size_t(NDRX_SAGA_PAYLOAD_MAX))
	C.memcpy(unsafe.Pointer(native.payload), unsafe.Pointer(&value[0]), C.size_t(len(value)))
	native.payload_len = C.long(len(value))
	return nil
}

// GetPayloadLen reads the payload length stored in the C record.
// @return Payload length in bytes
func (ctl *TpSagaCtl) GetPayloadLen() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl || nil == ctl.c_ctl {
		return 0
	}
	return int64(ctl.c_ctl.payload_len)
}

// SetPayloadLen changes the payload view; newly exposed bytes are zeroed.
// @param length New length in bytes, between zero and NDRX_SAGA_PAYLOAD_MAX
// @return Nil on success, otherwise an error without changing the payload
func (ctl *TpSagaCtl) SetPayloadLen(length int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl || 0 > length || NDRX_SAGA_PAYLOAD_MAX < length {
		return NewCustomATMIError(TPEINVAL, "invalid Saga payload length")
	}
	native := ctl.native()
	if 0 == length {
		C.free(unsafe.Pointer(native.payload))
		native.payload = nil
		native.payload_len = 0
		return nil
	}
	if nil == native.payload {
		data := C.calloc(C.size_t(NDRX_SAGA_PAYLOAD_MAX), 1)
		if nil == data {
			return NewCustomATMIError(TPEOS, "failed to allocate Saga payload")
		}
		native.payload = (*C.char)(data)
	}
	if native.payload_len < C.long(length) {
		C.memset(unsafe.Add(unsafe.Pointer(native.payload), int(native.payload_len)),
			0, C.size_t(length-int64(native.payload_len)))
	}
	native.payload_len = C.long(length)
	return nil
}

// GetPreviousPayload copies the data replaced by a TPSAGA_DATAONLY rejoin.
// @return Previous participant data, or nil when empty or ctl is nil
func (ctl *TpSagaCtl) GetPreviousPayload() []byte {
	defer runtime.KeepAlive(ctl)
	if nil == ctl || ctl.native().previous_payload_len == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(&ctl.native().previous_payload[0]),
		C.int(ctl.native().previous_payload_len))
}

// TpTxStatus holds one native transaction selector and status result.
// The zero value is ready to use. Copies of an initialized value share the native record; use a new
// value for independent state. Access a record serially.
type TpTxStatus struct {
	c_ctl *C.TPTXSTATUS
}

// native returns the single backing C record, allocating it on first use.
// @return Native record shared by this wrapper and its copies
func (ctl *TpTxStatus) native() *C.TPTXSTATUS {
	if nil == ctl.c_ctl {
		ctl.c_ctl = new(C.TPTXSTATUS)
	}
	return ctl.c_ctl
}

// GetFlags reads TPTXS identifier selector from the native record.
// @return Current TPTXS identifier selector
func (ctl *TpTxStatus) GetFlags() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().flags)
}

// SetFlags writes TPTXS identifier selector to the native record.
// @param value New TPTXS identifier selector
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetFlags(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Flags exceeds native long")
	}
	ctl.native().flags = C.long(value)
	return nil
}

// GetBranchXID reads participant branch identifier from the native record.
// @return Current participant branch identifier
func (ctl *TpTxStatus) GetBranchXID() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().branch_xid[0])
}

// SetBranchXID writes participant branch identifier to the native record.
// @param value New participant branch identifier
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetBranchXID(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	return copyCText(ctl.native().branch_xid[:], value, "BranchXID")
}

// GetStatus reads transaction status from the native record.
// @return Current transaction status
func (ctl *TpTxStatus) GetStatus() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().status)
}

// SetStatus writes transaction status to the native record.
// @param value New transaction status
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetStatus(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Status exceeds native long")
	}
	ctl.native().status = C.long(value)
	return nil
}

// GetTimeout reads timeout in seconds from the native record.
// @return Current timeout in seconds
func (ctl *TpTxStatus) GetTimeout() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().timeout)
}

// SetTimeout writes timeout in seconds to the native record.
// @param value New timeout in seconds
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetTimeout(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Timeout exceeds native long")
	}
	ctl.native().timeout = C.long(value)
	return nil
}

// GetTimeoutLeft reads remaining lifetime in seconds from the native record.
// @return Current remaining lifetime in seconds
func (ctl *TpTxStatus) GetTimeoutLeft() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().timeout_left)
}

// SetTimeoutLeft writes remaining lifetime in seconds to the native record.
// @param value New remaining lifetime in seconds
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetTimeoutLeft(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "TimeoutLeft exceeds native long")
	}
	ctl.native().timeout_left = C.long(value)
	return nil
}

// GetXID returns a value copy of the native global transaction identifier.
// @return Native global XID
func (ctl *TpSagaCtl) GetXID() XID {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return XID{}
	}
	return XID{c_xid: ctl.native().xid}
}

// SetXID stores a native global transaction identifier.
// @param value Global XID obtained from registration or TpXIDFromString
// @return Nil on success, otherwise TPEINVAL for a nil control
func (ctl *TpSagaCtl) SetXID(value XID) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	ctl.native().xid = value.c_xid
	return nil
}

// GetXID returns a value copy of the native global transaction identifier.
// @return Native global XID
func (ctl *TpTxStatus) GetXID() XID {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return XID{}
	}
	return XID{c_xid: ctl.native().xid}
}

// SetXID stores a native global transaction identifier.
// @param value Global XID obtained from registration or TpXIDFromString
// @return Nil on success, otherwise TPEINVAL for a nil control
func (ctl *TpTxStatus) SetXID(value XID) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	ctl.native().xid = value.c_xid
	return nil
}

// GetTranID returns a value copy of the native suspend token.
// @return Native TPTRANID, valid only for the matching runtime ABI
func (ctl *TpTxStatus) GetTranID() TPTRANID {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return TPTRANID{}
	}
	return TPTRANID{c_tptranid: ctl.native().tranid}
}

// SetTranID stores an existing native suspend token without converting fields.
// @param value Token returned by TpSuspend
// @return Nil on success, otherwise TPEINVAL for a nil control
func (ctl *TpTxStatus) SetTranID(value TPTRANID) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "transaction status is nil")
	}
	ctl.native().tranid = value.c_tptranid
	return nil
}

// GetTranIDString encodes the stored native token for TpResumeString.
// @return Base64 token for this runtime ABI
func (ctl *TpTxStatus) GetTranIDString() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return base64.StdEncoding.EncodeToString(C.GoBytes(
		unsafe.Pointer(&ctl.native().tranid), C.sizeof_TPTRANID))
}

// SetTranIDString decodes a TpSuspendString token directly into the native field.
// @param value Base64 native TPTRANID token
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpTxStatus) SetTranIDString(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "transaction status is nil")
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if nil != err || len(decoded) != C.sizeof_TPTRANID {
		return NewCustomATMIError(TPEINVAL, "invalid Base64 TPTRANID token for this runtime ABI")
	}
	C.memcpy(unsafe.Pointer(&ctl.native().tranid), unsafe.Pointer(&decoded[0]), C.sizeof_TPTRANID)
	return nil
}

// GetClientID reads the durable LRA client identity.
// @return Client identifier returned by TpTxStatus or TpTxSetInfo
func (st *TpTxStatus) GetClientID() string {
	defer runtime.KeepAlive(st)
	if nil == st {
		return ""
	}
	return C.GoString(&st.native().client_id[0])
}

// SetClientID selects the client identity to store through TpTxSetInfo.
// @param value Client identity, at most NDRX_LRA_CLIENT_ID_MAX UTF-8 bytes
// @return Nil on success, otherwise TPEINVAL without changing the value
func (st *TpTxStatus) SetClientID(value string) ATMIError {
	defer runtime.KeepAlive(st)
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	return copyCText(st.native().client_id[:], value, "ClientID")
}

// GetLRAURL reads the durable LRA coordinator base URL; LRA ids are
// <URL>/<global XID>.
// @return Coordinator URL returned by TpTxStatus or TpTxSetInfo, empty if none
func (st *TpTxStatus) GetLRAURL() string {
	defer runtime.KeepAlive(st)
	if nil == st {
		return ""
	}
	return C.GoString(&st.native().lra_url[0])
}

// SetLRAURL selects the coordinator base URL to store through TpTxSetInfo.
// The participant callbacks of the LRA carry ids under this URL.
// @param value Empty, or an http(s) URL of at most NDRX_SAGA_URL_MAX bytes
// @return Nil on success, otherwise TPEINVAL without changing the value
func (st *TpTxStatus) SetLRAURL(value string) ATMIError {
	defer runtime.KeepAlive(st)
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	return copyCText(st.native().lra_url[:], value, "LRAURL")
}

// GetStartTime reads the transaction start timestamp.
// @return Milliseconds since the Unix epoch, or zero for an unknown transaction
func (st *TpTxStatus) GetStartTime() int64 {
	defer runtime.KeepAlive(st)
	if nil == st {
		return 0
	}
	return int64(st.native().start_time)
}

// GetFinishTime reads the deadline while active, or latest progress timestamp.
// @return Milliseconds since the Unix epoch, or zero for an unknown transaction
func (st *TpTxStatus) GetFinishTime() int64 {
	defer runtime.KeepAlive(st)
	if nil == st {
		return 0
	}
	return int64(st.native().finish_time)
}

// IsLRA reports whether the selected transaction was begun with TPSAGALRA.
// @return True for LRA transactions
func (st *TpTxStatus) IsLRA() bool {
	defer runtime.KeepAlive(st)
	return nil != st && st.native().is_lra != 0
}

// IsRecovering reports pending transaction completion or after notifications.
// @return True while a known LRA has completion work remaining
func (st *TpTxStatus) IsRecovering() bool {
	defer runtime.KeepAlive(st)
	return nil != st && st.native().is_recovering != 0
}

// TpXIDToString converts a global XID through the native TPCONVXID codec.
// @param xid Native global XID obtained from Saga registration or
//
//	TpXIDFromString
//
// @return Native TPCONVXID string and nil error, or empty string and ATMIError
func (ac *ATMICtx) TpXIDToString(xid *XID) (string, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return "", err
	}
	if nil == xid {
		return "", NewCustomATMIError(TPEINVAL, "XID is nil")
	}
	var text [C.TPCONVMAXSTR]C.char
	if C.Otpconvert(&ac.c_ctx, &text[0],
		(*C.char)(unsafe.Pointer(&xid.c_xid)), C.TPCONVXID|C.TPTOSTRING) != SUCCEED {
		return "", ac.NewATMIError()
	}
	return C.GoString(&text[0]), nil
}

// TpXIDFromString decodes a global XID using the native TPCONVXID codec.
// Branch XID strings and TpSuspendString tokens use different representations.
// @param value NUL-free global XID string produced by native TPCONVXID
//
//	conversion
//
// @return Decoded native XID and nil error, or nil and ATMIError
func (ac *ATMICtx) TpXIDFromString(value string) (*XID, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if len(value) == 0 || len(value) >= C.TPCONVMAXSTR ||
		strings.IndexByte(value, 0) >= 0 {
		return nil, NewCustomATMIError(TPEINVAL, "invalid global XID string")
	}
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	var xid XID
	if C.Otpconvert(&ac.c_ctx, cValue, (*C.char)(unsafe.Pointer(&xid.c_xid)),
		C.TPCONVXID) != SUCCEED {
		return nil, ac.NewATMIError()
	}
	return &xid, nil
}

// TpTranIDFromXID converts an Enduro/X global XID string to a participant
// token for TpResume. It does not grant TpCommit or TpAbort initiator rights.
// After XA work, detach with TpSuspend; the initiator decides the outcome.
// On failure detach with TpSuspend(TPTXABORTONLY): the initiator's commit
// then rolls back with TPEABORT.
// Conversion does not contact the coordinator or check transaction state.
// Dynamic RMs enlist when used, and manual RMs require TpRMStart after resume.
// @param xid NUL-free global XID string, not a branch ID or an LRA URL
// @return Participant transaction identifier and nil error, or nil and ATMIError
func (ac *ATMICtx) TpTranIDFromXID(xid string) (*TPTRANID, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if len(xid) == 0 || len(xid) >= C.TPCONVMAXSTR ||
		strings.IndexByte(xid, 0) >= 0 {
		return nil, NewCustomATMIError(TPEINVAL, "invalid global XID string")
	}
	cValue := C.CString(xid)
	defer C.free(unsafe.Pointer(cValue))
	var tranid TPTRANID
	if C.Otpconvert(&ac.c_ctx, cValue, (*C.char)(unsafe.Pointer(&tranid.c_tptranid)),
		C.TPCONVXIDTRAN) != SUCCEED {
		return nil, ac.NewATMIError()
	}
	return &tranid, nil
}

// GetXIDString returns the global XID of the transaction identifier in the
// native TPCONVXID string form, e.g. for building an LRA id.
// @return Global XID string, or empty string for a nil or zero-value identifier
func (t *TPTRANID) GetXIDString() string {
	defer runtime.KeepAlive(t)
	if nil == t {
		return ""
	}
	return C.GoString(&t.c_tptranid.tmxid[0])
}

// TpSagaEnlist registers ctl in the current branch, or ctl.GetXID() with
// TPSAGA_JOINXID. Read the returned identifiers with GetXID and GetBranchXID.
// @param ctl Registration metadata, explicit callback URLs and payload up to
//
//	NDRX_SAGA_PAYLOAD_MAX bytes; identifiers are available through getters
//
// @param flags TPSAGA_JOINXID or TPSAGA_LRAJOIN select ctl.GetXID(); TPSAGA_UPDATE for
//
//	explicit replacement, with TPSAGA_DATAONLY to preserve URLs and return old data
//
// @return Nil on successful registration, otherwise ATMIError; plain duplicate
//
//	registration in a branch fails
func (ac *ATMICtx) TpSagaEnlist(ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if C.Otpsagaenlist(&ac.c_ctx, ctl.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpSagaEnd ends the participant's local working section, retaining its
// registration. If dlmspace is configured, it releases the lock provided
// by the Distributed Lock Manager (DLM) in Enduro/X Enterprise Edition
// (endurox-ee). Standard mode requires this call even without DLM.
// Only ctl.GetFlags() is input; flags must be zero. TPSAGA_NOLOCK
// registrations have no section to end.
// @param ctl Saga control; only registration mode Flags is read
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpSagaEnd(ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if C.Otpsagaend(&ac.c_ctx, ctl.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpSagaLeave withdraws ctl.GetBranchXID() from an active transaction,
// preventing its callbacks. It does not release the local working section;
// call TpSagaEnd for that. TPSAGA_BYKEY selects XID, SagaName and SagaKey.
// @param ctl Control with BranchXID identifying the registered participant
// @param flags 0 for BranchXID, TPSAGA_BYKEY for XID plus profile and key
// @return Nil on successful removal, otherwise ATMIError
func (ac *ATMICtx) TpSagaLeave(ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if C.Otpsagaleave(&ac.c_ctx, ctl.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpSagaGet reads registered metadata and URLs by ctl.GetBranchXID(), including
// during recovery. Payload is preserved; flags must be zero.
// @param ctl BranchXID selects the participant; receives metadata,
//
//	registration mode, global XID and URLs without returning payload bytes
//
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError; the native payload is
//
//	preserved
func (ac *ATMICtx) TpSagaGet(ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if C.Otpsagaget(&ac.c_ctx, ctl.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpSagaUpdate replaces an active participant's payload and URLs. With
// TPSAGA_URLSONLY it relocates URLs during recovery without changing callback
// presence. URL-only calls ignore Payload, ContentType and registration Flags.
// @param ctl Participant BranchXID and replacement callback URLs; full updates
//
//	also copy Payload, ContentType and registration Flags
//
// @param flags 0 for an active participant full update, or TPSAGA_URLSONLY to
//
//	relocate URLs during recovery
//
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpSagaUpdate(ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpSagaCtl is nil")
	}
	if C.Otpsagaupdate(&ac.c_ctx, ctl.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTxStatus reads the status and timeouts of the selected transaction.
// Flags selects one identifier; the call's flags argument must be zero.
// @param st Transaction selector: Flags is exactly one of TPTXS_XID,
//
//	TPTXS_BRXID or TPTXS_TRANID; receives Status, Timeout and TimeoutLeft
//
// @param flags Reserved; must be 0
// @return Nil on successful lookup, otherwise ATMIError; returned timeouts are
//
//	seconds
func (ac *ATMICtx) TpTxStatus(st *TpTxStatus, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(st)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if C.Otptxstatus(&ac.c_ctx, st.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTxSetInfo durably stores the client identity and coordinator URL of an
// active LRA.
// @param st Selector, ClientID and LRAURL to store; receives current information
// @param flags Reserved; must be 0
// @return Nil after the metadata is stored, otherwise ATMIError
func (ac *ATMICtx) TpTxSetInfo(st *TpTxStatus, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(st)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if C.Otptxsetinfo(&ac.c_ctx, st.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTxComplete closes or cancels an LRA without joining the calling context.
// An optional participant-data update is serialized with the completion decision.
// Query TpTxStatus afterward to distinguish pending completion from final failure.
// @param st Transaction selector as for TpTxStatus; input only
// @param ctl Nil for no update, or SagaName, SagaKey and Payload for the final update
// @param flags 0 for close, TPTXC_ABORT for cancel
// @return Nil on completion, otherwise the native commit/abort ATMIError
func (ac *ATMICtx) TpTxComplete(st *TpTxStatus, ctl *TpSagaCtl, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(st)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	var native *C.TPSAGACTL
	if nil != ctl {
		native = ctl.native()
	}
	if C.Otptxcomplete(&ac.c_ctx, st.native(), native, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTxRenew requests st.GetTimeout() seconds from now, capped by the earliest
// participant deadline. The returned status contains the effective timeout.
// Flags selects one identifier; the call's flags argument must be zero.
// @param st Transaction selector and requested Timeout in seconds from now;
//
//	receives Status, Timeout and TimeoutLeft
//
// @param flags 0, or TPTXR_EXTEND to ignore shortening under the TM lock
// @return Nil on successful renewal, otherwise ATMIError
func (ac *ATMICtx) TpTxRenew(st *TpTxStatus, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(st)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == st {
		return NewCustomATMIError(TPEINVAL, "TpTxStatus is nil")
	}
	if C.Otptxrenew(&ac.c_ctx, st.native(), C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

/* vim: set ts=4 sw=4 et smartindent: */
