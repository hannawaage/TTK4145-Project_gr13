package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	masterID := id
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	// Heihei
	var (
		numPeers           int
		currentMsgID       int
		numTimeouts        int
		elev               config.Elevator
		onlineIDs          []int
		receivedReceipt    []int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		timeStamps         [config.NumFloors]time.Timer
		allElevs           [config.NumElevs]config.Elevator
		online             bool
	)

	go func() {
		for {
			select {
			case elev = <-esmChns.Elev:
				allElevs[id] = elev
				if !online {
					if currentAllOrders[id] != elev.Orders {
						updatedLocalOrders[id] = elev.Orders
						esmChns.CurrentAllOrders <- updatedLocalOrders
						floor := setTimeStamps(&timeStamps, &currentAllOrders, &updatedLocalOrders)
						if !(floor < 0) {
							timeStamps[floor].Reset(5 * time.Second)
							fmt.Println("timer set for floor", floor)
						} else {
							fmt.Println("no timer set for floor", floor)
						}
						currentAllOrders = updatedLocalOrders
					}
				}
			}
		}
	}()

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()
	for i := 0; i < config.NumFloors; i++ {
		timeStamps[i] = *time.NewTimer(5 * time.Second)
		timeStamps[i].Stop()
	}

	go func() {
		for {
			currentMsgID = rand.Intn(256)
			msg := config.Message{elev, updatedLocalOrders, currentMsgID, false, localIP, id}
			syncCh.SendChn <- msg
			msgTimer.Reset(800 * time.Millisecond)
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			if id != recID {
				if !contains(onlineIDs, recID) {
					onlineIDs = append(onlineIDs, recID)
					numPeers = len(onlineIDs)
					online = true
					for i := 0; i < numPeers; i++ {
						theID := onlineIDs[i]
						if theID < masterID {
							masterID = theID
						}
					}
				}
				allElevs[recID] = incomming.Elev
				if id == masterID {
					updatedLocalOrders = CostFunction(id, allElevs, onlineIDs)
				} else if recID == masterID {
					updatedLocalOrders = incomming.AllOrders
				}
				if currentAllOrders != updatedLocalOrders {
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
					floor := setTimeStamps(&timeStamps, &currentAllOrders, &updatedLocalOrders)
					if !(floor < 0) {
						timeStamps[floor].Reset(5 * time.Second)
						fmt.Println("timer set for floor", floor)
					} else {
						fmt.Println("no timer set for floor", floor)
					}
				}
				if incomming.IsReceipt {
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				} else {
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		case <-msgTimer.C:
			numTimeouts++
			if numTimeouts > 2 {
				fmt.Println("Three timeouts in a row")
				numTimeouts = 0
				numPeers = 0
				onlineIDs = onlineIDs[:0]
				receivedReceipt = receivedReceipt[:0]
				masterID = id
				online = false
				updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
				esmChns.CurrentAllOrders <- updatedLocalOrders
				currentAllOrders = updatedLocalOrders
			}
		case <-timeStamps[0].C:
			updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
			esmChns.CurrentAllOrders <- updatedLocalOrders
			currentAllOrders = updatedLocalOrders
			fmt.Println("Order timeout")
		case <-timeStamps[1].C:
			updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
			esmChns.CurrentAllOrders <- updatedLocalOrders
			currentAllOrders = updatedLocalOrders
			fmt.Println("Order timeout")
		case <-timeStamps[2].C:
			updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
			esmChns.CurrentAllOrders <- updatedLocalOrders
			currentAllOrders = updatedLocalOrders
			fmt.Println("Order timeout")
		case <-timeStamps[3].C:
			updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
			esmChns.CurrentAllOrders <- updatedLocalOrders
			currentAllOrders = updatedLocalOrders
			fmt.Println("Order timeout")
		}
	}

}

func setTimeStamps(prevTime *[config.NumFloors]time.Timer, current *[config.NumElevs][config.NumFloors][config.NumButtons]bool, updated *[config.NumElevs][config.NumFloors][config.NumButtons]bool) int {
	new := -1
	for elev := 0; elev < config.NumElevs; elev++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons; btn++ {
				if updated[elev][floor][btn] && !current[elev][floor][btn] {
					new = floor
					return new
				}
			}
		}
	}
	return new
}
