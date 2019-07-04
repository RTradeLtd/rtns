# RTNS

RTNS (RTrade Name Service) is a stand-alone IPNS record management service, designed to facilitate secure publishing of IPNS records, leveraging an encrypted keystore known as [kaas](https://github.com/RTradeLtd/kaas). Internally it facilitates scheduled republishing of all published records.

It is essentially a modified and condensed version of [go-ipfs/namesys](https://github.com/ipfs/go-ipfs/tree/master/namesys) with minor optimizations.

## Development

### Using `$GOPATH`

Ref: https://splice.com/blog/contributing-open-source-git-repositories-go/

1. Fork the repository
2. Clone the repository by running `git clone git@github.com:RTradeLtd/rtns.git $GOPATH/src/github.com/RTradeLtd/rtns`
3. Run `cd $GOPATH/src/github.com/RTradeLtd/rtns`
3. Set up remotes.

```bash
git remote rename origin upstream
git remote add origin git@github.com:<your-github-username>/rtns.git
```
4. Add `export GO111MODULE=on` to `.bashrc` or `.bash_profile` (if you're on a Mac) or `.zshrc` (if you're using [zsh](https://github.com/robbyrussell/oh-my-zsh)). Make sure to reload the rc file of your choice by running `source <rc-file>`
5. Run `go mod download` to download the dependencies
6. To run the tests, use `go test ./...`

### Outside $GOPATH

1. Fork and clone the repository to any location on your machine
2. Run `cd rtns`
3. Set up a remote for the upstream repository

```bash
git remote add upstream git@github.com:RTradeLtd/rtns.git
```

3. Run `go mod download` to download the dependencies
4. To run the tests, use `go test ./...`

## Limitations

* When used within Temporal, any keys derived from the fail-over KaaS host are not eligible for automated republishing.

## Future Improvements

* DNSLink support
* Act as a gateway implementation to TNS (Temporal Name Server)
* Enable HA Kaas Backend
  * This will involve repeatedly iterating through all available KaaS hosts attempting to retrieve the private key, until we either find the key or we iterate through all available KaaS hosts without finding one, triggering an error
* Enable automatic topic subscription for IPNS pubsub
  * This would involve using rtfs to call an IPFS node, establishing a subscription for a given topic.