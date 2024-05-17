package rtlimit

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type LimitUnit struct {
	t     time.Time
	count uint64
	lock  sync.Mutex
}
type Limit struct {
	rate     uint64
	db       map[string]*LimitUnit
	ctime    time.Time
	duration time.Duration
	clean    time.Duration
	lock     sync.Mutex
}

// Create new limiter
// a: request
// b: time duration
// c: clean garbage timeout
func New(a uint64, b time.Duration, c time.Duration) Limit {
	return Limit{
		rate:     a,
		db:       make(map[string]*LimitUnit),
		ctime:    time.Now(),
		duration: b,
		clean:    c,
		lock:     sync.Mutex{},
	}
}

func Run(listen string, a uint64, b time.Duration, c time.Duration) {
	l := New(a, b, c)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if l.Check(r.FormValue("id")) {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})

	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Client(s string) bool {
	resp, err := http.Get(s)
	if err != nil {
		log.Println(err)
		return true
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

func (a *Limit) Check(s string) bool {
	t := time.Now()
	a.lock.Lock()
	r := a.rate
	d := a.duration
	if t.Sub(a.ctime) > a.clean {
		for k, v := range a.db {
			if t.Sub(v.t) > a.clean {
				delete(a.db, k)
			}
		}
		a.ctime = t
	}
	limitunit, ok := a.db[s]

	if !ok {
		a.db[s] = &LimitUnit{
			t:     t,
			count: r - 1,
			lock:  sync.Mutex{},
		}
		a.lock.Unlock()
		return true
	}
	a.lock.Unlock()

	limitunit.lock.Lock()
	if limitunit.count > 0 {
		limitunit.count--
		limitunit.lock.Unlock()
		return true
	}

	if t.Sub(limitunit.t) < d {
		limitunit.lock.Unlock()
		return false
	}

	limitunit.t = t
	limitunit.count = r - 1
	limitunit.lock.Unlock()
	return true
}
