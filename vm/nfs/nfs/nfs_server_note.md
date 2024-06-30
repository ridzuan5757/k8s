# preparation
```
sudo mkdir /src/nfs/kubedata
sudo chown nobody: /src/nfs/kubedata
```

# install nfs-utils
```
sudo apt install nfs-utils
sudo systemctl enable nfs-server
sudo systemctl start nfs-server
```

# modify `/etc/exports` file
```
sudo vi /etc/exports
```

```
# /etc/exports - directories exported to NFS clients
/src/nfs/kubeadata
*(rw,sync,no_subtree_check,no_root_squash,no_all_squash, insecure)
```

# exports 

```
sudo exportfs -rav
```

# verify th export data

```
sudo exportfs -v
```


