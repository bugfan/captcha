package captcha

import (
	"flag"
	"log"
	"testing"
)

func TestVerify(*testing.T) {
	go func() {
		var addr = flag.String("addr", ":9997", "assign bind addr")
		flag.Parse()
		log.Fatal(Start(*addr))
	}()
	<-(chan int)(nil)
}
