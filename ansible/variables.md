# Variables: Facts and Magic Variables

We can retrieve or discover certain variables containing information about the
remote system or about Ansible itself. Variables related to remote system are
called `facts`. With facts, we can use the behaviour or state of one system as a
configuration on other systems.

For example, we can use IP address of one system as a configuration value on
another system. Variables related to Ansible are called `magic variables`.

## Ansible facts

Ansible facts are data related to the remote systems, including:
- Operating systems
- IP addresses
- Attached filesystems
- etc.

We can access this data in the `ansible_facts` variable. By default, we can also
access some Ansible facts as top-level variables with the `ansible_` prefix. We
can disable this behaviour using this `INJECT_FACTS_AS_VARS` settings. To see
all available facts, add this task to a play:

```yaml
- name: Print all available facts
  ansible.builtin.debug:
    var: ansible_facts
```

To see the raw information as gathered, run this command at the command line:

```bash
ansible <hostname> -m ansible.builtin.setup
```

Facts include a large amount of variable data, which may look like this:

```javascript
{
    "ansible_all_ipv4_addresses": [
        "REDACTED IP ADDRESS"
    ],
    "ansible_all_ipv6_addresses": [
        "REDACTED IPV6 ADDRESS"
    ],
    "ansible_apparmor": {
        "status": "disabled"
    },
    "ansible_architecture": "x86_64",
    "ansible_bios_date": "11/28/2013",
    "ansible_bios_version": "4.1.5",
    "ansible_cmdline": {
        "BOOT_IMAGE": "/boot/vmlinuz-3.10.0-862.14.4.el7.x86_64",
        "console": "ttyS0,115200",
        "no_timer_check": true,
        "nofb": true,
        "nomodeset": true,
        "ro": true,
        "root": "LABEL=cloudimg-rootfs",
        "vga": "normal"
    },
    "ansible_date_time": {
        "date": "2018-10-25",
        "day": "25",
        "epoch": "1540469324",
        "hour": "12",
        "iso8601": "2018-10-25T12:08:44Z",
        "iso8601_basic": "20181025T120844109754",
        "iso8601_basic_short": "20181025T120844",
        "iso8601_micro": "2018-10-25T12:08:44.109968Z",
        "minute": "08",
        "month": "10",
        "second": "44",
        "time": "12:08:44",
        "tz": "UTC",
        "tz_offset": "+0000",
        "weekday": "Thursday",
        "weekday_number": "4",
        "weeknumber": "43",
        "year": "2018"
    },
    "ansible_default_ipv4": {
        "address": "REDACTED",
        "alias": "eth0",
        "broadcast": "REDACTED",
        "gateway": "REDACTED",
        "interface": "eth0",
        "macaddress": "REDACTED",
        "mtu": 1500,
        "netmask": "255.255.255.0",
        "network": "REDACTED",
        "type": "ether"
    },
    "ansible_default_ipv6": {},
    "ansible_device_links": {
        "ids": {},
        "labels": {
            "xvda1": [
                "cloudimg-rootfs"
            ],
            "xvdd": [
                "config-2"
            ]
        },
        "masters": {},
        "uuids": {
            "xvda1": [
                "cac81d61-d0f8-4b47-84aa-b48798239164"
            ],
            "xvdd": [
                "2018-10-25-12-05-57-00"
            ]
        }
    },
    "ansible_devices": {
        "xvda": {
            "holders": [],
            "host": "",
            "links": {
                "ids": [],
                "labels": [],
                "masters": [],
                "uuids": []
            },
            "model": null,
            "partitions": {
                "xvda1": {
                    "holders": [],
                    "links": {
                        "ids": [],
                        "labels": [
                            "cloudimg-rootfs"
                        ],
                        "masters": [],
                        "uuids": [
                            "cac81d61-d0f8-4b47-84aa-b48798239164"
                        ]
                    },
                    "sectors": "83883999",
                    "sectorsize": 512,
                    "size": "40.00 GB",
                    "start": "2048",
                    "uuid": "cac81d61-d0f8-4b47-84aa-b48798239164"
                }
            },
            "removable": "0",
            "rotational": "0",
            "sas_address": null,
            "sas_device_handle": null,
            "scheduler_mode": "deadline",
            "sectors": "83886080",
            "sectorsize": "512",
            "size": "40.00 GB",
            "support_discard": "0",
            "vendor": null,
            "virtual": 1
        },
        "xvdd": {
            "holders": [],
            "host": "",
            "links": {
                "ids": [],
                "labels": [
                    "config-2"
                ],
                "masters": [],
                "uuids": [
                    "2018-10-25-12-05-57-00"
                ]
            },
            "model": null,
            "partitions": {},
            "removable": "0",
            "rotational": "0",
            "sas_address": null,
            "sas_device_handle": null,
            "scheduler_mode": "deadline",
            "sectors": "131072",
            "sectorsize": "512",
            "size": "64.00 MB",
            "support_discard": "0",
            "vendor": null,
            "virtual": 1
        },
        "xvde": {
            "holders": [],
            "host": "",
            "links": {
                "ids": [],
                "labels": [],
                "masters": [],
                "uuids": []
            },
            "model": null,
            "partitions": {
                "xvde1": {
                    "holders": [],
                    "links": {
                        "ids": [],
                        "labels": [],
                        "masters": [],
                        "uuids": []
                    },
                    "sectors": "167770112",
                    "sectorsize": 512,
                    "size": "80.00 GB",
                    "start": "2048",
                    "uuid": null
                }
            },
            "removable": "0",
            "rotational": "0",
            "sas_address": null,
            "sas_device_handle": null,
            "scheduler_mode": "deadline",
            "sectors": "167772160",
            "sectorsize": "512",
            "size": "80.00 GB",
            "support_discard": "0",
            "vendor": null,
            "virtual": 1
        }
    },
    "ansible_distribution": "CentOS",
    "ansible_distribution_file_parsed": true,
    "ansible_distribution_file_path": "/etc/redhat-release",
    "ansible_distribution_file_variety": "RedHat",
    "ansible_distribution_major_version": "7",
    "ansible_distribution_release": "Core",
    "ansible_distribution_version": "7.5.1804",
    "ansible_dns": {
        "nameservers": [
            "127.0.0.1"
        ]
    },
    "ansible_domain": "",
    "ansible_effective_group_id": 1000,
    "ansible_effective_user_id": 1000,
    "ansible_env": {
        "HOME": "/home/zuul",
        "LANG": "en_US.UTF-8",
        "LESSOPEN": "||/usr/bin/lesspipe.sh %s",
        "LOGNAME": "zuul",
        "MAIL": "/var/mail/zuul",
        "PATH": "/usr/local/bin:/usr/bin",
        "PWD": "/home/zuul",
        "SELINUX_LEVEL_REQUESTED": "",
        "SELINUX_ROLE_REQUESTED": "",
        "SELINUX_USE_CURRENT_RANGE": "",
        "SHELL": "/bin/bash",
        "SHLVL": "2",
        "SSH_CLIENT": "REDACTED 55672 22",
        "SSH_CONNECTION": "REDACTED 55672 REDACTED 22",
        "USER": "zuul",
        "XDG_RUNTIME_DIR": "/run/user/1000",
        "XDG_SESSION_ID": "1",
        "_": "/usr/bin/python2"
    },
    "ansible_eth0": {
        "active": true,
        "device": "eth0",
        "ipv4": {
            "address": "REDACTED",
            "broadcast": "REDACTED",
            "netmask": "255.255.255.0",
            "network": "REDACTED"
        },
        "ipv6": [
            {
                "address": "REDACTED",
                "prefix": "64",
                "scope": "link"
            }
        ],
        "macaddress": "REDACTED",
        "module": "xen_netfront",
        "mtu": 1500,
        "pciid": "vif-0",
        "promisc": false,
        "type": "ether"
    },
    "ansible_eth1": {
        "active": true,
        "device": "eth1",
        "ipv4": {
            "address": "REDACTED",
            "broadcast": "REDACTED",
            "netmask": "255.255.224.0",
            "network": "REDACTED"
        },
        "ipv6": [
            {
                "address": "REDACTED",
                "prefix": "64",
                "scope": "link"
            }
        ],
        "macaddress": "REDACTED",
        "module": "xen_netfront",
        "mtu": 1500,
        "pciid": "vif-1",
        "promisc": false,
        "type": "ether"
    },
    "ansible_fips": false,
    "ansible_form_factor": "Other",
    "ansible_fqdn": "centos-7-rax-dfw-0003427354",
    "ansible_hostname": "centos-7-rax-dfw-0003427354",
    "ansible_interfaces": [
        "lo",
        "eth1",
        "eth0"
    ],
    "ansible_is_chroot": false,
    "ansible_kernel": "3.10.0-862.14.4.el7.x86_64",
    "ansible_lo": {
        "active": true,
        "device": "lo",
        "ipv4": {
            "address": "127.0.0.1",
            "broadcast": "host",
            "netmask": "255.0.0.0",
            "network": "127.0.0.0"
        },
        "ipv6": [
            {
                "address": "::1",
                "prefix": "128",
                "scope": "host"
            }
        ],
        "mtu": 65536,
        "promisc": false,
        "type": "loopback"
    },
    "ansible_local": {},
    "ansible_lsb": {
        "codename": "Core",
        "description": "CentOS Linux release 7.5.1804 (Core)",
        "id": "CentOS",
        "major_release": "7",
        "release": "7.5.1804"
    },
    "ansible_machine": "x86_64",
    "ansible_machine_id": "2db133253c984c82aef2fafcce6f2bed",
    "ansible_memfree_mb": 7709,
    "ansible_memory_mb": {
        "nocache": {
            "free": 7804,
            "used": 173
        },
        "real": {
            "free": 7709,
            "total": 7977,
            "used": 268
        },
        "swap": {
            "cached": 0,
            "free": 0,
            "total": 0,
            "used": 0
        }
    },
    "ansible_memtotal_mb": 7977,
    "ansible_mounts": [
        {
            "block_available": 7220998,
            "block_size": 4096,
            "block_total": 9817227,
            "block_used": 2596229,
            "device": "/dev/xvda1",
            "fstype": "ext4",
            "inode_available": 10052341,
            "inode_total": 10419200,
            "inode_used": 366859,
            "mount": "/",
            "options": "rw,seclabel,relatime,data=ordered",
            "size_available": 29577207808,
            "size_total": 40211361792,
            "uuid": "cac81d61-d0f8-4b47-84aa-b48798239164"
        },
        {
            "block_available": 0,
            "block_size": 2048,
            "block_total": 252,
            "block_used": 252,
            "device": "/dev/xvdd",
            "fstype": "iso9660",
            "inode_available": 0,
            "inode_total": 0,
            "inode_used": 0,
            "mount": "/mnt/config",
            "options": "ro,relatime,mode=0700",
            "size_available": 0,
            "size_total": 516096,
            "uuid": "2018-10-25-12-05-57-00"
        }
    ],
    "ansible_nodename": "centos-7-rax-dfw-0003427354",
    "ansible_os_family": "RedHat",
    "ansible_pkg_mgr": "yum",
    "ansible_processor": [
        "0",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "1",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "2",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "3",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "4",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "5",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "6",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz",
        "7",
        "GenuineIntel",
        "Intel(R) Xeon(R) CPU E5-2670 0 @ 2.60GHz"
    ],
    "ansible_processor_cores": 8,
    "ansible_processor_count": 8,
    "ansible_processor_nproc": 8,
    "ansible_processor_threads_per_core": 1,
    "ansible_processor_vcpus": 8,
    "ansible_product_name": "HVM domU",
    "ansible_product_serial": "REDACTED",
    "ansible_product_uuid": "REDACTED",
    "ansible_product_version": "4.1.5",
    "ansible_python": {
        "executable": "/usr/bin/python2",
        "has_sslcontext": true,
        "type": "CPython",
        "version": {
            "major": 2,
            "micro": 5,
            "minor": 7,
            "releaselevel": "final",
            "serial": 0
        },
        "version_info": [
            2,
            7,
            5,
            "final",
            0
        ]
    },
    "ansible_python_version": "2.7.5",
    "ansible_real_group_id": 1000,
    "ansible_real_user_id": 1000,
    "ansible_selinux": {
        "config_mode": "enforcing",
        "mode": "enforcing",
        "policyvers": 31,
        "status": "enabled",
        "type": "targeted"
    },
    "ansible_selinux_python_present": true,
    "ansible_service_mgr": "systemd",
    "ansible_ssh_host_key_ecdsa_public": "REDACTED KEY VALUE",
    "ansible_ssh_host_key_ed25519_public": "REDACTED KEY VALUE",
    "ansible_ssh_host_key_rsa_public": "REDACTED KEY VALUE",
    "ansible_swapfree_mb": 0,
    "ansible_swaptotal_mb": 0,
    "ansible_system": "Linux",
    "ansible_system_capabilities": [
        ""
    ],
    "ansible_system_capabilities_enforced": "True",
    "ansible_system_vendor": "Xen",
    "ansible_uptime_seconds": 151,
    "ansible_user_dir": "/home/zuul",
    "ansible_user_gecos": "",
    "ansible_user_gid": 1000,
    "ansible_user_id": "zuul",
    "ansible_user_shell": "/bin/bash",
    "ansible_user_uid": 1000,
    "ansible_userspace_architecture": "x86_64",
    "ansible_userspace_bits": "64",
    "ansible_virtualization_role": "guest",
    "ansible_virtualization_type": "xen",
    "gather_subset": [
        "all"
    ],
    "module_setup": true
}
```

