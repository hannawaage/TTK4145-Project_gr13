package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	masterID := id

	var (
		numPeers         int
		currentMsgID     int
		elev             config.Elevator
		onlineIDs        []int
		receivedReceipt  []int
		updatedAllOrders [config.NumElevs][config.NumFloors][config.NumButtons]int
		currentAllOrders [config.NumElevs][config.NumFloors][config.NumButtons]int
		allElevs         [config.NumElevs]config.Elevator
		orderTimeStamps    [config.NumFloors]int
		online           bool
		faultyElev int = -1
	)

	go func() {
		for {
			select {
			case elev = <-esmChns.Elev:
				allElevs[id] = elev
				if !online {
					if currentAllOrders[id] != elev.Orders {
						updatedAllOrders[id] = elev.Orders
						esmChns.CurrentAllOrders <- updatedAllOrders
						currentAllOrders = updatedAllOrders
					}
				}
			}
		}
	}()

	msgTimer := time.NewTimer(10 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			if id != faultyElev{
				UpdateTimeStamp(&orderTimeStamps, &currentAllOrders, &allElevs)
				if TimeStampTimeout(&orderTimeStamps) {
					go func() { syncCh.OrderTimeout <- true }()
				}
				currentMsgID = rand.Intn(256)
				msg := config.Message{elev, updatedAllOrders, currentMsgID, false, id}
				syncCh.SendChn <- msg
				msgTimer.Reset(200 * time.Millisecond)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			if id != recID {
				if (recID == faultyElev || id == faultyElev) {
					break
				}
				if !Contains(onlineIDs, recID) {
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
				if incomming.IsReceipt {
					if incomming.MsgId == currentMsgID {
						if !Contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				} else {
					allElevs[recID] = incomming.Elev
					for elevator := 0; elevator < config.NumElevs; elevator++ {
						if !Contains(onlineIDs, allElevs[elevator].Id) && (elevator != id){
							allElevs[elevator].Orders = [config.NumFloors][config.NumButtons]int{}
						}
					}
					if id == masterID {
						updatedAllOrders = CostFunction(id, allElevs, onlineIDs)
					} else if recID == masterID {
						updatedAllOrders = incomming.AllOrders
					}
					if currentAllOrders != updatedAllOrders {
						esmChns.CurrentAllOrders <- updatedAllOrders
						currentAllOrders = updatedAllOrders
					}
					msg := config.Message{elev, updatedAllOrders, incomming.MsgId, true, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		case <-msgTimer.C:
				fmt.Println("Offline")
				numPeers = 0
				onlineIDs = onlineIDs[:0]
				receivedReceipt = receivedReceipt[:0]
				masterID = id
				online = false
				updatedAllOrders = MergeAllOrders(id, updatedAllOrders)
				esmChns.CurrentAllOrders <- updatedAllOrders
				currentAllOrders = updatedAllOrders
		case timeout := <-syncCh.OrderTimeout:
            if timeout {
				faultyElev = FindFaultyElev(&currentAllOrders, &orderTimeStamps)
					fmt.Println("Faulty: ", faultyElev)
					if id != faultyElev {
						updatedAllOrders = MergeAllOrders(id, updatedAllOrders)
					}
					updatedAllOrders[faultyElev] = [config.NumFloors][config.NumButtons]int{}
					esmChns.CurrentAllOrders <- updatedAllOrders
					currentAllOrders = updatedAllOrders
				fmt.Println("Order  timeout")
				orderTimeStamps = [config.NumFloors]int{}
				numPeers = 0
				onlineIDs = onlineIDs[:0]
				receivedReceipt = receivedReceipt[:0]
				masterID = id
				online = false
            }
		}
	}
}
