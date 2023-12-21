data "hund_issues" "standing" {
  standing = true
}

data "hund_issues" "upcoming" {
  upcoming = true
}

output "standing_issues" {
  value = data.hund_issues.standing.issues
}

output "upcoming_issues" {
  value = data.hund_issues.upcoming.issues
}
