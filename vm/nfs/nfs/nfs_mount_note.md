# package install
```
sudo apt install nfs-common
```

# note on mounting nfs on worker nodes

```
mount -t nfs <ip_address>:/srv/nfs/kubedata /mnt
```

# verify the mount data

```
mount | grep kubedata
```

# to unmount the data

```
unmount /mnt
```
