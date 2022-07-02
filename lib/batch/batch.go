package batch

import (
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {

	usersChan := make(chan user)

	// Prepare async results collector
	var wgCollector sync.WaitGroup
	wgCollector.Add(1)
	go func() {
		defer wgCollector.Done()
		for {
			if v, ok := <-usersChan; ok {
				res = append(res, v)
			} else {
				break
			}
		}
	}()

	// Run queries batch within pool
	erg := new(errgroup.Group)
	erg.SetLimit(int(pool))
	var i int64
	for i = 0; i < n; i++ {
		id := i
		erg.Go(func() error {
			usersChan <- getOne(id)
			return nil
		})
	}

	erg.Wait()

	close(usersChan)

	wgCollector.Wait()
	return res
}
