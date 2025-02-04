/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import "mosn.io/proxy-wasm-go-host/proxywasm/common"

// Exports contains ABI that exported by wasm module.
type Exports interface {
	ProxyOnContextCreate(contextID int32, parentContextID int32) error
	ProxyOnDone(contextID int32) (int32, error)
	ProxyOnLog(contextID int32) error
	ProxyOnDelete(contextID int32) error
	ProxyOnMemoryAllocate(size int32) (int32, error)

	ProxyOnVmStart(rootContextID int32, vmConfigurationSize int32) (int32, error)
	ProxyOnConfigure(rootContextID int32, pluginConfigurationSize int32) (int32, error)

	ProxyOnTick(rootContextID int32) error

	ProxyOnNewConnection(contextID int32) (Action, error)

	ProxyOnDownstreamData(contextID int32, dataLength int32, endOfStream int32) (Action, error)
	ProxyOnDownstreamConnectionClose(contextID int32, closeType int32) error

	ProxyOnUpstreamData(contextID int32, dataLength int32, endOfStream int32) (Action, error)
	ProxyOnUpstreamConnectionClose(contextID int32, closeType int32) error

	ProxyOnRequestHeaders(contextID int32, headers int32, endOfStream int32) (Action, error)
	ProxyOnRequestBody(contextID int32, bodyBufferLength int32, endOfStream int32) (Action, error)
	ProxyOnRequestTrailers(contextID int32, trailers int32) (Action, error)
	ProxyOnRequestMetadata(contextID int32, nElements int32) (Action, error)

	ProxyOnResponseHeaders(contextID int32, headers int32, endOfStream int32) (Action, error)
	ProxyOnResponseBody(contextID int32, bodyBufferLength int32, endOfStream int32) (Action, error)
	ProxyOnResponseTrailers(contextID int32, trailers int32) (Action, error)
	ProxyOnResponseMetadata(contextID int32, nElements int32) (Action, error)

	ProxyOnHttpCallResponse(contextID int32, token int32, headers int32, bodySize int32, trailers int32) error

	ProxyOnQueueReady(rootContextID int32, token int32) error

	ProxyOnGrpcCallResponseHeaderMetadata(contextID int32, calloutID int32, nElements int32) error
	ProxyOnGrpcCallResponseMessage(contextID int32, calloutID int32, msgSize int32) error
	ProxyOnGrpcCallResponseTrailerMetadata(contextID int32, calloutID int32, nElements int32) error
	ProxyOnGrpcCallClose(contextID int32, calloutID int32, statusCode int32) error
}

type ImportsHandler interface {
	// utils
	Log(level LogLevel, msg string) WasmResult
	GetRootContextID() int32
	SetEffectiveContextID(contextID int32) WasmResult
	SetTickPeriodMilliseconds(tickPeriodMilliseconds int32) WasmResult
	GetCurrentTimeNanoseconds() (int32, WasmResult)
	Done() WasmResult

	// config
	GetVmConfig() common.IoBuffer
	GetPluginConfig() common.IoBuffer

	// metric
	DefineMetric(metricType MetricType, name string) (int32, WasmResult)
	IncrementMetric(metricID int32, offset int64) WasmResult
	RecordMetric(metricID int32, value int64) WasmResult
	GetMetric(metricID int32) (int64, WasmResult)
	RemoveMetric(metricID int32) WasmResult

	// property
	GetProperty(key string) (string, WasmResult)
	SetProperty(key string, value string) WasmResult

	// l4
	GetDownStreamData() common.IoBuffer
	GetUpstreamData() common.IoBuffer
	ResumeDownstream() WasmResult
	ResumeUpstream() WasmResult

	// http
	GetHttpRequestHeader() common.HeaderMap
	GetHttpRequestBody() common.IoBuffer
	GetHttpRequestTrailer() common.HeaderMap

	GetHttpResponseHeader() common.HeaderMap
	GetHttpResponseBody() common.IoBuffer
	GetHttpResponseTrailer() common.HeaderMap

	HttpCall(url string, headers common.HeaderMap, body common.IoBuffer, trailer common.HeaderMap, timeoutMilliseconds int32) (int32, WasmResult)
	GetHttpCallResponseHeaders() common.HeaderMap
	GetHttpCallResponseBody() common.IoBuffer
	GetHttpCallResponseTrailer() common.HeaderMap

	ResumeHttpRequest() WasmResult
	ResumeHttpResponse() WasmResult
	SendHttpResp(respCode int32, respCodeDetail common.IoBuffer, respBody common.IoBuffer, additionalHeaderMap common.HeaderMap, grpcCode int32) WasmResult

	// grpc
	OpenGrpcStream(grpcService string, serviceName string, method string) (int32, WasmResult)
	SendGrpcCallMsg(token int32, data common.IoBuffer, endOfStream int32) WasmResult
	CancelGrpcCall(token int32) WasmResult
	CloseGrpcCall(token int32) WasmResult

	GrpcCall(grpcService string, serviceName string, method string, data common.IoBuffer, timeoutMilliseconds int32) (int32, WasmResult)
	GetGrpcReceiveInitialMetaData() common.HeaderMap
	GetGrpcReceiveBuffer() common.IoBuffer
	GetGrpcReceiveTrailerMetaData() common.HeaderMap

	// foreign
	CallForeignFunction(funcName string, param []byte) ([]byte, WasmResult)
	GetFuncCallData() common.IoBuffer

	// shared
	GetSharedData(key string) (string, uint32, WasmResult)
	SetSharedData(key string, value string, cas uint32) WasmResult

	RegisterSharedQueue(queueName string) (uint32, WasmResult)
	RemoveSharedQueue(queueID uint32) WasmResult
	ResolveSharedQueue(queueName string) (uint32, WasmResult)
	EnqueueSharedQueue(queueID uint32, data string) WasmResult
	DequeueSharedQueue(queueID uint32) (string, WasmResult)

	// for golang host environment
	// Wait until async call return, eg. sync http call in golang
	Wait() Action

	// custom extension
	GetCustomBuffer(bufferType BufferType) common.IoBuffer
	GetCustomHeader(mapType MapType) common.HeaderMap

	//Runtime
	GetState(storeName string, key string) (string, WasmResult)
	InvokeService(id string, method string, param string) (string, WasmResult)
}
