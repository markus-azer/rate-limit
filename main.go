package main

import (
	"fmt"
	"strconv"
	"time"

	tokenbucket "github.com/markus-azer/rate-limit/pkg"
)

type test struct {
	store []tokenbucket.TokenBucket
}

func (h *test) FetchToken(id string) (*tokenbucket.TokenBucket, error) {
	for i := range h.store {
		if h.store[i].Id == id {
			return &h.store[i], nil
		}
	}
	return nil, nil
}

func (h *test) RemoveTokens(id string, tokens int) error {
	for i := range h.store {
		if h.store[i].Id == id {
			h.store[i].Tokens = h.store[i].Tokens - tokens
		}
	}
	return nil
}

func (h *test) CreateBucket(id string, tokens int, startTime time.Time) error {
	tb := tokenbucket.TokenBucket{Id: id, Version: 1, Tokens: tokens, StartTime: time.Now()}
	h.store = append(h.store, tb)
	return nil
}

func (h *test) ResetBucket(id string, version int, tokens int, startTime time.Time) error {
	for i := range h.store {
		if h.store[i].Id == id {
			h.store[i].Tokens = tokens // make sure max don't exceed
			h.store[i].StartTime = startTime
		}
	}
	return nil
}

func main() {

	ha := new(test)

	rt := tokenbucket.NewRateLimiter(ha, time.Duration(1000), 1)

	re, _ := rt.Take("123")
	res, _ := rt.Take("123")
	fmt.Println(re)
	fmt.Println(res)
	for i := range ha.store {
		if ha.store[i].Id == "123" {
			xyz := strconv.Itoa(ha.store[i].Tokens)
			fmt.Println(xyz)
		}
	}

	time.Sleep(1 * time.Second)
	rew, _ := rt.Take("123")
	resw, _ := rt.Take("123")
	fmt.Println(rew)
	fmt.Println(resw)
	for i := range ha.store {
		if ha.store[i].Id == "123" {
			xyz := strconv.Itoa(ha.store[i].Tokens)
			fmt.Println(xyz)
		}
	}

}
