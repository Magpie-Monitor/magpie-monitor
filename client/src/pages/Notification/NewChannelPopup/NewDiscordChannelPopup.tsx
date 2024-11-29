import {
  DiscordChannelForm,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import { NewChannelPopupProps } from 'pages/Notification/NewChannelPopup/NewChannelPopup';
import './NewChannelPopup.scss';
import discordIcon from 'assets/discord-icon.png';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { discordWebhookValidation } from 'lib/validators';
import { useToast } from 'providers/ToastProvider/ToastProvider';

const defaultDiscordChannel: DiscordChannelForm = {
  name: '',
  webhookUrl: '',
};

const NewDiscordChannelPopup = ({
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: NewChannelPopupProps) => {
  const [discordChannel, setDiscordChannel] = useState<DiscordChannelForm>(
    defaultDiscordChannel,
  );
  const { showMessage } = useToast();

  const createDiscordChannel = async () => {
    if (discordChannel === defaultDiscordChannel) return;

    try {
      await ManagmentServiceApiInstance.postDiscordChannel(discordChannel);
      showMessage({
        message: 'Notification channel was created',
        type: 'INFO',
      });
    } catch (error) {
      // eslint-disable-next-line no-console
      showMessage({
        message: `Failed to add discord channel ${error}`,
        type: 'ERROR',
      });
    } finally {
      setIsDisplayed(false);
      onSubmit();
    }
  };

  const handleSubmit = () => {
    createDiscordChannel();
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(discordChannel.name) == null &&
      discordWebhookValidation(discordChannel.webhookUrl) == null
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
          icon={<img src={discordIcon} />}
          title="Add new Discord channel"
        />
        <Form
          id="new-channel-form"
          onSubmit={handleSubmit}
          className="new-channel-popup__form"
        >
          <label className="new-channel-popup__form__header">
            Assign human-readable name and webhook for discord channel
            configuration
          </label>
          <LabelInput
            label="Name"
            placeholder={'My discord channel'}
            validationMessage={nonEmptyFieldValidation}
            value={discordChannel.name}
            onChange={(name) =>
              setDiscordChannel((data) => ({ ...data, name }))
            }
          />
          <LabelInput
            label="Webhook Url"
            placeholder={'https://discord.com/api/webhooks/xxxx'}
            value={discordChannel.webhookUrl}
            validationMessage={discordWebhookValidation}
            onChange={(webhookUrl) =>
              setDiscordChannel((data) => ({ ...data, webhookUrl }))
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

export default NewDiscordChannelPopup;
