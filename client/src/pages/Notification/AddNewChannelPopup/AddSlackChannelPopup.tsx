import { ManagmentServiceApiInstance } from 'api/managment-service';
import Popup from 'components/Popup/Popup';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import { AddNewChannelPopupProps } from 'pages/Notification/AddNewChannelPopup/AddNewChannelPopup';
import './AddSlackChannelPopup.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import slackIcon from 'assets/slack-icon.png';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';

export interface SlackChannel {
  name: string;
  webhookUrl: string;
}

const defaultSlackChannel: SlackChannel = {
  name: '',
  webhookUrl: '',
};

const AddSlackChannelPopup = ({
  isDisplayed,
  setIsDisplayed,
}: AddNewChannelPopupProps) => {
  const [slackChannel, setSlackChannel] =
    useState<SlackChannel>(defaultSlackChannel);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const name = event.target.name;
    const value = event.target.value;
    setSlackChannel((inputSlackChannel) => ({
      ...inputSlackChannel,
      [name]: value,
    }));
  };

  const postSlackChannel = async () => {
    console.log('working');
    if (slackChannel === defaultSlackChannel) return;
    return;
    try {
      await ManagmentServiceApiInstance.postSlackChannel(slackChannel);
    } catch (error) {
      console.error('Error posting slack channels: ', error);
    } finally {
      setIsDisplayed(false);
    }
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    postSlackChannel();
  };

  return (
    <Popup isDisplayed={isDisplayed} setIsDisplayed={setIsDisplayed}>
      <div className="add-slack-channel-popup">
        <HeaderWithIcon
          icon={<img src={slackIcon} />}
          title="Add new Slack channel"
        />
        <Form
          id="add-slack-channel-form"
          onSubmit={handleSubmit}
          className="add-slack-channel-popup__form"
        >
          <label className="add-slack-channel-popup__form__header">
            Slack properties
          </label>
          <label className="add-slack-channel-popup__form__row">
            Name
            <input
              type="text"
              value={slackChannel.name}
              onChange={handleChange}
            />
          </label>
          <label className="add-slack-channel-popup__form__row">
            Webhook url
            <input
              type="text"
              value={slackChannel.webhookUrl}
              onChange={handleChange}
            />
          </label>
          <button form="add-slack-channel-form" type="submit">
            <ActionButton
              onClick={() => {}}
              description="SUBMIT"
              color={ActionButtonColor.GREEN}
            />
          </button>
        </Form>
      </div>
    </Popup>
  );
};

export default AddSlackChannelPopup;
