/* queue.go - fifo queue    */
/*
modification history
--------------------
2014/3/10, by Zhang Miao, create
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/queue"
    
    var q queue.Queue
    
    q.Init()
    
    q.Append("abcd")
    len := q.Len()
    
    msg = q.Remove()
    // type convert is required here
    msgStr := msg.(string)
*/
package queue

import (
    "container/list"
    "errors"
    "sync"    
)

/* queue    */
type Queue struct {
    lock    sync.Mutex
    cond    *sync.Cond
    queue   *list.List 
    maxLen  int             // max queue length   
}

/* Initialize the queue */
func (q *Queue) Init() {
    q.cond = sync.NewCond(&q.lock)
    q.queue = list.New()
    q.maxLen = -1
}

/* set max queue length */
func (q *Queue) SetMaxLen(maxLen int) {
    q.lock.Lock()
    q.maxLen = maxLen
    q.lock.Unlock()    
}

/* Add to the queue */
func (q *Queue) Append(item interface{}) error {
    var err error
    
    q.cond.L.Lock()
    
    if q.maxLen != -1 && q.queue.Len() >= q.maxLen {
        err = errors.New("Queue is full")
    } else {
        q.queue.PushBack(item)
        q.cond.Signal()
        err = nil
    }
    
    q.cond.L.Unlock()
    return err
}

/* Remove from the queue */
func (q *Queue) Remove() interface{} {
    q.cond.L.Lock()

    for q.queue.Len() == 0 {
        q.cond.Wait()
    }
        
    item := q.queue.Front()
    q.queue.Remove(item)
    
    q.cond.L.Unlock()

    return item.Value
}

/* Get length of the queue */
func (q *Queue) Len() int {
    var len int
    
    q.lock.Lock()
    len = q.queue.Len()
    q.lock.Unlock()
    
    return len
}