We can reference the model of the first disk in the facts shown above in a
template or playbook as:

```yaml
{{ ansible_facts['devices']['xvda']['model'] }}
```

To reference the system hostname:

```yaml
{{ ansible_facts['nodename'] }}
```

We can use facts in conditionals, and also in templates. We can also use facts
to create dynamic groups of hosts that match particular criteria, see the
`group_by` module documentation for details.

> [!NOTE]
> Because `ansible_date_time` is created and cached when Ansible gathers facts
> before each playbook run, it can get state with long-running playbooks. If the
> playbook takes a long time to run, use the `pipe` filter or `now()` with
> Jinja2 tempalte isntead of `ansible_date_time`.

### Package requirements

On som distros, missing fact values or facts set to default values because the
packages that support gathering those facts are not installed by default. We can
install the necessary packages on remote hosts using the OS package manager.
Known dependencies include:
- Linux Network fact gathering - depends on `ip` binary, commonly included in
  the `iproute2` package.

### Caching facts

Like registered variables, facts are stored in memory by default. However,
unlike registered variables, facts can be gathered independently and cached for
repeated use. With cached facts, we can refer to facts from one system when
configuring a second system, even if Ansible executes the current play on the
second system first. For example:

```bash
{{ hostvars['asdf.example.com']['ansible_facts']['os_family'] }}
```

Caching is controlled by the cache plugins. By default, Ansible uses the memory
cache plugin, which stores facts in memory for the duration of the current
playbook run. To retain Ansible facts for repeated use, select a different
cache plugin.

Fact caching can improve performance. If we manage thousands of hosts, we can
configure fact caching to run nightly, and then manage configuration on a
smaller set of servers periodically throughout the day. With cached facts, we
have access to variables and information about all hosts even when we are only
managing a small number of servers.

### Disabling facts

By default, Ansible gathers facts at the beginning of each play. If we do not
need to gather facts ( for example, if we know everything about your systems
centrally), we can turn off fact gathering at the play level to improve
scalability. Disabling facts may paricularly improve performance in push mode
with very large numbers of systems, or if we are using Ansible on experimental
platforms. To disable fact gathering:

```yaml
- hosts: whatever
    gather_facts: false
```

### Adding custom facts

The setup module in Ansible automatically discovers a standard set of facts
about each host. If we want to add custom values to the facts, we can write a
custom facts module, set temporary facts with a `ansible.builtin.set_fact` task,
or provide permanent custom facts using the facts.d directory.

#### facts.d or local facts

We can add static custom facts by adding static files to fact       
