package remote

import (
	"www.baidu.com/golang-lib/remote/test_pb"
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	math_rand "math/rand"
)

/*
func TestNewClient(t *testing.T) {
	client := NewClient("udp", "localhost:8080", 3, NewPbClientCodec)

	request := new(test_pb.WafRequest)
	service := "news"
	request.Service = &service

	signature := "a947cf27d138c8935d743e3889dfa4d1"
	request.Signature = &signature

	clientip := uint32(12345)
	request.Clientip = &clientip
	request.Req = new(test_pb.WafRequest_Request)
	request.Req.Method = test_pb.WafRequest_Request_GET.Enum()
	request.Req.Version = test_pb.WafRequest_Request_HTTP_1_1.Enum()
	uri := "xxxxxxxxxxxxx"
	request.Req.Uri = &uri

	response := new(test_pb.WafResponse)

	var wg sync.WaitGroup
	fn := func(wg *sync.WaitGroup) {

		for i := 0; i < 3000; i++ {
			if err := client.Call(request, response, 2*time.Millisecond); err != nil {
				t.Errorf("client call failed %s", err)
			}

			//client.Call(request, response, nil)
			//t.Logf("code %v", response)

		}
		wg.Done()
	}

    wg.Add(3)
    print("start go routine\n")
	go fn(&wg)
	go fn(&wg)
	go fn(&wg)

	wg.Wait()
	// client.GoNoReturn(request)
	client.Close()
}
*/

func TestPbCall(t *testing.T) {
	popr := math_rand.New(math_rand.NewSource(616))
	request := test_pb.NewPopulatedWafRequest(popr, false)

	concurrency := 1
	client := NewClient("unix", "/tmp/dictserver.sock", 1*time.Second, concurrency,
		NewPbClientCodec, 1000)

	response := new(test_pb.WafResponse)
	time.Sleep(time.Second)
	err := client.Call(request, response, 20*time.Millisecond)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(err)
}

func BenchmarkPbCall(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	request := test_pb.NewPopulatedRequest(popr, true)

	total := 400
	concurrency := 2
	client := NewClient("unix", "/tmp/test.sock", 1*time.Second, concurrency,
		NewPbClientCodec, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		fn := func(wg *sync.WaitGroup) {

			response := new(test_pb.WafResponse)
			for i := 0; i < total/concurrency; i++ {
				if err := client.Call(request, response, 2*time.Millisecond); err != nil {
					b.Logf("client call failed %s", err)
				}

				//client.Call(request, response, nil)
				//t.Logf("code %v", response)

			}
			wg.Done()
		}

		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			go fn(&wg)
		}

		wg.Wait()
	}
}

/*
func BenchmarkGoNoReturn(b *testing.B) {
    popr := math_rand.New(math_rand.NewSource(616))
    request := test_pb.NewPopulatedRequest(popr, false)



    total := 400
    concurrency := 2
	client := NewClient("udp", "localhost:8080", 1*time.Second, concurrency,
        NewPbClientCodec, 1000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var wg sync.WaitGroup
	    fn := func(wg *sync.WaitGroup) {

		    for i := 0; i < total/concurrency; i++ {
			    client.GoNoReturn(request)
		    }
		    wg.Done()
	    }

        wg.Add(concurrency)
        for j:=0; j<concurrency; j++ {
	        go fn(&wg)
        }

	    wg.Wait()
    }

}
*/

type Header1 struct {
	Key   string
	Value string
}

type Req struct {
	Body   []byte
	Uri    string
	Method string

	H []Header1
}

type Request struct {
	R   Req
	Sig string
}

type Response struct {
	Ok bool
}

func BenchmarkMsgPackCall(b *testing.B) {
	res := new(Response)

	r := new(Request)
	r.R.Uri = "xxxxxxxxxxxxxxxxxxxx"
	r.R.Method = "Get"
	r.R.H = append(r.R.H, Header1{"abc", "def"})
	r.R.Body = bytes.Repeat([]byte("s"), 3900)
	r.Sig = "a947cf27d138c8935d743e3889dfa4d1"

	total := 400
	concurrency := 8
	client := NewClient("unix", "/tmp/test.sock", 1*time.Second, concurrency,
		nil, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		fn := func(wg *sync.WaitGroup) {

			for i := 0; i < total/concurrency; i++ {
				if err := client.Call(r, res, 2*time.Millisecond); err != nil {
					b.Logf("client call failed %s", err)
				} else {
					b.Logf("client call success")
				}

				//client.Call(request, response, nil)
				//t.Logf("code %v", response)

			}
			wg.Done()
		}

		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			go fn(&wg)
		}

		wg.Wait()
	}

}

