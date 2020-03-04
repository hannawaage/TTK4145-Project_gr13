package timer

import "time"

//const DoorOpenTime = 3000 * time.Millisecond

type TimerChns struct {
    StartTimer chan int
    StopTimer chan bool
    TimerTimeout chan bool
}



func RunTimer(timerChns TimerChns) {
  var timer time.Timer
  for {
    select {
    case start := <- timerChns.StartTimer:
      timer := time.NewTimer(time.Duration(start))
      _ = timer
    case <- timerChns.StopTimer:
        timer.Stop()
    case <- timer.C:
      timer.Stop()
      timerChns.TimerTimeout <- true
    }
  }

}
