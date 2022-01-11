# jw-webhook-reconciler

Reconsiler for managing webhook definitions in JW player

## JW API

The JW API has a 60 requests/minute per API token or IP. So if the reconciliation is run very frequently use a separate API key for the process.

There is a page size limit of 100 on the API and this version simply requests 100 and does not try to read any more pages so if you have more that 100 webhook definitions you can not use this version.

## Docker

with a spec JSON file in the current directory run the following remembering to supply a secret

```sh
docker run --rm -v `pwd`:`pwd` -w `pwd` reconsile /home/reconsiler/reconsile --secret= list spec.json
```

There is a pre-build docker image available at jlindstorff/jw-webhook-reconciler:beta1 so you could just run:

```sh
docker run --rm -v `pwd`:`pwd` -w `pwd` jlindstorff/jw-webhook-reconciler:beta2 /home/reconciler/reconcile --secret=<secret> list spec.json
```
