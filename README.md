# RTNS

RTNS (RTrade Name Service) is a stand-alone IPNS record management service, designed to facilitate secure publishing of IPNS records, leveraging an encrypted keystore known as [kaas](https://github.com/RTradeLtd/kaas). Internally it facilitates scheduled republishing of all published records.

It is essentially a modified and condensed version of [go-ipfs/namesys](https://github.com/ipfs/go-ipfs/tree/master/namesys) with minor optimizations

# Limitations

* When used within Temporal, any keys derived from the fail-over KaaS host are not eligible for automated republishing

# Future Improvements

* DNSLink support
* Act as a gateway implementation to TNS (Temporal Name Server)
* Enable HA Kaas Backend
  * This will involve repeatedly iterating through all available KaaS hosts attempting to retrieve the private key, until we either find the key or we iterate through all available KaaS hosts without finding one, triggering an error
* Enable automatic topic subscription for IPNS pubsub
  * This would involve using rtfs to call an IPFS node, establishing a subscription for a given topic