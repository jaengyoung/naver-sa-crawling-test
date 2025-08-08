provider "aws" {
  region  = "ap-northeast-2"
  profile = "skale"

  default_tags {
    tags = {
      Owner = "jaeyoung.lim"
    }
  }
}