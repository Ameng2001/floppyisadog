// Package smsserver comment
// This file was generated by tars2go 1.1.10
// Generated from Sms.tars
package smsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	m "github.com/TarsCloud/TarsGo/tars/model"
	"github.com/TarsCloud/TarsGo/tars/protocol/codec"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/basef"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/requestf"
	"github.com/TarsCloud/TarsGo/tars/protocol/tup"
	"github.com/TarsCloud/TarsGo/tars/util/current"
	"github.com/TarsCloud/TarsGo/tars/util/tools"
	"github.com/TarsCloud/TarsGo/tars/util/trace"
	"unsafe"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = fmt.Errorf
	_ = codec.FromInt8
	_ = unsafe.Pointer(nil)
	_ = bytes.ErrTooLarge
)

// Sms struct
type Sms struct {
	servant m.Servant
}

// QueueSend is the proxy function for the method defined in the tars file, with the context
func (obj *Sms) QueueSend(req *SmsRequest, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = req.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}
	tarsResp := new(requestf.ResponsePacket)
	tarsCtx := context.Background()

	err = obj.servant.TarsInvoke(tarsCtx, 0, "QueueSend", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	readBuf := codec.NewReader(tools.Int8ToByte(tarsResp.SBuffer))
	err = readBuf.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// QueueSendWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Sms) QueueSendWithContext(tarsCtx context.Context, req *SmsRequest, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = req.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	traceData, ok := current.GetTraceData(tarsCtx)
	if ok && traceData.TraceCall {
		traceData.NewSpan()
		var traceParam string
		traceParamFlag := traceData.NeedTraceParam(trace.EstCS, uint(buf.Len()))
		if traceParamFlag == trace.EnpNormal {
			value := map[string]interface{}{}
			value["req"] = req
			p, _ := json.Marshal(value)
			traceParam = string(p)
		} else if traceParamFlag == trace.EnpOverMaxLen {
			traceParam = "{\"trace_param_over_max_len\":true}"
		}
		tars.Trace(traceData.GetTraceKey(trace.EstCS), trace.TraceAnnotationCS, tars.GetClientConfig().ModuleName, obj.servant.Name(), "QueueSend", 0, traceParam, "")
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 0, "QueueSend", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	readBuf := codec.NewReader(tools.Int8ToByte(tarsResp.SBuffer))
	err = readBuf.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	if ok && traceData.TraceCall {
		var traceParam string
		traceParamFlag := traceData.NeedTraceParam(trace.EstCR, uint(readBuf.Len()))
		if traceParamFlag == trace.EnpNormal {
			value := map[string]interface{}{}
			value[""] = ret
			p, _ := json.Marshal(value)
			traceParam = string(p)
		} else if traceParamFlag == trace.EnpOverMaxLen {
			traceParam = "{\"trace_param_over_max_len\":true}"
		}
		tars.Trace(traceData.GetTraceKey(trace.EstCR), trace.TraceAnnotationCR, tars.GetClientConfig().ModuleName, obj.servant.Name(), "QueueSend", int(tarsResp.IRet), traceParam, "")
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// QueueSendOneWayWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Sms) QueueSendOneWayWithContext(tarsCtx context.Context, req *SmsRequest, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = req.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 1, "QueueSend", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// SetServant sets servant for the service.
func (obj *Sms) SetServant(servant m.Servant) {
	obj.servant = servant
}

// TarsSetTimeout sets the timeout for the servant which is in ms.
func (obj *Sms) TarsSetTimeout(timeout int) {
	obj.servant.TarsSetTimeout(timeout)
}

// TarsSetProtocol sets the protocol for the servant.
func (obj *Sms) TarsSetProtocol(p m.Protocol) {
	obj.servant.TarsSetProtocol(p)
}

// AddServant adds servant  for the service.
func (obj *Sms) AddServant(imp SmsServant, servantObj string) {
	tars.AddServant(obj, imp, servantObj)
}

// AddServantWithContext adds servant  for the service with context.
func (obj *Sms) AddServantWithContext(imp SmsServantWithContext, servantObj string) {
	tars.AddServantWithContext(obj, imp, servantObj)
}

type SmsServant interface {
	QueueSend(req *SmsRequest) (ret int32, err error)
}
type SmsServantWithContext interface {
	QueueSend(tarsCtx context.Context, req *SmsRequest) (ret int32, err error)
}

// Dispatch is used to call the server side implement for the method defined in the tars file. withContext shows using context or not.
func (obj *Sms) Dispatch(tarsCtx context.Context, val interface{}, tarsReq *requestf.RequestPacket, tarsResp *requestf.ResponsePacket, withContext bool) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	readBuf := codec.NewReader(tools.Int8ToByte(tarsReq.SBuffer))
	buf := codec.NewBuffer()
	switch tarsReq.SFuncName {
	case "QueueSend":
		var req SmsRequest

		if tarsReq.IVersion == basef.TARSVERSION {

			err = req.ReadBlock(readBuf, 1, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			reqTup := tup.NewUniAttribute()
			reqTup.Decode(readBuf)

			var tupBuffer []byte

			reqTup.GetBuffer("req", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = req.ReadBlock(readBuf, 0, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.JSONVERSION {
			var jsonData map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
			decoder.UseNumber()
			err = decoder.Decode(&jsonData)
			if err != nil {
				return fmt.Errorf("decode reqpacket failed, error: %+v", err)
			}
			{
				jsonStr, _ := json.Marshal(jsonData["req"])
				req.ResetDefault()
				if err = json.Unmarshal(jsonStr, &req); err != nil {
					return err
				}
			}

		} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}

		traceData, ok := current.GetTraceData(tarsCtx)
		if ok && traceData.TraceCall {
			var traceParam string
			traceParamFlag := traceData.NeedTraceParam(trace.EstSR, uint(readBuf.Len()))
			if traceParamFlag == trace.EnpNormal {
				value := map[string]interface{}{}
				value["req"] = req
				p, _ := json.Marshal(value)
				traceParam = string(p)
			} else if traceParamFlag == trace.EnpOverMaxLen {
				traceParam = "{\"trace_param_over_max_len\":true}"
			}
			tars.Trace(traceData.GetTraceKey(trace.EstSR), trace.TraceAnnotationSR, tars.GetClientConfig().ModuleName, tarsReq.SServantName, "QueueSend", 0, traceParam, "")
		}

		var funRet int32
		if !withContext {
			imp := val.(SmsServant)
			funRet, err = imp.QueueSend(&req)
		} else {
			imp := val.(SmsServantWithContext)
			funRet, err = imp.QueueSend(tarsCtx, &req)
		}

		if err != nil {
			return err
		}

		if tarsReq.IVersion == basef.TARSVERSION {
			buf.Reset()

			err = buf.WriteInt32(funRet, 0)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			rspTup := tup.NewUniAttribute()

			err = buf.WriteInt32(funRet, 0)
			if err != nil {
				return err
			}

			rspTup.PutBuffer("", buf.ToBytes())
			rspTup.PutBuffer("tars_ret", buf.ToBytes())

			buf.Reset()
			err = rspTup.Encode(buf)
			if err != nil {
				return err
			}
		} else if tarsReq.IVersion == basef.JSONVERSION {
			rspJson := map[string]interface{}{}
			rspJson["tars_ret"] = funRet

			var rspByte []byte
			if rspByte, err = json.Marshal(rspJson); err != nil {
				return err
			}

			buf.Reset()
			err = buf.WriteSliceUint8(rspByte)
			if err != nil {
				return err
			}
		}

		if ok && traceData.TraceCall {
			var traceParam string
			traceParamFlag := traceData.NeedTraceParam(trace.EstSS, uint(buf.Len()))
			if traceParamFlag == trace.EnpNormal {
				value := map[string]interface{}{}
				value[""] = funRet
				p, _ := json.Marshal(value)
				traceParam = string(p)
			} else if traceParamFlag == trace.EnpOverMaxLen {
				traceParam = "{\"trace_param_over_max_len\":true}"
			}
			tars.Trace(traceData.GetTraceKey(trace.EstSS), trace.TraceAnnotationSS, tars.GetClientConfig().ModuleName, tarsReq.SServantName, "QueueSend", 0, traceParam, "")
		}

	default:
		return fmt.Errorf("func mismatch")
	}
	var statusMap map[string]string
	if status, ok := current.GetResponseStatus(tarsCtx); ok && status != nil {
		statusMap = status
	}
	var contextMap map[string]string
	if ctx, ok := current.GetResponseContext(tarsCtx); ok && ctx != nil {
		contextMap = ctx
	}
	*tarsResp = requestf.ResponsePacket{
		IVersion:     tarsReq.IVersion,
		CPacketType:  0,
		IRequestId:   tarsReq.IRequestId,
		IMessageType: 0,
		IRet:         0,
		SBuffer:      tools.ByteToInt8(buf.ToBytes()),
		Status:       statusMap,
		SResultDesc:  "",
		Context:      contextMap,
	}

	_ = readBuf
	_ = buf
	_ = length
	_ = have
	_ = ty
	return nil
}
