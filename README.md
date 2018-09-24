# eph-music-micro
Suite of microservices supporting my main `eph-music` repo

## Infrastructure as Code
Terraform is used to manage cloud resources

### Google Cloud
[Guide](https://cloud.google.com/community/tutorials/managing-gcp-projects-with-terraform)
Follow the guide in order to stand up / tear down

Make sure to adjust the `GOOGLE_PROJECT` environment variable from the
terraform admin project to the `music` project
