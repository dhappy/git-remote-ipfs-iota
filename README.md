# Interplanetary Filesystem (IPFS) Git Remote Helper

Push and fetch commits to IPFS. Updates are broadcast to the IOTA tangle where others may retrieve the most recent version.

## Installation

`npm install --global git-remote-ipfs-mam`

## Usage

#### (Insecure) Cloud Backup

1. `git remote add ipfs ipfs+mam://myproject`
2. `git push ipfs --tags`
3. `git push ipfs --all`
4. Pin the resultant hash on a pinning service.

#### Push `master` with tags and get an IPFS CID back:

`git push --tags ipfs+mam:: master`

#### Pull the most recent version:

`git pull ipfs+mam::Qma5iwyvJqxzHqCT9aqyc7dxZXXGoDeSUyPYFqkCWGJw92`

#### Clone a repository:

`git clone ipfs::Qma5iwyvJqxzHqCT9aqyc7dxZXXGoDeSUyPYFqkCWGJw92 repo`

#### See debugging info:

`IGIS_DEBUG=t git push ipfs::`

## IPFS Data Structures

This program is an extension of [git-remote-ipfs](https://github.com/dhappy/git-remote-ipfs). For information about the IPFS file structure, see that repository.

## IOTA Data Structures

In IOTA, it is possible to write to any address. A Masked Authentication Message (MAM) channel is formed by including in each published message the next address that will be used.

When a repository is published, this remote publishes a signed JSON Linked Data object to the MAM channel for the format:

```javascript
{
  '@context': {
    schema: 'http://schema.org/',
    action: 'schema:action',
    agent: 'schema:name',
    repository: 'schema:url',
    publisher: 'schema:url',
    next_root: 'schema:url',
    published_at: 'schema:datetime',
  },
  action: 'RepositoryUpdate',
  repository: 'ipfs://QmThisIsTheCIDOfTheRepo',
  publisher: 'did:key:zABase58EncodedED25519Key',
  next_root: `iota://NEXT9TANGLE9ADDRESS9IN9THE9MAM9TREE:TAG9FROM9REPO9UUID`,
  agent: 'git-remote-ipfs+mam',
  published_at: new Date(),
  'https://w3id.org/security#proof': {…}
}
```

Additionally, a signed JSON-LD object is published to the tangle address `99IPFS9MAM9CHNL9LINK9VA99${multicodec_of_repo_cid}` of the format:

```javascript
{
  '@context': {
    schema: 'http://schema.org/',
    action: 'schema:action',
    agent: 'schema:name',
    publisher: 'schema:url',
    bundle: 'schema:url',
    published_at: 'schema:datetime',
  },
  action: 'MAMLink',
  publisher: 'did:key:zABase58EncodedED25519Key',
  bundle: `iota:bundle://IOTA9BUNDLE9HASH9In9MAM9TREE`,
  agent: 'git-remote-ipfs+mam',
  published_at: new Date(),
  'https://w3id.org/security#proof': {…}
}
```

If a different remote is asked to clone from that CID, it can check that address and get the bundle hash of a message in the channel. The `publisher` field is also present in the IPFS repository, so the user is able to verify the signatures (and differentiate between genuine messages and those inserted by an attacker).

# License
MIT
