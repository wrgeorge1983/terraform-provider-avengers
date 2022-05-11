terraform {
  required_version = ">= 0.12"
  required_providers {
    avengers = {
      source = "github.com/wrgeorge1983/avengers"
    }
  }

}

provider "avengers" {
  host = "http://localhost:8000"
}

resource "avengers_resource" "this" {
  name   = "some guy"
  alias  = "another guy"
  weapon = "stuff"
}

output "created_avenger" {
  value = avengers_resource.this
}