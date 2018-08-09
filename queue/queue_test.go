/* queue_test.go - test for queue.go  */
/*
modification history
--------------------
2014/3/10, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package queue

import (
    "fmt"
    "testing"
    "time"
)

func consumer(queue *Queue) {
    for {
        number := queue.Remove()
        fmt.Println("Read from queue: ", number)
    }
}

func producer(queue *Queue, t *testing.T) {
    for i := 0; i < 10; i = i + 1 {
        retVal := queue.Append(i)
        if retVal != nil {
            t.Error("queue.Append() should return nil")
        }        
        fmt.Println("write to queue: ", i)
    }
}

func TestSendQueue(t *testing.T) {
    var queue Queue
    queue.Init()
        
    go consumer(&queue)
    go producer(&queue, t)

    time.Sleep(2 * time.Second)
}

func TestQueueIsFull(t *testing.T) {
    var queue Queue
    queue.Init()
    queue.SetMaxLen(3)
        
    for i := 0; i < 10; i = i + 1 {
        retVal := queue.Append(i)
        
        if i < 3 {
            if retVal != nil {
                t.Error("queue.Append() should return nil")
            }
        } else {
            if retVal == nil {
                t.Error("queue.Append() should return error")
            }            
        }
    }    
}
