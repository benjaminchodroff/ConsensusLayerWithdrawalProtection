# CLWP Rescuer Bot

This is the CLWP rescuer bot. We do not take any responsibility for this code, and only advanced users should use it. If you wish to volunteer for CLWP, we strongly recommend using the Volunteer instructions which does not rely on this bot. This bot may be used with any beacon node, but was designed to handle beacon nodes such as Teku that cannot accept `BLSToExecutionChange` prior to Capella. You must run your own beacon node.

## Description

The bot fetches the change operations from a git repo and submits them to one or many beacon nodes. 

Edit the config.json file to set your beacon node endpoint, adjust git mainnet, polling period etc.

## Build

If you have go installed then add this folder to your GO_PATH and build it:

```
go build .
go run .
```

Or use Docker:

```
build . -t clwp.rescuer
docker run clwp.rescuer
```

