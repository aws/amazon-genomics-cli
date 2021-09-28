## Sarek Workflow

Due to the memory requirements of some steps of the sarek workflow you will need to run this on a deployed `bigMemCtx`

```shell
agc context deploy -c bigMemCtx
agc workflow run sarek -c bigMemCtx
```