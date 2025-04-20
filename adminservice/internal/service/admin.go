package service

import (
	"fmt"
	"net/http"
	"time"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dbhsgfsdhgshg"))
}

func Slow(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte(fmt.Sprintf("all done.\n")))

}
