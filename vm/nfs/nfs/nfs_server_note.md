# preparation
```
sudo mkdir /src/nfs/kubedata
sudo chown nobody:nogroup /src/nfs/kubedata
```

# install nfs-utils
```
pacman install nfs-utils                  # for arch

sudo apt install nfs-kernel-server
sudo systemctl enable nfs-kernel-server
sudo systemctl start nfs-kernel-server
```

# modify `/etc/exports` file
```
sudo vi /etc/exports
```

```
# /etc/exports - directories exported to NFS clients
/src/nfs/kubeadata      *(rw,sync,no_subtree_check)
```

```
sudo systemctl restart nfs-kernel-server
```

# exports 

```
sudo exportfs -rav
```

# verify th export data

```
sudo exportfs -v
```


