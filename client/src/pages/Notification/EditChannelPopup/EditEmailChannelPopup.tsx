import {
  EditEmailChannelForm,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import { useState } from 'react';
import { Form } from 'react-router-dom';
import './EditChannelPopup.scss';
import emailIcon from 'assets/mail-icon.svg';
import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import LabelInput, {
  nonEmptyFieldValidation,
} from 'components/LabelInput/LabelInput';
import NewChannelPopupHeader from 'pages/Notification/NewChannelPopupHeader/NewChannelPopupHeader';
import { EditChannelPopupProps } from './EditChannelPopup';
import { emailValidation } from 'lib/validators';
import { useToast } from 'providers/ToastProvider/ToastProvider';

interface EditEmailChannelPopupProps extends EditChannelPopupProps {
  id: string;
  name: string;
  email: string;
}

const EditEmailChannelPopup = ({
  id,
  name,
  email,
  isDisplayed,
  setIsDisplayed,
  onSubmit,
}: EditEmailChannelPopupProps) => {
  const [emailChannel, setEmailChannel] = useState<EditEmailChannelForm>({
    id,
    name,
    email,
  });
  const { showMessage } = useToast();

  const editEmailChannel = async () => {
    if (emailChannel.name === name && emailChannel.email === email) return;

    try {
      await ManagmentServiceApiInstance.editEmailChannel(emailChannel);
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
    editEmailChannel();
  };

  const isFormValid = () => {
    return (
      nonEmptyFieldValidation(emailChannel.name) == null &&
      emailValidation(emailChannel.email) == null
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
          icon={<img src={emailIcon} />}
          title="Edit Email channel"
        />
        <Form
          id="edit-channel-form"
          onSubmit={handleSubmit}
          className="edit-channel-popup__form"
        >
          <label className="edit-channel-popup__form__header">
            Update channel with new name or new receiver email.
          </label>
          <LabelInput
            label="Name"
            validationMessage={nonEmptyFieldValidation}
            value={emailChannel.name}
            onChange={(newName) =>
              setEmailChannel((data) => ({ ...data, name: newName }))
            }
          />
          <LabelInput
            label="Email"
            validationMessage={emailValidation}
            value={emailChannel.email}
            onChange={(newEmail) =>
              setEmailChannel((data) => ({
                ...data,
                email: newEmail,
              }))
            }
          />
        </Form>
        <div className="edit-channel-popup__buttons">
          <ActionButton
            onClick={handleSubmit}
            description="Save"
            disabled={!isFormValid()}
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

export default EditEmailChannelPopup;
