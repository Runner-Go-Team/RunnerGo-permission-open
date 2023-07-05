FROM runnergo/debian:stable-slim

ADD  permission  /data/permission/permission


CMD ["/data/permission/permission","-m", "1"]
