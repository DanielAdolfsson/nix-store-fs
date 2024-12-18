# nix-store-fs

`nix-store-fs` is a FUSE filesystem that creates a tailored view of /nix/store, limiting access to the exact items required by a specific path, ensuring clarity and security.

## Features

- Creates a filtered view of `/nix/store`.
- Grants read-only access based on specified paths.
- Automatically resolves paths using the Nix daemon.

## Installation

To install `go-nix-fs`, use the following commands:

```sh
git clone https://github.com/DanielAdolfsson/nix-store-fs.git
cd nix-store-fs/src
go build
```

## Usage

To mount the filesystem, use the following command:

```sh
./nix-store-fs <item> <mountpoint>
```

Replace `<item>` with a name that exists in `/nix/store`, and `<mountpoint>` with the directory where you want to mount the filesystem.

## Example

```sh
./go-nix-fs wlwyqvdalg32pdf20klxndhhqmra9mmh-bash-interactive-5.2p37 /mnt/nix-view

$ ls -1 /mnt/nix-view
1ci7cipl06rf3c8cr7vz2zzr36wpxms1-glibc-2.40-36
49l1b1x9aw21c00qmma7zxzbj7qa91pz-readline-8.2p13
60i9rqqaj1zzyspws5byblcc3gq8kp4v-libidn2-2.3.7
83mg03jav95qmj15qbw9c9l464brlrg3-libunistring-1.2
l2xykaacd3xgc1i83j7qvdc4k064w820-xgcc-13.3.0-libgcc
qirrh7cdr8vm2wg56mbx5slgcvdnhcb9-ncurses-6.4.20221231
wlwyqvdalg32pdf20klxndhhqmra9mmh-bash-interactive-5.2p37

$ nix path-info --recursive /nix/store/wlwyqvdalg32pdf20klxndhhqmra9mmh-bash-interactive-5.2p37
/nix/store/83mg03jav95qmj15qbw9c9l464brlrg3-libunistring-1.2
/nix/store/60i9rqqaj1zzyspws5byblcc3gq8kp4v-libidn2-2.3.7
/nix/store/l2xykaacd3xgc1i83j7qvdc4k064w820-xgcc-13.3.0-libgcc
/nix/store/1ci7cipl06rf3c8cr7vz2zzr36wpxms1-glibc-2.40-36
/nix/store/qirrh7cdr8vm2wg56mbx5slgcvdnhcb9-ncurses-6.4.20221231
/nix/store/49l1b1x9aw21c00qmma7zxzbj7qa91pz-readline-8.2p13
/nix/store/wlwyqvdalg32pdf20klxndhhqmra9mmh-bash-interactive-5.2p37
```

## Options

- `-daemon-socket-path`: Path to the Nix daemon socket (default `/nix/var/nix/daemon-socket/socket`).
- `-store-path`: Path to the Nix store (default `/nix/store`).

```