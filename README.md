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
