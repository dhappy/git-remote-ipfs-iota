# Interplanetary Git Service (IGiS) Remote Helper

Push and fetch commits to IPFS.

## Installation

`npm install --global git-remote-igis`

## Usage

Push:

`git push ipfs:: master`

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

## Generated File Structure

* `/`: the contents of the branch that was pushed
* `.git/blobs/`, `.git/trees/`, `.git/commits/`, `.git/tags/`: various Git objects stored by their SHA1 hash as filename
* `.git/refs/heads/*`: files containing the root hash of various Git branches
* `.git/HEAD`: the name of the branch contained in this repo

The virtual filesystem makes all the trees associated with all the commits available, but takes about twice as long to generate:

* `.git/vfs/messages/`: all the trees linked by commit message and sorted by date
* `.git/vfs/authors/#{name}/`: commits sorted by author
* `.git/vfs/rev/messages/`, `.git/vfs/rev/authors/#{name}/`: the commits as before, but prefaced with a count to reverse the order
* `.git/vfs/commits/`: vfs commits named by commit SHA1
* `.git/vfs/trees/`: content trees named by tree SHA1
* `.git/vfs/HEAD`: root vfs commit

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
