TTK4145 Elevator Project - spring 2020
======================================

Description
-----------
The goal of the Elevator project is to create software for controlling `n` elevators working in parallel across `m` floors. Where the main requirements can be summarized in bullet points.

- No orders are lost
- Multiple elevators should be more efficient than one
- An individual elevator should behave sensibly and efficiently
- The lights and buttons should function as expected

Our project is written for `n = 3` elevators, and `m == 4` floors. Where we have tried to avoid hard-coding these values: meaning that we are able to add a fourth elevator with no extra configuration, or change the number of floors with minimal configuration. But we have not tested for `n > 3` elevators and `m != 4` floors.
We have also decided how to handle some unspecified behaviour:

- Which orders are cleared when stopping at a floor?
 - Where we assume that everyone enters/exits the elevator when the door opens, meaning that we clear all orders, for that elevator, at that floor.

- How the elevator behaves when it cannot connect to the network (router) during initialization
 - We decided that our elevators should enter a "single-elevator" mode when this happens.

- How the hall (call up, call down) buttons work when the elevator is disconnected from the network
 - We choose to refuse to take these new orders.

- Stop button & obstruction switch can be disabled
  - We kept them disabled.

Lastly we also have that the following assumptions always will be true during testing:
  - At least one elevator is always working normally
  - No multiple simultaneous errors: Only one error happens at a time, but the system must still return to a fail-safe state after this error
    - (A network packet loss is *not* an error in this context, and must be considered regardless of any other (single) error that can occur)
  - No network partitioning: There will never be a situation where there are multiple sets of two or more elevators with no connection between them


Go to [TTK4145 - Project](https://github.com/TTK4145/Project) to read the full description of the project.

Our solution
------------
In this project we have solved the problem described under "Elevator Project". Our solution is based on a master-slave design. Where the elevators/nodes use UDP broadcasting to communicate. The project was written in `Google GO`.

Disclaimer
----------

The following modules/drivers where entirely copied from the TTK4145 resources:

- The [network module](https://github.com/maghauke/TTK4145-Project_gr13/tree/master/Project/network) can be found [here](https://github.com/TTK4145/Network-go).
- The [elevator driver](https://github.com/maghauke/TTK4145-Project_gr13/tree/master/Project/driver-go) can be found [here](https://github.com/TTK4145/driver-go).
