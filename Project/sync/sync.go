package sync

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id string, syncCh config.SyncChns, esmChns config.EsmChns) {
	const numPeers = config.NumElevs - 1
	idDig, _ := strconv.Atoi(id)
	idDig--
	masterID := idDig
	var (
		elev               config.Elevator
		onlineIPs          []string
		receivedReceipt    []string
		currentMsgID       int
		numTimeouts        int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		online             bool
		allElevs           [config.NumElevs]config.Elevator
	)
	go func() {
		for {
			select {
			case b := <-syncCh.Online:
				if b {
					online = true
					fmt.Println("Yaho, we are online!")
				} else {
					online = false
					fmt.Println("Boo, we are offline.")

				}
			case elev := <-esmChns.Elev:
				if updatedLocalOrders[idDig] != elev.Orders {
					if online {
						updatedLocalOrders = mergeLocalOrders(idDig, updatedLocalOrders, elev.Orders)
					} else {
						updatedLocalOrders[idDig] = elev.Orders
					}
				}
				allElevs[idDig] = elev
			}
		}
	}()

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	go func() {
		for {
			if currentAllOrders != updatedLocalOrders {
				if !online {
					updatedLocalOrders = mergeAllOrders(idDig, updatedLocalOrders)
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
				}
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()
	go func() {
		for {
			currentMsgID = rand.Intn(256)
			if masterID == idDig && online {
				if !updatedLocalOrders[0][0][0] {
					fmt.Println("Dette sendes fra master nå")
					fmt.Println(updatedLocalOrders)
				}
			}
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
			recIDDig, _ := strconv.Atoi(recID)
			recIDDig--
			if id != recID { //Hvis det ikke er fra oss selv, BYTTES TIL IP VED KJØRING PÅ FORSKJELLIGE MASKINER
				if !contains(onlineIPs, recID) {
					// Dersom heisen enda ikke er registrert, sjekker vi om vi nå er online og sjekker om vi er master
					onlineIPs = append(onlineIPs, recID)
					if len(onlineIPs) == numPeers {
						syncCh.Online <- true
						for i := 0; i < numPeers; i++ {
							theID, _ := strconv.Atoi(onlineIPs[i])
							theID--
							if theID < masterID {
								masterID = theID
							}
						}
					}
				}
				if !incomming.Receipt {
					if online {
						allElevs[recIDDig] = incomming.Elev
						allElevs[recIDDig].Orders = incomming.AllOrders[recIDDig]
						/*
							if (idDig == 2) && !incomming.AllOrders[0][0][0] {
								fmt.Println("Incomming for master er")
								fmt.Println(incomming.AllOrders[0])
								fmt.Println("Current for master er")
								fmt.Println(currentAllOrders[0])
							}*/
						if currentAllOrders != incomming.AllOrders {
							// Hvis vi mottar noe nytt
							if idDig == masterID {
								// Hvis jeg er master
								updatedLocalOrders = CostFunction(allElevs)
								// Lokale endringer tas med i elev uansett
							} else if recIDDig == masterID {
								updatedLocalOrders = incomming.AllOrders
								if currentAllOrders != updatedLocalOrders {
									esmChns.CurrentAllOrders <- updatedLocalOrders
									currentAllOrders = updatedLocalOrders
								}
							}
						}
					}
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				} else {
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
								if currentAllOrders != updatedLocalOrders {
									if masterID == idDig {
										esmChns.CurrentAllOrders <- updatedLocalOrders
										currentAllOrders = updatedLocalOrders
									}
								}
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
				masterID = idDig
			}
		}
	}
}

func contains(elevs []string, str string) bool {
	for _, a := range elevs {
		if a == str {
			return true
		}
	}
	return false
}

func costfcn(id int, current [config.NumElevs][config.NumFloors][config.NumButtons]bool, new [config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	current[id] = new
	allOrderMat := mergeAllOrders(0, current)
	return allOrderMat
}

func mergeAllOrders(id int, all [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged[id] = all[id]
	for elev := 0; elev < config.NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons; btn++ {
				if all[elev][floor][btn] && btn != config.NumButtons-1 {
					merged[id][floor][btn] = true
					merged[elev][floor][btn] = false
				}
			}
		}
	}
	return merged
}

func mergeLocalOrders(id int, inc [config.NumElevs][config.NumFloors][config.NumButtons]bool, local [config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged = inc
	for floor := 0; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if local[floor][btn] {
				merged[id][floor][btn] = true
			}
		}
	}
	return merged
}
