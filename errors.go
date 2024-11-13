package git

import "errors"

var (
	ErrNotEmptyRepository = errors.New("destination path already exists and is not an empty directory")
	ErrUnknownRevision    = errors.New("unknown revision or path not in the working tree")
)

/*
https://jvns.ca/blog/2024/04/10/notes-on-git-error-messages/


fatal: refusing to merge unrelated histories
---
fatal: remote origin already exists
---
fatal: not a git repository
---
fatal: repository not found
---
$ git push
To github.com:jvns/int-exposed
! [rejected]        main -> main (non-fast-forward)
error: failed to push some refs to 'github.com:jvns/int-exposed'
hint: Updates were rejected because the tip of your current branch is behind
hint: its remote counterpart. Integrate the remote changes (e.g.
hint: 'git pull ...') before pushing again.
hint: See the 'Note about fast-forwards' in 'git push --help' for details.
---
$ git pull
hint: You have divergent branches and need to specify how to reconcile them.
hint: You can do so by running one of the following commands sometime before
hint: your next pull:
hint:
hint:   git config pull.rebase false  # merge
hint:   git config pull.rebase true   # rebase
hint:   git config pull.ff only       # fast-forward only
hint:
hint: You can replace "git config" with "git config --global" to set a default
hint: preference for all repositories. You can also pass --rebase, --no-rebase,
hint: or --ff-only on the command line to override the configured default per
hint: invocation.
fatal: Need to specify how to reconcile divergent branches.
---
$ git checkout asdf
error: pathspec 'asdf' did not match any file(s) known to git
---
$ git switch asdf
fatal: invalid reference: asdf
---
$ git switch v0.1
fatal: a branch is expected, got tag 'v0.1'`
---
$ git checkout HEAD^
Note: switching to 'HEAD^'.

You are in 'detached HEAD' state. You can look around, make experimental
changes and commit them, and you can discard any commits you make in this
state without impacting any branches by switching back to a branch.

If you want to create a new branch to retain commits you create, you may
do so (now or later) by using -c with the switch command. Example:

  git switch -c

Or undo this operation with:

  git switch -

Turn off this advice by setting config variable advice.detachedHead to false

HEAD is now at 182cd3f add "swap byte order" button
---
$ git status
interactive rebase in progress; onto c694cf8
Last command done (1 command done):
   pick 0a9964d wip
No commands remaining.
You are currently rebasing branch 'main' on 'c694cf8'.
  (fix conflicts and then run "git rebase --continue")
  (use "git rebase --skip" to skip this patch)
  (use "git rebase --abort" to check out the original branch)

Unmerged paths:
  (use "git restore --staged ..." to unstage)
  (use "git add ..." to mark resolution)
  both modified:   index.html

no changes added to commit (use "git add" and/or "git commit -a")
---
$ git rebase main
CONFLICT (modify/delete): index.html deleted in 0ce151e (wip) and modified in HEAD.  Version HEAD of index.html left in tree.
error: could not apply 0ce151e... wip
---
CONFLICT (modify/delete): index.html deleted on `main` and modified on `mybranch`
---
$ git status
On branch master
You have unmerged paths.
  (fix conflicts and run "git commit")
  (use "git merge --abort" to abort the merge)

Unmerged paths:
(use “git add/rm …” as appropriate to mark resolution)
deleted by them: the_file

no changes added to commit (use “git add” and/or “git commit -a”)
---
$ git clean
fatal: clean.requireForce defaults to true and neither -i, -n, nor -f given; refusing to clean
---
error: Your local changes to the following files would be overwritten by merge
---
error: pathspec 'file.txt' did not match any file(s) known to git
---
error: failed to clone some remote refs
---

*/
