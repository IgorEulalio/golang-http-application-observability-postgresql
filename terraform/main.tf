resource "aws_s3_bucket" "b" {
 bucket = "dev_s3_bucket"
 acl    = "public" ## test

 tags = {
   Name        = "Environment"
   Environment = "Dev"
 }
}
