In order to SSH into nodes you need to exec into docker containers. Let's do it.

First, we will get list of nodes by running kubectl get nodes -o wide:

NAME                 STATUS   ROLES                  AGE     VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION    CONTAINER-RUNTIME
kind-control-plane   Ready    control-plane,master   5m5s    v1.21.1   172.18.0.2    <none>        Ubuntu 21.04   5.11.0-1017-gcp   containerd://1.5.2
kind-worker          Ready    <none>                 4m38s   v1.21.1   172.18.0.4    <none>        Ubuntu 21.04   5.11.0-1017-gcp   containerd://1.5.2
kind-worker2         Ready    <none>                 4m35s   v1.21.1   172.18.0.3    <none>        Ubuntu 21.04   5.11.0-1017-gcp   containerd://1.5.2
Let's suppose we want to SSH into kind-worker node.

Now, we will get list of docker containers (docker ps -a) and check if all nodes are here:

CONTAINER ID   IMAGE                  COMMAND                  CREATED          STATUS         PORTS                       NAMES
7ee204ad5fd1   kindest/node:v1.21.1   "/usr/local/bin/entr…"   10 minutes ago   Up 8 minutes                               kind-worker
434f54087e7c   kindest/node:v1.21.1   "/usr/local/bin/entr…"   10 minutes ago   Up 8 minutes   127.0.0.1:35085->6443/tcp   kind-control-plane
2cb2e9465d18   kindest/node:v1.21.1   "/usr/local/bin/entr…"   10 minutes ago   Up 8 minutes                               kind-worker2
Take a look at the NAMES column - here are nodes names used in Kubernetes.

Now we will use standard docker exec command to connect to the running container and connect to it's shell - docker exec -it kind-worker sh, then we will run ip a on the container to check if IP address matches the address from the kubectl get nodes command:

# ls
bin  boot  dev  etc  home  kind  lib  lib32  lib64  libx32  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
...
11: eth0@if12: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 172.18.0.4/16 brd 172.18.255.255 scope global eth0
    ...
# 
As can see, we successfully connected to the node used by Kind Kubernetes - the IP address 172.18.0.4 matches the IP address from the kubectl get nodes command.