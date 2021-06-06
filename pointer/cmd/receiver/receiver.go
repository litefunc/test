package main

import (
	"sync"
	"test/logger"
	"time"
)

type S1 struct {
	A int
	B *int
}

func (s *S1) SetA(i int) {
	s.A = i
	logger.Debug(s)
	logger.Debugf(`%p`, s)
}

func (s S1) SetA1(i int) {
	s.A = i
	logger.Debug(s)
	logger.Debugf(`%v`, s)
}

func SetA(s S1, i int) {
	s.A = i
	logger.Debugf(`%p`, &s)
}

func main() {
	var s1 S1
	s2 := &s1
	s3 := *s2
	s4 := *s2
	s5 := s2

	s1.SetA(1)
	s2.SetA1(2)

	SetA(s1, 2)

	logger.Debug(s1, s2)
	logger.Debugf(`%p, %p, %p, %p, %p, %p`, &s1, s2, &s3, &s4, s5, s1.B)

	now := time.Now()
	var wc sync.WaitGroup
	for i := 0; i < 10; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()
			s2.sleep()
			logger.Trace(i)
		}(i)
	}
	wc.Wait()

	logger.Debug(time.Since(now))

}

func (rec S1) sleep() {
	time.Sleep(time.Second * 1)
}
