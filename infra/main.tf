module "ecr_naver_sa_crawler_test" {
  source          = "git::https://github.com/NHNAD-wooyeon/hyper-infra-modules.git//terraform-aws-ecr"
  repository_name = "naver-sa-crawler-test"
}