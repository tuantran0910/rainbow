from datetime import datetime

import dagster as dg
from dagster_slack import SlackResource
from slack_sdk.web.client import WebClient

from shared.constants import DAGIT_BASE_URL
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS
from shared.constants import SLACK_BOT_TOKEN
from shared.constants import SLACK_CHANNEL


def build_and_post_slack_message_for_asset_run_failure(
    context: dg.RunFailureSensorContext, slack: SlackResource
) -> None:
    """
    Build and post a Slack message for asset run failures.

    Args:
        context: The run failure sensor context containing run information.
        slack: The Slack resource for posting messages.
    """
    # Check if run will be retried — skip alert if so
    tags = context.dagster_run.tags or {}
    if tags.get("dagster/will_retry", "false") == "true":
        context.log.info(f"Run {context.dagster_run.run_id} will retry, skipping alert.")
        return

    dagster_run: dg.DagsterRun = context.dagster_run
    job_name: str = dagster_run.job_name
    run_id: str = dagster_run.run_id

    # Get run failure information
    failure_message = None
    job_failure_data = getattr(context.failure_event, "event_specific_data", None)
    if job_failure_data and hasattr(job_failure_data, "first_step_failure_event"):
        step_failure_event = job_failure_data.first_step_failure_event
        step_failure_data = getattr(step_failure_event, "event_specific_data", None)
        if step_failure_data and hasattr(step_failure_data, "error"):
            error_info = step_failure_data.error
            if error_info and error_info.message:
                failure_message = f"{error_info.message.strip()}"
                if error_info.cause and error_info.cause.message:
                    failure_message += f"\nCaused by: {error_info.cause.message.strip()}"
    else:
        failure_message = "Asset run failed"

    # Build asset list
    asset_keys = []
    if hasattr(dagster_run, "asset_selection") and dagster_run.asset_selection:
        asset_keys = [key.to_user_string() for key in dagster_run.asset_selection]

    assets_text = "N/A"
    if asset_keys:
        assets_text = ", ".join(f"`{key}`" for key in asset_keys[:5])
        if len(asset_keys) > 5:
            assets_text += f" and {len(asset_keys) - 5} more"

    # Timestamps
    current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S UTC")
    dagster_ui_link = f"{DAGIT_BASE_URL}/runs/{run_id}"

    # Build the Slack message blocks
    slack_message_body = [
        {
            "color": "#FF0000",
            "blocks": [
                {
                    "type": "header",
                    "text": {"type": "plain_text", "text": f"🚨 Job {job_name} failed"},
                },
                {
                    "type": "section",
                    "text": {
                        "type": "mrkdwn",
                        "text": f"*📦 Assets:* {assets_text}\n*🛠 Job:* `{job_name}`\n*🆔 Run ID:* `{run_id}`\n*⏰ Time:* `{current_time}`",
                    },
                },
                {
                    "type": "section",
                    "text": {
                        "type": "mrkdwn",
                        "text": f"*🐛 Error Details:*\n```{failure_message}```",
                    },
                },
                {
                    "type": "actions",
                    "elements": [
                        {
                            "type": "button",
                            "text": {"type": "plain_text", "text": "🔗 View in Dagster UI"},
                            "url": dagster_ui_link,
                            "style": "primary",
                        }
                    ],
                },
            ],
        }
    ]

    # Send to Slack
    slack_client: WebClient = slack.get_client()
    slack_client.chat_postMessage(channel=SLACK_CHANNEL, attachments=slack_message_body)


@dg.run_failure_sensor(
    name="failed_asset_run_sensor",
    tags=DAGSTER_TAGS,
    metadata=DAGSTER_METADATA,
)
def failed_asset_run_sensor(context: dg.RunFailureSensorContext, slack: SlackResource) -> None:
    """
    Sensor that triggers when an asset run fails.
    It sends a Slack notification with details about the failed run.

    Args:
        context: The run failure sensor context.
        slack: The slack resource.
    """
    build_and_post_slack_message_for_asset_run_failure(context=context, slack=slack)


# Combine all monitoring definitions
defs = dg.Definitions(
    resources={
        "slack": SlackResource(
            token=SLACK_BOT_TOKEN,
        ),
    },
    sensors=[failed_asset_run_sensor],
)
