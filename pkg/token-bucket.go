package tokenbucket

import (
	"fmt"
	"time"
)

// type TokenBucket interface {

// }

// type bucket interface {
// }

//TokenBucket TokenBucket
type TokenBucket struct {
	Id        string
	Version   int
	Tokens    int
	StartTime time.Time
}

type bucketStore interface {
	FetchToken(id string) (*TokenBucket, error)
	RemoveTokens(id string, tokens int) error
	CreateBucket(id string, tokens int, startTime time.Time) error
	ResetBucket(id string, version int, tokens int, startTime time.Time) error
}

//RateLimiter RateLimiter
type RateLimiter struct {
	bs       bucketStore
	windowMs time.Duration
	max      int
}

//NewRateLimiter NewRateLimiter
func NewRateLimiter(bs bucketStore, windowMs time.Duration, max int) *RateLimiter {
	return &RateLimiter{
		bs:       bs,
		windowMs: windowMs,
		max:      max,
	}
}

//Take Take
func (rt *RateLimiter) Take(id string) (bool, error) {
	availableTokens := rt.refill(id)

	if availableTokens > 0 {
		rt.bs.RemoveTokens(id, 1)
		return true, nil
	}

	return false, nil
}

func (rt *RateLimiter) refill(id string) int {

	token, _ := rt.bs.FetchToken(id)

	//check if id exist if not create new one
	if token == nil {
		rt.bs.CreateBucket(id, rt.max, time.Now())
		return rt.max
	}

	//if now - startTime in ms  > windowMs
	//then check if so refill and substract older if exist
	diff := time.Now().Sub(token.StartTime) //token.StartTime.Sub(time.Now())

	fmt.Println(diff, rt.windowMs)
	if diff > rt.windowMs {
		fmt.Println("Here")
		tokensToBeAdded := rt.max + token.Tokens                          //todo check this with max and minus
		rt.bs.ResetBucket(id, token.Version, tokensToBeAdded, time.Now()) //token.StartTime.Add(rt.windowMs)) //handle time window
		return tokensToBeAdded
	}

	return token.Tokens
}
