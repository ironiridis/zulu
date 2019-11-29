# zulu
 Keep the clock correct on your old AV gear.

## Install
    go get github.com/ironiridis/zulu

## Example
Single device at 192.168.0.1, unauthenticated, update once an hour until error or terminated

    zulu -driver crestron2seriesctp 192.168.0.1

Three devices at 192.168.0.2-4, two with a password, quit after setting time

    zulu -driver crestron2seriesctp -oneshot 192.168.0.2 :secret@192.168.0.3 :letmein@192.168.0.4

Set time every 4 hours, ignore errors, connect on a non-default port

    zulu -driver crestron2seriesctp -permissive -rate 4h 192.168.0.5:8080

## Why
`zulu` is a modular application intended for managing the internal clock on devices as they age. The initial target of this is Crestron 2-series processors (discontinued September 2015) as these often rely on some time-of-day based logic, yet do not have any means for synchronizing their clock to an authoritative time source. These devices often run perfectly well otherwise, but experience clock drift and RTC reset on power loss. `zulu` gives these devices additional usable lifetime.

## Usage
The only mandatory flag is `-driver`. Use `-driver list` or run with no arguments to see the available drivers. At this time the only included driver is `crestron2seriesctp`, but I intend to write more (and pull requests are of course welcome).

Invoke using `zulu -driver {driver} [flags] [user:password@host:port]...`

### Flags
* `oneshot` Exit after setting time once on all specified hosts.
* `permissive` Continue (instead of exiting) after printing any errors to console.
* `sequential` Instead of trying to connect to as many hosts as possible, connect to each host individually in order and wait for the operation to complete before continuing.
* `rate {duration}` Specify the duration to sleep after all devices have been attempted before looping. Ignored if `-oneshot` is specified.
* `driver {driver}` Choose the device and connection types for all hosts specified in this invocation.

### Hosts
In general, `zulu` will be using TCP-IP to connect to each device. Unless otherwise specified by a driver, hosts are in the traditional ssh-style format:

    [[user][:password]@]hostname-or-ip[:non-default-port]

Only the hostname/ip is required, and hostnames should be resolvable given an appropriately configured network. Drivers will supply any unspecified parameters with typical defaults. Hence, these are all (theoretically) valid:

* `192.168.0.2`
* `192.168.0.3:8080`
* `user@192.168.0.4`
* `user:password@192.168.0.5`
* `:password@192.168.0.6`
* `ExamplePRO2`
* `:secret@ExampleCP2e:7777`

### Drivers
Drivers will generally be as simple as possible to accomplish the task. Robust error handling or recovery is not a priority in this application. Some drivers will be OS-specific (eg any drivers written with RS232 or USB connectivity in mind).

* `crestron2seriesctp` connects to the ethernet-enabled Crestron processors of the 2-series generation (ie AV2/PRO2/PAC2/RACK2 with a C2ENET card, or CNX-DVP4/CP2e/DIN-AP2/MC2e/MP2e/PAC2M/QM-RMC/QM-RMCRX[-BA], and others). Tested with firmware 4.008. Might work with other basic Crestron CTP devices, like the CEN-TVAV.

Drivers on my radar but not implemented yet (lacking hardware for the moment):

* `crestron2series232` connects to any Crestron 2-series processor with a basic serial console port (in particular this excludes the PAC2M and the DIN-AP2)
* `amxnix000` connects to most NI-x000 series controllers from AMX (eg NI-2000/NI-3000/NI-4000/NI-700) over IP

