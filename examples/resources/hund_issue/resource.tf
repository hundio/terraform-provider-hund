resource "hund_issue" "example" {
  archive_on_destroy = true

  title         = "a terraform issue"
  body          = "detailed description of the issue"
  component_ids = ["5d72d51f8fbb65b5d3a587e1"]

  # Back-date an Issue
  # began_at = "2023-09-23T13:00:00Z"

  # Make this Issue scheduled
  schedule = {
    starts_at = timeadd("2023-09-23T13:00:00Z", "2h")
    ends_at   = timeadd("2023-09-23T13:00:00Z", "10h")
  }

  # Retrospectively create an Issue
  # updates = [
  #   {
  #     label = "resolved"
  #     body = "description of an already-resolved Issue"
  #   }
  # ]
}
