/* l2queue.go - fifo queue with two levels   */
/*
modification history
--------------------
2014/3/11, by Zhang Miao, create
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/l2queue"
    
    var q l2queue.Queue
    
    q.Init()
    
    q.Append("abcd", "high")
    q.Append("efgh", "low")
    len := q.Len()
    
    msg = q.Remove()
    // type convert is required here
    msgStr := msg.(string)

*/
package l2queue

import (
    "container/list"
    "sync"    
)

/* queue    */
type Queue struct {
    lock    sync.Mutex
    cond    *sync.Cond
    q_high  *list.List  // queue with high priority
    q_low   *list.List  // queue with low priority
}

/* Initialize the queue */
func (q *Queue) Init() {
    q.cond = sync.NewCond(&q.lock)
    q.q_high = list.New()
    q.q_low = list.New()
}

/* Add to the queue */
func (q *Queue) Append(item interface{}, level string) {    
    q.cond.L.Lock()
    
    if level == "high" {
        q.q_high.PushBack(item)
    } else {
        q.q_low.PushBack(item)
    }
    
    q.cond.Signal()
    q.cond.L.Unlock()
}

/* Remove from the queue.  
   Try to get item from queue of high level, then of low level
*/
func (q *Queue) Remove() interface{} {
    q.cond.L.Lock()

    for q.q_high.Len() == 0 && q.q_low.Len() == 0 {
        q.cond.Wait()
    }
    
    var item *list.Element
    
    if q.q_high.Len() != 0 {
        item = q.q_high.Front()
        q.q_high.Remove(item)
    } else {
        item = q.q_low.Front()
        q.q_low.Remove(item)        
    }
    
    q.cond.L.Unlock()

    return item.Value
}

/* Get length of the queue */
func (q *Queue) Len() (int, int) {    
    q.lock.Lock()
    
    len_high := q.q_high.Len()
    len_low := q.q_low.Len()
    
    q.lock.Unlock()
    
    return len_high, len_low
}
