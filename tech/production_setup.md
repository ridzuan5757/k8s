# Disable SWAP memory

This is needed for kubelet to function properly.

```bash
sudo swapoff -a
```

Modify `/etc/fstab`

```
/dev/disk/,,,
#/swap.img
```

# Install container runtime

```bash
sudo apt install docker.io -y
```

# Necessary application

```bash
sudo apt install apt-transport-https curl -y
```

# Add repo for installing k8s


