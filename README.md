# Lotus

Lotus is a fork of [Tulip](https://github.com/OpenAttackDefenseTools/tulip), used by Team Asia for ICC 2023

## Configuration
Before starting the stack, edit `services/configurations.py`:

```
services = [{"ip": "10.60.3.1", "port": 3003, "name": "closedsea"},
            {"ip": "10.61.3.1", "port": 1337, "name": "rpn"},
]
```

You can also edit this during the CTF, just rebuild the `api` service:
```
docker-compose up --build -d api
```

## Usage

The stack can be started with docker-compose, after creating an `.env` file. See `.env.example` as an example of how to configure your environment.
```
cp .env.example .env
# < Edit the .env file with your favourite text editor >
docker-compose up -d --build
```
To ingest traffic, it is recommended to create a shared bind mount with the docker-compose. One convenient way to set this up is as follows:
1. On the vulnbox, start a rotating packet sniffer (e.g. tcpdump, suricata, ...) (tcpdump example given below)
2. Using rsync, copy complete captures to the machine running tulip (e.g. to /traffic)
3. Add a bind to the assembler service so it can read /traffic

The ingestor will use inotify to watch for new pcap's and suricata logs. No need to set a chron job.


## tcpdump
First, you must know the network interface that the challenge is using. In ICC 2023, the network interface is named "game". You can check the network interfaces with:
```
ip addr
```

Then run the traffic sniffer with tcpdump:
```
tcpdump -i <interface> 'port <port>' -Z <user> -G 60 -w <filename>_%H_%M.pcap
```

Usually, i name the file captures/<service_name> to easily track the traces. can also add port.

## rsync
run
```
rsync -a <vm server>:captures/ <cwd>/services/traffic
```
in a 1 minute cron job

Todo: Find a way to sync consistently without cronjob, but tbh cronjob works fine

### Metadata
Tags are read from the metadata field of a rule. For example, here's a simple rule to detect a path traversal:
```
alert tcp any any -> any any (msg: "Path Traversal-../"; flow:to_server; content: "../"; metadata: tag path_traversal; sid:1; rev: 1;)
```
Once this rule is seen in traffic, the `path_traversal` tag will automatically be added to the filters in Tulip.


### eve.json
Suricata alerts are read directly from the `eve.json` file. Because this file can get quite verbose when all extensions are enabled, it is recommended to strip the config down a fair bit. For example:
```yaml
# ...
  - eve-log:
      enabled: yes
      filetype: regular #regular|syslog|unix_dgram|unix_stream|redis
      filename: eve.json
      pcap-file: false
      community-id: false
      community-id-seed: 0
      types:
        - alert:
            metadata: yes
            # Enable the logging of tagged packets for rules using the
            # "tag" keyword.
            tagged-packets: yes
# ...
```

Sessions with matched alerts will be highlighted in the front-end and include which rule was matched.

# Security
Your Tulip instance will probably contain sensitive CTF information, like flags stolen from your machines. If you expose it to the internet and other people find it, you risk losing additional flags. It is recommended to host it on an internal network (for instance behind a VPN) or to put Tulip behind some form of authentication.

# Contributing
If you have an idea for a new feature, bug fixes, UX improvements, or other contributions, feel free to open a pull request or create an issue!      
When opening a pull request, please target the `devel` branch.

# Credits
Lotus improvements and deployment was handled by [Zafirr](https://github.com/zafirr31), thanks to the creators of [Tulip](https://github.com/OpenAttackDefenseTools/tulip) for opensourcing their fork of [flower](https://github.com/secgroup/flower).

# TODO
QoL things to add

* Improve frontend
    * Add more options for filtering
* Add login page for security
* Get a new VPS to run the service
* Figure out command to sniff packets on windows server
* for demo 3, we need to check if the flag can have different formats
* Add tag for heap address, libc address, etc