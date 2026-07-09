# Sandbox Performance Test Cluster

Wrapper scripts for provisioning disposable OCP clusters on AWS using
`openshift-install` IPI, for Dev Sandbox performance testing.

## Prerequisites

- [openshift-install](https://mirror.openshift.com/pub/openshift-v4/clients/ocp/) and `oc` CLI for your target OCP version
- `envsubst` (part of `gettext` — `brew install gettext` on macOS)
- A [pull secret](https://console.redhat.com/openshift/install/pull-secret) saved to `~/.openshift/pull-secret.json`
- AWS credentials with permissions to create EC2, VPC, ELB, Route53, IAM, and S3 resources
- A Route53 hosted zone in the target AWS account (`aws route53 list-hosted-zones` to check)

### Red Hat internal access

If using a Red Hat AWS account via SAML:

```bash
kinit YOUR_KERBEROS@IPA.REDHAT.COM
export $(rh-aws-saml-login --output env <your-aws-account-alias>)
```

Credentials expire after 1 hour — re-run if a long install times out.

## Quick start

```bash
cd infra/
cp config.example config
# Edit config — set CLUSTER_NAME and BASE_DOMAIN at minimum

./create-cluster.sh       # ~40 minutes
```

## Cluster access

```bash
export KUBECONFIG=clusters/<name>/auth/kubeconfig
oc whoami
```

The kubeadmin password is at `clusters/<name>/auth/kubeadmin-password`.

## Running performance tests

See [setup/README.adoc](../setup/README.adoc) for the full guide. Quick start:

```bash
# From the project root
export KUBECONFIG=infra/clusters/<name>/auth/kubeconfig

make dev-deploy-latest      # Install Sandbox operators
go run setup/main.go --users 1 --default 1 --custom 0 --username baseline --testname=baseline
```

## Teardown

```bash
cd infra/
./destroy-cluster.sh
```

Do NOT delete the `clusters/<name>/` directory manually — the destroy script
needs `metadata.json` to clean up AWS resources.

## Notes

- **Route53**: `openshift-install` creates DNS records under `api.<cluster>.<basedomain>`
  and `*.apps.<cluster>.<basedomain>`. Set `BASE_DOMAIN` in your config to match a
  hosted zone in your AWS account.
- **AWS credential expiry**: if the install fails partway due to expired creds, refresh
  them and re-run `./create-cluster.sh` — `openshift-install` can resume from the
  existing cluster directory.
- **Cluster state**: the `clusters/` directory (gitignored) holds kubeconfig, admin
  password, and metadata. Each developer manages their own local state.
