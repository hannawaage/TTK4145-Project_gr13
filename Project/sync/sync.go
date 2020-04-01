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

	var (
		numPeers           int
		elev               config.Elevator
		onlineIPs          []int
		receivedReceipt    []int
		currentMsgID       int
		numTimeouts        int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		allElevs           [config.NumElevs]config.Elevator
		masterAck          bool
	)

	go func() {
		for {
			select {
			case elev = <-esmChns.Elev:
				if updatedLocalOrders[id] != elev.Orders {
					updatedLocalOrders[id] = elev.Orders
				}
				allElevs[id] = elev
				masterAck = false
			}
		}
	}()

	go func() {
		for {
			if currentAllOrders != updatedLocalOrders {
				if !(len(onlineIPs) > 0) {
					updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
				} else {
					if newCabOrdersOnly(id, &currentAllOrders, &updatedLocalOrders) {
						esmChns.CurrentAllOrders <- updatedLocalOrders
						currentAllOrders = updatedLocalOrders
					}
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

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
				if !contains(onlineIPs, recID) {
					onlineIPs = append(onlineIPs, recID)
					numPeers = len(onlineIPs)
					for i := 0; i < numPeers; i++ {
						theID := onlineIPs[i]
						if theID < masterID {
							masterID = theID
						}
					}
				}
				if len(onlineIPs) == numPeers {
					if id == masterID {
						allElevs[recID] = incomming.Elev
						allElevs[recID].Orders = incomming.AllOrders[recID]
						updatedLocalOrders = CostFunction(id, allElevs, onlineIPs)
						if currentAllOrders != updatedLocalOrders {
							esmChns.CurrentAllOrders <- updatedLocalOrders
							currentAllOrders = updatedLocalOrders
						}
					} else if recID == masterID {
						if (currentAllOrders != updatedLocalOrders) && !masterAck {
							updatedLocalOrders = mergeLocalOrders(id, &elev.Orders, updatedLocalOrders)
						} else {
							updatedLocalOrders = incomming.AllOrders
							if currentAllOrders != updatedLocalOrders {
								esmChns.CurrentAllOrders <- updatedLocalOrders
								currentAllOrders = updatedLocalOrders
							}
						}
					}
				}
				if !incomming.Receipt {
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				} else {
					if incomming.MsgId == currentMsgID {
						if recID == masterID {
							masterAck = true
						}
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				}
			}
		case <-msgTimer.C:
			numTimeouts++
			if numTimeouts > 2 {
				syncCh.Online <- false
				fmt.Println("Three timeouts in a row")
				numTimeouts = 0
				onlineIPs = onlineIPs[:0]
				masterID = id
				receivedReceipt = receivedReceipt[:0]
				numPeers = 0
			}
		}
	}
}
