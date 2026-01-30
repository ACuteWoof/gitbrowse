# Gitbrowse: a simple web server for git.

Note that this project is constantly updated and that changes may occur to
routes or anything that your users may get used to. Try to avoid using 301
redirects that may be cached and prevent the user from accessing routes that
may change in the future.

## Routes

| Route                                               | Description                                                                                                      |
| --------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `/`                                                 | shows a list of users                                                                                            |
| `/{user}`                                           | shows a list of repositories for the user                                                                        |
| `/{user}/{repo}/`                                   | shows the readme and licenses for a repo                                                                         |
| `/{user}/{repo}/branch/`                            | shows the list of branches for a repo                                                                            |
| `/{user}/{repo}/branch/{branch}/tree/{filepath...}` | shows the tree (or file if filepath points to a file and not a directory) on the given path on the given branch |
| `/{user}/{repo}/branch/{branch}/commit`             | shows the commit log for the given branch                                                                        |
| `/{user}/{repo}/tag/`                               | shows the list of annotated tags for a repo                                                                      |
| `/{user}/{repo}/tag/{name}/{fileName}`              | downloads a tag with the given file name, you will usually be led to this by a link                              |
| `/{user}/{repo}/show/{hash}`                        | shows the output of `git show` for the given hash                                                                |
| `/{user}/{repo}/grep?q={regex}`                     | greps the repo with the regex and shows the lines that match                                                     |

Browsing [this](https://git.lewoof.xyz/gitbrowse) site should give you a good
idea of what the site looks like, as I stick to the defaults other than having
added remote fonts and treating `/` as `/{user}/`.

## Configuration

The configuration files are mainly in the `main.go` and the `config/config.go`
files. (The project was originally written with the configuration being entirely
dependent on `config/config.go`, but some temporary changes in `main.go` turned
out to be not temporary --- this may of course change at any time in the
future).

The configuration for the page that lists all users is in `config/config.go`.
All other configuration can be found in `main.go`.

The configuration is part of the binary insofar as it is not a stylesheet in
`static`, so I recommend you have a branch for your configuration that updates
when the master branch is updated and compile from there. I use
[`entr`](https://github.com/eradman/entr) for this.

### The `config.PageConfig` struct

A quick look at `main.go` reveals two configuration functions that return
functions. This is also a result of one of those temporary changes I mentioned
earlier. Both are meant for obtaining variables of the same type, `config.PageConfig`.

`getIndexConfigGetterUser` deals with the configuration for the page that lists
the repositories for a given user (`/{user}/...`). Here we have access to one variable,
`username`.

`getRepoConfigGetter` deals with the configuration for all the pages that deal
with repositories (`/{user}/{repo}/...`).

There is also a variable, `IndexPageConfig` in `config/config.go` used for the
page that lists all users. It should be treated similarly to the variable in
`getRepoConfigGetter`, but without access to the `username` variable. Ignore
`RepoPageConfig`.

The fields you really should change are `RootDir`, `Title`, and `Description`.

Now to what the fields of `config.PageConfig` represent:

| Field         | Explanation                                                                                                                                                                          |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `URLRoot`     | The URL for the the user page (for index config) or the repo page (for repo config). E.g, the `URLRoot` for `/{user}/{repo}/branch/{branch}/tree/{filepath...}` is `/{user}/{repo}`. |
| `RootDir`     | The directory in the unix filesystem where either the user's repositories are immediately found (for index config) or the repo's directory is immediately found (for repo config).   |
| `CloneURL`    | Not used in the user page. The URL used to clone the repo. This must be handled by some other service as Gitbrowse does not bother to; I use nginx with the git http server.         |
| `Title`       | The title of the page, used for `<head>` and for the `<header>`                                                                                                                      |
| `Description` | The description of the page, used for `<head>`                                                                                                                                       |
| `Thumbnail`   | The URL to the thumbnail image, used for `<head>` and the image in `<header>`, this can be an https URL if you have an API that provides your users with profile pictures            |
| `Favicon`     | The URL to the favicon                                                                                                                                                               |
| `Styles`      | A list of URLs to stylesheets                                                                                                                                                        |

The header is generated by a function in `template/common-header.go`, and the head tag in `template/common-head.go`.

A quick look at the `static` directory and the stylesheets should give you a basic idea of customizing the look of the site. The `template` module contains the functions responsible for producing the HTML that is rendered.

JavaScript is not used in the output of the server to allow full functionality for users who prefer not to turn on JavaScript on their browsers, and because it is quite unnecessary.

Remember always to build after making any changes outside `static`.

## Building

Gitbrowse requires a git binary that is accessible by calling `git`.

To build, install `go` and run `go build`. The resulting binary will be the
executable for the web server. Run it in the same directory as the `static`
directory, or change `./static` to your preferred location in `main.go`.

## Running

Executing `./gitbrowse` will start the server on port 8088. This can also be
changed in `main.go`.

Here's the systemd service file I use to run it as a service:

```
[Unit]
Description=Gitbrowse

[Service]
Restart=always
RestartSec=15
User=acutewoof
WorkingDirectory=/home/acutewoof/program/gitbrowse
ExecStart=/home/acutewoof/program/gitbrowse/gitbrowse

[Install]
WantedBy=network-online.target
```

Note that the working directory contains the `static` directory.

## Addresses, Links, and Contact

See [https://www.lewoof.xyz#pgp](https://www.lewoof.xyz#pgp) for my PGP key.

- Email: [contact@lewoof.xyz](mailto:contact@lewoof.xyz)
- Site: [www.lewoof.xyz](https://www.lewoof.xyz)
- Solana: `JDkK2kpBmPm6YyYnLYNHpw8FhKyZ9AQ2CQTqj6BQxFKY`
- Bitcoin (native segwit): `c1qh5m8v9xyd8l4yc7d3qfs5a83rqk3eukcjlw6sh`
- BuyMeACoffee: [acutewoof](https://www.buymeacoffee.com/acutewoof)

## License

Gitbrowse is licensed under the GNU General Public License v3.0.

See `LICENSE.txt` for more information, and `THIRD_PARTY_LICENSES.txt` for
licenses of the libraries used.
