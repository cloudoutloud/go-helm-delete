# go-helm-delete

Cron job to periodically check and deletes helm releases in set namespace, over a set threshold.

## Environment vars

`HELM_NAMESPACE` Namespace to action in.
`EXEMPT_RELEASES` Release names to skip.
`THRESHOLD_HOURS` Override hours, defaults to 24 hours.

## Example build image

`docker build . -t helm-clean-up:v1.0.1`

## Example deploy k8s resources

`kubectl apply -f k8s/cron-job.yaml`

K8s resources are just examples edit to suit your needs.

