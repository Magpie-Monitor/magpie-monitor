import {
  EditDiscordChannelForm,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import './EditChannelPopup.scss';
import discordIcon from 'assets/discord-icon.png';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { EditChannelPopupProps } from './EditChannelPopup';
import { discordWebhookValidation } from 'lib/validators';
import { useToast } from 'providers/ToastProvider/ToastProvider';

interface EditDiscordChannelPopupProps extends EditChannelPopupProps {
  id: string;
  name: string;
  webhookUrl: string;
}

const EditDiscordChannelPopup = ({
  id,
  name,
  webhookUrl,
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: EditDiscordChannelPopupProps) => {
  const [discordChannel, setDiscordChannel] = useState<EditDiscordChannelForm>({
    id,
    name,
    webhookUrl: '',
  });

  const { showMessage } = useToast();

  const editDiscordChannel = async () => {
    if (
      discordChannel.name === name &&
      discordChannel.webhookUrl === webhookUrl
    )
      return;

    try {
      await ManagmentServiceApiInstance.editDiscordChannel(discordChannel);
      setIsDisplayed(false);
      onSubmit();
      showMessage({ message: 'Channel was updated', type: 'INFO' });
    } catch (error) {
      showMessage({
        message: `Failed to update channel ${error}`,
        type: 'ERROR',
      });
    }
  };

  const handleSubmit = () => {
    editDiscordChannel();
  };

  const validateNonEmptyDiscordWebhook = (value: string) => {
    if (value === '') {
      return null;
    }

    return discordWebhookValidation(value);
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(discordChannel.name) == null &&
      validateNonEmptyDiscordWebhook(discordChannel.webhookUrl) == null
    );
  };

  return (
    <OverlayComponent
      isDisplayed={isDisplayed}
      onClose={() => {
        setIsDisplayed(false);
      }}
    >
      <div className="edit-channel-popup">
        <NewChannelPopupHeader
          icon={<img src={discordIcon} />}
          title="Edit Discord channel"
        />
        <Form
          id="edit-channel-form"
          onSubmit={handleSubmit}
          className="edit-channel-popup__form"
        >
          <label className="edit-channel-popup__form__header">
            Update channel with new name or new webhook url.
          </label>
          <LabelInput
            label="Name"
            value={discordChannel.name}
            validationMessage={nonEmptyFieldValidation}
            onChange={(newName) =>
              setDiscordChannel((data) => ({ ...data, name: newName }))
            }
          />

          <LabelInput
            label="Webhook Url"
            value={discordChannel.webhookUrl}
            placeholder={webhookUrl}
            validationMessage={validateNonEmptyDiscordWebhook}
            onChange={(newWebhookUrl) =>
              setDiscordChannel((data) => ({
                ...data,
                webhookUrl: newWebhookUrl,
              }))
            }
          />
        </Form>
        <div className="edit-channel-popup__buttons">
          <ActionButton
            disabled={!isFormValid()}
            onClick={handleSubmit}
            description="Save"
            color={ActionButtonColor.GREEN}
          ></ActionButton>
          <ActionButton
            onClick={() => setIsDisplayed(false)}
            description="Cancel"
            color={ActionButtonColor.RED}
          ></ActionButton>
        </div>
      </div>
    </OverlayComponent>
  );
};

export default EditDiscordChannelPopup;
