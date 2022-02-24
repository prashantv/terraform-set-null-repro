terraform {
  required_providers {
    tftest = {
      version = "0.1.0"
      source  = "local/prashantv/tftest"
    }
  }
}

resource "tftest_service" "s" {
  name    = "svc1"
  job {
    namespace = "NS1"
    search_tags = {
      "foo" = "bar"
    }
  }
}

