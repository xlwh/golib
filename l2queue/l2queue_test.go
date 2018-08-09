/* l2queue_test.go - test for l2queue.go  */
/*
modification history
--------------------
2014/3/11, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package l2queue

import (
    "testing"
    "time"
)


func TestWorkQueue(t *testing.T) {
    var queue Queue
    queue.Init()

    queue.Append(1, "low")
    queue.Append(2, "low")
    queue.Append(3, "high")        

    
    for i := 0; i < 3; i = i + 1 {
        msg := queue.Remove()
        number := msg.(int)
               
        switch i {
            case 0:
                if number != 3 {
                    t.Error("number[0] should be 3")
                }
            case 1:
                if number != 1 {
                    t.Error("number[1] should be 1")
                }
            case 2:
                if number != 2 {
                    t.Error("number[2] should be 2")
                }
        }
    }

    time.Sleep(1 * time.Second)
}