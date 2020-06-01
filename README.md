# Interplanetary Git Service (IGiS) Remote Helper

Push and fetch commits to IPFS.

## Installation

`npm install --global git-remote-igis`

## Usage

#### (Insecure) Cloud Backup

1. `git push igis:: --tags # you can't push all and tags at the same time`
2. `git push igis:: --all`
3. Pin the resultant hash on a pinning service.

#### Push `master` with tags and get an IPFS CID back:

`git push --tags igis:: master`

#### Push `master` with tags and get an IPFS CID back:

`git push --tags igis:: master`

#### Pull a commit:

`git pull igis://Qma5iwyvJqxzHqCT9aqyc7dxZXXGoDeSUyPYFqkCWGJw92`

#### See debugging info:

`IGIS_DEBUG=t git push ipfs::`

## Generated File Structure

* `/`: the contents of the branch that was pushed
* `.git/`: CBOR-DAG representing a git repository
* `.git/HEAD`: string entry denoting the current default branch
* `.git/channel`: The MAM channel state
* `.git/refs/(heads|tags)/*`: Pointers to commit objects

Each commit then has:

* `parents`: The commit's parent commits
* `(author|committer)`: The commits author and committer signatures
* `gpgsig`: Optional signature from the initial commit
* `tree`: The filesystem state at the time of this commit
* `modes`: `tree` is an IPFS Protobuffer-UnixFS DAG which is browsable through the web, but can't store the file mode information, so this is that info.

## Overview

This remote serializes a Git commit tree to a CBOR-DAG stored in IPFS.

### IPLD Git Remote

Integrating Git and IPFS has been on ongoing work with several solutions over the years. The [predecessor to this one](//github.com/ipfs-shipyard/git-remote-ipld) stored the raw blocks in the IPFS DAG using a multihash version of git's SHA1s.

The SHA1 keys used by Git aren't exactly for the hash of the object. Each serialized form is prefaced with a header of the format: "`#{type} #{size}\x00`". So a Blob in Git is this header plus the file contents.

Because the IPLD remote stores the raw Git blocks, the file data is fully present, but unreadable because of the header.

## Troubleshooting

It is safe to delete `.git/remote-igis/cache/`. If you remove `.git/remote-igis/config.json` you remove the key used to write to the MAM chain and you'll have to start a new one. People following your previous chain wouldn't be able to find updates.

# License
MIT
