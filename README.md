container_id=`docker create -v /var/lib/docker-volumes/cgroupfs/my-volume/meminfo:/proc/meminfo -m=15m chenchun/hello /hello`

docker volume create -d cgroupfs_volume --name my-volume

docker start $container_id