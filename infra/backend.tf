terraform {
  backend "s3" {
    bucket       = "skale-terraform-state-bucket"                   # bucket name
    key          = "skale/naver-sa-crawling-test/terraform.tfstate" # file name
    region       = "ap-northeast-2"                                 # region
    use_lockfile = true
    encrypt      = true # lock status file encrypt
  }
}