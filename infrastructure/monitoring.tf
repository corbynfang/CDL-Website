# monitoring.tf — CloudWatch alerts for API 500 errors.
#
# How it works end-to-end:
#   1. Gin logs every request: "[GIN] | 500 | 2ms | ... | GET /api/v1/players"
#   2. The metric filter scans the ECS log group for lines matching "| 500 |"
#      and increments a custom CloudWatch metric each time one appears.
#   3. The alarm watches that metric. If >=3 errors occur in any 5-minute window,
#      it fires and publishes a message to the SNS topic.
#   4. SNS delivers the message to your email address.
#
# After terraform apply, AWS sends a confirmation email — click "Confirm subscription"
# before alerts will start arriving.

# ── SNS topic (the notification bus) ────────────────────────────────────────

resource "aws_sns_topic" "api_alerts" {
  name = "${local.prefix}-api-alerts"
  tags = local.tags
}

# ── Email subscription (delivers alarm messages to your inbox) ───────────────

resource "aws_sns_topic_subscription" "alert_email" {
  topic_arn = aws_sns_topic.api_alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# ── CloudWatch metric filter ─────────────────────────────────────────────────
# Scans the ECS log group for Gin log lines containing "| 500 |" and
# counts them into a custom metric called "API500Errors".

resource "aws_cloudwatch_log_metric_filter" "api_500s" {
  name           = "${local.prefix}-api-500-errors"
  log_group_name = "/ecs/${local.prefix}-api"
  pattern        = "\"| 500 |\""

  metric_transformation {
    name      = "API500Errors"
    namespace = "CDLytics/API"
    value     = "1"
    # Each matching log line adds 1 to the metric
  }
}

# ── CloudWatch alarm ─────────────────────────────────────────────────────────
# Fires when 3 or more 500 errors appear within a single 5-minute period.
# Lower this threshold if you want to catch even a single 500.

resource "aws_cloudwatch_metric_alarm" "api_500s" {
  alarm_name          = "${local.prefix}-api-500-errors"
  alarm_description   = "API is returning HTTP 500 errors — check ECS logs at /ecs/${local.prefix}-api"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  metric_name         = "API500Errors"
  namespace           = "CDLytics/API"
  period              = 300 # 5 minutes
  statistic           = "Sum"
  threshold           = 3

  # Missing data means no 500s were logged — treat as OK, not INSUFFICIENT_DATA
  treat_missing_data = "notBreaching"

  alarm_actions = [aws_sns_topic.api_alerts.arn]
  ok_actions    = [aws_sns_topic.api_alerts.arn] # also notify when it recovers

  tags = local.tags
}
