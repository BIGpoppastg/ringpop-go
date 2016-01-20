// Autogenerated by Thrift Compiler (1.0.0-dev)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package tcollector

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type TCollector interface {
	// Parameters:
	//  - Span
	Submit(span *Span) (r *Response, err error)
	// Parameters:
	//  - Spans
	MultiSubmit(spans []*Span) (r []*Response, err error)
}

type TCollectorClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewTCollectorClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *TCollectorClient {
	return &TCollectorClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewTCollectorClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *TCollectorClient {
	return &TCollectorClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Span
func (p *TCollectorClient) Submit(span *Span) (r *Response, err error) {
	if err = p.sendSubmit(span); err != nil {
		return
	}
	return p.recvSubmit()
}

func (p *TCollectorClient) sendSubmit(span *Span) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("submit", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := TCollectorSubmitArgs{
		Span: span,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *TCollectorClient) recvSubmit() (value *Response, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "submit" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "submit failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "submit failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error2 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error3 error
		error3, err = error2.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error3
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "submit failed: invalid message type")
		return
	}
	result := TCollectorSubmitResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

// Parameters:
//  - Spans
func (p *TCollectorClient) MultiSubmit(spans []*Span) (r []*Response, err error) {
	if err = p.sendMultiSubmit(spans); err != nil {
		return
	}
	return p.recvMultiSubmit()
}

func (p *TCollectorClient) sendMultiSubmit(spans []*Span) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("multi_submit", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := TCollectorMultiSubmitArgs{
		Spans: spans,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *TCollectorClient) recvMultiSubmit() (value []*Response, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "multi_submit" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "multi_submit failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "multi_submit failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error4 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error5 error
		error5, err = error4.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error5
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "multi_submit failed: invalid message type")
		return
	}
	result := TCollectorMultiSubmitResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type TCollectorProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      TCollector
}

func (p *TCollectorProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *TCollectorProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *TCollectorProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewTCollectorProcessor(handler TCollector) *TCollectorProcessor {

	self6 := &TCollectorProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self6.processorMap["submit"] = &tCollectorProcessorSubmit{handler: handler}
	self6.processorMap["multi_submit"] = &tCollectorProcessorMultiSubmit{handler: handler}
	return self6
}

func (p *TCollectorProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x7 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x7.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x7

}

type tCollectorProcessorSubmit struct {
	handler TCollector
}

func (p *tCollectorProcessorSubmit) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := TCollectorSubmitArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("submit", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := TCollectorSubmitResult{}
	var retval *Response
	var err2 error
	if retval, err2 = p.handler.Submit(args.Span); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing submit: "+err2.Error())
		oprot.WriteMessageBegin("submit", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("submit", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type tCollectorProcessorMultiSubmit struct {
	handler TCollector
}

func (p *tCollectorProcessorMultiSubmit) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := TCollectorMultiSubmitArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("multi_submit", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := TCollectorMultiSubmitResult{}
	var retval []*Response
	var err2 error
	if retval, err2 = p.handler.MultiSubmit(args.Spans); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing multi_submit: "+err2.Error())
		oprot.WriteMessageBegin("multi_submit", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("multi_submit", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Span
type TCollectorSubmitArgs struct {
	Span *Span `thrift:"span,1" json:"span"`
}

func NewTCollectorSubmitArgs() *TCollectorSubmitArgs {
	return &TCollectorSubmitArgs{}
}

var TCollectorSubmitArgs_Span_DEFAULT *Span

func (p *TCollectorSubmitArgs) GetSpan() *Span {
	if !p.IsSetSpan() {
		return TCollectorSubmitArgs_Span_DEFAULT
	}
	return p.Span
}
func (p *TCollectorSubmitArgs) IsSetSpan() bool {
	return p.Span != nil
}

func (p *TCollectorSubmitArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TCollectorSubmitArgs) readField1(iprot thrift.TProtocol) error {
	p.Span = &Span{}
	if err := p.Span.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Span), err)
	}
	return nil
}

func (p *TCollectorSubmitArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("submit_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TCollectorSubmitArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("span", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:span: ", p), err)
	}
	if err := p.Span.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Span), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:span: ", p), err)
	}
	return err
}

func (p *TCollectorSubmitArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TCollectorSubmitArgs(%+v)", *p)
}

// Attributes:
//  - Success
type TCollectorSubmitResult struct {
	Success *Response `thrift:"success,0" json:"success,omitempty"`
}

func NewTCollectorSubmitResult() *TCollectorSubmitResult {
	return &TCollectorSubmitResult{}
}

var TCollectorSubmitResult_Success_DEFAULT *Response

func (p *TCollectorSubmitResult) GetSuccess() *Response {
	if !p.IsSetSuccess() {
		return TCollectorSubmitResult_Success_DEFAULT
	}
	return p.Success
}
func (p *TCollectorSubmitResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *TCollectorSubmitResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TCollectorSubmitResult) readField0(iprot thrift.TProtocol) error {
	p.Success = &Response{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *TCollectorSubmitResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("submit_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TCollectorSubmitResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *TCollectorSubmitResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TCollectorSubmitResult(%+v)", *p)
}

// Attributes:
//  - Spans
type TCollectorMultiSubmitArgs struct {
	Spans []*Span `thrift:"spans,1" json:"spans"`
}

func NewTCollectorMultiSubmitArgs() *TCollectorMultiSubmitArgs {
	return &TCollectorMultiSubmitArgs{}
}

func (p *TCollectorMultiSubmitArgs) GetSpans() []*Span {
	return p.Spans
}
func (p *TCollectorMultiSubmitArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TCollectorMultiSubmitArgs) readField1(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Span, 0, size)
	p.Spans = tSlice
	for i := 0; i < size; i++ {
		_elem8 := &Span{}
		if err := _elem8.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem8), err)
		}
		p.Spans = append(p.Spans, _elem8)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *TCollectorMultiSubmitArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("multi_submit_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TCollectorMultiSubmitArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("spans", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:spans: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Spans)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Spans {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:spans: ", p), err)
	}
	return err
}

func (p *TCollectorMultiSubmitArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TCollectorMultiSubmitArgs(%+v)", *p)
}

// Attributes:
//  - Success
type TCollectorMultiSubmitResult struct {
	Success []*Response `thrift:"success,0" json:"success,omitempty"`
}

func NewTCollectorMultiSubmitResult() *TCollectorMultiSubmitResult {
	return &TCollectorMultiSubmitResult{}
}

var TCollectorMultiSubmitResult_Success_DEFAULT []*Response

func (p *TCollectorMultiSubmitResult) GetSuccess() []*Response {
	return p.Success
}
func (p *TCollectorMultiSubmitResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *TCollectorMultiSubmitResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TCollectorMultiSubmitResult) readField0(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Response, 0, size)
	p.Success = tSlice
	for i := 0; i < size; i++ {
		_elem9 := &Response{}
		if err := _elem9.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem9), err)
		}
		p.Success = append(p.Success, _elem9)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *TCollectorMultiSubmitResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("multi_submit_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TCollectorMultiSubmitResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.LIST, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Success)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Success {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *TCollectorMultiSubmitResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TCollectorMultiSubmitResult(%+v)", *p)
}