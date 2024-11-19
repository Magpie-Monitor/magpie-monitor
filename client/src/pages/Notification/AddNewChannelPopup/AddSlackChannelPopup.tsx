import { ManagmentServiceApiInstance } from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
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

  const handleChange = (event: React.FormEvent<HTMLInputElement>) => {
    const name = (event.target as HTMLInputElement).name;
    const value = (event.target as HTMLInputElement).value;
    console.log(`name, value [${name}]: ${value}`);
    setSlackChannel((inputSlackChannel) => ({
      ...inputSlackChannel,
      [name]: value,
    }));
  };

  const postSlackChannel = async () => {
    if (slackChannel === defaultSlackChannel) return;

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
    <OverlayComponent
      isDisplayed={isDisplayed}
      onClose={() => {
        setIsDisplayed(false);
      }}
    >
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
              type="name"
              name="name"
              value={slackChannel.name}
              onChange={handleChange}
            />
          </label>
          <label className="add-slack-channel-popup__form__row">
            Webhook url
            <input
              type="webhookUrl"
              name="webhookUrl"
              value={slackChannel.webhookUrl}
              onChange={handleChange}
            />
          </label>
        </Form>
        <ActionButton
          onClick={()=>{}}
          description="Submit" 
          color={ActionButtonColor.GREEN}
        ></ActionButton>
      </div>
    </OverlayComponent>
  );
};

export default AddSlackChannelPopup;
