package timer

import "time"

type TimerChns struct {
    startTimer chan int
    stopTimer chan bool
    timerTimeout chan bool
}
var timer time.Timer

func RunTimer(timerChns TimerChns) {
  for {
    select {
    case start := <- timerChn.startTimer:
      timer := time.NewTimer(start)
    }
  case <- timerChns.stopTimer:
      timer.Stop()
    case <- timer.C:
      timer.Stop()
      timerTimeout <- true
  }
}
