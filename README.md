# rtns

rtns is a IPNS publishing service for use with [RTradeLtd/kaas](https://github.com/RTradeLtd/kaas), with republishing of published records.

# limitations

When used within Temporal, any keys derived from the fail-over KaaS host are not eligible for automated republishing

# Future Improvements

* DNSLink support
* Act as a gateway implementation to TNS (Temporal Name Server)
* Enable HA Kaas Backend
  * This will involve repeatedly iterating through all available KaaS hosts attempting to retrieve the private key, until we either find the key or we iterate through all available KaaS hosts without finding one, triggering an error