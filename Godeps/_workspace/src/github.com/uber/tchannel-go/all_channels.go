// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tchannel

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

// channelMap is used to ensure that applications don't create multiple channels with
// the same service name in a single process.
var channelMap = struct {
	sync.Mutex
	existing map[string][]*Channel
}{
	existing: make(map[string][]*Channel),
}

func getCallerStack(skip int) string {
	callers := make([]uintptr, 32)
	n := runtime.Callers(skip+2 /* skip Callers and self */, callers)
	callers = callers[:n]

	buf := &bytes.Buffer{}
	for _, pc := range callers {
		f := runtime.FuncForPC(pc)
		name := f.Name()
		file, line := f.FileLine(pc)

		buf.WriteString(name)
		buf.WriteByte('\n')
		buf.WriteString("   at ")
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte('\n')
	}

	return buf.String()
}

func registerNewChannel(ch *Channel) {
	serviceName := ch.ServiceName()
	ch.createdStack = getCallerStack(1 /* skip self */)

	channelMap.Lock()
	defer channelMap.Unlock()

	existing := channelMap.existing[serviceName]
	channelMap.existing[serviceName] = append(existing, ch)
}

func removeClosedChannel(ch *Channel) {
	channelMap.Lock()
	defer channelMap.Unlock()

	channels := channelMap.existing[ch.ServiceName()]
	for i, v := range channels {
		if v != ch {
			continue
		}

		// Replace current index with the last element, and truncate channels.
		channels[i] = channels[len(channels)-1]
		channels = channels[:len(channels)-1]
		break
	}

	channelMap.existing[ch.ServiceName()] = channels
}
