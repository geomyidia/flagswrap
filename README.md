# flagswrap

[![Build Status][gh-actions-badge]][gh-actions]

*Some convenience wrapping around go-flags, mostly for error handling*

[![][logo]][logo-large]

The intent of this project is to provide a band-aide until the go-flags project has a chance to address the following issues:

* <https://github.com/jessevdk/go-flags/issues/306>
* <https://github.com/jessevdk/go-flags/issues/361>
* <https://github.com/jessevdk/go-flags/issues/377>

One possible use for this library is to change the (currently broken) example error handling [here](https://github.com/jessevdk/go-flags/blob/206428b03a5152306e043b5c0b7c7575c24afc61/examples/main.go) from this:

```golang
func main() {
    if _, err := parser.Parse(); err != nil {
        switch flagsErr := err.(type) {
        case flags.ErrorType:
            if flagsErr == flags.ErrHelp {
                os.Exit(0)
            }
            os.Exit(1)
        default:
            os.Exit(1)
        }
    }
}
```

to this:

```golang
func main() {
    if _, err := parser.Parse(); err != nil {
        wrappedErr := flagswrap.WrapError(err)
        switch {
        case wrappedErr.IsHelp():
            os.Exit(0)
        // Self-documenting go-flags errors:
        case wrappedErr.IsVerbose():
            os.Exit(1)
        // go-flags errors that need more context:
        case wrappedErr.IsSilent():
            // TODO: if you see duplicate error messages here, then
            // you just need to move the error in question from the
            // goFlagsSilentErrors to the goFlagsVerboseErrors map
            // in ./errors.go -- and then submit a PR!
            fmt.Printf("Error: %v\n", wrappedErr)
            os.Exit(1)
        default:
            // TODO: anything here might justify a PR ...
            fmt.Printf("ERROR (unexpected): %+v\n", wrappedErr)
            os.Exit(1)
        }
    }
}
```

You can give this a shot with:

```shell
make
./bin/simple -h
./bin/simple
./bin/simple --bloop
./bin/dupe
```

Btw, this doesn't feel super Go-idiomatic -- any PRs making such improvements would be gladly reviewed and (ideally) accepted :-)

[//]: ---Named-Links---

[logo]: assets/images/project-avatar-small.png
[logo-large]: assets/images/project-avatar.png
[github]: https://github.com/geomyidia/flagswrap
[gh-actions-badge]: https://github.com/geomyidia/flagswrap/workflows/ci%2Fcd/badge.svg
[gh-actions]: https://github.com/geomyidia/flagswrap/actions
