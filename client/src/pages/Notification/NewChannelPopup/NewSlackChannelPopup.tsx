import {
  ManagmentServiceApiInstance,
  SlackChannelForm,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import { NewChannelPopupProps } from 'pages/Notification/NewChannelPopup/NewChannelPopup';
import './NewChannelPopup.scss';
import slackIcon from 'assets/slack-icon.png';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { slackWebhookValidation } from 'pages/Notification/NotificationTable/SlackTable';

const defaultSlackChannel: SlackChannelForm = {
  name: '',
  webhookUrl: '',
};

const NewSlackChannelPopup = ({
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: NewChannelPopupProps) => {
  const [slackChannel, setSlackChannel] =
    useState<SlackChannelForm>(defaultSlackChannel);

  const createSlackChannel = async () => {
    if (slackChannel === defaultSlackChannel) return;

    try {
      await ManagmentServiceApiInstance.postSlackChannel(slackChannel);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error posting slack channels: ', error);
    } finally {
      setIsDisplayed(false);
      onSubmit();
    }
  };

  const handleSubmit = () => {
    createSlackChannel();
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(slackChannel.name) == null &&
      slackWebhookValidation(slackChannel.webhookUrl) == null
    );
  };

  return (
    <OverlayComponent
      isDisplayed={isDisplayed}
      onClose={() => {
        setIsDisplayed(false);
      }}
    >
      <div className="new-channel-popup">
        <NewChannelPopupHeader
          icon={<img src={slackIcon} />}
          title="Add new Slack channel"
        />
        <Form
          id="new-channel-form"
          onSubmit={handleSubmit}
          className="new-channel-popup__form"
        >
          <label className="new-channel-popup__form__header">
            Assign human-readable name and webhook for slack channel
            configuration
          </label>
          <LabelInput
            value={slackChannel.name}
            label="Name"
            placeholder={'My slack channel'}
            validationMessage={nonEmptyFieldValidation}
            onChange={(name) => setSlackChannel((data) => ({ ...data, name }))}
          />

          <LabelInput
            value={slackChannel.webhookUrl}
            label="Webhook Url"
            placeholder={'https://hooks.slack.com/services/xxx'}
            validationMessage={slackWebhookValidation}
            onChange={(webhookUrl) =>
              setSlackChannel((data) => ({ ...data, webhookUrl }))
            }
          />
        </Form>
        <div className="new-channel-popup__buttons">
          <ActionButton
            onClick={handleSubmit}
            description="Submit"
            disabled={!isFormValid()}
            color={ActionButtonColor.GREEN}
          />
          <ActionButton
            onClick={() => setIsDisplayed(false)}
            description="Cancel"
            color={ActionButtonColor.RED}
          />
        </div>
      </div>
    </OverlayComponent>
  );
};

export default NewSlackChannelPopup;
