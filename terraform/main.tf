resource "aws_s3_bucket" "b" {
 bucket = "yo-happy-my-tf-test-bucket"
 acl    = "public"


 tags = {
   Name        = "Yo bucket"
   Environment = "Dev"
 }
}
