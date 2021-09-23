# AGC minimal permissions

This CDK project will create two policies that represent minimal permissions for:

* Admins - roles that perform
  * `agc account *`
* Users - roles that peform
  * `agc context *`
  * `agc workflow *`
  * `agc logs *`

To install dependencies:

```bash
npm install
```

To deploy this into your account run:

```bash
cdk deploy
```
