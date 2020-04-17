package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
	"../network/localip"
)

const (
	NumElevs   = config.NumElevs
	NumFloors  = config.NumFloors
	NumButtons = config.NumButtons
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	masterID := id
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	var (
		numPeers           int
		currentMsgID       int
		numTimeouts        int
		elev               config.Elevator
		onlineIDs          []int
		receivedReceipt    []int
		updatedLocalOrders [NumElevs][NumFloors][NumButtons]bool
		currentAllOrders   [NumElevs][NumFloors][NumButtons]bool
		orderTimeStamps    [NumFloors]int
		allElevs           [NumElevs]config.Elevator
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
						currentAllOrders = updatedLocalOrders
					}
				}
			}
		}
	}()

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			updateTimeStamp(&orderTimeStamps, &currentAllOrders, &updatedLocalOrders)
			if TimeStampTimeout(&orderTimeStamps) {
				go func() { syncCh.OrderTimeout <- true }()
			}
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
		case timeout := <-syncCh.OrderTimeout:
			if timeout {
				updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
				esmChns.CurrentAllOrders <- updatedLocalOrders
				currentAllOrders = updatedLocalOrders
				elev.Orders = currentAllOrders[id]
				fmt.Println("Order  timeout")
				orderTimeStamps = [config.NumFloors]int{}
			}
		}
	}
}
