# letsencrypt-cert-remaining-days
A certificate parser to get how many days you have left on your letsencrypt certificates

A simple way to check how many days your certificate remains valid.

By default it reads the configuration files located on `/etc/letsencrypt` (actually under `/etc/letsencrypt/renew`), search for
the `cert = ` entry to find the certificate.  Then check how many days it has left.

You can specify a single domain using parameter `--domain`.

## Example

```bash
❯ ./letsencrypt-cert-days --letsencryptdir=$PWD/testing/etc/letsencrypt --domain=helio.loureiro.eng.br
Using directory: /home/helio/DEVEL/letsencrypt-cert-remaining-days/testing/etc/letsencrypt
helio.loureiro.eng.br=50
```

```bash
❯ sudo ./letsencrypt-cert-days
helio.truta.org=51
```
