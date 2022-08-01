package pubsub // import "github.com/docker/docker/pkg/pubsub"

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSendToOneSub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	c := p.Subscribe()

	p.Publish("hi")

	msg := <-c
	if msg.(string) != "hi" {
		t.Fatalf("expected message hi but received %v", msg)
	}
}

type Server struct {
	Bucket map[string]*Publisher
	m      sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Bucket: make(map[string]*Publisher),
	}
}

func (s *Server) sub(subName string) {

	p := s.Bucket[subName]

	if p == nil {
		p = NewPublisher(100*time.Millisecond, 10)
		s.Bucket[subName] = p
	}

	c := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok && strings.Contains(s, subName) {
			return true
		}
		return false
	})

	p.Publish("hello 123")
	p.Publish("hello golang")

	for {
		select {
		case msg := <-c:
			fmt.Println(msg)
		}
	}
}

func (s *Server) pub(subName string, v any) {
	p := s.Bucket[subName]
	if p == nil {
		p = NewPublisher(100*time.Millisecond, 10)
		s.Bucket[subName] = p
	}

	p.Publish(fmt.Sprintf("hello %d", v))
}

func TestSub(t *testing.T) {
	s := NewServer()
	go s.sub("hello")

	var i int
	i = 0
	for {
		i++
		s.pub("hello", i)
	}

	//p := NewPublisher(100*time.Millisecond, 10)
	//defer p.Close()
	//c := p.SubscribeTopic(func(v interface{}) bool {
	//	if s, ok := v.(string); ok && strings.Contains(s, "hello") {
	//		return true
	//	}
	//	return false
	//})
	//p.Publish("hello world")
	//p.Publish("hello golang")
	//
	//for {
	//	select {
	//	case msg := <-c:
	//		fmt.Println(msg)
	//	default:
	//		break
	//	}
	//}
}

func TestPub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()
	_ = p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok && strings.Contains(s, "hello") {
			return true
		}
		return false
	})

	p.Publish("hello 123")
	//p.Publish("hello world")
	//p.Publish("hello golang")
	//
	//for {
	//	select {
	//	case msg := <-c:
	//		fmt.Println(msg)
	//	default:
	//		break
	//	}
	//}

	//go func() {
	//	for v := range c {
	//		fmt.Println("golang subscribe: ", v)
	//	}
	//}()
}

func TestSendToMultipleSubs(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	var subs []chan interface{}
	subs = append(subs, p.Subscribe(), p.Subscribe(), p.Subscribe())

	p.Publish("hi")

	for _, c := range subs {
		msg := <-c
		if msg.(string) != "hi" {
			t.Fatalf("expected message hi but received %v", msg)
		}
	}
}

func TestEvictOneSub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	s1 := p.Subscribe()
	s2 := p.Subscribe()

	p.Evict(s1)
	p.Publish("hi")
	if _, ok := <-s1; ok {
		t.Fatal("expected s1 to not receive the published message")
	}

	msg := <-s2
	if msg.(string) != "hi" {
		t.Fatalf("expected message hi but received %v", msg)
	}
}

func TestClosePublisher(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	var subs []chan interface{}
	subs = append(subs, p.Subscribe(), p.Subscribe(), p.Subscribe())
	p.Close()

	for _, c := range subs {
		if _, ok := <-c; ok {
			t.Fatal("expected all subscriber channels to be closed")
		}
	}
}

const sampleText = "test"

type testSubscriber struct {
	dataCh chan interface{}
	ch     chan error
}

func (s *testSubscriber) Wait() error {
	return <-s.ch
}

func newTestSubscriber(p *Publisher) *testSubscriber {
	ts := &testSubscriber{
		dataCh: p.Subscribe(),
		ch:     make(chan error),
	}
	go func() {
		for data := range ts.dataCh {
			s, ok := data.(string)
			if !ok {
				ts.ch <- fmt.Errorf("Unexpected type %T", data)
				break
			}
			if s != sampleText {
				ts.ch <- fmt.Errorf("Unexpected text %s", s)
				break
			}
		}
		close(ts.ch)
	}()
	return ts
}

// for testing with -race
func TestPubSubRace(t *testing.T) {
	p := NewPublisher(0, 1024)
	var subs []*testSubscriber
	for j := 0; j < 50; j++ {
		subs = append(subs, newTestSubscriber(p))
	}
	for j := 0; j < 1000; j++ {
		p.Publish(sampleText)
	}
	time.AfterFunc(1*time.Second, func() {
		for _, s := range subs {
			p.Evict(s.dataCh)
		}
	})
	for _, s := range subs {
		s.Wait()
	}
}

func BenchmarkPubSub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		p := NewPublisher(0, 1024)
		var subs []*testSubscriber
		for j := 0; j < 50; j++ {
			subs = append(subs, newTestSubscriber(p))
		}
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			p.Publish(sampleText)
		}
		time.AfterFunc(1*time.Second, func() {
			for _, s := range subs {
				p.Evict(s.dataCh)
			}
		})
		for _, s := range subs {
			if err := s.Wait(); err != nil {
				b.Fatal(err)
			}
		}
	}
}
