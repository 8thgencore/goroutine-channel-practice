package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type result struct {
	msg string
	err error
}

type search func() *result
type replicas []search

func fakeSearch(kind string) search {
	return func() *result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return &result{
			msg: fmt.Sprintf("%q result", kind),
		}
	}
}

func getFirstResult(ctx context.Context, replicas replicas) *result {
	if len(replicas) == 0 {
		return nil
	}

	resultChan := make(chan *result, 1)

	for _, replica := range replicas {
		go func(replica search) {
			select {
			case resultChan <- replica(): // Отправляем результат в канал
			case <-ctx.Done(): // Если контекст завершен, завершаем
			}
		}(replica)
	}

	select {
	case res := <-resultChan: // Возвращаем первый полученный результат
		return res
	case <-ctx.Done(): // Если контекст завершился раньше, возвращаем ошибку
		return &result{err: ctx.Err()}
	}
}

func getResults(ctx context.Context, replicaKinds []replicas) []*result {
	results := make([]*result, len(replicaKinds))

	var wg sync.WaitGroup

	for i, searchKind := range replicaKinds {
		wg.Add(1)
		go func(i int, searchKind replicas) {
			defer wg.Done()
			results[i] = getFirstResult(ctx, searchKind)
		}(i, searchKind)
	}

	wg.Wait()

	return results
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	replicaKinds := []replicas{
		{fakeSearch("web1"), fakeSearch("web2")},
		{fakeSearch("image1"), fakeSearch("image2")},
		{fakeSearch("video1"), fakeSearch("video2")},
	}

	for _, res := range getResults(ctx, replicaKinds) {
		fmt.Println(res.msg, res.err)
	}
}
