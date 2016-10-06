```
root@worker1:/home/docker# docker volume create -d cgroupfs_volume --name myvolume -o cidfile=/tmp/containerid
myvolume
root@worker1:/home/docker# docker volume ls
DRIVER              VOLUME NAME
cgroupfs_volume     myvolume
root@worker1:/home/docker# docker run -d --cidfile /tmp/containerid --volume-driver=cgroupfs_volume -v myvolume:/proc/meminfo -m=15m chenchun/hello /hello
0034bcc32c73655c6dc5ff2220194c44a62a0a71d9b1f0133fe335103716af11
root@worker1:/home/docker# docker exec 0034 cat /proc/meminfo
MemTotal:       15360 kB
MemFree:        14552 kB
MemAvailable:   14552 kB
Buffers:        0 kB
Cached:         4 kB
SwapCached:     0 kB
root@worker1:/home/docker# cat /proc/mounts | grep myvolume
cgroupfs /var/lib/docker-volumes/cgroupfs/myvolume fuse.cgroupfs rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other 0 0
```