func TestPendingCallNumber(t *testing.T) {
	concurrency := 5

	// client could not connect to server
	client := NewClient("unix", "/tmp/non-exist.sock", 1*time.Second, concurrency,
		NewPbClientCodec, 1000)

	// should not panic
	t.Log(client.PendingCallNum())
}

func prepareClient(addrInfo map[string]int32) *Client {
	// Note: address in addrInfo should be unaccessable
	clientPool := NewClientByAddr("unix", addrInfo, 1*time.Second, 1, nil, 100)

	// just mark all Sclient available
	for _, sclient := range clientPool.clients {
		sclient.Status = CONNECTED
	}

	// order clients by address
	sort.Sort(SClients(clientPool.clients))

	return clientPool
}

func updateClient(clientPool *Client, addrInfo map[string]int32) {
	clientPool.Update(addrInfo)

	for _, sclient := range clientPool.clients {
		sclient.Status = CONNECTED
	}
}

type SClients []*SClient

func (s SClients) Len() int {
	return len(s)
}

func (s SClients) Less(i, j int) bool {
	return s[i].Address < s[j].Address
}

func (s SClients) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getBalanceResult(clientPool *Client, count int) ([]string, error) {
	result := make([]string, 0)
	for i := 0; i < count; i++ {
		client, err := clientPool.balance()
		if err != nil {
			return nil, err
		}
		result = append(result, client.Address)
	}
	return result, nil
}

func TestWRR_case1(t *testing.T) {
	// prepare client
	addrInfo := map[string]int32{
		"/tmp/notexist/a": 3,
		"/tmp/notexist/b": 2,
		"/tmp/notexist/c": 1,
	}
	clientPool := prepareClient(addrInfo)

	expectResult := []string{
		"/tmp/notexist/a",
		"/tmp/notexist/b",
		"/tmp/notexist/c",
		"/tmp/notexist/a",
		"/tmp/notexist/b",
		"/tmp/notexist/a",
		"/tmp/notexist/a",
		"/tmp/notexist/b",
	}

	// check result
	actualResult, err := getBalanceResult(clientPool, len(expectResult))
	if err != nil {
		t.Errorf("should not catch error: %s", err.Error())
	}
	if !reflect.DeepEqual(expectResult, actualResult) {
		t.Errorf("Balance error (expect: %v, actual %v)", expectResult, actualResult)
	}
}

func TestWRR_case2(t *testing.T) {
	// prepare client
	addrInfo := map[string]int32{
		"/tmp/notexist/a": 3,
		"/tmp/notexist/b": 2,
		"/tmp/notexist/c": 1,
	}
	clientPool := prepareClient(addrInfo)

	// mark '/tmp/notexist/b' unavailable
	for _, client := range clientPool.clients {
		if client.Address == "/tmp/notexist/b" {
			client.Status = CONNECTING
		}
	}

	expectResult := []string{
		"/tmp/notexist/a",
		"/tmp/notexist/c",
		"/tmp/notexist/a",
		"/tmp/notexist/a",
		"/tmp/notexist/a",
		"/tmp/notexist/c",
		"/tmp/notexist/a",
		"/tmp/notexist/a",
	}

	// check result
	actualResult, err := getBalanceResult(clientPool, len(expectResult))
	if err != nil {
		t.Errorf("should not catch error: %s", err.Error())
	}
	if !reflect.DeepEqual(expectResult, actualResult) {
		t.Errorf("Balance error (expect: %v, actual %v)", expectResult, actualResult)
	}
}

func TestWRR_case3(t *testing.T) {
	// prepare client
	addrInfo := map[string]int32{
		"/tmp/notexist/a": 3,
		"/tmp/notexist/b": 2,
		"/tmp/notexist/c": 1,
	}
	clientPool := prepareClient(addrInfo)

	// update client
	addrInfo2 := map[string]int32{
		"/tmp/notexist/a": 3,
		"/tmp/notexist/c": 2,
		"/tmp/notexist/d": 2,
		"/tmp/notexist/e": -1,
	}
	updateClient(clientPool, addrInfo2)

	clients := clientPool.clients
	if len(clients) != 3 {
		t.Errorf("number of clients should be 3")
		return
	}
	sort.Sort(SClients(clientPool.clients))
	clientPool.index = 0

	expectResult := []string{
		"/tmp/notexist/a",
		"/tmp/notexist/c",
		"/tmp/notexist/d",
		"/tmp/notexist/a",
		"/tmp/notexist/c",
		"/tmp/notexist/d",
		"/tmp/notexist/a",
		"/tmp/notexist/a",
	}

	// check result
	actualResult, err := getBalanceResult(clientPool, len(expectResult))
	if err != nil {
		t.Errorf("should not catch error: %s", err.Error())
	}
	if !reflect.DeepEqual(expectResult, actualResult) {
		t.Errorf("Balance error (expect: %v, actual %v)", expectResult, actualResult)
	}
}
