# Process for the Horizon Development Team to Make Updates to the Helloworld Service

The instructions for the open-horizon developers that are maintaining this example code are:

- Perform the steps in the [README.md Preconditions](#preconditions) section, **except**:
  - export `HZN_EXCHANGE_URL` to the staging instance
  - Do **not** run `hzn dev service new ...` (use the git files in this directory instead)
  - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Change `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` to be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:

```bash
make publish-all-arches
```

Note: Building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation
