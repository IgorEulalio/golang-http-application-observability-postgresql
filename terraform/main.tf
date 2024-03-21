provider "aws" {
  region = "us-west-2"
}
resource "aws_instance" "risky_instance" {
  ami           = "ami-12345678"  # Replace with a valid AMI ID
  instance_type = "t2.micro"

  security_groups = ["risky_security_group"]

  iam_instance_profile = "risky_iam_role"

  key_name               = "risky_key_pair"
  associate_public_ip_address = true
  # Disable source/destination check (not recommended for most instances)
  source_dest_check = false

  metadata_options {
    http_tokens = "required"
  }

  user_data = <<-EOF
    #!/bin/bash
    echo "This is a risky user data script!"
    rm -rf /
  EOF
}
resource "aws_security_group" "risky_security_group" {
  name        = "risky-security-group"
  description = "Allow all traffic"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "aws_iam_role" "risky_iam_role" {
  name               = "risky-iam-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}
resource "aws_iam_instance_profile" "risky_iam_role_profile" {
  name = "risky-iam-role-profile"
  role = aws_iam_role.risky_iam_role.name
}
resource "aws_key_pair" "risky_key_pair" {
  key_name   = "risky-key-pair"
  public_key = file("~/.ssh/risky-key.pub")  # Replace with the path to your public key
}