package agent

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func mainloop() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done
}
