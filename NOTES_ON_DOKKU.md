For the git revision to be available at build time in dokku's environment it is essential to configure dokku to keep the .git directory:

```bash
# keep the .git directory during builds
dokku git:set <app> keep-git-dir true
```

Additionally, the config file needs to be mounted via:
```bash
# mount the config file (/app is the workdir according to Dockerfile)
dokku storage:mount change-check /apps/change-check/change-check.config.yaml:/app/change-check.config.yaml
```