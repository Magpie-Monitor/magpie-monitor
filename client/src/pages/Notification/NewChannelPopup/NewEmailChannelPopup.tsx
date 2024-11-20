import {
  EmailChannelForm,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import { NewChannelPopupProps } from 'pages/Notification/NewChannelPopup/NewChannelPopup';
import './NewChannelPopup.scss';
import emailIcon from 'assets/mail-icon.svg';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { emailValidation } from 'lib/validators';

const defaultEmailChannel: EmailChannelForm = {
  name: '',
  email: '',
};

const NewEmailChannelPopup = ({
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: NewChannelPopupProps) => {
  const [emailChannel, setEmailChannel] =
    useState<EmailChannelForm>(defaultEmailChannel);

  const createEmailChannel = async () => {
    if (emailChannel === defaultEmailChannel) return;

    try {
      await ManagmentServiceApiInstance.postEmailChannel(emailChannel);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error posting slack channels: ', error);
    } finally {
      setIsDisplayed(false);
      onSubmit();
    }
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(emailChannel.name) == null &&
      emailValidation(emailChannel.email) == null
    );
  };

  const handleSubmit = () => {
    createEmailChannel();
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
          icon={<img src={emailIcon} />}
          title="Add new Email channel"
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
            label="Name"
            value={emailChannel.name}
            placeholder={'My email channel'}
            validationMessage={nonEmptyFieldValidation}
            onChange={(name) => setEmailChannel((data) => ({ ...data, name }))}
          />

          <LabelInput
            value={emailChannel.email}
            label="Email"
            placeholder={'contact@company.com'}
            validationMessage={emailValidation}
            onChange={(email) =>
              setEmailChannel((data) => ({ ...data, email }))
            }
          />
        </Form>
        <div className="new-channel-popup__buttons">
          <ActionButton
            onClick={handleSubmit}
            disabled={!isFormValid()}
            description="Submit"
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

export default NewEmailChannelPopup;
