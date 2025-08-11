# module "ecr_naver_sa_crawler_test" {
#   source          = "git::https://github.com/NHNAD-wooyeon/hyper-infra-modules.git//terraform-aws-ecr"
#   repository_name = "naver-sa-crawler-test"
# }

# module "lambda_naver_sa_crawler_test" {
#   source        = "git::https://github.com/NHNAD-wooyeon/hyper-infra-modules.git//terraform-aws-lambda"
#   function_name = "naver-sa-crawler-test"
#   package_type  = "Image"
#   image_uri     = "${module.ecr_naver_sa_crawler_test.repository_url}:latest"
#   timeout       = 60
#   memory_size   = 128
#   tracing_config = {
#     mode = "PassThrough"
#   }
# }