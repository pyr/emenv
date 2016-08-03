emenv: emacs-lisp dependency fetcher
====================================

emenv is a dependency fetcher for emacs-lisp. It is meant as a
companion to recent emacsen and does the following:

- Fetching ELPA package information from common sources (MELPA-stable,
  MELPA, GNU, Marmalade, Org)
-

emenv does not rely on emacs-lisp at all, it is a standalone program
with a simple emacs-lisp reader, sufficient to parse the syntax of
the config file and the syntax of ELPA archives

Building
--------

emenv is written in go, as such, you will only need to
fetch this repository and then run `go build`. The
resulting file can be copied on your PATH.

Configuration
-------------

```clojure
(source my-repo "https://elpa.example.com/packages")
(prefer my-repo omelpa-stable gnu melpa)
(package ag)
(package projectile)
(package auto-complete (from melpa-stable))

```

Running
-------

You can fetch latest repository information with:

```
emenv sync
```

You can then install packages with:

```
emenv install
```