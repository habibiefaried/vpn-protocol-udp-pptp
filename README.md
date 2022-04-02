# vpn-protocol
Custom VPN protocol with UDP

# Docker Test

```
docker build . -f Dockerfile.test -t vpn-protocol
docker run --privileged --name vpn -dit vpn-protocol
docker exec -it vpn bash
```