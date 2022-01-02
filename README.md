# jw-webhook-reconciler
Reconsiler for managing webhook definitions in JW player

## JW API

The JW API has a 60 requests/minute per API token or IP. So if the reconciliation is run very frequently use a separate API key for the process. 

There is a page size limit of 100 on the API and this version simply requests 100 and does not try to read any more pages so if you have more that 100 webhook definitions you can not use this version.

## Docker

with a spec JSON file in the current directory run the following remembering to supply a secret

```sh
docker run --rm -v `pwd`:`pwd` -w `pwd` reconsile /home/reconsiler/reconsile --spec=spec.json --secret= list
```
