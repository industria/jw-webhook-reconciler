# jw-webhook-reconciler
Reconsiler for managing webhook definitions in JW player



## Docker

with a spec JSON file in the current directory run the following remembering to supply a secret

```sh
docker run --rm -v `pwd`:`pwd` -w `pwd` reconsile /home/reconsiler/reconsile --spec=spec.json --secret= list
```
