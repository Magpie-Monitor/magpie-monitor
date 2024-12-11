import {
  EditSlackChannelForm,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import './EditChannelPopup.scss';
import slackIcon from 'assets/slack-icon.png';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { EditChannelPopupProps } from './EditChannelPopup';
import { slackWebhookValidation } from 'lib/validators';
import { useToast } from 'providers/ToastProvider/ToastProvider';

interface EditSlackChannelPopupProps extends EditChannelPopupProps {
  id: string;
  name: string;
  webhookUrl: string;
}

const EditSlackChannelPopup = ({
  id,
  name,
  webhookUrl,
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: EditSlackChannelPopupProps) => {
  const [slackChannel, setSlackChannel] = useState<EditSlackChannelForm>({
    id,
    name,
    webhookUrl: '',
  });

  const { showMessage } = useToast();

  const editSlackChannel = async () => {
    if (slackChannel.name === name && slackChannel.webhookUrl === webhookUrl)
      return;

    try {
      await ManagmentServiceApiInstance.editSlackChannel(slackChannel);
      showMessage({
        message: 'Channel updated successfully',
        type: 'INFO',
      });
      onSubmit();
      setIsDisplayed(false);
    } catch (error) {
      showMessage({
        message: `Failed to update channel ${error}`,
        type: 'ERROR',
      });
    }
  };

  const handleSubmit = () => {
    editSlackChannel();
  };

  const validateNonEmptySlackWebhook = (value: string) => {
    if (value === '') {
      return null;
    }

    return slackWebhookValidation(value);
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(slackChannel.name) == null &&
      validateNonEmptySlackWebhook(slackChannel.webhookUrl) == null
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
          icon={<img src={slackIcon} />}
          title="Edit Slack channel"
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
            value={slackChannel.name}
            validationMessage={nonEmptyFieldValidation}
            onChange={(newName) =>
              setSlackChannel((data) => ({ ...data, name: newName }))
            }
          />
          <LabelInput
            label="Webhook Url"
            value={slackChannel.webhookUrl}
            placeholder={webhookUrl}
            validationMessage={validateNonEmptySlackWebhook}
            onChange={(newSlack) =>
              setSlackChannel((data) => ({
                ...data,
                slack: newSlack,
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

export default EditSlackChannelPopup;
