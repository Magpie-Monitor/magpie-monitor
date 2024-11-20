const DISCORD_WEBHOOK_FORMAT =
  'https://discord.com/api/webhooks/<webhookId>/<authToken>';

export const discordWebhookValidation = (webhook: string) => {
  const discordWebhookRegex =
    /^https:\/\/discord\.com\/api\/webhooks\/\d+\/[a-zA-Z0-9_-]+$/;
  if (!discordWebhookRegex.test(webhook)) {
    return `Discord webhook must be in the following format ${DISCORD_WEBHOOK_FORMAT}`;
  }

  return null;
};

const SLACK_WEBHOOK_FORMAT =
  'https://hooks.slack.com/services/<workspaceId>/<botIdentifier>/<authToken>';

export const slackWebhookValidation = (webhook: string) => {
  const slackWebhookRegex =
    /^https:\/\/hooks\.slack\.com\/services\/[A-Z0-9]+\/[A-Z0-9]+\/[a-zA-Z0-9]+$/;
  if (!slackWebhookRegex.test(webhook)) {
    return `Slack webhook must be in the following format ${SLACK_WEBHOOK_FORMAT}`;
  }

  return null;
};

export const emailValidation = (value: string) => {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  if (!emailRegex.test(value)) {
    return 'Invalid email';
  }
  return null;
};
