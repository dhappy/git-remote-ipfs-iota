package core

import (
	"container/list"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	ipfs "github.com/ipfs/go-ipfs-api"
	ipldgit "github.com/ipfs/go-ipld-git"
	git "gopkg.in/src-d/go-git.v4"
	plumbing "gopkg.in/src-d/go-git.v4/plumbing"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	mh "gx/ipfs/QmU9a9NV9RdPNwZQDYd5uKsm6N6LJLSvLbywDDYFbaaC6P/go-multihash"
)

type Push struct {
	objectDir string
	gitDir    string

	done    uint64
	todoc   uint64
	todo    *list.List
	log     *log.Logger
	tracker *Tracker
	repo    *git.Repository

	NewNode func(hash *cid.Cid, data []byte) error
}

func NewPush(gitDir string, tracker *Tracker, repo *git.Repository) *Push {
	return &Push{
		objectDir: path.Join(gitDir, "objects"),
		gitDir:    gitDir,
		todo:      list.New(),
		log:       log.New(os.Stderr, "push: ", 0),
		tracker:   tracker,
		repo:      repo,
		todoc:     1,
	}
}

func (p *Push) PushHash(hash string) error {
	p.todo.PushFront(hash)
	return <-p.doWork()
}

func (p *Push) doWork() <-chan error {
	api := ipfs.NewLocalShell()
	out := make(chan error)

	go func() {
		for e := p.todo.Front(); e != nil; e = e.Next() {
			hash := e.Value.(string)

			sha, err := hex.DecodeString(hash)
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			has, err := p.tracker.HasEntry(sha)
			if err != nil {
				out <- fmt.Errorf("push/process: %v", err)
				return
			}

			if has {
				p.todoc--
				continue
			}

			expectedCid, err := CidFromHex(hash)
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			obj, err := p.repo.Storer.EncodedObject(plumbing.AnyObject, plumbing.NewHash(hash))
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			rawReader, err := obj.Reader()
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			raw, err := ioutil.ReadAll(rawReader)
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			switch obj.Type() {
			case plumbing.CommitObject:
				raw = append([]byte(fmt.Sprintf("commit %d\x00", obj.Size())), raw...)
			case plumbing.TreeObject:
				raw = append([]byte(fmt.Sprintf("tree %d\x00", obj.Size())), raw...)
			case plumbing.BlobObject:
				raw = append([]byte(fmt.Sprintf("blob %d\x00", obj.Size())), raw...)
			case plumbing.TagObject:
				raw = append([]byte(fmt.Sprintf("tag %d\x00", obj.Size())), raw...)
			}

			p.done++
			p.log.Printf("%d/%d %s %s\r\x1b[A", p.done, p.todoc, hash, expectedCid.String())

			res, err := api.DagPut(raw, "raw", "git")
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			err = p.tracker.AddEntry(sha)
			if err != nil {
				out <- fmt.Errorf("push: %v", err)
				return
			}

			if expectedCid.String() != res {
				out <- fmt.Errorf("CIDs don't match: expected %s, got %s", expectedCid.String(), res)
				return
			}

			if p.NewNode != nil {
				if err := p.NewNode(expectedCid, raw); err != nil {
					out <- err
					return
				}
			}

			p.processLinks(raw)
		}
		p.log.Printf("\n")
		out <- nil
	}()
	return out
}

func (p *Push) processLinks(object []byte) error {
	nd, err := ipldgit.ParseObjectFromBuffer(object)
	if err != nil {
		return fmt.Errorf("push/process: %v", err)
	}

	links := nd.Links()
	for _, link := range links {
		mhash := link.Cid.Hash()
		decoded, err := mh.Decode(mhash)
		if err != nil {
			return fmt.Errorf("push/process: %v", err)
		}

		has, err := p.tracker.HasEntry(decoded.Digest)
		if err != nil {
			return fmt.Errorf("push/process: %v", err)
		}

		if has {
			continue
		}

		p.todoc++
		p.todo.PushBack(hex.EncodeToString(decoded.Digest))
	}
	return nil
}
