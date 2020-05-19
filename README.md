# Interplanetary Git Service (IGiS) Remote Helper

Push and fetch commits to IPFS.

## Installation

`npm install --global git-remote-igis`

## Usage

Push `master` with tags and get an IPFS CID back:

`git push --tags igis:: master`

Push all branches to the repository `repo`:

`git push --all igis::repo`

Clone an example repository:

`git clone ipns://git-remote-ipfs.dhappy.org git-remote-ipfs`

Pull a commit:

`git pull ipfs://Qma5iwyvJqxzHqCT9aqyc7dxZXXGoDeSUyPYFqkCWGJw92`

Push with the `.git/vfs/` directory:

`GIT_IPFS_VFS=t git push ipfs::`

Push to an IPNS remote:

* `ipfs key gen --type=rsa --size=2048 mysite`
* `git remote add ipns ipns::key:mysite`
* `git push ipns`
* `git pull`

## Workflow

1. Developer Arrives
2. Browses Open Issues
3. Select an Issue to Work On
4. Pull Down the Active Working Tree
5. Make Commits on the Tree

## Generated File Structure

* `/`: the contents of the branch that was pushed
* `.git/`: CBOR-DAG representing a git repository
* `.git/HEAD`: string entry denoting the current default branch
* `.git/refs/(heads|tags)/*`: Pointers to commit objects
* `[commit]/channel`: The MAM channel state
* `[commit]/parents`: The commits parent commits
* `[commit]/(author|committer)`: The commits author and committer signatures


## Overview

Git is at its core an object database. There are four types of objects: Commits, Trees, Tags, & Blobs.

When a remote helper is asked to push, it receives the key of a Commit object. That commit has an associated Tree and zero or more Commit parents.

Trees are lists of umasks, names, and keys for either nested Trees or Blobs.

Tags are named links to specific commits. They are frequently used to mark versions.

The helper traverses the tree and parents of the root Commit and stores them all on the remote.

### IPLD Git Remote

Integrating Git and IPFS has been on ongoing work with several solutions over the years. The predecessor to this one stored the raw blocks in the IPFS DAG using a multihash version of git's SHA1s.

The SHA1 keys used by Git aren't exactly for the hash of the object. Each serialized form is prefaced with a header of the format: "`#{type} #{size}\x00`". So a Blob in Git is this header plus the file contents.

Because the IPLD remote stores the raw Git blocks, the file data is fully present, but unreadable because of the header.

### v0.2

v0.2 of this project was based on the IPLD helper and was thus in Go. It created named directories for different the types of objects and removed the header for the stored version.

This allowed creating a checked out version of the repository for essentially free because all the files are already present in the repository.

Unfortunately, with that representation of the object store the entire thing has to be written every time you push. It is monumentally slow for an operation that's done so frequently.

### v0.3

v0.3 is a Node app which leverages the IPFS DAG to create a hybrid data structure / filesystem representing the object store in much the same way git does internally.

This should allow me to only calculate back to a previously inserted commit and simply include it in the chain.

The language change gives me access to [OrbitDB](//github.com/orbitdb/) which should be exciting.

# Troubleshooting
* `fetch: manifest has unsupported version: x (we support y)` on any command
  - This usually means that cache tracker data format has changed
  - Remove the cache with: `rm -rf .git/remote-ipfs`

* `panic: runtime error: invalid memory address or nil pointer dereference`
  - This dramatic message likely means the IPFS server isn't running.

# License
MIT
